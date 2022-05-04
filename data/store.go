// code for the data store implementation
package data

import (
	"github.com/go-qiu/passer-auth-service/data/avl"
	"github.com/go-qiu/passer-auth-service/data/models"
	"github.com/go-qiu/passer-auth-service/data/stack"
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

/*
	Wrapper function to Check-In a parcel job
*/
func (ds *DataStore) CheckIn(id string) (models.User, error) {

	// // find the parcel job record
	// found := ds.avl.Find(id)
	// if found == nil {
	// 	// not found
	// 	return ErrJobNotFound
	// }

	updated, err := ds.avl.UpdateStatus(id, "READY")
	if err != nil {
		return models.User{}, err
	}

	// ok. updated
	return updated.GetItem(), nil
}

/*
	Wrapper function to Check-Out a parcel job
*/
func (ds *DataStore) CheckOut(id string) (models.User, error) {
	updated, err := ds.avl.UpdateStatus(id, "COMPLETED")
	if err != nil {
		return models.User{}, err
	}

	// ok. updated
	return updated.GetItem(), nil
}

func (ds *DataStore) Remove(id string) error {
	err := ds.avl.Remove(id)
	if err != nil {
		return err
	}

	// ok. removed.
	return nil
}

func (ds *DataStore) ListSelectedJobsByStatus(status string) error {

	return ds.avl.ListSelectedNodesByStatus(status, false)
}
