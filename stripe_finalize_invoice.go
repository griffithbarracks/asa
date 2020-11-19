package main

import (
    // "encoding/json"
    "fmt"
    "flag"
    "github.com/joho/godotenv"
    "github.com/stripe/stripe-go"
    // "github.com/stripe/stripe-go/customer"
    "github.com/stripe/stripe-go/invoice"
    // "github.com/stripe/stripe-go/invoiceitem"
    "log"
    "os"
)

func setKey(key string) {
  stripe.LogLevel = int(1)

  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  livekey := os.Getenv("livekey")
  testkey := os.Getenv("testkey")
  stripe.Key = testkey
  if key == "live" {
    stripe.Key = livekey
  }
}

func main() {

  // Arg handling
  keyArg := flag.String("key","test","Key to use: Live or Test")
  invoiceArg := flag.String("invoice","","Invoice Id")

  flag.Parse()
  // fmt.Printf("### FLAGS ### key=%s, customerid=%s, token=%s\n", *keyArg, *customerArg, *tokenArg)
  setKey(*keyArg)

  i_invoiceid := *invoiceArg

  params := &stripe.InvoiceFinalizeParams{
    AutoAdvance: stripe.Bool(true),
  }

  final_invoice, invoiceerr := invoice.FinalizeInvoice(i_invoiceid, params)


  if invoiceerr != nil {
    fmt.Printf("Error finalizing invoice for [%i] %s\n", i_invoiceid, invoiceerr)
    log.Fatal(invoiceerr)
  }

  fmt.Printf("Finalized Invoice: %s, %s, \"%s\", %d, %s\n",
    final_invoice.ID,
    final_invoice.CustomerEmail,
    final_invoice.Description,
    final_invoice.AmountDue,
    final_invoice.Status,
  )
}
