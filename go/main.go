package main

import (
	"fmt"
	iv "interview/src"
	"sync"
	"time"
)

func main() {
	m := map[string]iv.Student{"people": {"marsonshine"}}
	// m1 := map[string]iv.Student{"people": {"marsonshine"}}
	// m["people"].Name = "marsonshine"
	fmt.Printf("%v", m)

	ret := iv.GoExec("111", func(n string) string {
		return n + "func1"
	}, func(n string) string {
		return n + "func2"
	}, func(n string) string {
		return n + "func3"
	}, func(n string) string {
		return n + "func4"
	})
	fmt.Println(ret)

	// iv.Dead()

	for _, i := range [5]int{1, 2, 3, 4, 5} {
		fmt.Print(i)
	}
	wg := sync.WaitGroup{}
	wg.Add(5)
	for _, i := range [5]int{1, 2, 3, 4, 5} {
		go func() {
			fmt.Println(i)
			wg.Done()
		}()
	}
	wg.Wait()
	wg.Add(5)
	fmt.Println("")
	for _, i := range [5]int{1, 2, 3, 4, 5} {
		go func(i int) {
			fmt.Println(i)
			wg.Done()
		}(i)
	}
	wg.Wait()

	// 组合
	t := iv.Teacher{}
	t.ShowA()

	// defer 调用时参数作用域
	a := 1
	b := 2
	defer calc("1", a, calc("10", a, b))
	a = 0
	defer calc("2", a, calc("20", a, b))
	b = 1

	s := make([]int, 5) // make 指定了长度，追加会从 len(length) 处填充数据
	s = append(s, 1, 2, 3)
	fmt.Println(s) // 所以会显示 0 0 0 0 0 1 2 3

	// RWMutx
	// go A()
	// time.Sleep(2 * time.Second)
	// mu.Lock()
	// defer mu.Unlock()
	// count++
	// fmt.Println(count)

	//
	// var mu MyMutex
	// mu.Lock()
	// var mu2 = mu // 死锁，因为复制的已经把 lock 的状态也复制过来了
	// mu.count++
	// mu.Unlock()
	// mu2.Lock() // 复制的时候已经把 mu.Lock 状态复制过来了。
	// mu2.count++
	// mu2.Unlock()
	// fmt.Println(mu.count, mu2.count)

	// pool
	iv.Pool()

	// 定时器
	go func() {
		t := time.NewTicker(time.Second * 1)
		for {
			select {
			case <-t.C:
				go func() {
					defer func() {
						if err := recover(); err != nil {
							fmt.Println(err)
						}
					}()
					proc()
				}()
			}
		}
	}()
	select {}
}
func proc() {
	panic("ok")
}
func calc(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}

var mu sync.RWMutex
var count int

func A() {
	mu.RLock()
	defer mu.RUnlock()
	B()
}
func B() {
	time.Sleep(5 * time.Second)
	C()
}
func C() {
	mu.RLock()
	defer mu.RUnlock()
}

type MyMutex struct {
	count int
	sync.Mutex
}
