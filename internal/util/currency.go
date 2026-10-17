package util

import "strings"

const (
	USD = "USD"
	EUR = "EUR"
	UGX = "UGX"
)

func IsSupportedCurrency(currency string) bool {
	currency = strings.ToUpper(currency)
	return currency == USD || currency == EUR || currency == UGX
}
