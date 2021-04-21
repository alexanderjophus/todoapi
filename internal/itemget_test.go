package internal_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/trelore/todoapi/internal"
	"github.com/trelore/todoapi/internal/datastores/mem"
)

func Test_server_GetItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		prep   func(m *mem.Memory) (idToLookup string, err error)
		status int
	}{
		{
			name: "happy path: get existing item",
			prep: func(m *mem.Memory) (idToLookup string, err error) {
				item, err := m.Insert("hello world")
				if err != nil {
					return "", err
				}
				return item.ID.String(), nil
			},
			status: http.StatusOK,
		},
		// ! TODO
		// {
		// 	name: "error: item doesn't exist",
		// 	prep: func(m *mem.Memory) (idToLookup string, err error) {
		// 		return uuid.NewString(), nil
		// 	},
		// 	status: http.StatusNotFound,
		// },
		{
			name: "error: id is rubbish",
			prep: func(m *mem.Memory) (idToLookup string, err error) {
				return "bad_id", nil
			},
			status: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mem.New()
			idToRetrieve, err := tt.prep(m)
			if err != nil {
				t.Fatal(err)
			}
			handler := internal.NewServer(m)
			router := httprouter.New()
			router.GET("/items/:id", handler.GetItem)

			req, _ := http.NewRequest("GET", fmt.Sprintf("/items/%s", idToRetrieve), nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.status, rr.Code)
		})
	}
}
