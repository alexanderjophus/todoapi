package internal

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* swagger:route PUT /items/{itemID} item putItem
Updates an item from the todo list given an ID.

Sets the description and done flag to the input.

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
func (s server) putItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = idParam{
		ID: ps.ByName("id"),
	}
	// parse the input (no validation is done here)
	var in Item
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	item, err := s.db.Upsert(params.ID, &in)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// marshals it into bytes to respond with
	b, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(b)
}
