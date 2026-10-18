# IV Blog

Dead simple personal site generator in Go, with Home, Articles, and Poetry sections.

## Build

```sh
go run ./cmd/site
```

Output goes to `build/`.

## Dev Server

Serve locally:

```sh
go run ./cmd/site -serve
```

Serve with hot reload:

```sh
go run ./cmd/site -serve -watch
```

Open `http://localhost:8080`. Use `-addr` to choose another bind address:

```sh
go run ./cmd/site -serve -watch -addr 127.0.0.1:3000
```

## Nix

Enter the development shell:

```sh
nix develop
```

Or run the development server directly:

```sh
nix develop -c go run ./cmd/site -serve -watch
```

## Write

### Home

Edit `content/home.md` to update the professional summary on `/`. This is plain
Markdown, with no frontmatter. Replace the initial placeholder with your bio and
any contact or professional links.

### Articles

Add Markdown files to `content/articles/` with frontmatter:

```md
---
title: Peru
date: 2026-06-07
slug: peru
---

# Post title

Post body.
```

Generated article URL:

```text
/articles/peru/
```

Browse all articles at `/articles/`.

### Poetry

Add Markdown files to `content/poetry/` with the same frontmatter:

```md
---
title: Morning
date: 2026-10-17
slug: morning
---

# Morning

First line of a poem
Second line of a poem

A new stanza begins here.
```

Poems appear at `/poetry/<slug>/` and are listed at `/poetry/`. Single line breaks
are preserved in poems; blank lines separate stanzas. Include a heading in the
body for the visible title, as with articles. Both lists show newest entries first.

Home, article, and poetry edits all trigger rebuilds in watch mode.

## Private page

Visit `/megan` manually and sign in with username `megan` and the configured
password using the browser's login prompt. No public navigation links to this page.
Edit `content/megan.md` to change its content.

Authentication runs server-side: Caddy protects production and the Go development
server protects local previews, including watch mode. Both protect `/megan`,
`/megan/`, and nested paths such as `/megan/index.html`. Only a bcrypt hash is stored,
in `Caddyfile` and `cmd/site/main.go`; update both when changing the password.

The generated `build/megan/` files contain the page content. Serve the build with
the supplied Caddy configuration or Go server, not an unprotected static host.
Production uses HTTPS; local HTTP previews are intended for trusted local use.

## Deploy

```sh
./deploy.sh
```

Deployment builds the site, validates and reloads Caddy's authentication rules,
then replaces the static files on the VPS so protection is active before private
content is published.
