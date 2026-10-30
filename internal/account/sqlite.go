package account

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"breaklab/internal/progress"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite" // registers the "sqlite" database/sql driver
)

// sessionTTL is how long a session cookie stays valid after login.
// Fixed rather than sliding for v1 (design.md Open Questions) — doesn't
// change the observable contract, just how often someone has to log
// back in.
const sessionTTL = 30 * 24 * time.Hour

const schema = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT NOT NULL UNIQUE,
	username TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	public_profile INTEGER NOT NULL DEFAULT 0,
	api_token_hash TEXT,
	created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
	token_hash TEXT PRIMARY KEY,
	user_id INTEGER NOT NULL,
	expires_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS synced_attempts (
	user_id INTEGER NOT NULL,
	id TEXT NOT NULL,
	lab TEXT NOT NULL,
	variant TEXT NOT NULL,
	started_at TIMESTAMP NOT NULL,
	solved_at TIMESTAMP,
	failed_checks INTEGER NOT NULL,
	hints_used INTEGER NOT NULL,
	postmortem TEXT NOT NULL,
	status TEXT NOT NULL,
	hypothesis TEXT NOT NULL,
	PRIMARY KEY (user_id, id)
);
`

type sqliteStore struct {
	db *sql.DB
}

// Open opens (creating if needed) a SQLite-backed Store at path,
// applying the schema. Closing the returned Store closes the
// underlying database connection.
func Open(path string) (Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("applying schema to %s: %w", path, err)
	}
	return &sqliteStore{db: db}, nil
}

func (s *sqliteStore) Close() error { return s.db.Close() }

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func newToken() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

func (s *sqliteStore) CreateUser(email, username, password string) (User, error) {
	if !usernamePattern.MatchString(username) {
		return User{}, ErrInvalidUsername
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	now := time.Now().UTC()
	res, err := s.db.Exec(
		`INSERT INTO users (email, username, password_hash, created_at) VALUES (?, ?, ?, ?)`,
		email, username, string(hash), now,
	)
	if err != nil {
		switch {
		case isUniqueViolation(err, "users.email"):
			return User{}, ErrEmailTaken
		case isUniqueViolation(err, "users.username"):
			return User{}, ErrUsernameTaken
		}
		return User{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return User{}, err
	}
	return User{ID: id, Email: email, Username: username, CreatedAt: now}, nil
}

func (s *sqliteStore) Authenticate(email, password string) (User, error) {
	var u User
	var hash string
	var publicProfile int
	row := s.db.QueryRow(
		`SELECT id, email, username, password_hash, public_profile, created_at FROM users WHERE email = ?`,
		email,
	)
	if err := row.Scan(&u.ID, &u.Email, &u.Username, &hash, &publicProfile, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrInvalidCredentials
		}
		return User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return User{}, ErrInvalidCredentials
	}
	u.PublicProfile = publicProfile != 0
	return u, nil
}

func (s *sqliteStore) userByID(id int64) (User, error) {
	var u User
	var publicProfile int
	row := s.db.QueryRow(
		`SELECT id, email, username, public_profile, created_at FROM users WHERE id = ?`,
		id,
	)
	if err := row.Scan(&u.ID, &u.Email, &u.Username, &publicProfile, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}
	u.PublicProfile = publicProfile != 0
	return u, nil
}

func (s *sqliteStore) CreateSession(userID int64) (string, time.Time, error) {
	token, err := newToken()
	if err != nil {
		return "", time.Time{}, err
	}
	expiresAt := time.Now().UTC().Add(sessionTTL)
	if _, err := s.db.Exec(
		`INSERT INTO sessions (token_hash, user_id, expires_at) VALUES (?, ?, ?)`,
		hashToken(token), userID, expiresAt,
	); err != nil {
		return "", time.Time{}, err
	}
	return token, expiresAt, nil
}

func (s *sqliteStore) UserBySession(token string) (User, error) {
	var userID int64
	var expiresAt time.Time
	row := s.db.QueryRow(
		`SELECT user_id, expires_at FROM sessions WHERE token_hash = ?`,
		hashToken(token),
	)
	if err := row.Scan(&userID, &expiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}
	if time.Now().UTC().After(expiresAt) {
		return User{}, ErrNotFound
	}
	return s.userByID(userID)
}

func (s *sqliteStore) DeleteSession(token string) error {
	_, err := s.db.Exec(`DELETE FROM sessions WHERE token_hash = ?`, hashToken(token))
	return err
}

func (s *sqliteStore) IssueToken(userID int64) (string, error) {
	token, err := newToken()
	if err != nil {
		return "", err
	}
	if _, err := s.db.Exec(`UPDATE users SET api_token_hash = ? WHERE id = ?`, hashToken(token), userID); err != nil {
		return "", err
	}
	return token, nil
}

func (s *sqliteStore) RevokeToken(userID int64) error {
	_, err := s.db.Exec(`UPDATE users SET api_token_hash = NULL WHERE id = ?`, userID)
	return err
}

func (s *sqliteStore) UserByToken(token string) (User, error) {
	var id int64
	row := s.db.QueryRow(`SELECT id FROM users WHERE api_token_hash = ?`, hashToken(token))
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}
	return s.userByID(id)
}

func (s *sqliteStore) SetPublicProfile(userID int64, enabled bool) error {
	_, err := s.db.Exec(`UPDATE users SET public_profile = ? WHERE id = ?`, boolToInt(enabled), userID)
	return err
}

func (s *sqliteStore) UserByUsername(username string) (User, error) {
	var id int64
	row := s.db.QueryRow(`SELECT id FROM users WHERE username = ?`, username)
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}
	return s.userByID(id)
}

func (s *sqliteStore) DeleteUser(userID int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, stmt := range []string{
		`DELETE FROM sessions WHERE user_id = ?`,
		`DELETE FROM synced_attempts WHERE user_id = ?`,
		`DELETE FROM users WHERE id = ?`,
	} {
		if _, err := tx.Exec(stmt, userID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *sqliteStore) PushAttempts(userID int64, attempts []progress.Attempt) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, a := range attempts {
		if _, err := tx.Exec(
			`INSERT OR IGNORE INTO synced_attempts
				(user_id, id, lab, variant, started_at, solved_at, failed_checks, hints_used, postmortem, status, hypothesis)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			userID, a.ID, a.Lab, a.Variant, a.StartedAt, a.SolvedAt, a.FailedChecks, a.HintsUsed, a.Postmortem, a.Status, a.Hypothesis,
		); err != nil {
			return fmt.Errorf("pushing attempt %s: %w", a.ID, err)
		}
	}
	return tx.Commit()
}

func (s *sqliteStore) PullAttempts(userID int64) ([]progress.Attempt, error) {
	rows, err := s.db.Query(
		`SELECT id, lab, variant, started_at, solved_at, failed_checks, hints_used, postmortem, status, hypothesis
		FROM synced_attempts WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []progress.Attempt
	for rows.Next() {
		var a progress.Attempt
		var solvedAt sql.NullTime
		if err := rows.Scan(&a.ID, &a.Lab, &a.Variant, &a.StartedAt, &solvedAt, &a.FailedChecks, &a.HintsUsed, &a.Postmortem, &a.Status, &a.Hypothesis); err != nil {
			return nil, err
		}
		if solvedAt.Valid {
			t := solvedAt.Time
			a.SolvedAt = &t
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// isUniqueViolation reports whether err is a UNIQUE constraint failure
// on the given "table.column" — modernc.org/sqlite surfaces this as a
// plain error whose message names the column, so match on that rather
// than a driver-specific error type.
func isUniqueViolation(err error, column string) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed: "+column)
}
