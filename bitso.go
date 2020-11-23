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
	apiKey string
	sling  *sling.Sling
}

func (srv *Service) doReq(ctx context.Context, slingReq *sling.Sling, success interface{}) (resp *http.Response, err error) {
	var (
		req *http.Request
	)
	withAuth := bitsoAuth(slingReq, srv.apiKey, srv.secret)
	req, err = withAuth.Request()
	if err != nil {
		return
	}
	req = req.WithContext(ctx)
	resp, err = withAuth.Do(req, &success, nil)
	return
}

// New setups a service struct to interact with bitso.com API.
func New(httpClient *http.Client, apiKey, secret string) *Service {
	return &Service{
		secret: []byte(secret),
		apiKey: apiKey,
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
	)
	resp, err = srv.doReq(ctx, srv.sling.New().Get("v3/ticker").QueryStruct(TickerParams{Book: book}), &r)
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

// OrderBook performs a request to a book with bitso:v3/order_book
func (srv *Service) OrderBook(ctx context.Context, params OrderBookParams) (r orderBookResponse, err error) {
	var (
		resp *http.Response
	)
	resp, err = srv.doReq(ctx, srv.sling.New().Get("v3/order_book").QueryStruct(params), &r)
	r.Http = resp
	return
}

type Balance struct {
	Currency  string  `json:"currency"`
	Total     float64 `json:"total,string"`
	Locked    float64 `json:"locked,string"`
	Available float64 `json:"available,string"`
}

type BalanceRes struct {
	List []Balance `json:"balances"`
}

type balanceResponse struct {
	Success  bool
	Http     *http.Response
	Balances BalanceRes `json:"payload"`
}

// Balance performs a request to a book with bitso:v3/balance
func (srv *Service) Balance(ctx context.Context) (r balanceResponse, err error) {
	var (
		resp *http.Response
	)
	resp, err = srv.doReq(ctx, srv.sling.New().Get("v3/balance"), &r)
	r.Http = resp
	return
}
