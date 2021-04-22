package internal

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* swagger:route POST /items item PostItem
Add a new item to the todo list.

The item is given an ID on insertion to the data store.

	Consumes:
	- application/json

	Produces:
	- application/json

	Schemes: https

	Security:
		basicAuth:
			type: http
			scheme: basic

	Responses:
		default: genericError
		200: item
*/
func (s server) PostItem(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// parse the input (no validation is done here)
	var in PostItemRequest
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// insert into the data store
	q, err := s.db.Insert(in.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// marshals it into bytes to respond with
	b, err := json.Marshal(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

// PostItemRequest allows a user to post a todo item.
// It doesn't allow a user to post a done todo item.
// To mark an item as done, use PutItem to updated it.
//
// swagger:response item
type PostItemRequest struct {
	// Required: true
	// Example: do the laundry
	Description string `json:"item"`
}
