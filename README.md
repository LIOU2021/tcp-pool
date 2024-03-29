# introduction
- 底层采用带有头尾指标的ring buffer循环利用记忆体
# ref
- [golang 环形队列](https://juejin.cn/s/golang%20%E7%8E%AF%E5%BD%A2%E9%98%9F%E5%88%97)
- [go语言数据结构 环形队列](https://cloud.tencent.com/developer/article/1165585)
# tcp connect pool

# example
```go
package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/LIOU2021/tcp-pool/pool"
)

func main() {
	address := "127.0.0.1:8000"

	pool := &pool.ConnPool{
		Dial: func() (net.Conn, error) {
			return net.Dial("tcp", address)
		},
		MaxIdle:  10,
		MinIdle:  2,
		IdleTime: 10 * time.Second,
	}

	pool.CreatePool()

	conn, err := pool.Get()
	if err != nil {
		log.Fatal(err)
	}

	conn.Write([]byte(fmt.Sprintf("echo %s", "hello world")))
	buf := make([]byte, 4096)

	cnt, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Fail to read from server: %s\n", err)
	}
	fmt.Println(string(buf[0:cnt]))

	pool.Put(conn)

	pool.Release()
}

```