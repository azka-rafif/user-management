package status

import (
	"strings"
)

type OrderStatus int

const (
	Pending OrderStatus = iota
	Paid
	Shipped
)

func GetStatusFromString(s string) OrderStatus {
	switch strings.ToLower(s) {
	case "pending":
		return Pending
	case "paid":
		return Paid
	case "shipped":
		return Shipped
	default:
		return Pending
	}
}

func GetStringFromStatus(s OrderStatus) string {
	switch s {
	case Paid:
		return "Paid"
	case Pending:
		return "Pending"
	case Shipped:
		return "Shipped"
	default:
		return "Pending"
	}
}
