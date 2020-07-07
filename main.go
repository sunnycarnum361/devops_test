package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

type config struct {
	DatabaseHostAddr string `env:"DATABASE_HOSTADDR"`
}

func main() {
	var conf config
	// Read environment variables
	err := env.Parse(&conf)
	if err != nil {
		log.Panic(err)
	}

	// Open database
	db, err = sql.Open("mysql", conf.DatabaseHostAddr)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	log.Println(db.Stats())
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
	}
	log.Fatal(srv.ListenAndServe())
}

func index(w http.ResponseWriter, r *http.Request) {
	p, err := getData()
	if err != nil {
		log.Println(err.Error())
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
	}
	writeJSON(w, p, http.StatusOK)
}

type Person struct {
	ID      int64
	Name    string
	Age     int64
	Address string
}

func getData() (*Person, error) {
	q := `SELECT id, name, age, address FROM person WHERE id = ?`
	var p Person
	id := "10"
	err := db.QueryRow(q, id).Scan(
		&p.ID,
		&p.Name,
		&p.Age,
		&p.Address,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func writeJSON(w http.ResponseWriter, val interface{}, code ...int) {

	b, err := json.Marshal(val)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(code) > 0 {
		w.WriteHeader(code[0])
	}

	w.Write(b)
}
