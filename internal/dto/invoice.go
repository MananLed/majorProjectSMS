package dto

type Invoice struct {
	Amount float64 `json:"amount"`
}

type InvoiceDDB struct {
	PK     string  `dynamodbav:"PK"`
	SK     string  `dynamodbav:"SK"`
	ID     string  `dynamodbav:"id"`
	Amount float64 `dynamodbav:"amount"`
	Month  string  `dynamodbav:"month"`
	Year   int     `dynamodbav:"year"`
}
