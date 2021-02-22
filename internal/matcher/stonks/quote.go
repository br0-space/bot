package stonks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kmptnz/bot/internal/config"
)

type ResponseQuote struct {
	ShortName                  string  `json:"shortName"`
	Symbol                     string  `json:"symbol"`
	MarketState                string  `json:"marketState"`
	Currency                   string  `json:"currency"`
	ExchangeName               string  `json:"fullExchangeName"`
	ExchangeDelay              float64 `json:"exchangeDataDelayedBy"`
	RegularMarketChange        float64 `json:"regularMarketChange"`
	RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
	RegularMarketPrice         float64 `json:"regularMarketPrice"`
	RegularMarketPreviousClose float64 `json:"regularMarketPreviousClose"`
	RegularMarketOpen          float64 `json:"regularMarketOpen"`
	RegularMarketDayRange      string  `json:"regularMarketDayRange"`
	PostMarketChange           float64 `json:"postMarketChange"`
	PostMarketChangePercent    float64 `json:"postMarketChangePercent"`
	PostMarketPrice            float64 `json:"postMarketPrice"`
	PreMarketChange            float64 `json:"preMarketChange"`
	PreMarketChangePercent     float64 `json:"preMarketChangePercent"`
	PreMarketPrice             float64 `json:"preMarketPrice"`
}

type Quote struct {
	ResponseQuote
	Price                   float64
	Change                  float64
	ChangePercent           float64
	IsActive                bool
	IsRegularTradingSession bool
}

type Response struct {
	QuoteResponse struct {
		Quotes []ResponseQuote `json:"result"`
		Error  interface{}     `json:"error"`
	} `json:"quoteResponse"`
}

func getQuotes(symbols []string) ([]Quote, error) {
	url := fmt.Sprintf(config.Cfg.StonksMatcher.QuotesUrl, strings.Join(symbols, ","))

	request, _ := http.NewRequest("GET", url, nil)
	response, _ := http.DefaultClient.Do(request)

	responseBody := &Response{}
	if err := json.NewDecoder(response.Body).Decode(responseBody); err != nil {
		return nil, err
	}

	return transformResponseQuotes(responseBody.QuoteResponse.Quotes), nil
}

func transformResponseQuotes(responseQuotes []ResponseQuote) []Quote {
	quotes := make([]Quote, 0)
	for _, responseQuote := range responseQuotes {
		quotes = append(quotes, transformResponseQuote(responseQuote))
	}
	return quotes
}

func transformResponseQuote(responseQuote ResponseQuote) Quote {
	if responseQuote.MarketState == "REGULAR" {
		return Quote{
			ResponseQuote:           responseQuote,
			Price:                   responseQuote.RegularMarketPrice,
			Change:                  responseQuote.RegularMarketChange,
			ChangePercent:           responseQuote.RegularMarketChangePercent,
			IsActive:                true,
			IsRegularTradingSession: true,
		}
	}

	if responseQuote.MarketState == "CLOSED" {
		return Quote{
			ResponseQuote:           responseQuote,
			Price:                   responseQuote.PostMarketPrice,
			Change:                  responseQuote.PostMarketChange + responseQuote.RegularMarketChange,
			ChangePercent:           responseQuote.PostMarketChangePercent + responseQuote.RegularMarketChangePercent,
			IsActive:                true,
			IsRegularTradingSession: false,
		}
	}

	if responseQuote.MarketState == "PRE" {
		return Quote{
			ResponseQuote:           responseQuote,
			Price:                   responseQuote.PreMarketPrice,
			Change:                  responseQuote.PreMarketChange,
			ChangePercent:           responseQuote.PreMarketChangePercent,
			IsActive:                true,
			IsRegularTradingSession: false,
		}
	}

	return Quote{
		ResponseQuote:           responseQuote,
		Price:                   responseQuote.RegularMarketPrice,
		Change:                  0.0,
		ChangePercent:           0.0,
		IsActive:                false,
		IsRegularTradingSession: false,
	}
}
