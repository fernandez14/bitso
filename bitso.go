package bitso

import (
	"context"
	"github.com/dghubble/sling"
	"net/http"
	"time"
)

var (
	BaseURL = "https://api.bitso.com"
)

// Service is the API gateway struct.
type Service struct {
	secret []byte
	sling  *sling.Sling
}

// New setups a service struct to interact with bitso.com API.
func New(httpClient *http.Client, secret string) *Service {
	return &Service{
		secret: []byte(secret),
		sling:  sling.New().Client(httpClient).Base(BaseURL),
	}
}

type Tick struct {
	Book     string `json:"book"`
	High     string `json:"high"`
	Vwap     string `json:"vwap"`
	Volume   string `json:"volume"`
	Last     string `json:"last"`
	Low      string `json:"low"`
	Ask      string `json:"ask"`
	Bid      string `json:"bid"`
	Change24 string `json:"change_24"`
}

type tickResponse struct {
	Success bool
	Http    *http.Response
	Tick    Tick `json:"payload"`
}

type TickerParams struct {
	Book string `url:"book"`
}

// Ticker performs a request to a book with bitso:v3/ticker
func (srv *Service) Ticker(ctx context.Context, book string) (r tickResponse, err error) {
	var (
		resp *http.Response
		req  *http.Request
	)
	get := srv.sling.New().Get("v3/ticker").QueryStruct(TickerParams{Book: book})
	withAuth := bitsoAuth(get, srv.secret)
	req, err = withAuth.Request()
	if err != nil {
		return
	}
	req = req.WithContext(ctx)
	resp, err = get.Do(req, &r, nil)
	if err != nil {
		return
	}
	r.Http = resp
	r.Tick.Book = book
	return
}

type orderBookResponse struct {
	Success   bool
	Http      *http.Response
	OrderBook OrderBook `json:"payload"`
}

type OrderBookParams struct {
	Book      string `url:"book"`
	Aggregate bool   `url:"aggregate,omitempty"`
}

type OpenOrder struct {
	Book   string `json:"book"`
	Price  string `json:"price"`
	Amount string `json:"amount"`
}

type OrderBook struct {
	Asks     []OpenOrder `json:"asks"`
	Bids     []OpenOrder `json:"bids"`
	Updated  time.Time   `json:"updated_at"`
	Sequence string      `json:"sequence"`
}

// Ticker performs a request to a book with bitso:v3/ticker
func (srv *Service) OrderBook(ctx context.Context, params OrderBookParams) (r orderBookResponse, err error) {
	var (
		resp *http.Response
		req  *http.Request
	)
	get := srv.sling.New().Get("v3/order_book").QueryStruct(params)
	withAuth := bitsoAuth(get, srv.secret)
	req, err = withAuth.Request()
	if err != nil {
		return
	}
	req = req.WithContext(ctx)
	resp, err = get.Do(req, &r, nil)
	if err != nil {
		return
	}
	r.Http = resp
	return
}
