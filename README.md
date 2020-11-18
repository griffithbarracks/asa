# stripe

This repo contains go (golang.org) code for connecting to stripe (stripe.com) and invoicing.
The initial use case was to programmatically send a list of invoices to parents to pay for After School Activities.
Invoices were sent out by stripe via email and provided a form for parents to pay by card.
Commands are run at the command line to operate this system. For example:

```
go run stripe_list_invoices.go
```

## Prerequisites
- [go](https://golang.org/doc/install)
- [stripe-go](https://github.com/stripe/stripe-go)
- [godotenv](https://github.com/joho/godotenv)

## Files

### Invoicing
- stripe_invoice.go
- stripe_finalize_draft_invoices.go

### Refunds
- stripe_refund.go

### Listing
- stripe_list_charges.go
= stripe_list_invoices.go

### Testing
- stripe_test_addcard.go
- stripe_test_invoice_pay.go

### Utility functions
- src/stripey/stripey.go

## Setup Configuration
- .env file containing the test and live keys
- In the top level directory run the following to source local packages: 
```
export GOPATH=$(go env GOPATH):`pwd`
```
