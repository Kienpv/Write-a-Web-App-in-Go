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
	
	r.HandleFunc("/{username}", midleware.AuthRequired(userGetHandler)).Methods("GET")
	return r;
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	update, err := models.GetAllUpdates()
	// due to some update -> not enough arguments in call to client.cmdable.LRange, we need to change
	// add "context" to the import list, then in the indexGetHandler add this:
	// ctx := context.TODO()
	// comments, err := client.LRange(ctx, "comments", 0, 10).Result()
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return 
	}
	utils.ExecuteTemplate(w, "index.html", struct {
		Tittle string
		Updates []*models.Update	
	} {
		Tittle: "all update",
		Updates: update,
	})
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	userId, ok := session.Values["user_id"].(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return 
	}
	r.ParseForm()
	body := r.PostForm.Get("update")
	
	err := models.PostUpdates(userId, body)
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return 
	}
	http.Redirect(w, r, "/", 302)		
}

func userGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	user, err := models.GetUserByUsername(username)
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return 
	}
	userId, err := user.GetUserId()
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return 
	}
	update, err := models.GetUpdates(userId)
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return 
	}
	utils.ExecuteTemplate(w, "index.html",  struct {
		Tittle string
		Updates []*models.Update	
	} {
		Tittle: username,
		Updates: update,
	})
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "login.html", nil)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	user, err := models.AuthenticateUser(username, password)

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
	session.Values["user_id"], err = user.GetUserId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
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
