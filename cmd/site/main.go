package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	urlpath "path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
)

const (
	contentDir  = "content"
	staticDir   = "static"
	templateDir = "templates"
	buildDir    = "build"
)

const liveReloadScript = `<script>
(() => {
  const events = new EventSource("/__reload");
  events.addEventListener("reload", () => window.location.reload());
})();
</script>`

// Keep this hash in sync with the /megan basic_auth entry in Caddyfile.
const meganPasswordHash = "$2a$14$lZuSbdcBQhiXWKPpuYd4IO1iXha32j/BPNTZaZwAUNR/VjWGvzFo2"

type config struct {
	serve bool
	watch bool
	addr  string
}

type frontmatter struct {
	Title string `yaml:"title"`
	Date  string `yaml:"date"`
	Slug  string `yaml:"slug"`
}

type article struct {
	Title string
	Date  string
	Slug  string
	URL   string
	HTML  template.HTML
}

type pageData struct {
	Title    string
	Section  string
	Heading  string
	Empty    string
	Home     template.HTML
	Articles []article
	Article  article
}

func main() {
	cfg := config{}
	flag.BoolVar(&cfg.serve, "serve", false, "serve the generated site")
	flag.BoolVar(&cfg.watch, "watch", false, "rebuild and reload when source files change")
	flag.StringVar(&cfg.addr, "addr", "0.0.0.0:8080", "address for the server")
	flag.Parse()

	var err error
	if cfg.serve {
		err = serve(cfg)
	} else {
		err = build()
	}
	if err != nil {
		log.Fatal(err)
	}
}

func serve(cfg config) error {
	if err := build(); err != nil {
		return err
	}

	hub := newReloadHub()
	if cfg.watch {
		go watchSources(hub)
	}

	log.Printf("serving %s on %s", buildDir, cfg.addr)

	return http.ListenAndServe(cfg.addr, devHandler{
		hub:        hub,
		liveReload: cfg.watch,
	})
}

func build() error {
	home, err := os.ReadFile(filepath.Join(contentDir, "home.md"))
	if err != nil {
		return err
	}
	var summary bytes.Buffer
	if err := goldmark.Convert(home, &summary); err != nil {
		return err
	}
	megan, err := os.ReadFile(filepath.Join(contentDir, "megan.md"))
	if err != nil {
		return err
	}
	var privateBody bytes.Buffer
	if err := goldmark.Convert(megan, &privateBody); err != nil {
		return err
	}
	type collection struct {
		slug, title, empty string
		posts              []article
	}
	collections := []collection{
		{slug: "articles", title: "Articles", empty: "No articles yet."},
		{slug: "poetry", title: "Poetry", empty: "No poems yet."},
	}
	for i := range collections {
		posts, err := loadArticles(filepath.Join(contentDir, collections[i].slug))
		if err != nil {
			return err
		}
		collections[i].posts = posts
	}

	tmpl, err := template.ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return err
	}

	if err := os.RemoveAll(buildDir); err != nil {
		return err
	}
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return err
	}
	if err := copyDir(staticDir, filepath.Join(buildDir, "static")); err != nil {
		return err
	}

	if err := render(tmpl, "index.html", filepath.Join(buildDir, "index.html"), pageData{
		Title:   "IV",
		Section: "home",
		Home:    template.HTML(summary.String()),
	}); err != nil {
		return err
	}

	if err := render(tmpl, "private.html", filepath.Join(buildDir, "megan", "index.html"), pageData{
		Title: "Megan - IV",
		Home:  template.HTML(privateBody.String()),
	}); err != nil {
		return err
	}

	for _, collection := range collections {
		if err := render(tmpl, "list.html", filepath.Join(buildDir, collection.slug, "index.html"), pageData{
			Title:    collection.title + " - IV",
			Section:  collection.slug,
			Heading:  collection.title,
			Empty:    collection.empty,
			Articles: collection.posts,
		}); err != nil {
			return err
		}
		for _, article := range collection.posts {
			out := filepath.Join(buildDir, collection.slug, article.Slug, "index.html")
			if err := render(tmpl, "article.html", out, pageData{
				Title:   article.Title + " - IV",
				Section: collection.slug,
				Heading: collection.title,
				Article: article,
			}); err != nil {
				return err
			}
		}
	}

	fmt.Printf("built home, %d article(s), and %d poem(s) into %s\n", len(collections[0].posts), len(collections[1].posts), buildDir)
	return nil
}

func loadArticles(dir string) ([]article, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.Typographer),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)
	if filepath.Base(dir) == "poetry" {
		md.Renderer().AddOptions(html.WithHardWraps())
	}

	var articles []article
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		post, err := loadArticle(path, md)
		if err != nil {
			return nil, err
		}
		post.URL = "/" + filepath.Base(dir) + "/" + post.Slug + "/"
		articles = append(articles, post)
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Date > articles[j].Date
	})

	return articles, nil
}

func loadArticle(path string, md goldmark.Markdown) (article, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return article{}, err
	}

	meta, body, err := splitFrontmatter(string(raw))
	if err != nil {
		return article{}, fmt.Errorf("%s: %w", path, err)
	}

	var fm frontmatter
	if err := yaml.Unmarshal([]byte(meta), &fm); err != nil {
		return article{}, fmt.Errorf("%s: %w", path, err)
	}
	if fm.Title == "" || fm.Date == "" || fm.Slug == "" {
		return article{}, fmt.Errorf("%s: title, date, and slug are required", path)
	}
	if _, err := time.Parse("2006-01-02", fm.Date); err != nil {
		return article{}, fmt.Errorf("%s: date must use YYYY-MM-DD", path)
	}

	var out bytes.Buffer
	if err := md.Convert([]byte(body), &out); err != nil {
		return article{}, fmt.Errorf("%s: %w", path, err)
	}

	return article{
		Title: fm.Title,
		Date:  fm.Date,
		Slug:  fm.Slug,
		URL:   "/articles/" + fm.Slug + "/",
		HTML:  template.HTML(out.String()),
	}, nil
}

func splitFrontmatter(raw string) (string, string, error) {
	if !strings.HasPrefix(raw, "---\n") {
		return "", "", errors.New("missing frontmatter")
	}

	parts := strings.SplitN(strings.TrimPrefix(raw, "---\n"), "\n---\n", 2)
	if len(parts) != 2 {
		return "", "", errors.New("unclosed frontmatter")
	}

	return parts[0], parts[1], nil
}

func render(tmpl *template.Template, name, out string, data pageData) error {
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}

	file, err := os.Create(out)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.ExecuteTemplate(file, name, data)
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, body, 0o644)
	})
}

type devHandler struct {
	hub        *reloadHub
	liveReload bool
}

func (h devHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clean := urlpath.Clean("/" + r.URL.Path)
	if clean == "/megan" || strings.HasPrefix(clean, "/megan/") {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Robots-Tag", "noindex, nofollow, noarchive")
		username, password, ok := r.BasicAuth()
		if !ok || username != "megan" || bcrypt.CompareHashAndPassword([]byte(meganPasswordHash), []byte(password)) != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Megan", charset="UTF-8"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}
	if r.URL.Path == "/__reload" {
		h.hub.handle(w, r)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, ok := buildPath(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}

	if h.liveReload && filepath.Ext(file) == ".html" {
		h.serveHTML(w, r, file)
		return
	}

	http.ServeFile(w, r, file)
}

func (h devHandler) serveHTML(w http.ResponseWriter, r *http.Request, file string) {
	body, err := os.ReadFile(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := string(body)
	if strings.Contains(html, "</body>") {
		html = strings.Replace(html, "</body>", liveReloadScript+"\n  </body>", 1)
	} else {
		html += liveReloadScript
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Method == http.MethodHead {
		return
	}
	_, _ = w.Write([]byte(html))
}

func buildPath(requestPath string) (string, bool) {
	clean := urlpath.Clean("/" + requestPath)
	rel := strings.TrimPrefix(clean, "/")
	target := filepath.Join(buildDir, filepath.FromSlash(rel))

	info, err := os.Stat(target)
	if err != nil {
		return "", false
	}
	if info.IsDir() {
		target = filepath.Join(target, "index.html")
		info, err = os.Stat(target)
		if err != nil || info.IsDir() {
			return "", false
		}
	}

	return target, true
}

type fileState struct {
	modTime int64
	size    int64
}

func watchSources(hub *reloadHub) {
	paths := []string{contentDir, staticDir, templateDir}
	previous, err := snapshot(paths)
	if err != nil {
		log.Printf("watch: %v", err)
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		next, err := snapshot(paths)
		if err != nil {
			log.Printf("watch: %v", err)
			continue
		}
		if sameSnapshot(previous, next) {
			continue
		}

		previous = next
		if err := build(); err != nil {
			log.Printf("rebuild failed: %v", err)
			continue
		}
		hub.broadcast()
	}
}

func snapshot(paths []string) (map[string]fileState, error) {
	state := map[string]fileState{}
	for _, root := range paths {
		if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}

			info, err := entry.Info()
			if err != nil {
				return err
			}
			state[path] = fileState{
				modTime: info.ModTime().UnixNano(),
				size:    info.Size(),
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}
	return state, nil
}

func sameSnapshot(a, b map[string]fileState) bool {
	if len(a) != len(b) {
		return false
	}
	for path, left := range a {
		if right, ok := b[path]; !ok || left != right {
			return false
		}
	}
	return true
}

type reloadHub struct {
	mu      sync.Mutex
	clients map[chan struct{}]struct{}
}

func newReloadHub() *reloadHub {
	return &reloadHub{clients: map[chan struct{}]struct{}{}}
}

func (h *reloadHub) handle(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := h.subscribe()
	defer h.unsubscribe(ch)

	fmt.Fprint(w, ": connected\n\n")
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ch:
			fmt.Fprintf(w, "event: reload\ndata: %d\n\n", time.Now().UnixNano())
			flusher.Flush()
		}
	}
}

func (h *reloadHub) subscribe() chan struct{} {
	ch := make(chan struct{}, 1)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *reloadHub) unsubscribe(ch chan struct{}) {
	h.mu.Lock()
	delete(h.clients, ch)
	h.mu.Unlock()
}

func (h *reloadHub) broadcast() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.clients {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
