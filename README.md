# IV Blog

Dead simple static blog generator in Go.

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

## Deploy

```sh
./deploy.sh
```

Deployment builds the site, replaces the static files on the VPS, and reloads Caddy.
