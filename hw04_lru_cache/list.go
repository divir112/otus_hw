package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len  int
	head *ListItem
	last *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	defer func() {
		l.len++
	}()

	if l.head == nil && l.last == nil {
		node := &ListItem{
			Value: v,
		}
		l.head = node

		l.last = node

		return l.last
	}

	node := &ListItem{
		Value: v,
		Prev:  nil,
		Next:  l.head,
	}
	l.head.Prev = node
	l.head = node

	return node
}

func (l *list) PushBack(v interface{}) *ListItem {
	defer func() {
		l.len++
	}()

	if l.head == nil && l.last == nil {
		node := &ListItem{
			Value: v,
		}
		l.head = node

		l.last = node

		return l.last
	}

	node := &ListItem{
		Value: v,
		Prev:  l.last,
		Next:  nil,
	}
	l.last.Next = node
	l.last = node

	return node
}

func (l *list) Remove(i *ListItem) {
	if l.head == i && l.last == i {
		l.head = nil
		l.last = nil
		l.len--
		return
	}

	for node := l.head; node != nil; node = node.Next {
		if node == i {
			defer func() {
				l.len--
			}()
			switch i {
			case l.head:
				l.head = i.Next
				l.head.Prev = nil
				return
			case l.last:
				l.last = i.Prev
				l.last.Next = nil
				return
			}
			prevNode := node.Prev
			nextNode := node.Next

			prevNode.Next, nextNode.Prev = nextNode, prevNode
			return
		}
	}
}

func (l *list) MoveToFront(i *ListItem) {
	switch i {
	case l.last:
		i.Next = l.head
		l.last = i.Prev
		l.last.Next = nil
		l.head = i
		l.head.Prev = nil
		if l.Len() == 2 {
			l.last.Prev = l.head
		}
	case l.head:
		return
	default:
		l.Remove(i)
		i.Next = l.head
		i.Prev = nil
		l.head = i
	}
}
