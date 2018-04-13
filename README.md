# DataPool
DataPool is library to get data asynchronously and use it thread-safe.

### Example
Getting data from pool
```go
package main

import (
    "github.com/vbereschenko/dpool"
    "context"
    "time"
)

func main() {
    ctx, done := context.WithCancel(context.Background())
    pool := dpool.NewDataPool(ctx, func() (interface{}, error) {
        return "test", nil
    }, time.Second * 5)
    pool.Get()
    done()
}
```

Filling by pointer

```go
package main

import (
    "github.com/vbereschenko/dpool"
    "context"
    "time"
)

func main() {
    ctx, done := context.WithCancel(context.Background())
    pool := dpool.NewDataPool(ctx, func() (interface{}, error) {
        return "test", nil
    }, time.Second * 5)
    var result string 
    pool.FetchInto(&result)
    done()
}
```

### Benchmark

| calls    |ns/op |
|----------|------|
| 20000000 | 54.4 |
| 20000000 | 115  |