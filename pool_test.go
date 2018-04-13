package dpool

import (
    "testing"
    "context"
    "time"
)

type User struct {
    Name string
}

func TestMemoryDataPool_FetchInto(t *testing.T) {
    pool := NewDataPool(context.Background(), func() (interface{}, error) {
        return &User{Name: "UserName"}, nil
    }, time.Second * 1)

    var result *User

    if err := pool.FetchInto(&result); err != nil {
        t.Fatal(err)
    }
    if result.Name != "UserName" {
        t.Fail()
    }
}
