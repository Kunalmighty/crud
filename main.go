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
	Balance int64
}

type Totals struct {
	Worth       int64
	Assets      int64
	Liabilities int64
}

/**
* Connect to the database
* @param none
* @return pointer to database
 */
func dbConn() (db *sql.DB) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	return db
}

//Template variable initialization
var tmpl = template.Must(template.ParseGlob("form/*"))

/**
* Delivers the table contents to be displayed
* @param HTTP interfaces
* @return none
 */
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
		var balance int64
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

/**
* Delivers the table row to be displayed
* @param HTTP interfaces
* @return none
 */
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
		var balance int64
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

/**
* Delivers the table totals to be displayed
* @param HTTP interfaces
* @return none
 */
func Show2(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	tot := Totals{}

	selDB, err := db.Query("SELECT SUM(balance) AS atotal FROM al WHERE asslia = 'asset'")
	if err != nil {
		panic(err.Error())
	}
	lelDB, err := db.Query("SELECT SUM(balance) AS ltotal FROM al WHERE asslia = 'liability'")
	if err != nil {
		panic(err.Error())
	}

	var assets int64
	for selDB.Next() {
		selDB.Scan(&assets)
	}
	var liabilities int64
	for lelDB.Next() {
		lelDB.Scan(&liabilities)
	}
	tot.Worth = assets - liabilities
	tot.Assets = assets
	tot.Liabilities = liabilities

	tmpl.ExecuteTemplate(w, "Show2", tot)
	defer db.Close()
}

/**
* Works with the Insert function and New.tmpl to insert new row in table
* @param HTTP interfaces
* @return none
 */
func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

/**
* Works with the Update function and Edit.tmpl to edit a table entry
* @param HTTP interfaces
* @return none
 */
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
		var balance int64
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

/**
* Works with the New function and New.tmpl to insert new row in table
* @param HTTP interfaces
* @return none
 */
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

/**
* Works with the Edit function and Edit.tmpl to edit a table entry
* @param HTTP interfaces
* @return none
 */
func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		typeo := r.FormValue("asslia")
		balance := r.FormValue("balance")
		id := r.FormValue("uid")
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

/**
* Delete a row from table, as selected through Index.tmpl
* @param HTTP interfaces
* @return none
 */
func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	rec := r.URL.Query().Get("name")
	delForm, err := db.Prepare("DELETE FROM al WHERE name=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(rec)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

/**
* Get Heroku port
* @param none
* @return port string
 */
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
	http.HandleFunc("/show2", Show2)

	fmt.Println("listening...")
	err := http.ListenAndServe(GetPort(), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
