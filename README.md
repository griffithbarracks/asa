# stripe

This repo contains go (golang.org) code for connecting to stripe (stripe.com) and invoicing.
The initial use case was to programmatically send a list of invoices to parents to pay for After School Activities.
Invoices were sent out by stripe via email and provided a form for parents to pay by card.

## Prerequisites
- [go](https://golang.org/doc/install)
- [stripe-go](https://github.com/stripe/stripe-go)
- [godotenv](https://github.com/joho/godotenv)
- .env file containing the test and live keys

### Invoicing
- stripe_invoice.go
- stripe_finalize_draft_invoices.go

### Refunds
- stripe_refund.go

### Listing
- stripe_list_charges.go
= stripe_list_invoices.go

### Testing
= stripe_test_addcard.go
= stripe_test_invoice_pay.go
