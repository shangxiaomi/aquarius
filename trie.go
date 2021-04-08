package aquarius

import (
	"fmt"
	"strings"
)
// node methods is not concurrent safe
type node struct {
	pattern  string  // 存放此节点对应的模式串
	part     string  // 此节点的匹配部分
	children []*node // 子节点
	isWild   bool    // 是否为模糊匹配

}


func (root *node) String() string {
	return fmt.Sprintf("(part[%s], pattern[%s], ChildLen[%d] ,isWild[%t])", root.part, root.pattern, len(root.children), root.isWild)
}

/*
insert add route path , if parts include part of contains "*", the part is the last element
*/
func (root *node) insert(pattern string, parts []string, height int) {
	// 达到模式串指定的深度
	if len(parts) == height {
		if root.pattern != "" {
			// 路由添加短路逻辑1：如果某个路由重复添加，直接panic
			panic(fmt.Errorf("[trie node] insert error , pattern of node is not empty, pattern of node :%s, insert pattern :%s", root.pattern, pattern))
		}
		// 赋值模式串
		root.pattern = pattern
		return
	}
	part := parts[height]
	childNode := matchChildNode(root, part)

	checkPattern(pattern, part, root, height)

	if childNode == nil {
		childNode = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		root.children = append(root.children, childNode)
	}

	childNode.insert(pattern, parts, height+1)
}

// checkPattern 检查插入的节点是否合法
func checkPattern(pattern string, part string, root *node, height int) {
	if part[0] == ':' { // 要插入的是个部分通配符
		f := false
		for _, child := range root.children {
			if child.part[0] == '*' {
				panic(fmt.Errorf("插入param节点时，父节点不能有matchAll的孩子节点，父节点[%s]，孩子节点[%s]，要插入的路径[%s]，冲突的部分[%s]，第[%d]层", root, child, pattern, part, height))
			}
			if child.part[0] == ':' {
				f = true
				if child.part != part {
					panic(fmt.Errorf("插入param节点时，父节点不能有通配符但是参数名称不同的的孩子节点，父节点[%s]，孩子节点[%s]，要插入的路径[%s]，冲突的部分[%s]，第[%d]层", root, child, pattern, part, height))
				}
			}
		}
		if !f && len(root.children) > 0 {
			panic(fmt.Errorf("插入param节点时，父节点不能有非通配符的孩子节点，父节点[%s]，要插入的路径[%s]，冲突的部分[%s]，第[%d]层", root, pattern, part, height))
		}
	} else if part[0] == '*' { // 要插入的是个*通配符
		if len(root.children) > 0 {
			panic(fmt.Errorf("在插入catchAll节点时，父节点的children必须为非空，父节点[%s],要插入的路径[%s]，冲突的部分[%s]", root, pattern, part))
		}
	} else { // 要插入的是个静态路径
		for _, child := range root.children {
			if len(child.part) > 0 && (child.part[0] == '*' || child.part[0] == ':') {
				panic(fmt.Errorf("插入static路径时，父节点不能有通配符的孩子节点，父节点[%s]，孩子节点[%s]，要插入的路径[%s]，冲突的部分[%s]，第[%d]层", root, child, pattern, part, height))
			}
		}
	}
}

func isWild(part string) bool {
	return part[0] == ':' || part[0] == '*'
}

// search get the matched trie node
func (root *node) search(parts []string, height int) *node {
	// 达到指定深度，或者trie数的part部分为"*"匹配方式
	if len(parts) == height || strings.HasPrefix(root.part, "*") {
		if root.pattern != "" {
			return root
		}
		return nil
	}
	part := parts[height]
	children := matchChildren(root, part)
	for _, child := range children {
		res := child.search(parts, height+1)
		if res != nil {
			return res
		}
	}
	return nil
}

// 只匹配到第一个符合条件的孩子节点
func matchChildNode(root *node, part string) *node {
	for _, child := range root.children {
		if child == nil {
			// 所有的子节点不应给为空
			panic("child should not be nil")
		}
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 找到所有匹配的节点
func matchChildren(root *node, part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range root.children {
		if child == nil {
			// 所有的子节点不应给为空
			panic("child should not be nil")
		}
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func travel(root *node) {
	que := make([]*node, 0)
	tail := 0
	que = append(que, root)
	idx := 1
	for tail < len(que) {
		n := que[tail]
		fmt.Print(n, " ")
		que = append(que, n.children...)
		tail++
		if tail == idx {
			fmt.Println()
			idx += len(que) - tail
		}
	}
}
