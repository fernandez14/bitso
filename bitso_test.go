package bitso

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"testing"
)

var (
	apiKey string
	secret string
)

func testClient() *Service {
	if v, exists := os.LookupEnv("BITSO_KEY"); exists {
		apiKey = v
	}
	if v, exists := os.LookupEnv("BITSO_SECRET"); exists {
		secret = v
	}
	return New(http.DefaultClient, apiKey, secret)
}

func TestNew(t *testing.T) {
	n := testClient()
	if n.sling == nil {
		t.Error("sling http client pointer is nil")
	}
	if bytes.Compare(n.secret, []byte(secret)) != 0 {
		t.Error("secret bytes are unexpectedly unequal")
	}
}

func TestService_Ticker(t *testing.T) {
	c := testClient()
	ticker, err := c.Ticker(context.Background(), "btc_mxn")
	if err != nil {
		t.Errorf("ticker failed	err=%v", err)
	}
	if ticker.Success == false {
		t.Error("tick did not succeed")
	}
	if ticker.Tick.Book != "btc_mxn" {
		t.Errorf("tick invalid book received 	value=%v", ticker.Tick.Book)
	}
	t.Logf("tick read	tick=%+v", ticker.Tick)
}

func TestService_AvailableBooks(t *testing.T) {
	c := testClient()
	books, err := c.AvailableBooks(context.Background())
	if err != nil {
		t.Errorf("order book failed	err=%v", err)
	}
	if books.Success == false {
		t.Error("tick did not succeed")
	}
	if len(books.List) == 0 {
		t.Errorf("empty available_books 	value=%v", books)
	}
	t.Logf("available_books read	boosk=%+v", books.List)
}

func TestService_OrderBook(t *testing.T) {
	c := testClient()
	book, err := c.OrderBook(context.Background(), OrderBookParams{
		Book:      "btc_mxn",
		Aggregate: false,
	})
	if err != nil {
		t.Errorf("order book failed	err=%v", err)
	}
	if book.Success == false {
		t.Error("tick did not succeed")
	}
	if book.OrderBook.Sequence == "" {
		t.Errorf("tick invalid book sequence received 	value=%v", book.OrderBook.Sequence)
	}
	t.Logf("order_book read	book=%+v", book.OrderBook)
}

func TestService_Balance(t *testing.T) {
	c := testClient()
	balance, err := c.Balance(context.Background())
	if err != nil {
		t.Errorf("balance failed	err=%v", err)
	}
	if balance.Success == false {
		t.Error("balance did not succeed")
	}
	if len(balance.Balances.List) == 0 {
		t.Error("balance response empty received")
	}
	t.Logf("balances read	balance=%+v", balance.Balances)
}

func TestService_OpenOrders(t *testing.T) {
	c := testClient()
	open, err := c.OpenOrders(context.Background(), OpenOrderParams{Limit: 100})
	if err != nil {
		t.Errorf("open_orders failed	err=%v", err)
	}
	if open.Success == false {
		t.Error("open_orders did not succeed")
	}
	t.Logf("open_orders read	len=%+v", len(open.List))
}
