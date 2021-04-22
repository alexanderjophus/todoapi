package internal_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/trelore/todoapi/internal"
	"github.com/trelore/todoapi/internal/datastores/mem"
)

func Test_server_PostItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "happy path: post item",
			status: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := internal.NewServer(mem.New())
			router := httprouter.New()
			router.POST("/items", handler.PostItem)

			b := bytes.NewReader([]byte(`{"item": "do thing"}`))
			req, _ := http.NewRequest("POST", "/items", b)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.status, rr.Code)
		})
	}
}
