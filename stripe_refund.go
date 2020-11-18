package main

import (
  "flag"
  "fmt"
  "github.com/stripe/stripe-go"
  "github.com/stripe/stripe-go/refund"
  "strconv"
  "strings"
  "stripey"
)

func main() {

  // Arg handling
  keyArg := flag.String("key","test","Key to use: Live or Test")
  chargeArg := flag.String("charge","","Charge Id")
  amtArg := flag.String("amount","","Amount in cents")
  flag.Parse()
  stripey.SetKey(*keyArg)

  amount, _ := strconv.ParseInt(*amtArg, 0, 64)
  if strings.Compare(*chargeArg,"") == 0 {
    fmt.Printf("No charge id specified. Exiting.\n")
    return
  }

  if amount <= 0 {
    fmt.Printf("Zero or negative or unspecified *amount*. Exiting.\n")
    return
  }

  if amount > 2000 {
    fmt.Printf("€20 is the maximum *amount*. Exiting.\n")
    return
  }

  params := &stripe.RefundParams{
    Charge: chargeArg,
    Amount: stripe.Int64(amount),
  }
  r, _ := refund.New(params)

  fmt.Printf("Refund %s, amt=%d, status=%s\n", r.ID, r.Amount, r.Status)
}
