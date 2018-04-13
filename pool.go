// data pool example
// ---------
//    ctx, done := context.WithCancel(context.TODO())
//    pool := dpool.NewDataPool(ctx, func() (interface{}, error) {
//        return "test", nil
//    }, time.Second * 5)
//
//    pool.Get()
//    done()
package dpool

import (
    "context"
    "time"
    "sync"
    "fmt"
)

type DataProvider func() (interface{}, error)

type DataPool interface {
    Get() (interface{}, error)
}

func NewDataPool(ctx context.Context, provider DataProvider, repeat time.Duration) *memoryDataPool {
    pool := memoryDataPool{
        provider: provider,
        repeat:   time.NewTicker(repeat),
    }
    pool.result, pool.err = provider()

    go pool.run(ctx)

    return &pool
}

func (dataPool *memoryDataPool) Get() (interface{}, error) {
    dataPool.RLock()
    defer dataPool.RUnlock()

    return dataPool.result, dataPool.err
}

type memoryDataPool struct {
    sync.RWMutex

    provider DataProvider
    repeat   *time.Ticker

    result interface{}
    err    error
}

func (dataPool *memoryDataPool) run(ctx context.Context) {
    for {
        select {
        case <- ctx.Done():
            fmt.Println("Background task finished")
            return

        case <- dataPool.repeat.C:
            dataPool.Lock()
            dataPool.result, dataPool.err = dataPool.provider()
            dataPool.Unlock()
        }
    }
}
