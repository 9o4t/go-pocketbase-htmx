package main

import (
	"html"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	dataDir := os.Getenv("PB_DATA_DIR")
	if dataDir == "" {
		if st, err := os.Stat("/data"); err == nil && st.IsDir() {
			dataDir = "/data"
		} else {
			dataDir = "./pb_data"
		}
	}

	app := pocketbase.NewWithConfig(pocketbase.Config{DefaultDataDir: dataDir})

	app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		if err := ensureNotes(e.App); err != nil {
			return err
		}
		return ensureAdmin(e.App, os.Getenv("PB_ADMIN_EMAIL"), os.Getenv("PB_ADMIN_PASSWORD"))
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.BindFunc(func(e *core.RequestEvent) error {
			p := e.Request.URL.Path
			if strings.HasPrefix(p, "/api") || strings.HasPrefix(p, "/_/") || p == "/_/" {
				return e.Next()
			}
			if se.App.IsDev() && strings.HasPrefix(p, "/_") {
				return e.Next()
			}
			handler(se.App).ServeHTTP(e.Response, e.Request)
			return nil
		})
		return se.Next()
	})

	if len(os.Args) == 1 {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8090"
		}
		os.Args = append(os.Args, "serve", "--http=0.0.0.0:"+port)
	}
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func ensureNotes(app core.App) error {
	if _, err := app.FindCollectionByNameOrId("notes"); err == nil {
		return nil
	}
	col := core.NewBaseCollection("notes")
	col.Fields.Add(&core.TextField{Name: "title", Required: true, Presentable: true})
	col.Fields.Add(&core.TextField{Name: "body"})
	col.Fields.Add(&core.AutodateField{Name: "created", OnCreate: true})
	col.Fields.Add(&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true})
	return app.Save(col)
}

func ensureAdmin(app core.App, email, password string) error {
	if strings.TrimSpace(email) == "" || password == "" {
		return nil
	}
	if _, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, email); err == nil {
		return nil
	}
	col, err := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
	if err != nil {
		return err
	}
	rec := core.NewRecord(col)
	rec.SetEmail(email)
	rec.SetPassword(password)
	return app.Save(rec)
}

func handler(app core.App) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(pageHTML(listNotesHTML(app))))
	})
	mux.HandleFunc("GET /notes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(listNotesHTML(app)))
	})
	mux.HandleFunc("POST /notes", func(w http.ResponseWriter, r *http.Request) {
		title := strings.TrimSpace(r.FormValue("title"))
		body := strings.TrimSpace(r.FormValue("body"))
		if title != "" {
			if err := addNote(app, title, body); err != nil {
				log.Printf("note: %v", err)
			}
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(listNotesHTML(app)))
	})
	return mux
}

func addNote(app core.App, title, body string) error {
	col, err := app.FindCollectionByNameOrId("notes")
	if err != nil {
		return err
	}
	rec := core.NewRecord(col)
	rec.Set("title", title)
	rec.Set("body", body)
	return app.Save(rec)
}

func listNotesHTML(app core.App) string {
	col, err := app.FindCollectionByNameOrId("notes")
	if err != nil {
		return `<p class="empty">Notes collection is not ready.</p>`
	}
	var recs []*core.Record
	if err := app.RecordQuery(col).OrderBy("created DESC").Limit(50).All(&recs); err != nil || len(recs) == 0 {
		return `<p class="empty">No notes yet. Add one above.</p>`
	}
	var b strings.Builder
	b.WriteString("<ul>")
	for _, rec := range recs {
		b.WriteString("<li><strong>")
		b.WriteString(html.EscapeString(rec.GetString("title")))
		b.WriteString("</strong>")
		if body := rec.GetString("body"); body != "" {
			b.WriteString("<span>")
			b.WriteString(html.EscapeString(body))
			b.WriteString("</span>")
		}
		b.WriteString("</li>")
	}
	b.WriteString("</ul>")
	return b.String()
}

func pageHTML(list string) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1"/>
  <title>Go + PocketBase + HTMX</title>
  <script src="https://unpkg.com/htmx.org@2.0.4" defer></script>
  <style>
    :root { color-scheme: dark; }
    body { font-family: ui-sans-serif, system-ui, sans-serif; margin: 0; background: #0b0d10; color: #e8edf5; }
    main { max-width: 40rem; margin: 3rem auto; padding: 0 1.25rem; }
    h1 { font-size: 1.75rem; }
    p { color: #8b95a8; line-height: 1.5; }
    form { display: grid; gap: .75rem; margin: 1.5rem 0; }
    input, textarea { background: #151922; border: 1px solid #2a3140; color: inherit; padding: .6rem .7rem; border-radius: 8px; }
    button { background: #c6ff4a; color: #0b0d10; border: 0; padding: .65rem 1rem; border-radius: 999px; font-weight: 700; cursor: pointer; width: fit-content; }
    ul { list-style: none; padding: 0; display: grid; gap: .75rem; }
    li { background: #151922; border: 1px solid #2a3140; border-radius: 10px; padding: .9rem 1rem; display: grid; gap: .35rem; }
    .empty { font-style: italic; }
    code { color: #7ee0ff; }
  </style>
</head>
<body>
  <main>
    <h1>Go + PocketBase + HTMX</h1>
    <p>One binary, one volume at <code>/data</code>, one <code>/healthz</code>. Notes below are a PocketBase collection, swapped in with HTMX.</p>
    <form hx-post="/notes" hx-target="#notes" hx-swap="innerHTML">
      <input name="title" required placeholder="Note title"/>
      <textarea name="body" rows="3" placeholder="Optional body"></textarea>
      <button type="submit">Save note</button>
    </form>
    <div id="notes">` + list + `</div>
  </main>
</body>
</html>`
}
