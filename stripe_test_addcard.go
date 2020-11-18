package main

import (
    "flag"
    "fmt"
    "github.com/stripe/stripe-go"
    "github.com/stripe/stripe-go/card"
    "strings"
    "stripey"
)

func main() {

  // Arg handling
  customerArg := flag.String("cus","","Customer Id")
  tokenArg := flag.String("token","tok_visa","Card Token")
  flag.Parse()

  if strings.Compare(*customerArg,"") == 0 {
    fmt.Printf("No customer id specified. Exiting.\n")
    return
  }

  stripey.SetKey("test")

  params := &stripe.CardParams{
    Customer: stripe.String(*customerArg),
    Token: stripe.String(*tokenArg),
  }
  c, _ := card.New(params)

  fmt.Printf("Card = %s\n", c.ID)
}
