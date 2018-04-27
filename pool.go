// data pool example
// ---------
//    ctx, done := context.WithCancel(context.Background())
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
    "reflect"
    "errors"
)

type DataProvider func() (interface{}, error)

type DataPool interface {
    Get() (interface{}, error)
    FetchInto(interface{}) error
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

// provides actual values of DataProvider
// thread-safe function
func (dataPool *memoryDataPool) Get() (interface{}, error) {
    dataPool.RLock()
    defer dataPool.RUnlock()

    return dataPool.result, dataPool.err
}

// fills object by pointer provided
// this method checks if pointer is provided and if types of stored value
// equals to value that is requested to fill
func (dataPool *memoryDataPool) FetchInto(result interface{}) error {
    if dataPool.err != nil {
        return dataPool.err
    }

    dataPool.RLock()
    defer dataPool.RUnlock()

    if reflect.TypeOf(result).Kind() != reflect.Ptr {
        return errors.New("argument should be pointer")
    }

    if reflect.TypeOf(result).Elem().Name() != reflect.TypeOf(dataPool.result).Name() {
        return errors.New("types don't match")
    }

    reflect.ValueOf(result).Elem().Set(reflect.ValueOf(dataPool.result))

    return dataPool.err
}

type memoryDataPool struct {
    sync.RWMutex

    provider DataProvider
    repeat   *time.Ticker

    result interface{}
    err    error
}

func (dataPool *memoryDataPool) run(ctx context.Context) {
    var resultBuffer interface{}
    var errBuffer error
    for {
        select {
        case <- ctx.Done():
            return

        case <- dataPool.repeat.C:
            // fetching before locking
            resultBuffer, errBuffer = dataPool.provider()

            dataPool.Lock()
            dataPool.result, dataPool.err = resultBuffer, errBuffer
            dataPool.Unlock()
        }
    }
}
