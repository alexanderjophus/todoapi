package internal_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/trelore/todoapi/internal"
	"github.com/trelore/todoapi/internal/datastores/mem"
)

func Test_server_ListItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		prep   func(m *mem.Memory) (items []*internal.Item, err error)
		status int
	}{
		{
			name: "happy path: list existing items",
			prep: func(m *mem.Memory) (items []*internal.Item, err error) {
				item, err := m.Insert("hello world")
				if err != nil {
					return nil, err
				}
				return []*internal.Item{item}, nil
			},
			status: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mem.New()
			_, err := tt.prep(m)
			if err != nil {
				t.Fatal(err)
			}
			handler := internal.NewServer(m)
			router := httprouter.New()
			router.GET("/items", handler.ListItems)

			req, _ := http.NewRequest("GET", "/items", nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.status, rr.Code)
		})
	}
}
