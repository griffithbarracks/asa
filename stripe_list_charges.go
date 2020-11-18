package main

import (
    "flag"
    "github.com/stripe/stripe-go"
    // "github.com/stripe/stripe-go/customer"
    "github.com/stripe/stripe-go/charge"
    "fmt"
    "time"
    "strconv"
    "strings"
    "stripey"
)

func main() {
  keyArg := flag.String("key","test","Key to use: Live or Test")
  startdateArg := flag.String("startdate","2020-01-01","Earliest date for item retrieval yyyy-mmm-dd")
  emailArg := flag.String("email","","Email")
  flag.Parse()

  stripey.SetKey(*keyArg)

  const shortFormDate = "2006-01-02"
  startdate, _ := time.Parse(shortFormDate,*startdateArg)
  createdTimeUnix := strconv.FormatInt(startdate.Unix(),10)

  email := *emailArg
  customerid := stripey.GetCustomerId(email)

  if strings.Compare(email,"") == 0 ||
    (strings.Contains(email,"@") && strings.Contains(customerid,"cus_")) {
    chargeparams := &stripe.ChargeListParams{}
    chargeparams.Filters.AddFilter("limit", "", "20")
    if (strings.Contains(email,"@") && strings.Contains(customerid,"cus_")) {
      chargeparams.Filters.AddFilter("customer", "", customerid)
    }
    chargeparams.Filters.AddFilter("created", "gt", createdTimeUnix)
    j := charge.List(chargeparams)

    for j.Next() {
      c := j.Charge()
      fmt.Printf("%s, %s, %s, Last4: %s, Amt: €%3.2f, Refund: €%3.2f, Stat: %s\n",
        c.Customer.ID,
        c.ReceiptEmail,
        c.ID,
        c.PaymentMethodDetails.Card.Last4,
        float64(c.Amount)/100.0,
        float64(c.AmountRefunded)/100.0,
        c.Status,
      )
    }
  }
}
