package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	_ "github.com/mattn/go-sqlite3"
)

var form = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Check-In</title>
		<style>
			html, body {
				min-height: 100vh;
			}
			input[type="text"], input[type="submit"] {
				font-size: 1.3em;
			}
			#content {
				min-height: 100vh;
				max-width: 600px;
				margin: auto;
				display: flex;
				flex-direction: column;
				justify-content: center;
				align-items: center;
			}
			#form {
				text-align: center;
			}
			#form p {
				margin-bottom: 15px;
				font-size: 2em;
			}
		</style>
    </head>
    <body>
	<div id="content">
		<div id="form">
			<p>%s Check-In</p>
			<form method="post">
				<input type="text" id="id" name="id" placeholder="Scan your ID badge" autofocus>
				<input type="hidden" id="type" name="type" value="%s">
				<input type="submit" value="Submit">
			</form>
			<p>%s</p>
		</div>
	</div>
    </body>
</html>
`

// FormTemplate is used to fill in the form template
type FormTemplate struct {
	Type string
}

// Server is the quickscan server
type Server struct {
	db *sql.DB
}

// FormHandler handles form submissions
func (s *Server) FormHandler(w http.ResponseWriter, r *http.Request) {
	feedback := ""
	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
		typ := r.FormValue("type")
		id := r.FormValue("id")
		log.Printf("inserting: type=%s, id=%s\n", typ, id)
		if _, err := s.db.Exec("insert into checkin(typ, id, time) values(?, ?, ?)", typ, id, time.Now().String()); err != nil {
			log.Println("could not execute sql:", err)
		}
		feedback = fmt.Sprintf("ID %s submitted", id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	typ := r.URL.Query().Get("type")
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(form, strings.ToTitle(strings.ReplaceAll(typ, "-", " ")), typ, feedback)))
}

func main() {
	dbpath := os.Getenv("DB_PATH")
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", dbpath))
	if err != nil {
		panic(fmt.Errorf("could not create db: %w", err))
	}
	log.Println("opened database:", dbpath)

	if _, err := db.Exec("create table if not exists checkin(internal_id integer primary key autoincrement, typ text not null, id text not null, time text not null)"); err != nil {
		panic(fmt.Errorf("could not create table: %w", err))
	}

	s := &Server{db: db}

	hndlrs := handlers.CombinedLoggingHandler(os.Stdout, http.HandlerFunc(s.FormHandler))
	if os.Getenv("PROXY_HEADERS") == "true" {
		hndlrs = handlers.ProxyHeaders(hndlrs)
	}

	addr := os.Getenv("LISTEN_ADDR")
	log.Println("starting:", addr)
	if err := http.ListenAndServe(addr, hndlrs); err != nil {
		log.Println("http server failed:", err)
	}
}
