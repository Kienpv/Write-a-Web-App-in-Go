package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"github.com/go-redis/redis"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var client *redis.Client
var store = sessions.NewCookieStore([]byte(""))		// create an object to configure how sessions are stored
													// byte array used as a key to sign our cookies - any data we store in our sessions
													// gorilla sessions package ensure that our application only accept cookies were signed with our key
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
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
	// r.HandleFunc("/test-login", test_loginGetHandler).Methods("GET")

	fs := http.FileServer(http.Dir("./static/"))		//Directory to server files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs)) 			// to use these files server for all paths that start with the static prefix. 
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	_, ok := session.Values["username"]
	if !ok {
		http.Redirect(w, r, "/login", 302)
		return
	}

	comment, err := client.LRange("comments", 0, 10).Result()
	// due to some update -> not enough arguments in call to client.cmdable.LRange, we need to change
	// add "context" to the import list, then in the indexGetHandler add this:
	// ctx := context.TODO()
	// comments, err := client.LRange(ctx, "comments", 0, 10).Result()
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return 
	}
	templates.ExecuteTemplate(w, "index.html", comment)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment := r.PostForm.Get("comment_text")
	err := client.LPush("comments", comment).Err()
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return 
	}
	http.Redirect(w, r, "/", 302)		
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	hash, err := client.Get("user:" + username).Bytes()
	if err == redis.Nil {
		templates.ExecuteTemplate(w, "login.html", "Unknown user")
		return
	}
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return 
	}

	password := r.PostForm.Get("password")
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		templates.ExecuteTemplate(w, "login.html", "Invalid password")
		return
	}
	
	session, _ := store.Get(r, "session")
	session.Values["username"] = username
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)		
}

// func test_loginGetHandler(w http.ResponseWriter, r *http.Request) {
// 	session, _ := store.Get(r, "session")
// 	untyped, ok := session.Values["username"]
// 	if !ok { 
// 		return
// 	}
// 	username, ok := untyped.(string)
// 	if !ok {
// 		return
// 	}
// 	w.Write([]byte(username))
// }

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register.html", nil)
}

func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return 
	}
	err = client.Set("user:" + username, hash, 0).Err()
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return 
	}
	http.Redirect(w, r, "/login", 302)
}
