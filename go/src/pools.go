package interview

import (
	"fmt"
	"sync"
)

var pool = &sync.Pool{
	New: func() interface{} {
		fmt.Println("init from pool")
		return 0
	},
}

func Pool() {
	init := pool.Get()
	fmt.Println(init)

	pool.Put(1)
	pool.Put("marsonshine")
	num := pool.Get()
	fmt.Println(num)

	num = pool.Get()
	fmt.Println(num)

	for i := 0; i < 5; i++ {
		go func() {
			pool.Get()
		}()
	}
}
