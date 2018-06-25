package deals

import (
    "../items"
    "../orders"
//    "../utils/avl"
    "fmt"
    "math"
)

type Deal struct{
    item *items.Item
    buyOrder, sellOrder *orders.Order
}

var deals []Deal

func (d *Deal) Key() (int64, int64) {
    return d.sellOrder.LocationID, d.buyOrder.LocationID
}

func (d *Deal) SellLocID() int64{
    return d.sellOrder.LocationID
}

func (d *Deal) BuyLocID() int64{
    return d.buyOrder.LocationID
}

func (d *Deal) Pm3() float64 {
    prf := d.profitPerUnit()
    vol := float64(d.item.Volume)
    out := prf/vol

    return out
}

func min(a, b int) int {
    if a < b { return a }
    return b
}

func (d *Deal) amount() int {
    out := min(d.buyOrder.OrderRemain(), d.sellOrder.OrderRemain())

    if out < d.buyOrder.MinVolume { out = 0 }
    if out < d.sellOrder.MinVolume { out = 0 }

    return out
}

func (d *Deal) amountCargo(cargo float64) int {
    out := d.amount()
    out = min(out, int(math.Floor(cargo/d.item.Volume)))
    return out
}

func tax() float64 {
    return 1-0.01
}

func (d *Deal) profitPerUnit() float64 {
    ppu := (d.buyOrder.Price*tax()) - d.sellOrder.Price
    return ppu
}

func (d *Deal) profitQnt(qnt int) float64 {
    ppu := d.profitPerUnit()
    out := float64(qnt)*ppu

    return out
}

func (d *Deal) Execute(cargo float64) (float64, float64, string) {
    itmVol := d.item.Volume
    itmName := d.item.Name

    qnt := d.amountCargo(cargo)

    d.buyOrder.Execute(qnt)
    d.sellOrder.Execute(qnt)

    vol := float64(qnt)*itmVol
    cargo -= vol

    bFor := d.sellOrder.Price
    sFor := d.buyOrder.Price
    profit := d.profitQnt(qnt)

    strg := fmt.Sprintf("\n%d\tx %s \tbuy for: %.2f \tsell for: %.2f \tprofit: %.2f",
                        qnt,
                        itmName,
                        bFor,
                        sFor,
                        profit)

    return cargo, profit, strg
}

func (d *Deal) Reset() {
    d.buyOrder.Reset()
    d.sellOrder.Reset()
}

func makeDeal(item *items.Item, bOrder *orders.Order, sOrder *orders.Order, cDeals chan *Deal) bool {
    d := Deal{item, bOrder, sOrder}
    if d.profitPerUnit() > 0.0 {
        deals = append(deals, d)
        cDeals <- &d
        return true
    }
    return false
}

func computeBuyOrder(bOrder *orders.Order, cDeals chan *Deal) {
    item := items.GetItem(bOrder.ItemID)
    //iter := item.Buy_orders.GetIterator()

    for _, sOrder := range item.Sell_orders {
        //sOrder := (*iter.Value()).(items.OrderAvlData).Order
//        println()
  //      println(sOrder.ItemID)
    //    println(bOrder.ItemID)
        makeDeal(item, bOrder, sOrder, cDeals)
        /*if !makeDeal(item, bOrder, sOrder, cDeals) {
            break
        }*/
    }
}

func computeSellOrder(sOrder *orders.Order, cDeals chan *Deal) {
    item := items.GetItem(sOrder.ItemID)
    //iter := item.Sell_orders.GetIterator()

    //for iter.Next() {
    for _, bOrder := range item.Buy_orders {
        //bOrder := (*iter.Value()).(items.OrderAvlData).Order
        makeDeal(item, bOrder, sOrder, cDeals)
        /*if !makeDeal(item, bOrder, sOrder, cDeals) {
            break
        }*/
    }
}

func Cleanup() {
    deals = []Deal{}
}

func ComputeDeals(o *orders.Order, cDeals chan *Deal) {
    if o.IsBuyOrder {
        computeBuyOrder(o, cDeals)
    } else {
        computeSellOrder(o, cDeals)
    }

}
