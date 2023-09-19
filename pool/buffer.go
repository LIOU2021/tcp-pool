package pool

type circularQueue struct {
	items []*connWithTime // 环形队列的存储空间
	n     int             // 环形队列的容量
	head  int             // 队首指针
	tail  int             // 队尾指针
}

func newCircularQueue(n int) *circularQueue {
	return &circularQueue{
		items: make([]*connWithTime, n),
		n:     n,
		head:  0,
		tail:  0,
	}
}

func (q *circularQueue) enqueue(item *connWithTime) bool {
	if q.isFull() {
		item.Close()
		return false
	}
	q.items[q.tail] = item
	q.tail = (q.tail + 1) % q.n
	return true
}

func (q *circularQueue) dequeue() *connWithTime {
	if q.isEmpty() {
		return nil
	}
	item := q.items[q.head]
	q.items[q.head] = nil
	q.head = (q.head + 1) % q.n
	return item
}

func (q *circularQueue) isFull() bool {
	return (q.tail+1)%q.n == q.head
}

func (q *circularQueue) isEmpty() bool {
	return q.head == q.tail
}

func (q *circularQueue) size() int {
	return (q.tail - q.head + q.n) % q.n
}

// 遍历
func (q *circularQueue) each(fn func(node *connWithTime)) {
	for i := q.head; i < q.head+q.size(); i++ {
		fn(q.items[i%q.n])
	}
}

// 清空
func (q *circularQueue) clear() bool {
	q.n = 0
	q.head = 0
	q.tail = 0
	q.items = nil
	return true
}
