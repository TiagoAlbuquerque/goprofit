package avl

import (
//    "fmt"
)

type Data interface {
    Less (d Data) bool
}
type tAvl struct {
    data Data
    height int
    lAvl, rAvl *tAvl
}
type Avl struct {
    treeRoot *tAvl
}

func (a *tAvl) inLine(cData chan Data) {
    if a == nil { return }
    (*a).lAvl.inLine(cData)
    cData <- (*a).data
    (*a).rAvl.inLine(cData)
}

func (a *tAvl) inLineR(cData chan Data) {
    if a == nil { return }
    (*a).rAvl.inLineR(cData)
    cData <- (*a).data
    (*a).lAvl.inLineR(cData)
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

func (a *tAvl) balance(d Data) *tAvl {
    if (*a).lAvl.getHeight() - (*a).rAvl.getHeight() == 2 {
        if !d.Less((*a).lAvl.data) {
            (*a).lAvl = (*a).lAvl.lRotate()
        }
        a = a.rRotate()
    } else if (*a).rAvl.getHeight() - (*a).lAvl.getHeight() == 2 {
        if d.Less((*a).rAvl.data) {
            (*a).rAvl = (*a).rAvl.rRotate()
        }
        a = a.lRotate()
    }
    return a
}

func (a *tAvl) put(d Data) *tAvl{
    if a == nil {
        return &tAvl{d, 0, nil, nil}
    }
    println("not nil")
    if d.Less((*a).data) {
        (*a).lAvl = (*a).lAvl.put(d)
    } else {
        (*a).rAvl = (*a).rAvl.put(d)
    }
    a.updateHeight()
    a = a.balance(d)
    return a
}

func (a *tAvl) iter(reversed bool) chan Data{
    cOut := make (chan Data)
    defer close(cOut)
    if reversed { defer a.inLineR(cOut)
    } else { defer a.inLine(cOut) }
    return cOut
}

func (a *Avl) Put(d Data) {
    (*a).treeRoot = (*a).treeRoot.put(d)
}
