# IV Blog

Dead simple personal site generator in Go, with Home, Writing, and Poetry sections.

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

### Writing

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

Browse the Writing section at `/articles/`.

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

## Private documents

Visit `/megan` manually and enter the configured password in the login form.
No username is required. No public navigation links to this page.
The page lists Markdown documents from `content/megan/`, newest first, just like
Writing and Poetry. Use the same `title`, `date`, and `slug` frontmatter. Each
document gets its own protected page at `/megan/<slug>/`, with a link back to the
Megan index. Edit `content/megan/what-i-love-about-you.md` to update the love list.
Edits trigger rebuilds in watch mode.

Authentication runs server-side in Go, both locally (including watch mode) and in
production. Caddy proxies `/megan` and all nested paths to the Go service on
`127.0.0.1:8081`; it never serves those files directly. Only a bcrypt password hash
is stored, in `cmd/site/auth.go`. Successful login sets a signed, HttpOnly,
SameSite=Strict cookie, marked Secure over HTTPS, valid for 12 hours. Restarting
the server invalidates sessions. A direct document link returns to that document
after login.

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

The script also cross-compiles and installs the Go server outside the public web
root, then starts `iv-blog.service` as a systemd user service. The VPS needs a
working `systemctl --user` session and `curl`. Enable lingering for the deploy user
(`loginctl enable-linger trbiv`, with administrator authorization if required) so
the service starts at boot and stays running after logout. If the backend stops,
private URLs return a gateway error rather than exposing static content.

To serve an existing build without rebuilding:

```sh
go run ./cmd/site -serve -dir build -addr 127.0.0.1:8081
```
