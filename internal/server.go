package internal

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/trelore/todoapi/internal/middlewares"
)

// Datastore is an interface allowing the storage of
type Datastore interface {
	// Insert an item description into the datastore
	Insert(description string) (*Item, error)
	// list all the todo items in the datastore
	List() ([]*Item, error)
	// Get a singular item from the todo list
	Get(id string) (*Item, error)
	// Delete a todo item by id
	Delete(id string) error
	// Upsert a todo item - set it to done true/false & rename it
	Upsert(id string, item *Item) (*Item, error)
}

// item represents a todo item within the service
//
// swagger:response item
type Item struct {
	// in: id
	// Example: 49830a50-6e63-4435-91b1-632607ba56bd
	ID uuid.UUID `json:"id"`
	// Required: true
	// Example: do the laundry
	Description string `json:"item"`
	// Required: true
	// Example: false
	Done bool `json:"done"`
}

// server to hold the data store and the routing for the server
type server struct {
	db     Datastore
	router *httprouter.Router
}

// NewServer creates a new server struct, initialised with the routing set
func NewServer(db Datastore) server {
	s := server{
		router: httprouter.New(),
		db:     db,
	}

	s.router.POST("/items", middlewares.BasicAuth(s.PostItem))
	s.router.GET("/items", s.ListItems)
	s.router.GET("/items/:id", s.GetItem)
	s.router.DELETE("/items/:id", middlewares.BasicAuth(s.DeleteItem))
	s.router.PUT("/items/:id", middlewares.BasicAuth(s.putItem))
	return s
}

// ServeHTTP implements the handler interface
func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// swagger:parameters getItem deleteItem putItem
type idParam struct {
	/* ID of todo item that needs to be fetched
	Required: true
	In: path
	*/
	ID string `json:"id"`
}

// APIError example
//
// swagger:response genericError
type APIError struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}
