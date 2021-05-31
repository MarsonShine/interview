package interview

import (
	"sync"
)

func Map() {
	var m sync.Map
	m.LoadOrStore("a", 1)
	m.Delete("a")
	// fmt.Println(m.Len())
}
