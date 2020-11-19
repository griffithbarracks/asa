# stripe

This repo contains go (golang.org) code for connecting to stripe (stripe.com) and invoicing.
The initial use case was to programmatically send a list of invoices to parents to pay for After School Activities.
Invoices were sent out by stripe via email and provided a form for parents to pay by card.
Commands are run at the command line to operate this system. For example:

```
> go run stripe_list_invoices.go
```

## Prerequisites
- [go](https://golang.org/doc/install)
- [stripe-go](https://github.com/stripe/stripe-go)
- [godotenv](https://github.com/joho/godotenv)

## Files

### Invoicing
- stripe_invoice.go
  - _create a single invoice_
- stripe_invoice_offers.go
  - _create invoices from an offers csv file (see example csv)_
- stripe_finalize_draft_invoices.go
  - _finalize invoices to trigger email for payment_

### Listing
  - stripe_list_charges.go
    - _list all charges - charge ids can be used for refunds_
  - stripe_list_invoices.go
    - _list of invoices and their details including paid status_

### Refunds
- stripe_refund.go
  - _refund a previous charge (payment taken by card)_

### Testing
- stripe_test_addcard.go
  - _add a tokenized card for testing payments and refunds_
- stripe_test_invoice_pay.go
  - _test paying an invoice with a card previously added_
- offers-202001119.csv
  - _example offers file used for invoicing offers_

### Utility functions
- src/stripey/stripey.go
  - _common utility functions for stripe_

## Setup Configuration
- .env file containing the test and live keys
- In the top level directory run the following to source local packages:
```
> export GOPATH=$(go env GOPATH):`pwd`
```
