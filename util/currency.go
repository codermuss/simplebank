package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	TL  = "TL"
)

var SupportedCurrencies = []string{USD, EUR, CAD, TL}

// * Note [codermuss]: Returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, TL:
		return true
	}
	return false
}
