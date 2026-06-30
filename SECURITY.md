# Security Policy

## Scanning Pipeline

CI runs automated security checks on every push and daily via `.github/workflows/security.yml`:

**Gosec** — Go source static analysis  
**Nancy + govulncheck** — Dependency CVE scanning  
**Trivy** — Filesystem & misconfiguration scanning  
**StaticCheck + Go-Critic** — Go linters  

Findings are uploaded to GitHub Code Scanning as SARIF.

## Secrets

**Never commit secrets.** The `config.yaml` in this repo contains development credentials only.

In production, load all secrets via environment variables or a secrets manager.

### At-Risk Fields in `config.yaml`

| Field | Risk |
|---|---|
| `auth.secret` | Full auth bypass |
| `postgres.connections[].password` | Database access |
| `mongo.connections[].uri` | DB access |
| `encryption.key` | Data decryption |

### Production Checklist

- `app.env: production`, `debug: false`
- JWT auth enabled (`middleware.jwt: true`)
- Rate limiting and audit logging **on**
- CORS locked to known origins (no `*`)
- `sslmode: require` or `verify-full` on Postgres
- TLS/SCRAM on MongoDB
- HSTS headers on (provided by `security` middleware)
- Frontend served over HTTPS with appropriate CSP headers

## Frontend Security

The embedded Astro.js SPA is served by the Go binary. Ensure:
- CSP headers are set via the `security` middleware
- No secrets are embedded in the frontend build (`web/dist/`)
- Sanitize any user-generated content rendered in the UI
