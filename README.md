# Go + PocketBase + HTMX starter

A one-service Railway starter: Go embeds PocketBase, serves a single HTMX page, and stores notes in a `notes` collection.

[c17.io](https://c17.io) is the dogfood version of this stack. Users never pay c17; Railway usage is what the referral / template kickback is for.

**Zero required environment variables.** Railway injects `PORT`. Attach a volume at `/data`. Health check is `GET /healthz`. Optional PocketBase admin: `PB_ADMIN_EMAIL` + `PB_ADMIN_PASSWORD`.

## Cost note

On a small Railway service plus a 1 GB volume this is usually a few dollars a month. The SaaS path it replaces (Firebase or Supabase for auth/data, plus Vercel or Netlify for the site) is the ~$2,400/year line on [c17.io/stacks/go-pocketbase-htmx](https://c17.io/stacks/go-pocketbase-htmx). Trial credits cover the first deploy for most people.

## Run locally

Needs Go 1.27+.

```bash
go test .
go run . serve --http=127.0.0.1:8090
```

- App: http://127.0.0.1:8090
- PocketBase admin: http://127.0.0.1:8090/_/
- Data: `./pb_data` locally, `/data` when that directory exists (Railway volume)

```bash
docker build -t go-pocketbase-htmx .
docker run --rm -p 8090:8090 -e PORT=8090 go-pocketbase-htmx
```

## Railway

One service. Dockerfile builder. Health check `GET /healthz`. Attach a volume at `/data`.

No variables are required to boot. Railway injects `PORT`. The image listens on `0.0.0.0:$PORT` and sets `PB_DATA_DIR=/data`.

| | |
|---|---|
| Health check | `GET /healthz` |
| Volume | `/data` — PocketBase SQLite + files. Without it, notes and the admin user reset on every deploy. |
| `PORT` | Injected by Railway. Do not hard-code. |

`railway.toml` and `railway.json` already set the Dockerfile builder and `/healthz` health check.

### Optional

| Variable | Why |
|---|---|
| `PB_ADMIN_EMAIL` | Create the first PocketBase superuser on boot (`/_/`) |
| `PB_ADMIN_PASSWORD` | Pair for `PB_ADMIN_EMAIL`. Ignored if email is unset. |
| `PB_DATA_DIR` | Override the data directory. Default `/data` in the image, `./pb_data` locally |

See `TEMPLATE.md` for the marketplace / composer shape (source = this public repo, root directory).
