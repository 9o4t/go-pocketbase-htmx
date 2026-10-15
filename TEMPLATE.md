# Railway template config (Go + PocketBase + HTMX)

Publish from Dave's Railway account using this public repo as the source. Then put the published template code in `9o4t/c17` `internal/railway/templates.go` as `"go-pocketbase-htmx": "<code>"`. Do not invent a code.

Source: GitHub [`9o4t/go-pocketbase-htmx`](https://github.com/9o4t/go-pocketbase-htmx), **repository root**. Do not point the template at the c17.io site root or at `templates/go-pocketbase-htmx` inside `9o4t/c17`.

## Service

| Field | Value |
|---|---|
| Name | `web` |
| Source | GitHub `9o4t/go-pocketbase-htmx`, root directory `.` |
| Builder | Dockerfile (`Dockerfile` at repo root) |
| Start | image `CMD` — `app serve --http=0.0.0.0:$PORT` |
| Public | Yes, one domain |
| Healthcheck path | `/healthz` |
| Healthcheck timeout | 30s |
| Restart | on failure |

`railway.toml` and `railway.json` already set builder + healthcheck. The composer can import them.

## Volume

| Mount | Service | Why |
|---|---|---|
| `/data` | `web` | PocketBase SQLite + files. The image sets `PB_DATA_DIR=/data`. Without this volume, notes and the admin user reset on every deploy. |

## Variables Railway provides

| Variable | Service | Why |
|---|---|---|
| `PORT` | `web` | Listen port. Injected. Do not hard-code. |

## Variables the operator may provide

None are required to boot. The app serves `/` and `/healthz` with an empty `notes` collection.

| Variable | Required | Why |
|---|---|---|
| `PB_ADMIN_EMAIL` | No | First PocketBase superuser (`/_/`) |
| `PB_ADMIN_PASSWORD` | No | Pair for `PB_ADMIN_EMAIL`. Ignored if email is unset. |
| `PB_DATA_DIR` | No | Leave unset so the image default `/data` is used |

Do not collect a referral code, site URL, or API keys in the template. This starter has no outbound SaaS.

## Composer / dashboard publish checklist

1. Sign in as Dave on [railway.com](https://railway.com).
2. New project from this GitHub repo (`9o4t/go-pocketbase-htmx`, root, Dockerfile builder). Confirm it boots with **zero required env vars** (`PORT` injected).
3. Attach a volume at `/data` on service `web`.
4. Healthcheck `GET /healthz` (already in `railway.toml` / `railway.json`). Timeout 30s. Restart on failure.
5. Generate a template from that project (dashboard **Generate Template**, or `railway templates create --project <id> --json`).
6. In the template editor: one service, GitHub source = this public repo (root). Mark `PB_ADMIN_EMAIL` and `PB_ADMIN_PASSWORD` optional. No other variables.
7. Description: one-click Go + PocketBase + HTMX. Mention the `/data` volume.
8. Publish to the marketplace (category **Starters** is a good fit): dashboard publish, or `railway templates publish <template-id> --category Starters --description "One-click Go + PocketBase + HTMX. Volume at /data. Health check /healthz. Zero required env vars." --readme-file README.md --json`.
9. Copy the published **template code** (the short code in `https://railway.com/new/template/<code>`).
10. In `9o4t/c17`, set `internal/railway/templates.go`:

    `"go-pocketbase-htmx": "<code>"`

    That single entry switches `/stacks/go-pocketbase-htmx`, related process pages (auth/backend), newsletter issue `001-go-pocketbase-htmx`, and `/replace/{firebase,supabase,vercel,netlify}` from **Start on Railway** to the official Deploy button with `referralCode=4xHDyE`.
11. Run `go test ./...` in c17. Merge only when the code is the real published one.
