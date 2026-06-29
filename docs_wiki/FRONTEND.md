# Frontend — Phase 1: Astro Shell + Go Embed

## Done

- [x] Create `web/` Astro project (React + Tailwind CSS v3 + shadcn/ui-ready config)
- [x] Write global CSS with design tokens (binding, paper, leather, gold-leaf, book-cloth, ink, reading-lamp)
- [x] Set up Google Fonts (Cormorant display, Inter body, JetBrains Mono data)
- [x] Build `BaseLayout.astro` (meta, fonts, global CSS shell)
- [x] Build `SidebarLayout.astro` (sidebar + main content slot)
- [x] Build `Sidebar.tsx` React island (8 nav items + login link, gold-leaf branding)
- [x] Build 15 placeholder pages (dashboard, catalog, patrons, circulation, fines, reservations, reports, settings, login + sub-pages)
- [x] Create `web/embed.go` (`//go:embed dist/*`, package `web`)
- [x] Create `web/dist/placeholder.txt` for Go compile-time guarantee when `dist/` is empty
- [x] Create `web/.gitignore` (`dist/* !dist/placeholder.txt`, `node_modules/`)
- [x] Add `FrontendConfig` to `config/config.go` (`Enabled`, `Path` fields + viper defaults)
- [x] Add `frontend:` block to `config.yaml`
- [x] Add `startWithFrontend()` method in `internal/server/server.go`:
  - [x] API/health/swagger/metrics routes → Gin
  - [x] Serve exact file matches from embedded `dist/` (direct `fs.ReadFile` + `w.Write`)
  - [x] SPA fallback: `fs.ReadFile(subFS, "index.html")` + `w.Write`
- [x] Remove unused config fields (`app.quiet_startup`, `encryption.rotate_keys`, `encryption.key_rotation_interval`)
- [x] Remove dead struct fields from `config/config.go` (`QuietStartup`, `RotateKeys`, `KeyRotationInterval`, `ExternalConfig`, `ExternalService`, `SwaggerConfig.BasePath`)

## Verified

- [x] `go build ./cmd/app` compiles
- [x] `go vet` passes
- [x] `/` → 200 (index.html)
- [x] `/login` → 200 (index.html)
- [x] `/favicon.svg` → 200 (244 bytes)
- [x] `/_astro/Sidebar.Cj2WZahx.js` → 200 (30293 bytes)
- [x] `/api/v1/bogus` → 404 JSON
- [x] `/health` → 200

## In Progress

- [ ] Fix 301 redirect for `/index.html` and `/catalog`:
  - `/index.html`: `http.FileServer` redirects → `./`. Replace serve-exact-file block with direct `fs.ReadFile` + `w.Write`.
  - `/catalog`: `dist/catalog/` is a directory (Astro sub-pages). Skip directories in serve-exact-file check, let them fall through to SPA fallback.
- [ ] Clean up debug logging from `server.go`
- [ ] Restore `config.yaml` production values after testing

## Blocked

- [ ] `http.FileServer` 301 redirect for exact file hits (`/index.html`) and directory paths (`/catalog` → `/catalog/`)

## Next Steps

- [ ] Remove `http.FileServer` from `startWithFrontend()`, use only `fs.ReadFile` for both file serving and SPA fallback
- [ ] Run final `go build ./cmd/app` after frontend server fix
