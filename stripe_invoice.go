package main

import (
  "fmt"
  // "github.com/stripe/stripe-go"
  // StripeInvoice "github.com/stripe/stripe-go/invoice"
  // "github.com/stripe/stripe-go/invoiceitem"
  "strings"
  "strconv"
  // "time"
  "flag"
  "stripey"
)

func main() {
  keyArg := flag.String("key","test","Key to use: Live or Test")
  emailArg := flag.String("email","","Email of customer")
  amountArg := flag.String("amount","","Amount of invoice")
  descriptionArg := flag.String("desc","GBMDS After School Activity","Description")
  offerArg := flag.String("offer_id","","Offer Id")
  flag.Parse()

  email := *emailArg

  amount,err1 := strconv.Atoi(*amountArg)
  if err1 != nil {
    fmt.Printf("Error converting amount [%s]. Exiting. \n", *amountArg)
    return
  }
  description := *descriptionArg
  offerid := *offerArg

  if !(strings.Contains(email,"@") || strings.Contains(email,".")) {
    fmt.Printf("Invalid or unspecified *-email=*. Exiting.\n")
    return
  }

  if amount <= 0 {
    fmt.Printf("Zero or negative or unspecified *-amount=*. Exiting.\n")
    return
  }

  stripey.SetKey(*keyArg)

  stripey.CreateInvoice(email, description, int64(amount), offerid)
}
