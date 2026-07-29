# IV Blog

Dead simple static blog generator in Go.

## Build

From this directory:

```sh
go run ./cmd/site
```

Output goes to `build/`.

## Dev Server

Serve locally:

```sh
go run ./cmd/site -serve
```

Serve locally with hot reload:

```sh
go run ./cmd/site -serve -watch
```

Open `http://localhost:8080`.

The server binds to `0.0.0.0:8080` by default, so it is reachable through any network IP on this machine if firewall/router rules allow it.

Use another port:

```sh
go run ./cmd/site -serve -watch -addr :3000
```

Bind to all IPv4 addresses on another port:

```sh
go run ./cmd/site -serve -watch -addr 0.0.0.0:3000
```

## Nix

From repo root:

```sh
nix develop
cd go-blog
go run ./cmd/site -serve -watch
```

One-shot without entering shell:

```sh
nix develop -c bash -lc 'cd go-blog && go run ./cmd/site -serve -watch'
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
