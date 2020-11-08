package server

type node struct {
	children []*node
	path     string
	indices  string
	handler  []HandlerFunc
}

type methodTree struct {
	method string
	root   *node // 这里的root表示 每个method下router的root
}

type methodTrees []methodTree // 每种method 在这里都会有一个元素

func (m methodTrees) getRoot(method string) *node {
	for i := 0; i < len(m); i++ {
		if m[i].method == method {
			return m[i].root
		}
	}
	return nil
}

func min(num1, num2 int) int {
	if num1 > num2 {
		return num2
	}
	return num1
}

func (n *node) addRouter(path string, handler ...HandlerFunc) {
	fullPath := path
	rootPathLen := len(n.path)
	if rootPathLen > 0 {
		// 查看 当前path 和 待插入path 是否有相同部分，
		// 若有相同部分 则分割 待插入path，
		minLen := min(rootPathLen, len(path))
		var i int
		for i < minLen && path[i] == n.path[i] {
			i++
		}

		if i < rootPathLen {

			/*
					这里表示 当前插入的path 和root节点的path有不同的地方，
					这时就该改变当前的root节点，将path[:i]作为新root节点的path
					将裁剪过后的root节点插入 新的root节点
					eg : /HelloWorld  这里是原来root 节点的path
						新插入节点为 /HelloMe
						这里树的结构就会变为
						/Hello
				         /   \
				      World  Me
			*/
			// 这里就是在 将原来root的path 修改为 World,然后将该World变成 /Hello的一个儿子节点
			child := &node{
				children: n.children,
				path:     n.path[i:],
				indices:  n.indices,
				handler:  n.handler,
			}
			// 改变root节点的path
			n.children = []*node{child}
			n.indices = string([]byte{n.path[i]})
			n.handler = nil
			n.path = path[:i]
		}

		if i == len(path) {
			n.handler = handler
		}

		if i < len(path) {
			// 这里表示待插入的path和 root的不一样 需要插入
			path = path[i:]
			c := path[0]

			// 若分割后的path 和当前 root的索引中有相同的。那么表示该path该成为已有节点的儿子
			for j := 0; j < len(n.indices); j++ {
				if c == n.indices[j] {
					n = n.children[j]
					n.addRouter(path, handler...)
					return
				}
			}

			// 若没有则表示 待插入节点该成为 root的儿子
			n.indices += string([]byte{c})
			child := &node{}
			n.children = append(n.children, child)
			n = child

			n.insertNode(path, handler...)
			return
		}

	} else {
		// 这里表示创建的是method的root节点
		// 直接插入
		n.insertNode(fullPath, handler...)
	}

}

func (n *node) getHandler(path string) []HandlerFunc {

	if n.path == path {
		return n.handler
	}

	if len(path) > len(n.path) {
		/*
			若path 是 当前node 的儿子 那么path 一定和node.path 有共同部分，
			共同部分为 path[0~ len(n.path)]
		*/
		if path[:len(n.path)] == n.path {

			path = path[len(n.path):]
			c := path[0]

			for i := 0; i < len(n.indices); i++ {
				if n.indices[i] == c {
					n = n.children[i]
					return n.getHandler(path)
				}
			}
		}
	}
	return nil
}

func (n *node) insertNode(path string, handler ...HandlerFunc) {
	n.path = path
	n.handler = handler
}
