package mem

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/trelore/todoapi/internal"
)

// Memory is an in Memory implementation of Datastore
type Memory struct {
	items map[uuid.UUID]*internal.Item
}

// New creates a new in memory data store
func New() *Memory {
	return &Memory{
		items: map[uuid.UUID]*internal.Item{},
	}
}

// Insert implements the interface
func (m *Memory) Insert(description string) (*internal.Item, error) {
	i := &internal.Item{
		ID:          uuid.New(),
		Description: description,
		Done:        false,
	}
	m.items[i.ID] = i

	return i, nil
}

// List implements the interface
func (m *Memory) List() ([]*internal.Item, error) {
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
func (m *Memory) Get(id string) (*internal.Item, error) {
	item := m.items[uuid.MustParse(id)]
	if item == nil {
		return nil, ErrNoData
	}

	return item, nil
}

// Delete implements the interface
func (m *Memory) Delete(id string) error {
	i, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	m.items[i] = nil
	return nil
}

// Upsert implements the interface
func (m *Memory) Upsert(id string, item *internal.Item) (_ *internal.Item, err error) {
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
