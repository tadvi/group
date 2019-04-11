// Package group implements a runner with deterministic teardown. It is
// somewhat similar to package errgroup but collects all the errors.
package group

import "context"

// Tasks group collects runner functions and runs them concurrently.
type Tasks struct {
    runners []task
}

// Add function to run as part of a group of tasks.
func (ts *Tasks) Add(execute func(ctx context.Context) error) {
    ts.runners = append(ts.runners, task{execute})
}

// Run all the tasks.
func (ts *Tasks) Run(ctx context.Context) []error {
    if len(ts.runners) == 0 {
        return nil
    }

    // Run each task.
    errors := make(chan error, len(ts.runners))
    var errs []error

    for _, a := range ts.runners {
        go func(ctx context.Context, a task) {
            errors <- a.execute(ctx)
        }(ctx, a)
    }

    // Wait for all runners to stop.
    for i := 0; i < cap(errors); i++ {
        err := <-errors
        if err != nil {
            errs = append(errs, err)
        }
    }

    // Return collected errors.
    return errs
}

type task struct {
    execute func(ctx context.Context) error
}
