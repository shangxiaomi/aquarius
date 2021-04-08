package aquarius

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrie(t *testing.T) {
	n := new(node)
	addPath(n, "/hello/world")
	addPath(n, "/hello/shp")
	addPath(n, "/nishishui")
	addPath(n, "/nishishui/cxh")
	addPath(n, "/hello/:name")
	travel2(n)
}

func TestTriePanic(t *testing.T) {
	/*
		1.在插入wildcard节点时，父节点的children数组非空且wildChild被设置为false。例如：GET /user/getAll和GET /user/:id/getAddr，或者GET /user/*aaa和GET /user/:id。
		2.在插入wildcard节点时，父节点的children数组非空且wildChild被设置为true，但该父节点的wildcard子节点要插入的wildcard名字不一样。例如：GET /user/:id/info和GET /user/:name/info。
		3.在插入catchAll节点时，父节点的children非空。例如：GET /src/abc和GET /src/*filename，或者GET /src/:id和GET /src/*filename。
		4.在插入static节点时，父节点的wildChild字段被设置为true。
		5.在插入static节点时，父节点的children非空，且子节点nType为catchAll。
		6.某个路由重复插入「已解决」
	*/
	cases := []*trieTestCase{
		{
			desc: "插入Param节点时，父节点的children数组非空，且没有模糊匹配",
			paths: []string{
				"/user/getAll",
				"/user/:id/getAddr",
			},
			isPanic: true,
		},
		{
			desc: "插入Param节点时，父节点的children数组非空，有模糊匹配，但是param的名称不同",
			paths: []string{
				"/user/:name/getAll",
				"/user/:id/getAddr",
			},
			isPanic: true,
		},
		{
			desc: "插入Param节点时，父节点的children数组非空，有模糊匹配，但是param的名称相同",
			paths: []string{
				"/user/:id/getAll",
				"/user/:id/getAddr",
			},
			isPanic: false,
		},
		{
			desc: "插入cacthAll节点时，父节点有孩子节点",
			paths: []string{
				"/user/shp/",
				"/user/*shp",
			},
			isPanic: true,
		},
		{
			desc: "插入static节点时，父节点没有有孩子节点",
			paths: []string{
				"/user/",
				"/user/*shp",
			},
			isPanic: false,
		},
		{
			desc: "插入静态节点时,父亲节点不能有param的孩子节点",
			paths: []string{
				"/user/:shp/123",
				"/user/233/233",
			},
			isPanic: true,
		},
		{
			desc: "插入静态节点时,父亲节点不能有cacthAll的孩子节点",
			paths: []string{
				"/user/shp/*123",
				"/user/shp/233",
			},
			isPanic: true,
		},
		{
			desc: "插入静态节点时,父亲节不能为cacthAll节点",
			paths: []string{
				"/user/*trx/",
				"/user/shp/233",
			},
			isPanic: true,
		},
		{
			desc: "插入静态节点时,父亲节不能为cacthAll节点2",
			paths: []string{
				"/*trx/",
				"shp/233",
			},
			isPanic: true,
		},
		{
			desc: "不能重复插入路由",
			paths: []string{
				"/user/shp",
				"/user/shp",
			},
			isPanic: true,
		},
	}
	runTrieTestCase(t, cases)
}

func addPath(n *node, path string) {
	parts := parsePattern(path)
	n.insert(path, parts, 0)
}

func runTrieTestCase(t *testing.T, testCases []*trieTestCase) {
	for _, kase := range testCases {
		t.Run(kase.desc, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					assert.Equal(t, true, kase.isPanic)
					t.Log("[panic]", err)
				}
			}()
			n := new(node)
			for _, p := range kase.paths {
				addPath(n, p)
			}
			assert.Equal(t, false, kase.isPanic)
		})
	}
}

type trieTestCase struct {
	desc    string
	paths   []string
	isPanic bool
}
