package services

import (
	"errors"
	"fmt"
	"portfolio-tracker/config"
	"portfolio-tracker/models"
)

// AddTrade inserts a new trade into DB
func AddTrade(symbol, side string, qty int, price float64) error {
	if side != "BUY" && side != "SELL" {
		return errors.New("invalid side: must be BUY or SELL")
	}
	if qty <= 0 || price <= 0 {
		return errors.New("quantity and price must be positive")
	}

	trade := models.Trade{
		Symbol:   symbol,
		Side:     side,
		Quantity: qty,
		Price:    price,
	}

	if err := config.DB.Create(&trade).Error; err != nil {
		return err
	}
	fmt.Println("✅ Trade added:", trade)
	return nil
}

// GetNetPosition calculates holdings per symbol
func GetNetPosition() error {
	var trades []models.Trade
	if err := config.DB.Find(&trades).Error; err != nil {
		return err
	}

	positions := make(map[string]struct {
		Qty   int
		Invest float64
	})

	for _, t := range trades {
		pos := positions[t.Symbol]
		if t.Side == "BUY" {
			pos.Qty += t.Quantity
			pos.Invest += float64(t.Quantity) * t.Price
		} else if t.Side == "SELL" {
			pos.Qty -= t.Quantity
			pos.Invest -= float64(t.Quantity) * t.Price
		}
		positions[t.Symbol] = pos
	}

	fmt.Println("📊 Net Positions:")
	for sym, p := range positions {
		fmt.Printf("%s: %d shares, Net Investment: ₹%.2f\n", sym, p.Qty, p.Invest)
	}
	return nil
}
