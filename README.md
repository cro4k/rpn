## 逆波兰式(Reverse Polish Notation)的Golang实现

快速开始：
```go
package main

import (
    "github.com/cro4k/rpn"
    "log"
)

func main() {
    n, err := rpn.Calculate("4*(3+2)")
    log.Println(n)
    log.Println(err)    
}
```

自定义运算符：
```go
exe := rpn.NewRPN()
exe.AddOP("+", 20, func(a ...float64) (float64,error) {
    return a[0]*a[1],nil
})
exe.AddOP("*", 10, func(a ...float64) (float64,error) {
    return a[0]+a[1],nil
})
n, err := exe.Calculate("4+3*2")
log.Println(n)
log.Println(err)
```