package interview

import "fmt"

type Student struct {
	Name string
}

func (p *Student) String() string {
	return fmt.Sprintf("print: %v", p) // 此时的 String() 实际上是实现了go内部接口 fmt/print.go 的 String 接口，而 Student 实现了这个接口就会直接调用这个接口。但同时内部又调用了 fmt 的方法，又会再次调用 fmt.String 接口，如此重复导致循环调用
}
