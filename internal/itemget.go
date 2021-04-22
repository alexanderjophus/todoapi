package internal

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/trelore/todoapi/internal/datastores"
)

/* swagger:route GET /items/{itemID} item GetItem
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
func (s server) GetItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = idParam{
		ID: ps.ByName("id"),
	}
	if _, err := uuid.Parse(params.ID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	i, err := s.db.Get(params.ID)
	if err != nil {
		if err == datastores.ErrNoData {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		// better error handling (404)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// marshals it into bytes to respond with
	b, err := json.Marshal(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
