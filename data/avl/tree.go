package avl

import (
	"errors"
	"strings"

	"github.com/go-qiu/passer-auth-service/data/models"
	"github.com/go-qiu/passer-auth-service/data/stack"
)

/*
	avl tree struct
*/
type AVL struct {
	root *BinaryNode
}

var (
	ErrEmptyTree             error = errors.New("[AVL]: tree is empty")
	ErrNodeNotFound          error = errors.New("[AVL]: node not found")
	ErrDuplicatedNode        error = errors.New("[AVL]: found duplicated parcel job record")
	ErrEmptyNodeItemStatus   error = errors.New("[AVL]: item status code cannot be empty")
	ErrInvalidNodeItemStatus error = errors.New("[AVL]: item status code is invalid")
)

/*
	Wrapper function to instantiate an AVL tree (in-memory)
*/
func New() *AVL {
	return &AVL{root: nil}
}

/*
	Wrapper function to insert a new node into the AVL tree.
*/
func (tree *AVL) InsertNode(item models.User) error {

	root, err := insertNode(tree.root, item)
	if err != nil {
		return err
	}

	tree.root = root
	return nil
}

/*
	Wrapper function for traversing the AVL tree, in-order manner.
*/
// func (tree *AVL) Traverse() {

// 	// ok. tree is not empty

// 	// declare a stack to cache
// 	// the Job record (in each node)
// 	// in in-order traversal sequence.

// 	// stackItems := stack{top: nil, size: 0}
// 	stackItems := stack.New()

// 	traverse(tree.root, &stackItems)

// 	// items (i.e. Job records) will be listed in desc order
// 	// stackItems.ListAllNodes()

// 	// stackItemsAsc := stack{top: nil, size: 0}
// 	stackItemsAsc := stack.New()

// 	// reverse the order of the stack to get asc order
// 	for i := 0; i < stackItems.GetSize(); i++ {
// 		item, _ := stackItems.Pop()
// 		stackItemsAsc.Push(item)
// 	}

// 	// print out all the nodes in the stack
// 	stackItemsAsc.ListAllNodes()
// }

func (tree *AVL) ListAllNodes(s *stack.Stack) error {

	// ok. tree is not empty

	// declare a stack to cache
	// the Job record (in each node)
	// in in-order traversal sequence.

	// stackItems := stack{top: nil, size: 0}
	// stackItems := stack.New()

	// traverse all the nodes on the tree (in-order sequence)
	// and cache the outcome into the stack passed in.
	traverse(tree.root, s)
	return nil
}

func (tree *AVL) ListSelectedNodesByStatus(status string, descending bool) error {

	// declare a stack to cache
	// the Job record (in each node)
	// in in-order traversal sequence.

	// stackItems := stack{top: nil, size: 0}
	stackItems := stack.New()

	// traverse all the nodes on the tree (in-order sequence)
	traverse(tree.root, &stackItems)

	// items (i.e. Job records) will be listed in desc order
	if descending {

		// to list in descending order
		err := stackItems.ListSelectedNodesByStatus(status)
		if err != nil {
			return err
		}

	} else {

		// to list in ascending order
		// stackItemsAsc := stack{top: nil, size: 0}
		stackItemsAsc := stack.New()

		// reverse the order of the stack to get asc order
		for stackItems.GetSize() > 0 {
			item, _ := stackItems.Pop()
			// if item.Status == status {
			// 	stackItemsAsc.Push(item)
			// }
			stackItemsAsc.Push(item)
		}

		// print out all the nodes in the stack
		// err := stackItemsAsc.ListAllNodes()
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}

/*
	Wrapper function to find a specific node by id
*/
func (tree *AVL) Find(id string) *BinaryNode {

	return findNode(tree.root, id)
}

func (tree *AVL) Remove(id string) error {
	root, err := removeNode(tree.root, id)
	if err != nil {
		return err
	}
	tree.root = root
	return nil
}

/*
	Wrapper function to update a specific node (identified by id) on the tree
*/
func (tree *AVL) UpdateStatus(id string, status string) (BinaryNode, error) {

	// exceptions handling

	switch strings.ToUpper(status) {
	case "NEW", "READY", "ARRIVED", "COMPLETED":

		// valid status code

		// find the node
		found := findNode(tree.root, id)
		if found == nil {
			// not found
			return *found, ErrNodeNotFound
		}

		// ok. found node.
		// update the node.
		item := found.item
		item.IsActive = true

		err := updateNode(&found, item)
		if err != nil {
			return BinaryNode{}, err
		}

		return *found, nil
	default:
		return BinaryNode{}, errors.New("[AVL]: invalid item status")
	}

}
