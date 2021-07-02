package common

type node struct {
	element interface{}
	Next    *node
}

type Queue struct {
	size  int
	First *node
	Last  *node
}

// Returns the size of the queue
func (q *Queue) Size() int {
	return q.size
}

// Adds an element to the queue
func (q *Queue) Add(element interface{}) {
	e := &node{
		element: element,
		Next:    nil,
	}
	if q.First == nil {
		q.First = e
		q.Last = e
	} else {
		q.Last.Next = e
		q.Last = e
	}
	q.size++
}

// Extracts an element from the queue
func (q *Queue) Remove() interface{} {
	if q.First != nil {
		temp := q.First

		q.First = q.First.Next
		q.size--
		return temp.element
	}
	return nil
}

// Checks whether the queue has another element
func (q *Queue) HasNext() bool {
	return q.size > 0
}
