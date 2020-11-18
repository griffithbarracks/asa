package main

import (
    // "encoding/json"
    "fmt"
    "flag"
    // "github.com/stripe/stripe-go"
    // "github.com/stripe/stripe-go/customer"
    "github.com/stripe/stripe-go/invoice"
    "strings"
    "stripey"
    "log"
)

func main() {
  invoiceArg := flag.String("invoice","","Invoice Id")
  // amountArg := flag.String("amount","","Amount of invoice")
  flag.Parse()

  if strings.Compare(*invoiceArg,"") == 0 {
    fmt.Printf("No *invoice_id* specified. Exiting.\n")
    return
  }
  // amount,err1 := strconv.Atoi(*amountArg)
  // if err1 != nil {
  //   fmt.Printf("Error converting *amount* [%s]. Exiting. \n", *amountArg)
  //   return
  // }
  // if amount <= 0 {
  //   fmt.Printf("Zero or negative or unspecified *amount*. Exiting.\n")
  //   return
  // }

  stripey.SetKey("test")

  i_invoiceid := *invoiceArg

  final_invoice, invoiceerr := invoice.Pay(i_invoiceid, nil)

  if invoiceerr != nil {
    fmt.Printf("Error Paying invoice for [%i] %s\n", i_invoiceid, invoiceerr)
    log.Fatal(invoiceerr)
  }

  fmt.Printf("Paid Invoice: %s, %s, \"%s\", %d, %s\n",
    final_invoice.ID,
    final_invoice.CustomerEmail,
    final_invoice.Description,
    final_invoice.AmountPaid,
    final_invoice.Status,
  )
}
