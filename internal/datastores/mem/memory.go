package mem

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/trelore/todoapi/internal"
)

// memory is an in memory implementation of Datastore
type memory struct {
	items map[uuid.UUID]*internal.Item
}

// New creates a new in memory data store
func New() *memory {
	return &memory{
		items: map[uuid.UUID]*internal.Item{},
	}
}

// Insert implements the interface
func (m *memory) Insert(description string) (*internal.Item, error) {
	i := &internal.Item{
		ID:          uuid.New(),
		Description: description,
		Done:        false,
	}
	m.items[i.ID] = i

	return i, nil
}

// List implements the interface
func (m *memory) List() ([]*internal.Item, error) {
	items := []*internal.Item{}
	for _, v := range m.items {
		if v.ID == uuid.Nil {
			continue
		}
		items = append(items, v)
	}

	return items, nil
}

// Get implements the interface
func (m *memory) Get(id string) (*internal.Item, error) {
	item := m.items[uuid.MustParse(id)]
	if item == nil {
		return nil, ErrNoData
	}

	return item, nil
}

// Delete implements the interface
func (m *memory) Delete(id string) error {
	i, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	m.items[i] = nil
	return nil
}

// Upsert implements the interface
func (m *memory) Upsert(id string, item *internal.Item) (_ *internal.Item, err error) {
	item.ID, err = uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	m.items[item.ID] = item
	return item, nil
}

var (
	ErrNoData = fmt.Errorf("no data")
)
