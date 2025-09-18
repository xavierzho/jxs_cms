package queue

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func test(value string) error {
	return fmt.Errorf("test error")
}
func TestQueue(t *testing.T) {
	q := NewQueueWorker(context.Background(), nil, nil, nil)
	q.AddQueueJob([]*QueueJob{
		{
			Name:  "test",
			Run:   test,
			Retry: false,
		},
	})
	err := q.work("test", "1111")

	fmt.Printf("%v\n", err)
	fmt.Printf("%v\n", errors.Unwrap(err))
	fmt.Printf("%v\n", err.Error())
}
