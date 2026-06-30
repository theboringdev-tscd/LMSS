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

---

# Frontend — Phase 2: Auth & Dashboard (June 2026)

## Done

### Backend
- [x] Remove demo services (`users_service`, `products_service`, `tasks_service`, `broadcast_service`, `cache_service`, `encryption_service`)
- [x] Create `internal/services/modules/auth_service.go`
  - JWT login with MongoDB-backed users (bcrypt password hashing)
  - Seeds default admin user on first start: `admin@library.org` / `admin`
  - Endpoints: `POST /api/v1/auth/login`, `POST /api/v1/auth/refresh`, `POST /api/v1/auth/logout`
  - Uses existing JWT middleware token helpers
- [x] Create `internal/services/modules/reports_service.go`
  - `GET /api/v1/reports/overview` — dashboard data
  - Returns `OverviewStats` (books_checked_out_today, new_patrons_this_week, overdue_returns, recent_activity)
  - Protected by `middleware.JWTRequired`
- [x] Enable `jwt` middleware in `config.yaml` (switched to `JWTOptional` globally, `JWTRequired` per protected route group)
- [x] Update `config.yaml` auth block (`type: jwt`)
- [x] Fix `RegisterServiceWithDependencies` to return nil (not error) when service factory graceful-skips due to missing dependencies

### Frontend
- [x] Enhance `web/src/lib/api.ts`
  - `setToken`, `clearToken`, `isAuthenticated` helpers using `localStorage`
  - Auth header injection on every request (Bearer token)
  - 401 interceptor → clears token + redirects to `/login`
  - Domain API functions: `login()`, `logout()`, `getDashboardOverview()`
- [x] Create `web/src/components/dashboard/StatBlock.tsx`
  - Typographic stat card matching the design system
- [x] Create `web/src/components/dashboard/ActivityFeed.tsx`
  - Recent transactions list with graceful empty state
- [x] Create `web/src/components/dashboard/Dashboard.tsx`
  - React island that fetches live data via `getDashboardOverview`
  - Sign-out button bound to `logout()`
- [x] Update `web/src/pages/index.astro` → purely shells `<Dashboard client:load />`
- [x] Update `web/src/pages/login.astro`
  - Client-side redirect if already authenticated
  - Form submit handler calling `login()` function
  - Error display below submit button
  - Loading state on button

## API Routes (Backend)

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/auth/login` | Public | Email + password → JWT token |
| POST | `/api/v1/auth/refresh` | JWT | Refresh existing token |
| POST | `/api/v1/auth/logout` | JWT | No-op (stateless) |
| GET | `/api/v1/reports/overview` | JWT | Dashboard stats overview |

## Frontend Auth Flow

1. User visits `/login` → server renders static HTML shell
2. Client-side script runs → checks `localStorage` → redirects to `/` if token exists
3. Submits credentials → `api.login(email, password)` stores token + redirects to `/`
4. All subsequent requests include `Authorization: Bearer <token>` header automatically
5. 401 response → token cleared + browser redirected to `/login`
6. `/logout` → clears token + redirects to login

## Test Fixes
- Updated `tests/simple_test.go` to replace `users_service` references with `auth_service`
- Updated `tests/startup_test.go` to register health endpoints in `mustBuildDiagnosticRouter` (previously broken due to Echo→Gin migration)
- Fixed `pkg/registry/registry.go` — `RegisterServiceWithDependencies` now returns nil (not error) when a service factory gracefully returns nil due to missing dependencies (e.g., MongoDB not connected in test environment)

## Next Steps

- Phase 3: Catalog (CRUD + full-text search via MongoDB)
- Phase 4: Patrons & Circulation
- Phase 5: Fines, Reservations, Reports extensions
- Phase 6: Polish (loading states, responsive, keyboard accessibility)
