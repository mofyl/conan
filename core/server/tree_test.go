package server

import (
	"testing"
	"conan/log"
)

func Hello(c *Context) {
	log.Info("HelloWorld")
}

func Me(c *Context) {

	log.Info("HelloMe")
}

func TestTree(t *testing.T) {
	tree := make(methodTrees, 9)

	root := &node{}

	tree = append(tree, methodTree{
		method: "GET",
		root:   root,
	})

	root.addRouter("/HelloWorld", Hello)
	root.addRouter("/HelloMe", Me)

	h := tree.getRoot("GET").getHandler("/HelloMe")
	if h == nil {
		log.Error("Can not Find Router")
		return
	}
}
