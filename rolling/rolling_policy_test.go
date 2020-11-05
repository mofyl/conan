package rolling

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func getPolicy() *Policy {
	w := NewWindow(WindowOpt{Size: 10})
	return NewPolicy(w, PolicyOpts{BucketDuration: 300 * time.Millisecond})
}

func TestRollingPolicy(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	p := getPolicy()
	//time.Sleep(400 * time.Millisecond)
	//p.Add(1)
	time.Sleep(3300 * time.Millisecond)
	p.Add(1)

	for _, v := range p.window.window {
		fmt.Println(v)
	}
}
