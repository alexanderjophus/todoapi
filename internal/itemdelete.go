package internal

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* swagger:route DELETE /items/{itemID} item deleteItem
Deletes an item from the todo list given an ID.

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
		204:
*/
func (s server) deleteItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = idParam{
		ID: ps.ByName("id"),
	}
	s.db.Delete(params.ID)
	w.WriteHeader(http.StatusNoContent)
}
