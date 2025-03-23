package main

import (
	//"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"binai.net/internal/models"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	lots           models.LotModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	companyName    string
}

// Фильтр дата окончания прием заявок (статус прием заявок)
// хостинг запустить

func main() {
	addr := flag.String("addr", ":443", "HTTP network address")
	// dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	//dsn := fmt.Sprintf("user=%s dbname=%s host=%s password=%s sslmode=disable", "postgres", "TamaqQr", "localhost", "349349")
	dsn := fmt.Sprintf("user=%s dbname=%s host=%s password=%s sslmode=disable port=%d", "baha", "binai", "postgres", "adminadmin1", 5433)

	// dsn := fmt.Sprintf("user=%s dbname=%s host=%s password=%s sslmode=disable",
	// 	os.Getenv("DB_USER"),
	// 	os.Getenv("DB_NAME"),
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_PASSWORD"))

	flag.Parse()
	// fmt.Println("!й")
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// log.Printf("DB_USER: %s, DB_NAME: %s, DB_HOST: %s, DB_PORT: %s", os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	db, err := openDB(dsn)
	if err != nil {
		fmt.Println("ERROR", err)
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	sessionManager.Cookie.Secure = true

	companyName := "aetz"

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		lots:           &models.LotModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		companyName:    companyName,
	}

	// tlsConfig := &tls.Config{
	// 	CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	// }

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
		// TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.users.Insert("aetz", "test@gmail.com", "qwe")
	if err != nil {
		app.infoLog.Print("not signup")
	}

	infoLog.Printf("Starting server on https://binai.kz LO%s")
	err = srv.ListenAndServeTLS("./tls/fullchain.pem", "./tls/privkey.pem")
	infoLog.Printf("Starting server on http://binai.kz%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
