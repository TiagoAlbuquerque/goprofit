package avl

import (
	"sync"
)

type Data interface {
	Less(d *Data) bool
}

type stack struct {
	next *stack
	val  *tAvl
}

type Iterator struct {
	top *stack
	val *tAvl
}

type tAvl struct {
	data       *Data
	height     int
	lAvl, rAvl *tAvl
}

type Avl struct {
	rev   bool //reversed
	root  *tAvl
	mutex *sync.Mutex
}

const REVERSED = true
const DIRECT = false

func (itp *Iterator) Next() bool {
	if itp.top == nil {
		return false
	}

	itp.val = itp.top.val
	itp.top = itp.top.next

	if itp.val.rAvl != nil {
		next := stack{itp.top, itp.val.rAvl}
		itp.top = &next
		tree := itp.val.rAvl
		for tree.lAvl != nil {
			itp.top = &stack{itp.top, tree.lAvl}
			tree = tree.lAvl
		}
	}
	return true
}

func (itp *Iterator) Value() *Data {
	return itp.val.data
}

func (a *tAvl) getHeight() int {
	if a == nil {
		return -1
	}
	return a.height
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (a *tAvl) updateHeight() {
	a.height = max(a.lAvl.getHeight(), a.rAvl.getHeight()) + 1
}

func (a *tAvl) rRotate() *tAvl {
	node := a.lAvl
	a.lAvl = node.rAvl
	node.rAvl = a
	a.updateHeight()
	node.updateHeight()
	return node
}

func (a *tAvl) lRotate() *tAvl {
	node := a.rAvl
	a.rAvl = node.lAvl
	node.lAvl = a
	a.updateHeight()
	node.updateHeight()
	return node
}

func (a *tAvl) balance() *tAvl {
	if a.lAvl.getHeight()-a.rAvl.getHeight() == 2 {
		if a.lAvl.rAvl.getHeight() > a.lAvl.lAvl.getHeight() {
			a.lAvl = a.lAvl.lRotate()
		}
		a = a.rRotate()
	} else if a.rAvl.getHeight()-a.lAvl.getHeight() == 2 {
		if a.rAvl.lAvl.getHeight() > a.rAvl.rAvl.getHeight() {
			a.rAvl = a.rAvl.rRotate()
		}
		a = a.lRotate()
	}
	return a
}

func (a *tAvl) put(d *Data, rev bool) *tAvl {
	if a == nil {
		return &tAvl{d, 0, nil, nil}
	}
	if rev != (*d).Less(a.data) {
		a.lAvl = a.lAvl.put(d, rev)
	} else {
		a.rAvl = a.rAvl.put(d, rev)
	}
	a.updateHeight()
	a = a.balance()
	return a
}

func NewAvl(reversed bool) *Avl {
	out := new(Avl)
	out.rev = reversed
	out.mutex = new(sync.Mutex)
	return out
}

func (a *Avl) Put(d *Data) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.root = a.root.put(d, a.rev)
}

func (a *Avl) GetIterator() *Iterator {
	out := new(Iterator)

	if a.root != nil {
		out.top = &stack{out.top, a.root}
		tree := a.root
		for tree.lAvl != nil {
			out.top = &stack{out.top, tree.lAvl}
			tree = tree.lAvl
		}
	}
	return out
}
