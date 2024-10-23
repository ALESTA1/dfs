package directree

import "sync"

type Node struct {
	Name     string
	Children map[string]*Node
	index    int
	Hosts    []string
	Is_Dir   bool
	RwMutex  sync.RWMutex
}

func NewNode(name string) *Node {
	return &Node{
		Name:     name,
		Children: make(map[string]*Node),
		Hosts:    []string{},
		index:    0,
		Is_Dir:   true,
	}
}

func Insert(node *Node, i int, path []string, host string) bool {

	if i == len(path) {
		node.Is_Dir = false
		return true
	}
	currentNode := path[i]

	nextNode, exists := node.Children[currentNode]
	f := true
	if !exists {
		temp := NewNode(currentNode)
		node.Children[currentNode] = temp
		f = f && Insert(temp, i+1, path, host)
	} else {

		if i == len(path)-1 {
			return false
		}
		f = f && Insert(nextNode, i+1, path, host)
	}
	if f {
		node.Hosts = append(node.Hosts, host)
	}
	return f
}

func IsValidPath(node *Node, i int, path []string) bool {

	if i == len(path) {
		return true
	}

	f := false
	nextNode, exists := node.Children[path[i]]
	if exists {
		f = f || IsValidPath(nextNode, i+1, path)
	}

	return f
}

func IsDir(node *Node, i int, path []string) int {

	if i == len(path) {

		if node.Is_Dir {
			return 0
		} else {
			return 1
		}
	}

	f := 2

	nextNode, exists := node.Children[path[i]]
	if exists {

		return IsDir(nextNode, i+1, path)
	}

	return f
}

func GetHost(node *Node, i int, path []string) string {

	if i == len(path) {

		host := node.Hosts[node.index]
		node.index += 1
		node.index %= len(node.Hosts)
		return host
	}

	nextNode, exists := node.Children[path[i]]
	if exists {

		return GetHost(nextNode, i+1, path)
	}

	return ""
}

func CreateDir(node *Node, i int, path []string, host string) {

	if i == len(path) {
		return
	}
	currentNode := path[i]
	nextNode, exists := node.Children[currentNode]
	if !exists {
		temp := NewNode(currentNode)
		node.Children[currentNode] = temp
		Insert(temp, i+1, path, host)
	} else {
		Insert(nextNode, i+1, path, host)
	}
	node.Hosts = append(node.Hosts, host)

}

func Delete(node *Node, i int, path []string) []string {

	if i == len(path) {
		return node.Hosts
	}
	currentNode := path[i]
	nextNode := node.Children[currentNode]
	if i == len(path)-1 {
		delete(node.Children, currentNode)
	}
	return Delete(nextNode, i+1, path)
}

func Lock(node *Node, i int, path []string, exclusive bool) {

	if i == len(path)-1 {
		return
	}
	if exclusive {
		node.RwMutex.Lock()
	} else {
		node.RwMutex.RLock()
	}
	currentNode := path[i]
	nextNode := node.Children[currentNode]
	Lock(nextNode, i+1, path, exclusive)

}

func Unlock(node *Node, i int, path []string, exclusive bool) {

	if i == len(path)-1 {
		return
	}
	if exclusive {
		node.RwMutex.Unlock()
	} else {
		node.RwMutex.RUnlock()
	}
	currentNode := path[i]
	nextNode := node.Children[currentNode]
	Lock(nextNode, i+1, path, exclusive)

}

func FindNode(node *Node, i int, path []string) *Node {
	if i == len(path) {
		if node.Is_Dir {
			return node
		}

		return nil
	}
	currentNode := path[i]
	nextNode, exists := node.Children[currentNode]
	if !exists {
		return nil
	} else {
		return FindNode(nextNode, i+1, path)
	}

}

func List(node *Node, temp string, paths *[]string) {

	if !node.Is_Dir {
		*paths = append(*paths, temp)
		return
	}

	for key, val := range node.Children {
		next := temp
		next = next + "/" + key
		List(val, next, paths)
	}
}
