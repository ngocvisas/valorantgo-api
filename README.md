# Valorant API (Encore + Claude Code Demo)

A small Encore (Go) backend generated & refined with **vibe-coding** using Claude Code.

## ✨ Features
- Public APIs: list agents/weapons, health, stats
- Auth APIs: create/list user loadouts (dev auth via `Bearer dev-<name>`)
- Postgres via `sqldb` with automatic migration (table: `loadouts`)
- Ready for Encore Cloud deploy

## ▶️ Run locally
**Requirements:** Encore CLI + Docker Desktop

```bash
encore run
# API: http://localhost:4000
# Dev Dashboard: http://localhost:9400
