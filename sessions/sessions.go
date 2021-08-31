package sessions

import (
	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("k3y_5ecreT"))		// create an object to configure how sessions are stored
													// byte array used as a key to sign our cookies - any data we store in our sessions
													// gorilla sessions package ensure that our application only accept cookies were signed with our key
