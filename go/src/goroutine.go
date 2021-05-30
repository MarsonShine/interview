package interview

import (
	"fmt"
	"runtime"
	"time"
)

func gochannel() {
	ch := make(chan int, 1000)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()
	go func() {
		for {
			a, ok := <-ch
			if !ok {
				fmt.Println("close")
				return
			}
			fmt.Println("a: ", a)
		}
	}()
	close(ch) // goroutine 开启协程启动的时机不固定，所以写入 channel 的时机可能在 close 后，就会报错。
	fmt.Println("ok")
	time.Sleep(time.Second * 100)
}

type query func(string) string

func GoExec(name string, vs ...query) string {
	ch := make(chan string)
	fn := func(i int) {
		ch <- vs[i](name)
	}
	for i, _ := range vs {
		go fn(i)
	}
	return <-ch
}

func Dead() {
	var i byte
	go func() {
		for i = 0; i <= 255; i++ { // 重点：因为 byte 最大 256，所以这个判断永远成立
		}
	}()
	fmt.Println("Dropping mic")
	// Yield execution to force executing other goroutines runtime.Gosched()
	runtime.Gosched() // yield 让出自己的执行权
	runtime.GC()      // 无法回收，因为 goroutine 无法让出自己的执行权，所以 go 内部会将这个长时间运行的 g 标记成可抢占状态（preemt），而 gc 是需要将所有 preemt 的程序停止后才能进行的。所以程序运行到这里会卡死
	fmt.Println("Done")
}
