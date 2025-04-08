package main

import (
	"fmt"
	"sort"
	"time"
)

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}

type Order struct {
	Size      float64 // Amount to buy/sell
	Bid       bool    // true = buying, false = selling
	Limit     *Limit  // Price point this order is at
	Timestamp int64   // When the order was created
}

type Orders []*Order

func (o Orders) Len() int           { return len(o) }
func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool { return o[i].Timestamp < o[j].Timestamp }

// NewOrder function creates a pointer to a new order with given size and type of order
func NewOrder(bid bool, size float64) *Order {
	return &Order{
		Size:      size,
		Bid:       bid,
		Timestamp: time.Now().UnixNano(),
	}
}

func (o *Order) String() string {
	return fmt.Sprintf("[size %.2f]", o.Size)
}

type Limit struct {
	Price       float64 // The price point to buy/sell at
	Orders      Orders  // List of orders at this price
	TotalVolume float64 // Combined size of all orders bids or asks
}

type Limits []*Limit

type ByBestAsk struct{ Limits }

func (a ByBestAsk) Len() int           { return len(a.Limits) }
func (a ByBestAsk) Swap(i, j int)      { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }
func (a ByBestAsk) Less(i, j int) bool { return a.Limits[i].Price < a.Limits[j].Price }

type ByBestBid struct{ Limits }

func (b ByBestBid) Len() int           { return len(b.Limits) }
func (b ByBestBid) Swap(i, j int)      { b.Limits[i], b.Limits[j] = b.Limits[j], b.Limits[i] }
func (b ByBestBid) Less(i, j int) bool { return b.Limits[i].Price > b.Limits[j].Price }

// NewLimit creates and initializes a new Limit struct at the specified price
// A limit represents a price level where multiple orders can be placed
func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,      // The price point this limit represents
		Orders: []*Order{}, // Initialize an empty slice to store orders at this price
		// Note: TotalVolume starts at 0 by default
	}
}

// AddOrder method on Limit puts a new order at this price and updates the total volume.
func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size
}

// DeleteOrder function removes an order and reduces the total volume.
func (l *Limit) DeleteOrder(o *Order) {
	for i := range l.Orders {
		if l.Orders[i] == o {
			// Shift all elements after i one position to the left
			copy(l.Orders[i:], l.Orders[i+1:])
			l.Orders = l.Orders[:len(l.Orders)-1]
			break
		}
	}
	o.Limit = nil
	l.TotalVolume -= o.Size

	sort.Sort(l.Orders)
}

type Orderbook struct {
	Asks      []*Limit           // Sorted array of sell orders by price (lowest first)
	Bids      []*Limit           // Sorted array of buy orders by price (highest first)
	AskLimits map[float64]*Limit // Fast lookup table for sell price levels
	BidLimits map[float64]*Limit // Fast lookup table for buy price levels
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		Asks:      []*Limit{},               // Initialize empty slice for sell orders
		Bids:      []*Limit{},               // Initialize empty slice for buy orders
		AskLimits: make(map[float64]*Limit), // Initialize empty map for ask price lookups, map[ask_price] *BuyLimit
		BidLimits: make(map[float64]*Limit), // Initialize empty map for bid price lookups, map[bid_price] *SellLimit
	}
}

func (ob *Orderbook) PlaceOrder(price float64, o *Order) []Match {
	// 1. try to match the order
	// matching logic

	// 2. add the rest of the order to the books
	if o.Size > 0.0 {
		ob.add(price, o)
	}
	return []Match{}
}

func (ob *Orderbook) add(price float64, o *Order) {
	var limit *Limit
	// limit.AddOrder(o)

	if o.Bid {
		limit = ob.BidLimits[price]
	} else {
		limit = ob.AskLimits[price]
	}

	if limit == nil {
		limit = NewLimit(price)
		if o.Bid {
			ob.Bids = append(ob.Bids, limit)
			ob.BidLimits[price] = limit
		} else {
			ob.Asks = append(ob.Asks, limit)
			ob.AskLimits[price] = limit
		}
	}

}
