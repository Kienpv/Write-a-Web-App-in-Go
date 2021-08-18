package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"github.com/go-redis/redis"
)

var client *redis.Client
var templates *template.Template

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost: 6379", 	// where the redis serve is. In this example, we're running it on the same machine 
									// as our web application on port 6379 which is the default port for Redis.
	})
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", indexGetHandler).Methods("GET")
	r.HandleFunc("/", indexPostHandler).Methods("POST")

	fs := http.FileServer(http.Dir("./static/"))		//Directory to server files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs)) 			// to use these files server for all paths that start with the static prefix. 
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	comment, err := client.LRange("comments", 0, 10).Result()
	// due to some update -> not enough arguments in call to client.cmdable.LRange, we need to change
	// add "context" to the import list, then in the indexGetHandler add this:
	// ctx := context.TODO()
	// comments, err := client.LRange(ctx, "comments", 0, 10).Result()
	 if (err != nil) { return }
	templates.ExecuteTemplate(w, "index.html", comment)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment := r.PostForm.Get("comment_text")
	client.LPush("comments", comment)
	http.Redirect(w, r, "/", 302)		
}
