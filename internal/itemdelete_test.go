package internal_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/trelore/todoapi/internal"
	"github.com/trelore/todoapi/internal/datastores/mem"
)

func Test_server_deleteItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		prep   func(m *mem.Memory) (idToLookup string, err error)
		status int
	}{
		{
			name: "happy path: delete existing item",
			prep: func(m *mem.Memory) (idToLookup string, err error) {
				item, err := m.Insert("hello world")
				if err != nil {
					return "", err
				}
				return item.ID.String(), nil
			},
			status: http.StatusNoContent,
		},
		{
			name: "happy path: delete non-existent item", // this behaviour was chosen for idempotency reasons
			prep: func(m *mem.Memory) (idToLookup string, err error) {
				return uuid.NewString(), nil
			},
			status: http.StatusNoContent,
		},
		{
			name: "error: invalid id",
			prep: func(m *mem.Memory) (idToLookup string, err error) {
				return "hello world", nil
			},
			status: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mem.New()
			idToDelete, err := tt.prep(m)
			if err != nil {
				t.Fatal(err)
			}
			handler := internal.NewServer(m)
			router := httprouter.New()
			router.DELETE("/items/:id", handler.DeleteItem)

			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/items/%s", idToDelete), nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.status, rr.Code)
		})
	}
}
