# stripe

This repo contains go (golang.org) code for connecting to stripe (stripe.com) and invoicing.
The initial use case was to programmatically send a list of invoices to parents to pay for After School Activities.
Invoices were sent out by stripe via email and provided a form for parents to pay by card.
Commands are run at the command line to operate this system. For example:

```
> go run stripe.go ls -startdate=2022-01-10 -key=test -status=draft
```

## Prerequisites
- [go](https://golang.org/doc/install)
- [stripe-go](https://github.com/stripe/stripe-go)
- [godotenv](https://github.com/joho/godotenv)

## Commands

### Listing
- ls: list invoices
- charges: list all charges for customer - charge ids can be used for refunds
- getcustomer: lookup customer and show payment source details

### Invoicing
- invoice: create a single invoice
- offers: create invoices from an offers csv file (see example csv in ./test)
- finalize: finalize invoices to trigger email for payment
- void: void an invoice

### Testing
- `offers-202001119.csv`: example offers file used for invoicing offers

### Utility functions
- `./stripey/stripey.go`: common utility functions for stripe

## Setup Configuration
- .env file containing the test and live keys
- In the go.mod file add the following replacement if the packages are not found:

```
replace github.com/griffithbarracks/utils/stripey v0.0.1 => ./stripey
```

- In the top level directory run the following to source local packages:

```go get github.com/griffithbarracks/utils/stripey@v0.0.1```
