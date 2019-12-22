package utils

// Queue ...
type Queue struct {
	items items
}

type items []interface{}

// NewQueue get the a new instance from queue datastructure
func NewQueue() Queue {
	return Queue{
		items: make([]interface{}, 0),
	}
}

// Push pushs a new item to the queue
func (q Queue) Push(item interface{}) Queue {
	q.items = append(q.items, item)
	return q
}

// IsEmpty checks if a queue is empty or not
func (q Queue) IsEmpty() bool {
	return len(q.items) == 0
}

// Pop pops an item from the queue
func (q *Queue) Pop() (e interface{}) {
	e, q.items = q.items[0], q.items[1:len(q.items)]
	return
}
