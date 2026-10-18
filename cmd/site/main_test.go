package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMeganAuthentication(t *testing.T) {
	for _, liveReload := range []bool{false, true} {
		h := devHandler{hub: newReloadHub(), liveReload: liveReload}
		for _, path := range []string{"/megan", "/megan/", "/megan/index.html", "/megan/nested/file", "/%6degan/", "/articles/../megan/", "//megan/"} {
			for _, method := range []string{http.MethodGet, http.MethodHead} {
				r := httptest.NewRequest(method, path, nil)
				w := httptest.NewRecorder()
				h.ServeHTTP(w, r)
				if w.Code != http.StatusUnauthorized || !strings.HasPrefix(w.Header().Get("WWW-Authenticate"), "Basic ") {
					t.Fatalf("%s %s: expected Basic Auth challenge, got %d", method, path, w.Code)
				}
				if w.Header().Get("Cache-Control") != "no-store" {
					t.Fatalf("%s: missing no-store", path)
				}
			}
		}
		r := httptest.NewRequest(http.MethodGet, "/megan/", nil)
		r.SetBasicAuth("megan", "wrong-password")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != http.StatusUnauthorized || len(w.Result().Cookies()) != 0 {
			t.Fatalf("wrong password accepted: %d", w.Code)
		}
	}
}

func TestMeganAuthenticatedPage(t *testing.T) {
	password := os.Getenv("MEGAN_TEST_PASSWORD")
	if password == "" {
		t.Skip("set MEGAN_TEST_PASSWORD to verify the configured password")
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
		h := devHandler{hub: newReloadHub(), liveReload: liveReload, root: root}
		for _, path := range []string{"/megan/", "/megan/what-i-love-about-you/"} {
			r := httptest.NewRequest(http.MethodGet, path, nil)
			unauthorized := httptest.NewRecorder()
			h.ServeHTTP(unauthorized, r)
			if unauthorized.Code != http.StatusUnauthorized || strings.Contains(unauthorized.Body.String(), "Private message") {
				t.Fatalf("private content exposed at %s", path)
			}
			r.SetBasicAuth("other-user", password)
			wrongUser := httptest.NewRecorder()
			h.ServeHTTP(wrongUser, r)
			if wrongUser.Code != http.StatusUnauthorized {
				t.Fatalf("non-static user accepted at %s", path)
			}
			r.SetBasicAuth("megan", password)
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
