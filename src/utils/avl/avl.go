package avl

type Data interface {
    (d Data)Less(l Data) bool
}
type Avl struct {
    data Data
    height int
    lAvl, rAvl Avl
}
var root Avl

func (avl Avl) inLine(cData chan Data) {
    if avl == nil { return }
    avl.lData.inline(cData)
    cData <- avl.data
    avl.rData.inline(cData)
}

func (avl Avl) inLineR(cData chan Data) {
    if avl == nil { return }
    avl.rData.inlineR(cData)
    cData <- avl.data
    avl.lData.inlineR(cData)
}

func height(avl Avl) {
    if avl == nil
        return -1
    return avl.height
}
func max(a, b int) int{
    if a > b { return a }
    return b
}

func (avl Avl) updateHeight() {
    avl.height = max(height(avl.lAvl), height(avl.rAvl)) +1
}

func (avl Avl) rRotate() Avl{
    node := avl.lAvl
    avl.lAvl = node.rAvl
    node.rAvl = avl
    avl.updateHeight()
    node.updateHeight()
    return node
}

func (avl Avl) lRotate() Avl{
    node := avl.rAvl
    avl.rAvl = node.lAvl
    node.lAvl = avl
    avl.updateHeight()
    node.updateHeight()
    return node
}

func (avl Avl) balance(d Data) *Avl {
    if height(avl.lAvl) - height(avl.rAvl) == 2 {
        if !d.Less(avl.lAvl.data) {
            avl.lAvl = lRotate(avl.lAvl)
        }
        alv = rRotate(avl)
    } else if height(avl.rAvl) - height(avl.lAvl) == 2 {
        if d.Less(avl.rAvl.data) {
            avl.rAvl = rRotate(avl.rAvl)
        }
        avl = avl.lRotate()
    }
    return avl
}

func (avl Avl) Insert(data Data) *Avl{
    if avl == nil {
        avl = &Avl{Data, 0, nil, nil}
        return avl
    }
    if data.Less(avl.data) {
        avl = avl.lAvl.Insert(data)
    } else {
        avl = avl.rAvl.Insert(data)
    }
    avl.updateHeight()
    avl = avl.balance(data)
    return avl
}

func (avl Avl) Iter(reversed bool) chan Data{
    cOut := make (chan Data)
    defer cout.Close()
    if reversed { defer avl.inLineR(cOut)
    } else { defer alv.inLine(cOut) }
    return cOut
}
