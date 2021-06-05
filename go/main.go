package main

import (
	"fmt"
	src "interview/src"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

func firstMissingPositive(nums []int) int {
	n := len(nums)
	for i := 0; i < n; i++ {
		for nums[i] > 0 && nums[i] <= n && nums[i] != nums[nums[i]-1] {
			nums[nums[i]-1], nums[i] = nums[i], nums[nums[i]-1]
		}
	}
	for i := 0; i < n; i++ {
		if nums[i] != i+1 {
			return i + 1
		}
	}
	return n + 1
}
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
func main() {
	// fmt.Println(firstMissingPositive([]int{3, 4, -1, 1}))
	// fmt.Println(src.HasCycle(src.CreateListNode()))
	// fmt.Println(src.MergeTowList(&src.ListNode{Val: 1, Next: &src.ListNode{Val: 4, Next: &src.ListNode{Val: 5, Next: nil}}}, &src.ListNode{Val: 1, Next: &src.ListNode{Val: 3, Next: &src.ListNode{Val: 4}}}))
	fmt.Println(src.FirstMissingPositive([]int{3, 4, -1, 1}))
	// m := map[string]iv.Student{"people": {"marsonshine"}}
	// // m1 := map[string]iv.Student{"people": {"marsonshine"}}
	// // m["people"].Name = "marsonshine"
	// fmt.Printf("%v", m)

	// ret := iv.GoExec("111", func(n string) string {
	// 	return n + "func1"
	// }, func(n string) string {
	// 	return n + "func2"
	// }, func(n string) string {
	// 	return n + "func3"
	// }, func(n string) string {
	// 	return n + "func4"
	// })
	// fmt.Println(ret)

	// // iv.Dead()

	// for _, i := range [5]int{1, 2, 3, 4, 5} {
	// 	fmt.Print(i)
	// }
	wg := sync.WaitGroup{}
	// ch := make(chan int)
	// defer close(ch)
	wg.Add(5)
	for _, i := range [5]int{1, 2, 3, 4, 5} {
		go func() {
			fmt.Println(i)
			wg.Done()
		}()
		// go func(i int) {
		// 	ch <- i
		// 	// wg.Done()
		// }(i)
	}
	// select {
	// case x := <-ch:
	// 	fmt.Printf("recvdq %d", x)
	// default:
	// 	fmt.Println("nothing to do")
	// }
	// for {
	// 	fmt.Println(<-ch)
	// }
	wg.Wait()
	// wg.Add(5)
	// fmt.Println("")
	// for _, i := range [5]int{1, 2, 3, 4, 5} {
	// 	go func(i int) {
	// 		fmt.Println(i)
	// 		wg.Done()
	// 	}(i)
	// }
	// wg.Wait()

	// // 组合
	// t := iv.Teacher{}
	// t.ShowA()

	// // defer 调用时参数作用域
	// a := 1
	// b := 2
	// defer calc("1", a, calc("10", a, b))
	// a = 0
	// defer calc("2", a, calc("20", a, b))
	// b = 1

	// s := make([]int, 5) // make 指定了长度，追加会从 len(length) 处填充数据
	// s = append(s, 1, 2, 3)
	// fmt.Println(s) // 所以会显示 0 0 0 0 0 1 2 3

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
	// iv.Pool()

	// channel
	// var ch chan int
	// go func() {
	// 	// ch = make(chan int, 1)
	// 	ch <- 1
	// }()
	// go func(ch chan int) {
	// 	time.Sleep(time.Second)
	// 	<-ch
	// }(ch)
	// c := time.Tick(1 * time.Second)
	// for range c {
	// 	fmt.Printf("#goroutines: %d\n", runtime.NumGoroutine())
	// }

	//
	// go f()
	// c <- 0
	// print(a)

	// // 定时器
	// go func() {
	// 	t := time.NewTicker(time.Second * 1)
	// 	for {
	// 		select {
	// 		case <-t.C:
	// 			go func() {
	// 				defer func() {
	// 					if err := recover(); err != nil {
	// 						fmt.Println(err)
	// 					}
	// 				}()
	// 				proc()
	// 			}()
	// 		}
	// 	}
	// }()
	// select {}

	// 内存逃逸分析
	a := foo("hello")
	b := a.s + " world"
	c := b + "!"
	fmt.Println(c)

	// 零拷贝 string - byte 转换
	aaa := "aaa"
	ssh := *(*reflect.StringHeader)(unsafe.Pointer(&aaa)) // 可以把字符串 aaa 转成底层结构的形式。
	bbb := *(*[]byte)(unsafe.Pointer(&ssh))               // 可以把 ssh 底层结构体转成 byte 的切片的指针。
	fmt.Printf("%v", bbb)                                 // 再通过 * 转为指针指向的实际内容。
}

type AA struct {
	s string
}

func foo(s string) *AA {
	a := new(AA)
	a.s = s
	return a
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

var c = make(chan int)
var a int

func f() {
	a = 1
	<-c
}
