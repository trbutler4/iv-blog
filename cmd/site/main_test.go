package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMeganAuthentication(t *testing.T) {
	auth, err := newMeganAuth()
	if err != nil {
		t.Fatal(err)
	}
	for _, liveReload := range []bool{false, true} {
		h := devHandler{hub: newReloadHub(), liveReload: liveReload, auth: auth}
		for _, path := range []string{"/megan", "/megan/", "/megan/index.html", "/megan/nested/file", "/%6degan/", "/articles/../megan/", "//megan/"} {
			for _, method := range []string{http.MethodGet, http.MethodHead} {
				r := httptest.NewRequest(method, path, nil)
				w := httptest.NewRecorder()
				h.ServeHTTP(w, r)
				if w.Code != http.StatusUnauthorized || w.Header().Get("WWW-Authenticate") != "" {
					t.Fatalf("%s %s: expected password form without browser auth challenge, got %d", method, path, w.Code)
				}
				if method == http.MethodGet && (!strings.Contains(w.Body.String(), `name="password"`) || strings.Contains(w.Body.String(), `name="username"`)) {
					t.Fatal("expected password-only form")
				}
				if w.Header().Get("Cache-Control") != "no-store" {
					t.Fatalf("%s: missing no-store", path)
				}
			}
		}
		r := httptest.NewRequest(http.MethodPost, "/megan/login", strings.NewReader("password=wrong-password"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != http.StatusUnauthorized || len(w.Result().Cookies()) != 0 {
			t.Fatalf("wrong password accepted: %d", w.Code)
		}
	}
	for _, value := range []string{"forged", auth.token(time.Now().Add(-time.Hour)), auth.token(time.Now().Add(time.Hour)) + "tampered"} {
		r := httptest.NewRequest(http.MethodGet, "/megan/", nil)
		r.AddCookie(&http.Cookie{Name: sessionCookie, Value: value})
		w := httptest.NewRecorder()
		devHandler{auth: auth}.ServeHTTP(w, r)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("invalid session accepted: %s", value)
		}
	}
}

func TestMeganAuthenticatedPage(t *testing.T) {
	password := os.Getenv("MEGAN_TEST_PASSWORD")
	if password == "" {
		t.Skip("set MEGAN_TEST_PASSWORD to verify the configured password")
	}
	auth, err := newMeganAuth()
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	for _, relative := range []string{"megan", "megan/what-i-love-about-you"} {
		dir := filepath.Join(root, relative)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("<body>Private message</body>"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, liveReload := range []bool{false, true} {
		h := devHandler{hub: newReloadHub(), liveReload: liveReload, root: root, auth: auth}
		for _, path := range []string{"/megan/", "/megan/what-i-love-about-you/"} {
			r := httptest.NewRequest(http.MethodGet, path, nil)
			unauthorized := httptest.NewRecorder()
			h.ServeHTTP(unauthorized, r)
			if unauthorized.Code != http.StatusUnauthorized || strings.Contains(unauthorized.Body.String(), "Private message") {
				t.Fatalf("private content exposed at %s", path)
			}
			form := url.Values{"password": {password}, "next": {path}}
			login := httptest.NewRequest(http.MethodPost, "https://example.com/megan/login", strings.NewReader(form.Encode()))
			login.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			response := httptest.NewRecorder()
			h.ServeHTTP(response, login)
			if response.Code != http.StatusSeeOther || response.Header().Get("Location") != path {
				t.Fatalf("login failed: %d %s", response.Code, response.Body.String())
			}
			cookies := response.Result().Cookies()
			if len(cookies) != 1 || !cookies[0].HttpOnly || !cookies[0].Secure || cookies[0].SameSite != http.SameSiteStrictMode || cookies[0].Path != "/megan" {
				t.Fatalf("incorrect session cookie: %+v", cookies)
			}
			r.AddCookie(cookies[0])
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "Private message") {
				t.Fatalf("authenticated page failed: %d %s", w.Code, w.Body.String())
			}
			if strings.Contains(w.Body.String(), "EventSource") != liveReload {
				t.Fatal("incorrect live reload behavior")
			}
		}
	}
}

func TestMeganDestination(t *testing.T) {
	for _, raw := range []string{"https://evil.example/megan/", "//evil.example/megan/", "/articles/", "/megan/../../", "/megan/login", "/megan/..%2farticles/"} {
		if got := meganDestination(raw); got != "/megan/" {
			t.Errorf("unsafe destination %q accepted as %q", raw, got)
		}
	}
	if got := meganDestination("/megan/what-i-love-about-you/"); got != "/megan/what-i-love-about-you/" {
		t.Errorf("lost private document destination: %s", got)
	}
}

func TestCollections(t *testing.T) {
	for _, section := range []string{"articles", "poetry", "megan"} {
		t.Run(section, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), section)
			if err := os.Mkdir(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			for _, post := range []struct{ slug, date string }{
				{"older", "2026-01-01"},
				{"newer", "2026-02-01"},
			} {
				body := "---\ntitle: Example\ndate: " + post.date + "\nslug: " + post.slug + "\n---\n\n# Example\n\nFirst line\nSecond line\n\nNext stanza.\n"
				if err := os.WriteFile(filepath.Join(dir, post.slug+".md"), []byte(body), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			posts, err := loadArticles(dir)
			if err != nil {
				t.Fatal(err)
			}
			if len(posts) != 2 || posts[0].Slug != "newer" {
				t.Fatalf("expected newest post first, got %+v", posts)
			}
			if posts[0].URL != "/"+section+"/newer/" {
				t.Fatalf("wrong collection URL: %s", posts[0].URL)
			}
			html := string(posts[0].HTML)
			if strings.Contains(html, "<br>") != (section == "poetry") {
				t.Fatalf("incorrect line break rendering: %s", html)
			}
			if !strings.Contains(html, "<p>Next stanza.</p>") {
				t.Fatalf("missing separate stanza: %s", html)
			}
		})
	}
}
