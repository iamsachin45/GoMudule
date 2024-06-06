package hello

import (
	"sync"
)

// Node represents a node in the B-tree.
type Node struct {
	keys     []string
	children []*Node
}

// Database represents a simple key-value store with B-tree indexing.
type Database struct {
	root *Node
	mu   sync.RWMutex // Mutex for thread safety
}

// NewNode creates a new instance of a B-tree node.
func NewNode() *Node {
	return &Node{}
}

// NewDatabase creates a new instance of the Database.
func NewDatabase() *Database {
	return &Database{
		root: NewNode(),
	}
}

// Insert inserts a key-value pair into the B-tree.
func (db *Database) Insert(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	insert(db.root, key, value)
}

func insert(node *Node, key, value string) {
	// Find the appropriate position to insert the key.
	i := len(node.keys) - 1
	for i >= 0 && node.keys[i] > key {
		i--
	}
	i++

	// If the node is a leaf, insert the key directly.
	if len(node.children) == 0 {
		node.keys = append(node.keys, "")
		copy(node.keys[i+1:], node.keys[i:])
		copy(node.children[i+1:], node.children[i:])
		node.keys[i] = key
		return
	}

	// If the node is not a leaf, recursively insert into the appropriate child.
	if i < len(node.children) {
		insert(node.children[i], key, value)
	}
}

// Search searches for a key in the B-tree.
func (db *Database) Search(key string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return search(db.root, key)
}

func search(node *Node, key string) (string, bool) {
	i := 0
	for i < len(node.keys) && node.keys[i] < key {
		i++
	}

	if i < len(node.keys) && node.keys[i] == key {
		return node.keys[i], true
	}

	if len(node.children) > 0 {
		return search(node.children[i], key)
	}

	return "", false
}
