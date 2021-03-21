package main

import (
	"abc/dbb"
	"database/sql"
	"fmt"
	"html/template"

	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

var db *sql.DB

// Person type for export
type Person struct {
	firstname string
	lastname  string
	email     string
	gender    string
	country   string
}

var (
	templates = template.Must(template.ParseGlob("templates/*"))
)

func main() {

	var err error
	db, err = dbb.Connect()

	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", insertcreate)
	r.HandleFunc("/process", insertcreateProcess)
	http.Handle("/", r)
	//r.HandleFunc("/notfound", notfound)
	fmt.Println("Your website is now online.......")
	http.ListenAndServe(":3000", nil)

}

func insertcreate(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.New("form").ParseFiles("insert.html"))
	templates.Execute(w, nil)
}

func insertcreateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		fmt.Println("insertcreateProcess okay")

		return
	}

	p := Person{}
	p.firstname = r.FormValue("firstname")
	p.lastname = r.FormValue("lastname")
	p.email = r.FormValue("email")
	p.gender = r.FormValue("gender")
	p.country = r.FormValue("country")

	// validate form values
	if p.firstname == "" || p.lastname == "" || p.email == "" || p.gender == "" || p.country == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		fmt.Println("validate okay")

		return
	}

	// insert values
	insertstmt := `insert into person (firstname, lastname, email, gender, country) values ($1, $2, $3, $4, $5)`
	_, err := db.Exec(insertstmt, p.firstname, p.lastname, p.email, p.gender, p.country)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// confirm insertion
	templates.ExecuteTemplate(w, "insertprocess.html", p)
}

/*

func insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		templates := template.Must(template.New("result").ParseFiles("insert.html"))

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p := Person{}
		p.firstname = r.FormValue("firstname")
		p.lastname = r.FormValue("lastname")
		p.email = r.FormValue("email")
		p.gender = r.FormValue("gender")
		p.country = r.FormValue("country")

		// validate form values
		if p.firstname == "" || p.lastname == "" || p.email == "" || p.gender == "" || p.country == "" {
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}

		var err error

		_, err = db.Exec(`insert into "person" ("firstname", "lastname", "email", "gender", "country") values ($1, $2, $3, $4, $5)`, p.firstname, p.lastname, p.email, p.gender, p.country)
		if err != nil {
			fmt.Println("okay")
		}

		if err := templates.Execute(w, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}
		return

	}
	http.Error(w, "", http.StatusBadRequest)

}

*/
