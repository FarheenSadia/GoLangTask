package main

import "fmt"

type Order interface {
	Execute() error
}

type MarketOrder struct {
	Quantity    int
	Stock       string
	MarketPrice float64
}

func (m MarketOrder) Execute() error {
	total := float64(m.Quantity) * m.MarketPrice
	fmt.Printf("%s: %d shares, Net Investment: ₹%.2f\n", m.Stock, m.Quantity, total)
	return nil
}

type LimitOrder struct {
	Quantity    int
	Stock       string
	LimitPrice  float64
	MarketPrice float64
}

func (l LimitOrder) Execute() error {
	if l.LimitPrice < l.MarketPrice {
		return fmt.Errorf("limit price %.2f is below market price %.2f", l.LimitPrice, l.MarketPrice)
	}
	total := float64(l.Quantity) * l.LimitPrice
	fmt.Printf("%s: %d shares, Net Investment: ₹%.2f\n", l.Stock, l.Quantity, total)
	return nil
}

func ProcessOrder(o Order) {
	err := o.Execute()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func main() {
	m1 := MarketOrder{60, "INFY", 1400.00}
	l1 := LimitOrder{100, "AAPL", 174.25, 170.00}
	l2 := LimitOrder{200, "TSLA", 680.00, 700.00}

	ProcessOrder(m1)
	ProcessOrder(l1)
	ProcessOrder(l2)
}
