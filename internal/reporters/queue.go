package reporters

type Queue[T any] struct {
	data        []T
	head, tail  int
	size, limit int
}

func New[T any](limit int) *Queue[T] {
	return &Queue[T]{
		data:  make([]T, limit),
		limit: limit,
	}
}

func (q *Queue[T]) Push(v T) bool {
	if q.size == q.limit {
		return false // очередь полна
	}
	q.data[q.tail] = v
	q.tail++
	if q.tail == q.limit {
		q.tail = 0
	}
	q.size++
	return true
}

// Pop извлекает элемент из начала очереди
func (q *Queue[T]) Pop() (v T, ok bool) {
	if q.size == 0 {
		return v, false // пусто
	}
	v = q.data[q.head]
	q.head++
	if q.head == q.limit {
		q.head = 0
	}
	q.size--
	return v, true
}

func (q *Queue[T]) Peek() (v T, ok bool) {
	if q.size == 0 {
		return v, false
	}
	return q.data[q.head], true
}

func (q *Queue[T]) Len() int { return q.size }

func (q *Queue[T]) Cap() int { return q.limit }

func (q *Queue[T]) Reset() {
	q.head, q.tail, q.size = 0, 0, 0
}
