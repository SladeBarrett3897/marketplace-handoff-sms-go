package main

import "fmt"

type Order struct{ ID, SellerName, BuyerPhone, Asset string }

func HandoffMessage(o Order) string {
	return fmt.Sprintf("Order %s: %s is ready for handoff. Seller asset: %s.", o.ID, o.SellerName, o.Asset)
}

func NotifyHandoff(c *SMSClient, o Order) (string, error) {
	return c.Send(o.BuyerPhone, HandoffMessage(o), "order-handoff-"+o.ID)
}
