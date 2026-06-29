package students

type FinanceSnapshot struct {
	Currency     string  `json:"currency" bson:"currency"`
	TotalDue     float64 `json:"totalDue" bson:"totalDue"`
	TotalPaid    float64 `json:"totalPaid" bson:"totalPaid"`
	Balance      float64 `json:"balance" bson:"balance"`
	Overpayment  float64 `json:"overpayment" bson:"overpayment"`
	InvoiceCount int64   `json:"invoiceCount" bson:"invoiceCount"`
}
