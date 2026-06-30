
<div align="center">
  <img src=".github/assets/stackyrd-banner.png" alt="stackyrd" style="width: 100%; max-width: 700px;"/>
</div>
<div align="center">
  <img src="https://www.shieldcn.dev/github/license/theboringdev-tscd/LMSS.svg?variant=secondary&size=xs" alt="License"/>
  <img src="https://shieldcn.dev/badge/Go-00ADD8.svg?logo=go&logoColor=fff&variant=branded&size=xs" alt="Go Version"/>
  <img src="https://www.shieldcn.dev/github/ci/theboringdev-tscd/LMSS.svg?variant=secondary&size=xs" alt="Build Status"/>
  <img src="https://shieldcn.dev/github/theboringdev-tscd/LMSS/ci.svg?variant=secondary&size=xs&logo=lu%3AShield&label=Security+scan&name=security.yml" alt="Security Status"/>
  <img src="https://www.shieldcn.dev/github/release/theboringdev-tscd/LMSS.svg?size=xs&variant=secondary" alt="Release"/>
  <img src="https://www.shieldcn.dev/badge/Agent--friendly-AGENTS.md-D97757.svg?variant=secondary&size=xs" alt="Agents Friendly"/>
</div>
<br>

**LMSS (Library Management System Stackyard)** is a full-featured library management application built on the stackyrd Go framework - designed for librarians and library staff who need to find, manage, and move library assets (books, patrons, loans) as quickly as they can walk to a shelf.

### Core Architecture

| Layer | Technology |
|-------|------------|
| **Backend** | Go with [Gin](https://github.com/gin-gonic/gin) - auto-discovered services, middleware, infrastructure |
| **Frontend** | Astro.js SPA with React islands, shadcn/ui, Tailwind CSS - embedded in the Go binary |
| **Database** | MongoDB (flexible schema for library metadata) |
| **Design** | Library-themed palette (binding, paper, leather, gold-leaf), typographic dashboard, Cormorant/Inter/JetBrains Mono |

### What it does

- **Catalog** - Browse, search, add, and manage the book collection with full-text search
- **Patrons** - Patron management with profiles and loan history
- **Circulation** - Check out, return, and track loans
- **Fines** - Fine assessment and payment management
- **Reservations** - Book holds and reservation management
- **Reports** - Usage reports and analytics
- **Auth** - JWT-based authentication for staff

## Quick Start

### Installation & Run

```bash
# Clone the repository
git clone https://github.com/theboringdev-tscd/LMSS.git
cd stackyrd

# Install dependencies
go mod download

# Run the application
go run cmd/app/main.go

# To build the application
go run scripts/build/build.go

# To download package
go run scripts/pkg/pkg.go

```

## Preview

![Console](.github/assets/console.png)

## Documentation

- **[Full Documentation](docs_wiki/)** - Comprehensive guides and references
- **[Plugin System Guide](PLUGIN_GUIDE.md)** - Creating and managing TypeScript, Lua, Python, and Go plugins
- **[Contributing Guide](CONTRIBUTING.md)** - Development workflow and guidelines

## License

Distributed under the Apache License Version 2.0. See `LICENSE` for full information.
