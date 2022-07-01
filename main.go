package main

import (
	"fmt"
	"net/http"
	"html/template"
	"database/sql"

	// go get -u github.com/gorilla/mux
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"

)

type Article struct {
	Id uint16
	Title, Anons, Full_text string
}

var posts = []Article{}
var showPost = Article{}

func index(w http.ResponseWriter, r *http.Request){

	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if(err != nil){
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/go-data")
	if err != nil {
		panic(err)
	}
	
	defer db.Close()

	res, err := db.Query("SELECT * FROM `articles` ")
	if err != nil {
		panic(err)
	}
	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_text)
		if err != nil {
			panic(err)
		}
		// fmt.Println(fmt.Sprintf("Post %s - id %d", post.Title, post.Id))
		posts = append(posts, post)
	}

	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request){
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")

	if(err != nil){
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request){
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Ne wse dannie sapolneny")
	} else {	
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/go-data")
		if err != nil {
			panic(err)
		}
		
		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES('%s', '%s', '%s')", title, anons, full_text))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
	
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}


}

func show_post(w http.ResponseWriter, r *http.Request){

	vars := mux.Vars(r)

	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")
	
	// w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "ID: %v\n", vars["id"])

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/go-data")
		if err != nil {
			panic(err)
		}
		
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM `articles` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}
	showPost = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_text)
		if err != nil {
			panic(err)
		}
		showPost = post
	}
	// fmt.Println(fmt.Sprintf("Post %s - id %d", post.Title, post.Id))
	t.ExecuteTemplate(w, "show", showPost)
}

func handleFunc(){

	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}

