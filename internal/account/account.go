// Package account implements BreakLab's optional hosted identity layer:
// signup/login, sessions, API tokens, and the data an opted-in public
// profile reads. It is entirely separate from the local, unauthenticated
// `breaklab serve` UI — nothing in this package runs unless a hosted
// deployment explicitly constructs a Store.
package account

import (
	"errors"
	"regexp"
	"time"

	"breaklab/internal/progress"
)

// usernamePattern matches this project's existing lab/variant id rule
// (internal/catalog/validate.go, internal/docker/docker.go): lowercase
// letters, digits, and hyphens only, so a username is always safe to
// drop straight into a URL path (/u/<name>).
var usernamePattern = regexp.MustCompile(`^[a-z0-9-]+$`)

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrInvalidUsername    = errors.New("username must match ^[a-z0-9-]+$")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrNotFound           = errors.New("account not found")
)

// User is one hosted account.
type User struct {
	ID            int64
	Email         string
	Username      string
	PublicProfile bool
	CreatedAt     time.Time
}

// Store is the hosted account persistence layer. Every lookup method is
// scoped by the caller-supplied credential/session/token it's given —
// it never returns or mutates another account's data.
type Store interface {
	// CreateUser registers a new account. Returns ErrEmailTaken,
	// ErrUsernameTaken, or ErrInvalidUsername without creating anything
	// on a conflict or malformed username.
	CreateUser(email, username, password string) (User, error)

	// Authenticate verifies email/password and returns the matching
	// User, or ErrInvalidCredentials for any mismatch (including an
	// unknown email) — the error never reveals which part was wrong.
	Authenticate(email, password string) (User, error)

	// CreateSession starts a new session for userID, returning the
	// plaintext token to set as a cookie; only its hash is stored.
	CreateSession(userID int64) (token string, expiresAt time.Time, err error)

	// UserBySession resolves a session cookie token to its account, or
	// ErrNotFound if the token is unknown, revoked, or expired.
	UserBySession(token string) (User, error)

	// DeleteSession invalidates token (logout). Deleting an
	// already-invalid token is a no-op, not an error.
	DeleteSession(token string) error

	// IssueToken generates a new long-lived API token for userID,
	// replacing any token issued before it, and returns the plaintext
	// token; only its hash is stored.
	IssueToken(userID int64) (token string, err error)

	// RevokeToken clears userID's current API token, if any.
	RevokeToken(userID int64) error

	// UserByToken resolves a bearer API token to its account, or
	// ErrNotFound if the token is unknown or has been revoked.
	UserByToken(token string) (User, error)

	// SetPublicProfile turns userID's /u/<username> page on or off.
	SetPublicProfile(userID int64, enabled bool) error

	// UserByUsername looks up an account by its public username,
	// regardless of its PublicProfile setting — callers check that
	// field themselves to decide whether to render or 404.
	UserByUsername(username string) (User, error)

	// DeleteUser permanently removes userID's account, every session,
	// its API token, and every synced attempt belonging to it.
	DeleteUser(userID int64) error

	// PushAttempts stores attempts for userID, ignoring any whose ID is
	// already recorded for that account — idempotent, so re-pushing an
	// already-synced attempt is a no-op rather than a duplicate.
	PushAttempts(userID int64, attempts []progress.Attempt) error

	// PullAttempts returns every attempt synced for userID.
	PullAttempts(userID int64) ([]progress.Attempt, error)

	// Close releases the underlying database connection.
	Close() error
}
