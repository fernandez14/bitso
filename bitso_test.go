package bitso

import (
	"bytes"
	"context"
	"net/http"
	"testing"
)

func testClient() *Service {
	return New(http.DefaultClient, "test")
}

func TestNew(t *testing.T) {
	n := testClient()
	if n.sling == nil {
		t.Error("sling http client pointer is nil")
	}
	if bytes.Compare(n.secret, []byte("test")) != 0 {
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
