# String 相关知识

此时的 String() 实际上是实现了go内部接口 fmt/print.go 的 String 接口，而 Student 实现了这个接口就会直接调用这个接口。但同时内部又调用了 fmt 的方法，又会再次调用 fmt.String 接口，如此重复导致循环调用

```go
type Student struct {
	Name string
}

func (p *Student) String() string {
	return fmt.Sprintf("print: %v", p) // 结构体默认实现 String() 方法，所以当显示实现 String 时会发生循环递归调用
}
```

此时编译器会提示：`Sprintf format %v with arg p causes recursive String method call`

# Panic / Recover 相关知识

- `panic` 能中断程序控制流，并**递归调用程序中 `defer` 代码段**
- `recover` 可以中止 `panic` 造成的程序崩溃。**它是一个只能在 `defer` 中发挥作用的函数，在其他作用域中调用不会发挥作用；**

# for 循环

```go
wg := sync.WaitGroup{}
// 1
wg.Add(5)
for _, i := range [5]int{1, 2, 3, 4, 5} {
  go func() {
    fmt.Println(i)
    wg.Done()
  }()
}
wg.Wait()
// 2
wg.Add(5)
fmt.Println("")
for _, i := range [5]int{1, 2, 3, 4, 5} {
  go func(i int) {
    fmt.Println(i)
    wg.Done()
  }(i)
}
wg.Wait()
```

第 2 段代码的输出这里有个误区，就是以为开启了 5 个协程运行，肯定顺序是随机的。其实通过调用我们发现**第一个运行输出永远都是 5，也就是最后一个循环体**。为什么会这样呢？

其实输出顺序是受调度器决定到底执行哪一个的。其实我们通过观察源码以及《Go语言设计与实现》中的[调度器](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-goroutine/)介绍中我们可知。我们在开启 G 的时候 go 是会将这些 G push 到一个队列中的，由 P 来控制 M 与 G 的运行。当创建一 个 G 时，会优先放入到下一个调度的 runnext 字段上作为下一次优先调度的 G。因此， 最先输出的是最后创建的G，也就是 5。

```go
func newproc(siz int32, fn *funcval) {
	argp := add(unsafe.Pointer(&fn), sys.PtrSize)
	gp := getg()
	pc := getcallerpc()
	systemstack(func() {
    // 这里就是初始化 g
		newg := newproc1(fn, argp, siz, gp, pc)

		_p_ := getg().m.p.ptr()
    // 调度 g
		runqput(_p_, newg, true)

		if mainStarted {
			wakep()
		}
	})
}

func runqput(_p_ *p, gp *g, next bool) {
	if randomizeScheduler && next && fastrand()%2 == 0 {
		next = false
	}

	if next {
	retryNext:
		oldnext := _p_.runnext
    // 当 next 是 true 时总会将新进来的 G 放入下一次调度字段中
		if !_p_.runnext.cas(oldnext, guintptr(unsafe.Pointer(gp))) {
			goto retryNext
		}
		if oldnext == 0 {
			return
		}
		// Kick the old runnext out to the regular run queue.
		gp = oldnext.ptr()
	}

retry:
	h := atomic.LoadAcq(&_p_.runqhead) // load-acquire, synchronize with consumers
	t := _p_.runqtail
	if t-h < uint32(len(_p_.runq)) {
		_p_.runq[t%uint32(len(_p_.runq))].set(gp)
		atomic.StoreRel(&_p_.runqtail, t+1) // store-release, makes the item available for consumption
		return
	}
	if runqputslow(_p_, gp, h, t) {
		return
	}
	// the queue is not full, now the put above must succeed
	goto retry
}
```

# 关于组合继承

# 锁

```go
type UserAges struct {
	ages map[string]int
	mu   sync.Mutex
	rw   sync.RWMutex
}

func (ua *UserAges) Add(name string, age int) {
	ua.mu.Lock()
	defer ua.mu.Unlock()
	ua.ages[name] = age
}
func (ua *UserAges) Get(name string) int {
	if age, ok := ua.ages[name]; ok { // 并发时会爆异常，fatal error: concurrent map read and map write，应该该用读写锁
		return age
	}
	return -1
}
func (ua *UserAges) Add2(name string, age int) {
	ua.rw.Lock()
	defer ua.rw.Unlock()
	ua.ages[name] = age
}

func (ua *UserAges) Get2(name string) int {
	ua.rw.RLock()
	defer ua.rw.RUnlock()
	if age, ok := ua.ages[name]; ok {
		return age
	}
	return -1
}
```

