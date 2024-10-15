package directree

type Node struct {
	Name     string
	Children map[string]*Node
	Hosts    []string
}

func NewNode(name string) *Node {
	return &Node{
		Name:     name,
		Children: make(map[string]*Node),
		Hosts:    []string{},
	}
}

func Insert(node *Node, i int, path []string, host string) bool {

	if i == len(path) {
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

	if i == len(path)-1 {
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

	if i == len(path)-1 {

		if len(node.Children) == 0 {
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
