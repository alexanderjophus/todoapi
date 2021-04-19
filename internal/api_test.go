package internal_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/trelore/todoapi/internal"
	"github.com/trelore/todoapi/internal/datastores/mem"
)

// basic API test to show how to use the API
func TestAPI(t *testing.T) {
	s := httptest.NewServer(internal.NewServer(mem.New()))
	defer s.Close()

	// List todo items should be ok with no items
	req, err := http.NewRequest("GET", s.URL+"/items", nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("expected no error, got: %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("status does not match")
	}
	// todo add body check

	// List todo items should be ok with no items
	req, err = http.NewRequest("POST", s.URL+"/items", bytes.NewBuffer([]byte(`{"description":"Do a backflip!"}`)))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("steve", "netherite")

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("expected no error, got: %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("status does not match")
	}
	// todo add body check

	// todo get item

	// todo update item

	// todo delete item
}
