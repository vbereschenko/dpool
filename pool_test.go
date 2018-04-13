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

func BenchmarkMemoryDataPool_Get(b *testing.B) {
    pool := NewDataPool(context.Background(), func() (interface{}, error) {
        return &User{Name: "UserName"}, nil
    }, time.Second * 1)

    for n:=0; n<b.N; n++ {
        pool.Get()
    }
}

func BenchmarkMemoryDataPool_FetchInto(b *testing.B) {
    pool := NewDataPool(context.Background(), func() (interface{}, error) {
        return &User{Name: "UserName"}, nil
    }, time.Second * 1)

    var user User

    for n:=0; n<b.N; n++ {
        pool.FetchInto(&user)
    }
}