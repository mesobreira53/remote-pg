package main

import (
	"net/http"
)

// The serverError helper writes a log entry at the Error level (including the request method and URIas attributes)
// and then sends a generic 500 Internal Server Error response to the client.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.logger.Error(err.Error(), "method:", method, "uri:", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper Sends a specific status code and corrresponding description
// to the user. For instance, we will send resposnes like "400 Bad Request" or "404 Not Found"
// when there is a problem with the request that the user sent
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
