package stripey
import (
  "fmt"
  "strings"
  "strconv"
  "time"
  "github.com/stripe/stripe-go"
  "github.com/stripe/stripe-go/charge"
  "github.com/stripe/stripe-go/refund"
)

func Charges (emailArg *string, startdateArg *string) {

  const shortFormDate = "2006-01-02"
  startdate, _ := time.Parse(shortFormDate,*startdateArg)
  createdTimeUnix := strconv.FormatInt(startdate.Unix(),10)

  email := *emailArg
  customerid := GetCustomerId(email)

  if strings.Compare(email,"") == 0 ||
    (strings.Contains(email,"@") && strings.Contains(customerid,"cus_")) {
    chargeparams := &stripe.ChargeListParams{}
    chargeparams.Filters.AddFilter("limit", "", "20")
    if (strings.Contains(email,"@") && strings.Contains(customerid,"cus_")) {
      chargeparams.Filters.AddFilter("customer", "", customerid)
    }
    chargeparams.Filters.AddFilter("created", "gt", createdTimeUnix)

    j := charge.List(chargeparams)

    header := true
    count := 0

    for j.Next() {
      c := j.Charge()
      if (header) {
        fmt.Println("CustID, Email, ID, Last4, Amt, Refunded, Status")
        header = false
      }
      count += 1
      fmt.Printf("%s, %s, %s, %s, €%3.2f, €%3.2f, %s\n",
        c.Customer.ID,
        c.ReceiptEmail,
        c.ID,
        c.PaymentMethodDetails.Card.Last4,
        float64(c.Amount)/100.0,
        float64(c.AmountRefunded)/100.0,
        c.Status,
      )
    }
    fmt.Printf("-- %d charge(s) found\n", count)
  }
}

func Refund (chargeArg *string, amtArg *string) {
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
