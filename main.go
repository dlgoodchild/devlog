package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

const (
	DBHost     = "127.0.0.1"
	DBPort     = ":3306"
	DBUser     = "root"
	DBPass     = ""
	DBDatabase = "devlog"

	ServePort = ":8080"

	SpecificPostBySlugQuery = "SELECT `guid`, `title`, `content`, `created_at` FROM `post` WHERE `guid`=?"
	LatestPostsQuery = "SELECT `guid`, `title`, `content`, `created_at` `date` FROM `post` ORDER BY `created_at` DESC LIMIT 5"
)

var database *sql.DB

type Post struct {
	Title      string
	RawContent string
	Content    template.HTML
	Date       string
	Slug       string
}

func redirectHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", 302)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	var post = []Post{}
	rows, err := database.Query(LatestPostsQuery)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	defer rows.Close();
	for rows.Next() {
		currPost := Post{}
		rows.Scan(&currPost.Slug, &currPost.Title, &currPost.RawContent, &currPost.Date)
		currPost.Content = template.HTML(currPost.RawContent)
		post = append(post, currPost)
	}

	// is there a better way to temporarily pass empty data to a template?
	/*
	type Empty struct {
	}
	empty := Empty{}

	t, _ := template.ParseFiles("templates/home.html")
	t.Execute(w, empty)
	*/

	t, _ := template.ParseFiles("templates/home.html")
	t.Execute(w, post)
}

func servePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGuid := vars["guid"]

	currPost := Post{}
	fmt.Println("Requesting post: " + postGuid)

	err := database.QueryRow(SpecificPostBySlugQuery, postGuid).Scan(&currPost.Slug, &currPost.Title, &currPost.RawContent, &currPost.Date)
	currPost.Content = template.HTML(currPost.RawContent)

	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)

		log.Println("Couldn't get post: " + postGuid)
		log.Println(err.Error())
		return
	}

	t, _ := template.ParseFiles("templates/post.html")
	t.Execute(w, currPost)
}

func main() {
	dbConn := fmt.Sprintf("%s:%s@tcp(%s%s)/%s", DBUser, DBPass, DBHost, DBPort, DBDatabase)
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Println("Couldn't connect to DB!")
		log.Println(err.Error())
	}
	database = db

	routes := mux.NewRouter()
	routes.HandleFunc("/post/{guid:[a-zA-Z0-9]+[a-zA-Z0-9\\-]+}", servePost)
	routes.HandleFunc("/", redirectHome)
	routes.HandleFunc("/home", serveHome)

	http.Handle("/", routes)
	http.ListenAndServe(ServePort, nil)
}
