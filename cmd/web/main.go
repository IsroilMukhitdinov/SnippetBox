package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/IsroilMukhitdinov/snippetbox/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	infoLog       *log.Logger
	errLog        *log.Logger
	staticDir     string
	templateCache map[string]*template.Template
	snippetModel  *models.SnippetModel
}

func main() {
	addr := flag.Int("addr", 9808, "network address. default: 9808")
	htmlDir := flag.String("html", "./ui/html", "html files directory. default: .ui/html")
	staticDir := flag.String("static", "./ui/static", "static files directory. default: .ui/static")
	driver := flag.String("driver", "mysql", "datbase driver name. default: mysql")
	dsn := flag.String("dsn", "web:Snippetbox_9808@/snippetbox?parseTime=true", "data source name for the database")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	templateCache, err := newTemplateCache(*htmlDir)
	if err != nil {
		errLog.Fatal(err)
	}

	db, err := openDB(*driver, *dsn)
	if err != nil {
		errLog.Printf("could not establish a connection with the %s database\n%s\n", *driver, err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := &application{
		infoLog:       infoLog,
		errLog:        errLog,
		staticDir:     *staticDir,
		templateCache: templateCache,
		snippetModel: &models.SnippetModel{
			DB: db,
		},
	}

	srv := http.Server{
		Addr:     fmt.Sprintf(":%d", *addr),
		Handler:  app.routes(),
		ErrorLog: errLog,
	}

	infoLog.Printf("server started on port %d\n", *addr)

	err = srv.ListenAndServe()
	if err != nil {
		errLog.Fatal(err)
	}
}

func openDB(driver string, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
