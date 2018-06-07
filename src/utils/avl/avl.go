package avl

import (
//    "fmt"
)

type Data interface {
    Less (d *Data) bool
}

type Iterator struct {
    stack *Iterator
    tree *tAvl
}

type tAvl struct {
    data *Data
    height int
    lAvl, rAvl *tAvl
}

type Avl struct {
    Reversed bool
    root *tAvl
}

func (itp *Iterator) Next() bool {
    if itp.stack == nil {
        return false
    }

    itp.tree = itp.stack.tree
    itp.stack = itp.stack.stack

    if itp.tree.rAvl != nil {
        next := Iterator{itp.stack, itp.tree.rAvl}
        itp.stack = &next
        tree := itp.tree.rAvl
        for tree.lAvl != nil {
            next = Iterator{itp.stack, tree.lAvl}
            itp.stack = &next
            tree = tree.lAvl
        }
    }
    return true
}

func (itp *Iterator) Value() *Data {
    return itp.tree.data
}

func (a *tAvl)getHeight() int{
    if a == nil { return -1 }
    return (*a).height
}

func max(a, b int) int{
    if a > b { return a }
    return b
}

func (a *tAvl) updateHeight() {
    (*a).height = max((*a).lAvl.getHeight(), (*a).rAvl.getHeight()) +1
}

func (a *tAvl) rRotate() *tAvl{
    node := (*a).lAvl
    (*a).lAvl = node.rAvl
    node.rAvl = a
    (*a).updateHeight()
    node.updateHeight()
    return node
}

func (a *tAvl) lRotate() *tAvl{
    node := (*a).rAvl
    (*a).rAvl = node.lAvl
    node.lAvl = a
    (*a).updateHeight()
    node.updateHeight()
    return node
}

func (a *tAvl) balance(d *Data, rev bool) *tAvl {
    if (*a).lAvl.getHeight() - (*a).rAvl.getHeight() == 2 {
        if rev != (!(*d).Less((*a).lAvl.data)) {
            (*a).lAvl = (*a).lAvl.lRotate()
        }
        a = a.rRotate()
    } else if (*a).rAvl.getHeight() - (*a).lAvl.getHeight() == 2 {
        if rev != ((*d).Less((*a).rAvl.data)) {
            (*a).rAvl = (*a).rAvl.rRotate()
        }
        a = a.lRotate()
    }
    return a
}

func (a *tAvl) put(d *Data, rev bool) *tAvl{
    if a == nil {
        return &tAvl{d, 0, nil, nil}
    }
    if rev != (*d).Less((*a).data) {
        (*a).lAvl = (*a).lAvl.put(d, rev)
    } else {
        (*a).rAvl = (*a).rAvl.put(d, rev)
    }
    a.updateHeight()
    a = a.balance(d, rev)
    return a
}

func NewAvl(reversed bool) Avl {
    return Avl{reversed, nil}
}

func (a *Avl) Put(d *Data) {
    (*a).root = (*a).root.put(d, a.Reversed)
}

func (a *Avl) GetIterator() *Iterator{
    out := Iterator{nil, nil}


    return &out
}
