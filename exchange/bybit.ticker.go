package exchange

func (e *Exchange) GetTicker(s string) (ticker Ticker, err error) {
	var tr TickerResponse
	params := make(RequestParameters)

	params["symbol"] = s
    _, err = e.PublicRequest("GET", "/v2/public/tickers", params, &tr)

	return tr.Results[0], err
}
