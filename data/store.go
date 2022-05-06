// code for the data store implementation
package data

import (
	"errors"

	"github.com/go-qiu/passer-auth-service/data/avl"
	"github.com/go-qiu/passer-auth-service/data/models"
	"github.com/go-qiu/passer-auth-service/data/stack"
)

var (
	ErrEmptyTree             error = errors.New("[AVL]: tree is empty")
	ErrNodeNotFound          error = errors.New("[AVL]: node not found")
	ErrDuplicatedNode        error = errors.New("[AVL]: found duplicated parcel job record")
	ErrEmptyNodeItemStatus   error = errors.New("[AVL]: item status code cannot be empty")
	ErrInvalidNodeItemStatus error = errors.New("[AVL]: item status code is invalid")
)

type DataStore struct {
	avl *avl.AVL
}

func New() *DataStore {
	return &DataStore{avl: avl.New()}
}

func (ds *DataStore) InsertNode(item models.User) error {

	err := ds.avl.InsertNode(item)
	if err != nil {
		return err
	}

	// ok.
	return nil
}

func (ds *DataStore) ListAllNodes(s *stack.Stack, requireDesc bool) error {

	err := ds.avl.ListAllNodes(s)
	if err != nil {
		return err
	}

	// ok.
	if !requireDesc {
		// need it to be ascending (smallest value at top)
		stackAsc := stack.New()
		for s.GetSize() > 0 {
			item, _ := s.Pop()
			stackAsc.Push(item)
		}
		size := stackAsc.GetSize()
		top := stackAsc.GetTop()
		s.SetTop(top)
		s.SetSize(size)
	}
	return nil
}

// wrapper function to find a specific data point by id
func (ds *DataStore) Find(id string) (*avl.BinaryNode, error) {
	found := ds.avl.Find(id)
	if found == nil {
		// not found
		return found, ErrNodeNotFound
	}

	return found, nil
}

// wrapper function to remove a specific data point by id (i.e. email)
func (ds *DataStore) Remove(id string) error {

	err := ds.avl.Remove(id)
	if err != nil {
		return err
	}

	return nil
}
