// Binary Node code

package avl

import (
	"strings"

	"github.com/go-qiu/passer-auth-service/data/stack"
)

// binary node struct
type BinaryNode struct {
	id     string
	item   interface{}
	left   *BinaryNode
	right  *BinaryNode
	height int
}

func (n *BinaryNode) GetItem() interface{} {
	return n.item
}

/*
	Wrapper function to get the height of the sub-tree, wrt a specific node.
*/
func (n *BinaryNode) Height() int {
	if n == nil {
		return 0
	}

	return n.height
}

/*
	Private function to return the greater height, between left and right sub-tree height, at a specific node of interest (NoI)
*/
func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func (n *BinaryNode) updateHeight() {

	// comparasion of the maximum children node height and include the height (in the avl tree)
	// where the specific node is at.
	ht := max(n.left.Height(), n.right.Height()) + 1

	n.height = ht
}

/*
	Private function to execute a right rotation operation
	at a specific node
*/
func rightRotate(x *BinaryNode) *BinaryNode {
	y := x.left
	t := y.right

	y.right = x
	x.left = t

	x.updateHeight()
	y.updateHeight()

	return y
}

/*
	Private function to execute a left rotate operation on a specific node, x
*/
func leftRotate(x *BinaryNode) *BinaryNode {
	y := x.right
	t := y.left

	y.left = x
	x.right = t

	x.updateHeight()
	y.updateHeight()

	return y
}

func (n *BinaryNode) balanceFactor() int {
	if n == nil {
		return 0
	}

	return n.left.Height() - n.right.Height()
}

func newNode(item interface{}, id string) *BinaryNode {
	return &BinaryNode{
		id:    id,
		item:  item,
		left:  nil,
		right: nil,
	}
}

func (n *BinaryNode) Item() interface{} {
	return n.item
}

func (n *BinaryNode) Left() *BinaryNode {
	return n.left
}

func (n *BinaryNode) Right() *BinaryNode {
	return n.right
}

func insertNode(node *BinaryNode, item interface{}, id string) (*BinaryNode, error) {
	//
	if node == nil {
		// reached a leaf node
		return newNode(item, id), nil
	}

	if id == node.id {
		// a duplicate node (i.e. Job record). not allowed.
		return nil, ErrDuplicatedNode
	}

	if id > node.id {
		// new Job record id is greater than node's value
		right, err := insertNode(node.right, item, id)
		if err != nil {
			return nil, err
		}

		node.right = right

	} else if id < node.id {
		// new Job record id is lesser than node's value
		left, err := insertNode(node.left, item, id)
		if err != nil {
			return nil, err
		}

		node.left = left
	}

	// recursive. execute tree balancing at this node (if needed)
	return rotateInsert(node), nil
}

func rotateInsert(node *BinaryNode) *BinaryNode {
	// update the height on every insertion
	node.updateHeight()

	// calculate the balance factor
	bf := node.balanceFactor()

	// nodes lined-up to the left
	if bf > 1 && node.id < node.left.id {
		return rightRotate(node)
	}

	// nodes lined-up to the right
	if bf < -1 && node.id > node.right.id {
		return leftRotate(node)
	}

	// nodes lined-up to a 'less than' shape
	if bf > 1 && node.id > node.left.id {
		node.left = leftRotate(node.left)
		return rightRotate(node)
	}

	// nodes lined-up to a 'greater than' shape
	if bf < -1 && node.id < node.right.id {
		node.right = rightRotate(node.right)
		return leftRotate(node)
	}

	return node
}

/*
	Private function to traverse the avl tree in an in-order manner
*/
func traverse(node *BinaryNode, s *stack.Stack) {
	// exit condition
	if node == nil {
		return
	}

	// in-order traversing.
	// use a stack to cache the contents of the nodes as the tree is traversed.
	traverse(node.left, s)
	s.Push(node.item)
	traverse(node.right, s)
}

/*
	Private function to find a specific node by id (recursively), in the avl tree.
*/
func findNode(node *BinaryNode, id string) *BinaryNode {

	if node == nil {
		// end of search.  not found.
		return nil
	}

	if node.id == id {
		// found.
		return node
	}

	if node.id < id {
		// target id is greater than current node id
		return findNode(node.right, id)
	}

	if node.id > id {
		return findNode(node.left, id)
	}

	return nil
}

/*
	Private function to find the least valueable child node
	of a current node.
*/
func least(node *BinaryNode) *BinaryNode {
	if node == nil {
		return nil
	}

	if node.left == nil {
		return node
	}

	// recursive call
	return least(node.left)

}

func removeNode(node *BinaryNode, id string) (*BinaryNode, error) {

	if node == nil {
		return nil, ErrNodeNotFound
	}

	if id > node.id {

		// target id is greater than current node id
		right, err := removeNode(node.right, id)
		if err != nil {
			return nil, err
		}
		node.right = right

	} else if id < node.id {
		// target id is lesser than the current node id
		left, err := removeNode(node.left, id)
		if err != nil {
			return nil, err
		}

		node.left = left
	} else {
		// found.
		if node.left != nil && node.right != nil {
			// has 2 children nodes

			// find successor
			successor := least(node.right)
			// successor := greatest(node.left)
			item := successor.item

			// remove the successor
			right, err := removeNode(node.right, id)
			// left, err := removeNode(node.left, id)
			if err != nil {
				return nil, err
			}
			node.right = right
			// node.left = left

			node.item = item

		} else if node.left != nil || node.right != nil {
			// has 1 child node (left or right)
			// move the child node position to the current node

			if node.left != nil {
				node = node.left
			} else {
				// node.right is not nil
				node = node.right
			}

		} else if node.left == nil && node.right == nil {
			//  current node is a leaf node
			node = nil
		}

	}

	// return node, nil
	return rotateDelete(node), nil
}

func rotateDelete(node *BinaryNode) *BinaryNode {

	if node == nil {
		// exception handling, for the 'removal' of the
		// successor (that has no children nodes)
		return node
	}

	node.updateHeight()
	bf := node.balanceFactor()

	// nodes lined-up to the left
	if bf > 1 && node.left.balanceFactor() >= 0 {
		return rightRotate(node)
	}

	// nodes lined-up like a 'less than' shape
	if bf > 1 && node.left.balanceFactor() < 0 {
		node.left = leftRotate(node.left)
		return rightRotate(node)
	}

	// nodes lined-up to the right
	if bf < -1 && node.right.balanceFactor() <= 0 {
		return leftRotate(node)
	}

	// nodes linked-up to a 'greater than' shape
	if bf < -1 && node.right.balanceFactor() > 0 {
		node.right = rightRotate(node.right)
		return leftRotate(node)
	}

	return node
}

/*
	function to update a specific node
*/
func updateNode(n **BinaryNode, updated interface{}, id string) error {

	if len(strings.TrimSpace(id)) == 0 {
		return ErrEmptyNodeItemStatus
	}

	if (*n).id == id {
		(*n).item = updated
	}

	return nil
}
