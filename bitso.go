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

// Ticker performs a request to bitso:v3/ticker
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

// OrderBook performs a request to bitso:v3/order_book
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

// Balance performs a request to bitso:v3/balance
func (srv *Service) Balance(ctx context.Context) (r balanceResponse, err error) {
	var (
		resp *http.Response
	)
	resp, err = srv.doReq(ctx, srv.sling.New().Get("v3/balance"), &r)
	r.Http = resp
	return
}

type AvailableBook struct {
	Book          string `json:"currency"`
	MinimumAmount string `json:"minimum_amount"`
	MaximumAmount string `json:"maximum_amount"`
	MinimumPrice  string `json:"minimum_price"`
	MaximumPrice  string `json:"maximum_price"`
	MinimumValue  string `json:"minimum_value"`
	MaximumValue  string `json:"maximum_value"`
}

type availableBooksResponse struct {
	Success bool
	Http    *http.Response
	List    []AvailableBook `json:"payload"`
}

// AvailableBooks performs a request to bitso:v3/available_books
func (srv *Service) AvailableBooks(ctx context.Context) (r availableBooksResponse, err error) {
	var (
		resp *http.Response
	)
	resp, err = srv.doReq(ctx, srv.sling.New().Get("v3/available_books"), &r)
	r.Http = resp
	return
}

type cancelOrderResponse struct {
	Success bool
	Http    *http.Response
	List    []string `json:"payload"`
}

type CancelOrderParams struct {
	OrderIDs  string `url:"oids,omitempty"`
	OriginIDs string `url:"origin_ids,omitempty"`
	All       bool   `url:"-"`
}

// CancelOrder performs a request to bitso:v3/orders
func (srv *Service) CancelOrder(ctx context.Context, params CancelOrderParams) (r cancelOrderResponse, err error) {
	var (
		resp *http.Response
	)
	req := srv.sling.New()
	if params.All {
		req = req.Delete("v3/orders/all")
	} else {
		req = req.QueryStruct(params)
	}
	resp, err = srv.doReq(ctx, req, &r)
	r.Http = resp
	return
}

type placeOrderResponse struct {
	Success bool
	Http    *http.Response
	List    []string `json:"payload"`
}

type PlaceOrderParams struct {
	Book        string `json:"book"`
	Side        string `json:"side"`
	Type        string `json:"type"`
	Major       string `json:"type,omitempty"`
	Minor       string `json:"minor,omitempty"`
	Price       string `json:"price,omitempty"`
	Stop        string `json:"stop,omitempty"`
	TimeInForce string `json:"time_in_force,omitempty"`
	OriginID    string `json:"origin_id,omitempty"`
}

// PlaceOrder performs a request to bitso:v3/orders
func (srv *Service) PlaceOrder(ctx context.Context, params PlaceOrderParams) (r placeOrderResponse, err error) {
	var (
		resp *http.Response
	)
	req := srv.sling.New().Post("v3/orders").BodyJSON(&params)
	resp, err = srv.doReq(ctx, req, &r)
	r.Http = resp
	return
}

func (srv *Service) Bid(ctx context.Context, amount, price, book string) (r placeOrderResponse, err error) {
	r, err = srv.PlaceOrder(ctx, PlaceOrderParams{
		Book:  book,
		Side:  "buy",
		Type:  "limit",
		Major: amount,
		Price: price,
	})
	return
}

func (srv *Service) Ask(ctx context.Context, amount, price, book string) (r placeOrderResponse, err error) {
	r, err = srv.PlaceOrder(ctx, PlaceOrderParams{
		Book:  book,
		Side:  "sell",
		Type:  "limit",
		Major: amount,
		Price: price,
	})
	return
}
