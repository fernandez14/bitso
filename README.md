# Bitso API golang client ![Go](https://github.com/fernandez14/bitso/workflows/Go/badge.svg) [![GoDoc](https://godoc.org/github.com/fernandez14/bitso?status.svg)](https://pkg.go.dev/github.com/fernandez14/bitso)  

<img align="right" width="100" src="https://assets.bitso.com/static/media/logo_dark.93562fe3.svg">

The Bitso API in a gopher friendly way

This library uses [dghubble/sling](https://github.com/dghubble/sling) internally: A Go HTTP client library for creating and sending API requests.

### Features

* V3 partial coverage
* Params and responses with custom structs.
* HMAC authentication implemented.
* Context for requests available.

## Install

```
go get github.com/fernandez14/bitso
```

## Documentation

Read [GoDoc](https://pkg.go.dev/github.com/fernandez14/bitso)


## Usage

Use bitso to set path, method, header, query, or body properties and create an `http.Request`.

```go
// A example request to ticker endpoint... 
srv := bitso.New(http.DefaultClient, "account_api_secret")
res, err := srv.Ticker(context.Background(), "btc_mxn")
if err != nil {
    log.Error(err)
}
log.Printf("success=%+v", res.Success)
log.Printf("tick=%+v", res.Tick)
// tick={Book:btc_mxn High:381830.40 Vwap:375965.0682188447 Volume:174.82970689 Last:377297.63 Low:369472.00 Ask:377297.59 Bid:376409.35 Change24:4097.63}
```

