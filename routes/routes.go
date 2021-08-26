package routes

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/Write-a-Web-App-in-Go/models"
	"github.com/Write-a-Web-App-in-Go/sessions"
	"github.com/Write-a-Web-App-in-Go/utils"
	"github.com/Write-a-Web-App-in-Go/midleware"
)

func NewRoute() *mux.Router{
	r := mux.NewRouter()
	r.HandleFunc("/", midleware.AuthRequired(indexGetHandler)).Methods("GET")
	r.HandleFunc("/", midleware.AuthRequired(indexPostHandler)).Methods("POST")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
	// r.HandleFunc("/test-login", test_loginGetHandler).Methods("GET")

	fs := http.FileServer(http.Dir("./static/"))		//Directory to server files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs)) 			// to use these files server for all paths that start with the static prefix.
	return r;
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	comment, err := models.GetComments()
	// due to some update -> not enough arguments in call to client.cmdable.LRange, we need to change
	// add "context" to the import list, then in the indexGetHandler add this:
	// ctx := context.TODO()
	// comments, err := client.LRange(ctx, "comments", 0, 10).Result()
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return 
	}
	utils.ExecuteTemplate(w, "index.html", comment)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment := r.PostForm.Get("comment_text")
	err := models.PostComments(comment)
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return 
	}
	http.Redirect(w, r, "/", 302)		
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "login.html", nil)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	err := models.AuthenticateUser(username, password)

	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			utils.ExecuteTemplate(w, "login.html", "Unknown user")
		case models.ErrInvalidLogin:
			utils.ExecuteTemplate(w, "login.html", "Invalid password")
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
		}
		return 
	}
	
	session, _ := sessions.Store.Get(r, "session")
	session.Values["username"] = username
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)		
}

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "register.html", nil)
}

func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")	
	err := models.RegisterUser(username, password)
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return 
	}
	http.Redirect(w, r, "/login", 302)
}
