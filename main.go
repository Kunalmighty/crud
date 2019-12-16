package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Record struct {
	Id      int
	Name    string
	Type    string
	Balance uint64
}

func dbConn() (db *sql.DB) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM al ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	rec := Record{}
	res := []Record{}
	for selDB.Next() {
		var id int
		var name string
		var asslia string
		var balance uint64
		err = selDB.Scan(&id, &asslia, &balance, &name)
		if err != nil {
			panic(err.Error())
		}
		rec.Id = id
		rec.Name = name
		rec.Type = asslia
		rec.Balance = balance
		res = append(res, rec)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM al WHERE id= $1", nId)
	if err != nil {
		panic(err.Error())
	}
	rec := Record{}
	for selDB.Next() {
		var id int
		var name string
		var asslia string
		var balance uint64
		err = selDB.Scan(&id, &asslia, &balance, &name)
		if err != nil {
			panic(err.Error())
		}
		rec.Id = id
		rec.Name = name
		rec.Type = asslia
		rec.Balance = balance
	}
	tmpl.ExecuteTemplate(w, "Show", rec)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM al WHERE id=$1", nId)
	if err != nil {
		panic(err.Error())
	}
	rec := Record{}
	for selDB.Next() {
		var id int
		var name string
		var asslia string
		var balance uint64
		err = selDB.Scan(&id, &asslia, &balance, &name)
		if err != nil {
			panic(err.Error())
		}
		rec.Id = id
		rec.Name = name
		rec.Type = asslia
		rec.Balance = balance
	}
	tmpl.ExecuteTemplate(w, "Edit", rec)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		nameo := r.FormValue("name")
		typeo := r.FormValue("asslia")
		balanceo := r.FormValue("balance")
		insForm, err := db.Prepare("INSERT INTO al(name, asslia, balance) VALUES( $1 , $2 , $3 )")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(nameo, typeo, balanceo)
		log.Println("INSERT: Name: " + nameo + " | Type: " + typeo)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		typeo := r.FormValue("type")
		balance := r.FormValue("balance")
		id := r.FormValue("id")
		insForm, err := db.Prepare("UPDATE al SET name=$1, asslia=$2, balance=$3 WHERE id=$4")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, typeo, balance, id)
		log.Println("UPDATE: Name: " + name + " | Type: " + typeo)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	rec := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM al WHERE id=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(rec)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello. This is our first Go web app on Heroku!")
}

func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func main() {

	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)

	fmt.Println("listening...")
	err := http.ListenAndServe(GetPort(), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	db := dbConn()
	_, err = db.Exec(`CREATE TYPE types AS ENUM ('asset', 'liability');`)
	if err != nil {
		log.Printf("Error creating type enum: %q \n", err)
	}

	/*
			_, err = db.Exec(`
		    CREATE TABLE IF NOT EXISTS al (
			  id SERIAL,
			  asslia TYPES,
			  balance MONEY,
			  name VARCHAR(64) NOT NULL UNIQUE,
		      CHECK (CHAR_LENGTH(TRIM(name)) > 0)
		    );`)
			if err != nil {
				log.Printf("Error creating table: %q \n", err)
			}
	*/

}
