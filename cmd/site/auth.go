package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const meganPasswordHash = "$2a$14$lZuSbdcBQhiXWKPpuYd4IO1iXha32j/BPNTZaZwAUNR/VjWGvzFo2"
const sessionCookie = "megan_session"
const sessionLifetime = 12 * time.Hour

//go:embed login.html
var loginHTML string

var loginTemplate = template.Must(template.New("login").Parse(loginHTML))

type meganAuth struct {
	key [32]byte
}

func newMeganAuth() (*meganAuth, error) {
	a := &meganAuth{}
	_, err := rand.Read(a.key[:])
	return a, err
}

func (a *meganAuth) signature(value string) string {
	mac := hmac.New(sha256.New, a.key[:])
	mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (a *meganAuth) token(expires time.Time) string {
	value := strconv.FormatInt(expires.Unix(), 10)
	return value + "." + a.signature(value)
}

func (a *meganAuth) authenticated(r *http.Request) bool {
	cookie, err := r.Cookie(sessionCookie)
	if err != nil {
		return false
	}
	value, signature, ok := strings.Cut(cookie.Value, ".")
	if !ok || !hmac.Equal([]byte(signature), []byte(a.signature(value))) {
		return false
	}
	expires, err := strconv.ParseInt(value, 10, 64)
	return err == nil && time.Now().Unix() < expires
}

func meganDestination(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.IsAbs() || u.Host != "" || strings.Contains(u.Path, "\\") {
		return "/megan/"
	}
	clean := path.Clean("/" + u.Path)
	if clean == "/megan/login" || (clean != "/megan" && !strings.HasPrefix(clean, "/megan/")) {
		return "/megan/"
	}
	if strings.HasSuffix(u.Path, "/") || clean == "/megan" {
		clean += "/"
	}
	return (&url.URL{Path: clean, RawQuery: u.RawQuery}).String()
}

func showLogin(w http.ResponseWriter, r *http.Request, destination, message string, status int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Frame-Options", "DENY")
	w.WriteHeader(status)
	if r.Method != http.MethodHead {
		_ = loginTemplate.Execute(w, struct{ Next, Error string }{destination, message})
	}
}

func (a *meganAuth) authorize(w http.ResponseWriter, r *http.Request, clean string) bool {
	if clean == "/megan/login" && r.Method == http.MethodPost {
		r.Body = http.MaxBytesReader(w, r.Body, 4096)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid login request", http.StatusBadRequest)
			return false
		}
		next := meganDestination(r.PostForm.Get("next"))
		if bcrypt.CompareHashAndPassword([]byte(meganPasswordHash), []byte(r.PostForm.Get("password"))) != nil {
			showLogin(w, r, next, "That password isn't right. Try again.", http.StatusUnauthorized)
			return false
		}
		expires := time.Now().Add(sessionLifetime)
		http.SetCookie(w, &http.Cookie{
			Name: sessionCookie, Value: a.token(expires), Path: "/megan",
			HttpOnly: true, Secure: r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https",
			SameSite: http.SameSiteStrictMode, MaxAge: int(sessionLifetime.Seconds()), Expires: expires,
		})
		http.Redirect(w, r, next, http.StatusSeeOther)
		return false
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	if clean == "/megan/login" || !a.authenticated(r) {
		showLogin(w, r, meganDestination(r.URL.RequestURI()), "", http.StatusUnauthorized)
		return false
	}
	return true
}
