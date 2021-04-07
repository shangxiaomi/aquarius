package aquarius

import "testing"

func TestTrie(t *testing.T) {
	n := newNode()
	addPath(n, "/hello/world")
	addPath(n, "/hello/shp")
	addPath(n, "/nishishui")
	addPath(n, "/nishishui/cxh")
	addPath(n, "/hello/:name")
	travel2(n)
}

func addPath(n *node, path string) {
	parts := parsePattern(path)
	n.insert(path, parts, 0)
}
