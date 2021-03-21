package main

import (
	"database/sql"
	"fmt"
	"html/template"

	"abc/dbb"

	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

// Person type for export
type Person struct {
	Id        int
	Firstname string
	Lastname  string
	Email     string
	Gender    string
	Country   string
}

var (
	templates = template.Must(template.ParseGlob("templates/*"))
)

var db *sql.DB

func main() {
	var err error
	db, err = dbb.Connect()
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.NotFoundHandler = http.HandlerFunc(notfound)
	r.HandleFunc("/", home)
	r.HandleFunc("/insert", insertcreate)
	r.HandleFunc("/insert/process", insertcreateprocess)
	r.HandleFunc("/delete", delete)
	r.HandleFunc("/update", updatestart)
	r.HandleFunc("/updateprocess", updateend)
	http.Handle("/", r)
	//r.HandleFunc("/notfound", notfound)
	fmt.Println("Your website is now online.......")
	http.ListenAndServe(":3000", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	// inserting values
	updtstmt := `SELECT * FROM person`
	rows, err := db.Query(updtstmt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	//
	var persons []Person
	for rows.Next() {
		var p Person
		err = rows.Scan(&p.Id, &p.Firstname, &p.Lastname, &p.Email, &p.Gender, &p.Country)
		if err != nil {
			panic(err)
		}
		persons = append(persons, p)
	}
	templates.ExecuteTemplate(w, "browseall.html", persons)

}

// get request
func insertcreate(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles("home.html"))
	templates.Execute(w, nil)
}

func insertcreateprocess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		fmt.Println("get issues")
		return
	}
	p := Person{}
	p.Firstname = r.FormValue("firstname")
	p.Lastname = r.FormValue("lastname")
	p.Email = r.FormValue("email")
	p.Gender = r.FormValue("gender")
	p.Country = r.FormValue("country")

	if p.Firstname == "" || p.Lastname == "" || p.Email == "" || p.Gender == "" || p.Country == "" {
		templates.ExecuteTemplate(w, "home.html", "Data was wrongly entered or omitted")
		fmt.Println("some data were not entered or wrongly inputed")
		return
	}

	// inserting values
	insertstmt := `insert into person (firstname, lastname, email, gender, country) values ($1, $2, $3, $4, $5)`
	_, err := db.Exec(insertstmt, p.Firstname, p.Lastname, p.Email, p.Gender, p.Country)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// confirm insertion
	templates.ExecuteTemplate(w, "home.html", "Data successfully entered!")

}

func delete(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "insert.html", nil)
}

//
func updatestart(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}
	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	row := db.QueryRow(`SELECT * FROM person WHERE id = $1;`, id)
	p := Person{}
	err := row.Scan(&p.Id, &p.Firstname, &p.Lastname, &p.Email, &p.Gender, &p.Country)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		http.Redirect(w, r, "/", 307)
		return
	}
	templates.ExecuteTemplate(w, "update.html", p)

}

func updateend(w http.ResponseWriter, r *http.Request) {
	fmt.Println("part 2")
	if r.Method != "POST" {
		fmt.Println("part 2a")
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	p := Person{}
	p.Firstname = r.FormValue("firstname")
	p.Lastname = r.FormValue("lastname")
	p.Email = r.FormValue("email")
	p.Gender = r.FormValue("gender")
	p.Country = r.FormValue("country")

	if p.Firstname == "" || p.Lastname == "" || p.Email == "" || p.Gender == "" || p.Country == "" {
		fmt.Println("All the fields have not been filled.")
		templates.ExecuteTemplate(w, "updateend.html", "Your data was NOT updated. Check back to do that")
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return

	}
	fmt.Println("part 33333")

	_, err := db.Exec("UPDATE person SET id = $1, firstname = $2, lastname = $3, email = $4, gender = $5, country = $6 WHERE id = $1;", p.Id, p.Firstname, p.Lastname, p.Email, p.Gender, p.Country)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		templates.ExecuteTemplate(w, "updateend.html", "it did not update records")

		return
	}
	templates.ExecuteTemplate(w, "updateend.html", "you have successfully updated your data")

}

func notfound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	templates.ExecuteTemplate(w, "404page.html", nil)
}
