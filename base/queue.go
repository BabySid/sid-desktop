package base

import (
	"container/list"
	"errors"
	"sync"
)

const (
	discardFIFO = iota
)

type Queue struct {
	container       *list.List
	capacity        int
	discardStrategy int
	mutex           *sync.Mutex
}

func NewQueue() *Queue {
	list.New()
	return &Queue{
		container:       list.New(),
		capacity:        -1,
		discardStrategy: discardFIFO,
		mutex:           &sync.Mutex{},
	}
}

func (q *Queue) SetCapacity(cap int) {
	q.capacity = cap
}

func (q *Queue) PushBack(item interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if q.capacity > 0 && q.container.Len() >= q.capacity {
		if q.discardStrategy == discardFIFO {
			if it := q.container.Front(); it != nil {
				q.container.Remove(it)
			}
		} else {
			panic(errors.New("cannot run here"))
		}
	}

	q.container.PushBack(item)
}

func (q *Queue) PopFront() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if it := q.container.Front(); it != nil {
		return q.container.Remove(it)
	}

	return nil
}

func (q *Queue) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return q.container.Len()
}

type QueueTraversalHandle func(interface{}) error

func (q *Queue) Traversal(handle QueueTraversalHandle) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for e := q.container.Front(); e != nil; e = e.Next() {
		if err := handle(e.Value); err != nil {
			break
		}
	}
}
