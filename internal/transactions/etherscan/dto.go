package etherscan

type txListResponse struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Result  []txListItem `json:"result"`
}

type txListItem struct {
	BlockNumber string `json:"blockNumber"`
	TimeStamp   string `json:"timeStamp"`
	Hash        string `json:"hash"`

	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`

	Input        string `json:"input"`
	MethodID     string `json:"methodId"`
	FunctionName string `json:"functionName"`

	TxReceiptStatus string `json:"txreceipt_status"`
	IsError         string `json:"isError"`
}
