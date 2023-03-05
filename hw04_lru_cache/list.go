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
	len   int
	first *ListItem
	last  *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	if l.first != nil {
		l.first.Prev = &ListItem{
			Next: l.first,
		}
		l.first = l.first.Prev
	} else {
		l.first = &ListItem{}
		l.last = l.first
	}

	l.first.Value = v
	l.len++

	return l.first
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.last != nil {
		l.last.Next = &ListItem{
			Prev: l.last,
		}
		l.last = l.last.Next
	} else {
		l.last = &ListItem{}
		l.first = l.last
	}

	l.last.Value = v
	l.len++

	return l.last
}

func (l *list) Remove(i *ListItem) {
	if i != nil {
		l.extractElement(i)
		l.len--
		switch l.len {
		case 0:
			l.first = nil
			l.last = nil
		case 1:
			if i == l.first {
				l.first = l.last
			}
			if i == l.last {
				l.last = l.first
			}
		}
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i != l.first && i != nil {
		l.extractElement(i)
		l.first.Prev = i
		i.Next = l.first
		i.Prev = nil
		l.first = i
	}
}

func (l *list) extractElement(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
}

func NewList() List {
	return &list{}
}
