package internal

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* swagger:route GET /items/{itemID} item getItem
Gets an item from the todo list given an ID.

	Consumes:
	- application/json

	Produces:
	- application/json

	Schemes: https

	Responses:
		default: genericError
		200: item
*/
func (s server) getItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = idParam{
		ID: ps.ByName("id"),
	}
	// better error handling (400)
	i, err := s.db.Get(params.ID)
	if err != nil {
		// better error handling (404)
		http.Error(w, err.Error(), 500)
		return
	}
	// marshals it into bytes to respond with
	b, err := json.Marshal(i)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(b)
}
