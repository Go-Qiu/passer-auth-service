package stack

import (
	"errors"
)

type node struct {
	item interface{}
	next *node
}

type Stack struct {
	top  *node
	size int
}

var ErrEmptyStack error = errors.New("[stack]: stack is empty")

// wrapper function to instantiate a stack
func New() Stack {
	newStack := Stack{top: nil, size: 0}
	return newStack
}

// wrapper function to get the size of the instantiated stack
func (s *Stack) GetSize() int {
	return s.size
}

// wrapper function to set the top node pointer o f the stack
func (s *Stack) SetTop(n *node) {
	s.top = n
}

// wrapper function to set the size of the stack
func (s *Stack) SetSize(size int) {
	s.size = size
}

// wrapper function to get the top node pointer of the stack
func (s *Stack) GetTop() *node {
	return s.top
}

// wrapper function to push an item  (job) into the instantiated stack.
func (s *Stack) Push(item interface{}) error {
	newNode := &node{item: item, next: nil}

	if s.top == nil {

		// empty stack
		s.top = newNode
	} else {

		// stack is not empty
		// associate the node.next attribute of the new node to
		// the current Top node in the stack
		newNode.next = s.top

		// change the Top node in the stack to the new node.
		s.top = newNode
	}

	s.size++
	return nil
}

/*
	function to list all the node in the Stack
*/
func (s *Stack) ListAllNodesV2() ([]interface{}, error) {

	if s.top == nil {
		return nil, ErrEmptyStack
	}

	// ok. stack is not empty.
	currentNode := s.top
	users := []interface{}{}

	for currentNode != nil {
		users = append(users, currentNode.item)
		currentNode = currentNode.next
	}

	return users, nil
}

// func (s *Stack) ListSelectedNodesByStatus(status string) error {

// 	if s.top == nil {
// 		return ErrEmptyStack
// 	}

// 	// ok. stack is not empty
// 	currentNode := s.top
// 	count := 0
// 	for currentNode != nil {
// 		count++
// 		fmt.Printf("=================< Parcel Job Record # %04d - STARTS >=========\n", count)
// 		fmt.Printf("|\n")
// 		fmt.Printf("| User Email       : %s\n", currentNode.item.Id)
// 		fmt.Printf("|--------------------------------------------------------------\n")
// 		fmt.Printf("| First Name       : %s\n", currentNode.item.Name.First)
// 		fmt.Printf("| Last Name        : %s\n", currentNode.item.Name.Last)
// 		fmt.Printf("|\n")
// 		fmt.Printf("<<<<<<<<<<<<<<<<<< Parcel Job Record # %04d ENDS >>>>>>>>>>>>>>>>\n", count)
// 		fmt.Println()
// 		currentNode = currentNode.next
// 	}

// 	if count == 0 {
// 		// empty after the filtering
// 		return ErrEmptyStack
// 	}

// 	return nil
// }

/*
	function to pop a node from the top of the stack
*/
func (s *Stack) Pop() (interface{}, error) {

	if s.top == nil {
		return nil, ErrEmptyStack
	}

	// ok. the stack is not empty
	item := s.top.item
	if s.top.next == nil {
		// no more node
		s.top = nil
	} else {
		// there are still nodes beneath top
		s.top = s.top.next
	}
	s.size--
	return item, nil
}

/*
	function to peek the Top noode of the Stack
*/
func (s *Stack) Peek() (interface{}, error) {

	if s.top == nil {
		// stack is empty
		return nil, ErrEmptyStack
	}

	// ok. stack is not empty.
	item := s.top.item
	return item, nil
}
