package group_test

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/tadvi/group"
)

func TestZero(t *testing.T) {
    var g group.Tasks
    res := make(chan []error)

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() { res <- g.Run(ctx) }()
    select {
    case err := <-res:
        if err != nil {
            t.Errorf("%v", err)
        }
    case <-time.After(100 * time.Millisecond):
        t.Error("timeout")
    }
}

func TestOne(t *testing.T) {
    myError := errors.New("foobar")
    var g group.Tasks

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    g.Add(func(ctx context.Context) error { return myError })

    res := make(chan []error)
    go func() { res <- g.Run(ctx) }()
    select {
    case errs := <-res:
        if want, got := myError, errs; len(got) != 1 || want != got[0] {
            t.Errorf("got %v, want %v", got, want)
        }
    case <-time.After(100 * time.Millisecond):
        t.Error("timeout")
    }
}

func TestMany(t *testing.T) {
    myError := errors.New("foobar")
    var g group.Tasks

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    g.Add(func(ctx context.Context) error { return myError })
    g.Add(func(ctx context.Context) error { return myError })
    g.Add(func(ctx context.Context) error { return nil })

    res := make(chan []error)
    go func() { res <- g.Run(ctx) }()
    select {
    case errs := <-res:
        if want, got := myError, errs; len(got) != 2 || want != got[0] || want != got[1] {
            t.Errorf("got %v, want %v", got, want)
        }
    case <-time.After(100 * time.Millisecond):
        t.Errorf("timeout")
    }
}
