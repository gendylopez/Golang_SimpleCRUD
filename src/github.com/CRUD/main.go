package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Student struct{
	Id int
	FirstName string
	LastName string
}

func dbConn() (db *sql.DB){
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "rootroot"
	dbName := "studentdb"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName);
	if err != nil{
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    selDB, err := db.Query("SELECT * FROM student ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
    }
    stud := Student{}
    res := []Student{}
    for selDB.Next() {
        var id int
        var firstName, lastName string
        err = selDB.Scan(&id, &firstName, &lastName)
        if err != nil {
            panic(err.Error())
        }
        stud.Id = id
        stud.FirstName = firstName
        stud.LastName = lastName
        res = append(res, stud)
    }
    tmpl.ExecuteTemplate(w, "Index", res)
    defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM student WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    stud := Student{}
    for selDB.Next() {
        var id int
        var firstName, lastName string
        err = selDB.Scan(&id, &firstName, &lastName)
        if err != nil {
            panic(err.Error())
        }
        stud.Id = id
        stud.FirstName = firstName
        stud.LastName = lastName
    }
    tmpl.ExecuteTemplate(w, "Show", stud)
    defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "New", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        firstName := r.FormValue("fname")
        lastName := r.FormValue("lname")
        insForm, err := db.Prepare("INSERT INTO Student(firstName, lastName) VALUES(?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(firstName, lastName)
        log.Println("INSERT: First Name: " + firstName + " | Last Name: " + lastName)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM Student WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    stud := Student{}
    for selDB.Next() {
        var id int
        var firstName, lastName string
        err = selDB.Scan(&id, &firstName, &lastName)
        if err != nil {
            panic(err.Error())
        }
        stud.Id = id
        stud.FirstName = firstName
        stud.LastName = lastName
    }
    tmpl.ExecuteTemplate(w, "Edit", stud)
    defer db.Close()
}

func Update(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        firstName := r.FormValue("fname")
        lastName := r.FormValue("lname")
        id := r.FormValue("uid")
        insForm, err := db.Prepare("UPDATE Student SET firstName=?, lastName=? WHERE id=?")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(firstName, lastName, id)
        log.Println("UPDATE: First Name: " + firstName + " | Last Name: " + lastName)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    emp := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM Student WHERE id=?")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(emp)
    log.Println("DELETE")
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func main(){
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
    http.HandleFunc("/new", New)
    http.HandleFunc("/insert", Insert)
    http.HandleFunc("/edit", Edit)
    http.HandleFunc("/update", Update)
    http.HandleFunc("/delete", Delete)
    log.Fatal(http.ListenAndServe(":8080", nil))
}