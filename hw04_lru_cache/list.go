package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v any) *ListItem
	PushBack(v any) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value any
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func NewList() *list {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v any) *ListItem {
	prevFront := l.front
	f := &ListItem{
		Value: v,
		Next:  prevFront,
	}

	if prevFront != nil {
		prevFront.Prev = f
	}

	l.front = f
	if l.back == nil {
		l.back = f
	}
	l.len++

	return f
}

func (l *list) PushBack(v any) *ListItem {
	prevBack := l.back
	b := &ListItem{
		Value: v,
		Prev:  prevBack,
	}

	if prevBack != nil {
		prevBack.Next = b
	}

	l.back = b
	if l.front == nil {
		l.front = b
	}
	l.len++

	return b
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}
	l.shrink(i)
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || i == l.front {
		return
	}

	f := i
	l.shrink(i)

	prevFront := l.front
	f.Next = prevFront
	if prevFront != nil {
		prevFront.Prev = f
	}
	l.front = f
}

func (l *list) shrink(i *ListItem) {
	switch {
	case i == l.front && i == l.back:
		l.front = l.front.Next
		l.back = l.back.Prev
		return
	case i == l.front:
		l.front = l.front.Next
	case i == l.back:
		l.back = l.back.Prev
	}

	p, n := i.Prev, i.Next
	if p != nil {
		p.Next = n
	}
	if n != nil {
		n.Prev = p
	}
}
