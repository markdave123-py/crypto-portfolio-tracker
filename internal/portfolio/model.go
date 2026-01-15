package portfolio

// Holding represents an owned asset in a portfolio
type Holding struct {
	Chain           string  // ethereum, polygon, etc
	ContractAddress string  // empty for native asset
	Amount          float64 // the amount owned
}

// Portfolio represents a wallet portfolio snapshot
type Portfolio struct {
	Wallet   string
	Holdings []Holding
}

//
type HoldingView struct {
	Chain           string
	ContractAddress string
	Amount          float64
	PriceUSD        float64
	ValueUSD        float64
}

// portfolio to be returned with computed field TotalValueUSD
type PortfolioView struct {
	Wallet        string
	Holdings      []HoldingView
	TotalValueUSD float64
}
