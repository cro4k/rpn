#逆波兰式Golang实现

usage:
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