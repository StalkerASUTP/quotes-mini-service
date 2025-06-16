package quote

type CreateQuoteRequest struct {
	Author string `json:"author"`
	Quote  string `json:"quote"`
}
type GetWithParamResponse struct {
	Quotes []Quote `json:"quotes"`
	Count  int64   `json:"count" example:"100"`
}
