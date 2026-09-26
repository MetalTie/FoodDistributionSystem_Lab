// The GateKeeper utilizes a circular queue for storage and retrieval. The package has been
// designed to be generic to allow use in future systems with respect to the data type and queue size.
//
// The software suite consists of main.go,
// food, stats, gatekeeper,
// producer, circularque, and sales packages.
//
// This represents the software to manage an "embedded" planetary system
// food receiving and distribution system.

package circularque

import (
    "errors"
    //"fmt"
)

type CircularQue[T any] struct {
    capacity    int
    box         []T
    front, rear int
    mesnum      int
}

func NewCircularQue[T any](capacity int) *CircularQue[T] {
    return &CircularQue[T]{
        capacity: capacity,
        box:      make([]T, capacity),
        front:    0,
        rear:     0,
        mesnum:   0,
    }
}

func (q *CircularQue[T]) AcceptMessage(msg T) error {
    if q.mesnum >= q.capacity {
        return errors.New("ERROR - Message rejected - queue is full!")
    }
    q.rear = (q.rear + 1) % q.capacity
    q.box[q.rear] = msg
    q.mesnum++
    return nil
}

func (q *CircularQue[T]) RetrieveMessage() (T, error) {
    var zero T
    if q.mesnum <= 0 {
        return zero, errors.New("ERROR - No message in the queue to retrieve!")
    }
    q.front = (q.front + 1) % q.capacity
    msg := q.box[q.front]
    q.mesnum--
    return msg, nil
}

func (q *CircularQue[T]) IsEmpty() bool {
    return q.mesnum == 0
}

func (q *CircularQue[T]) IsFull() bool {
    return q.mesnum >= q.capacity
}