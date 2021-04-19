package internal

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* swagger:route GET /items item listItems
Retrieves a list of todo items

Retrieves a list of todo items, the list is not filtered nor paginated.

	Consumes:
	- application/json

	Produces:
	- application/json

	Schemes: https

	Responses:
		default: genericError
		200: listItemsResponse
*/
func (s server) listItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// lists everything in the datastore
	is, err := s.db.List()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// marshals it into bytes to respond with
	b, err := json.Marshal(&is)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(b)
}

// ListItemsResponse is the response from ListItems
//
// swagger:response listItemsResponse
type ListItemsResponse struct {
	// A list of todo items and their done state
	Items []Item `json:"items"`
}
