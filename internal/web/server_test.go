package web

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func testConfig(t *testing.T, port int) Config {
	t.Helper()
	return Config{
		Host:         "127.0.0.1",
		Port:         port,
		TemplatesDir: t.TempDir(),
		Out:          &strings.Builder{},
	}
}

func TestServePortAlreadyInUse(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = Serve(ctx, testConfig(t, port))
	if err == nil || !strings.Contains(err.Error(), strconv.Itoa(port)) || !strings.Contains(err.Error(), "already in use") {
		t.Fatalf("err = %v, want a message naming port %d as already in use", err, port)
	}
}

func TestServeShutsDownOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	cfg := testConfig(t, 0) // port 0: let the OS pick a free one
	go func() { done <- Serve(ctx, cfg) }()

	// Give Serve a moment to bind before cancelling.
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Serve returned %v after shutdown", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Serve did not shut down after ctx was cancelled")
	}
}

func TestServeOpensBrowserOnlyForLoopback(t *testing.T) {
	cases := []struct {
		name string
		host string
	}{
		{"loopback", "127.0.0.1"},
		{"non-loopback", "0.0.0.0"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var opened string
			out := &strings.Builder{}
			cfg := testConfig(t, 0)
			cfg.Host = tc.host
			cfg.Out = out
			cfg.OpenBrowser = func(url string) error {
				opened = url
				return nil
			}

			ctx, cancel := context.WithCancel(context.Background())
			done := make(chan error, 1)
			go func() { done <- Serve(ctx, cfg) }()
			time.Sleep(50 * time.Millisecond)
			cancel()
			<-done

			if tc.host == "127.0.0.1" {
				if opened == "" {
					t.Error("expected OpenBrowser to be called for a loopback host")
				}
				if out.String() != "" {
					t.Errorf("expected no warning for a loopback host, got %q", out.String())
				}
			} else {
				if opened != "" {
					t.Errorf("OpenBrowser was called for a non-loopback host: %q", opened)
				}
				if !strings.Contains(out.String(), "unauthenticated") {
					t.Errorf("expected a one-time warning for a non-loopback host, got %q", out.String())
				}
			}
		})
	}
}

func TestStaticAssetRoutesServeVendoredFiles(t *testing.T) {
	cfg := testConfig(t, 0)
	for _, dir := range []string{"fonts", "js"} {
		if err := os.MkdirAll(cfg.TemplatesDir+"/"+dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	fixtures := map[string]string{
		"fonts/JetBrainsMono-Regular.woff2": "regular-stub",
		"fonts/JetBrainsMono-Medium.woff2":  "medium-stub",
		"js/cytoscape.min.js":               "cytoscape-stub",
	}
	for rel, content := range fixtures {
		if err := os.WriteFile(cfg.TemplatesDir+"/"+rel, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	for _, route := range []string{
		"/fonts/JetBrainsMono-Regular.woff2",
		"/fonts/JetBrainsMono-Medium.woff2",
		"/js/cytoscape.min.js",
	} {
		req := httptest.NewRequest("GET", route, nil)
		rec := httptest.NewRecorder()
		Routes(cfg).ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("%s: status = %d, want 200", route, rec.Code)
		}
	}
}
