package main

import (
  "fmt"
  "github.com/stripe/stripe-go"
  StripeInvoice "github.com/stripe/stripe-go/invoice"
  "github.com/stripe/stripe-go/invoiceitem"
  "strings"
  "strconv"
  "time"
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
    fmt.Printf("Invalid or specified *email*. Exiting.\n")
    return
  }

  if amount <= 0 {
    fmt.Printf("Zero or negative or unspecified *amount*. Exiting.\n")
    return
  }

  stripey.SetKey(*keyArg)

  customerid := stripey.GetCustomerId(email)

  ii_params := &stripe.InvoiceItemParams{
    Customer: stripe.String(customerid),
    Amount: stripe.Int64(int64(amount)),
    Currency: stripe.String(string(stripe.CurrencyEUR)),
    Description: stripe.String(description),
  }
  ii, ii_err := invoiceitem.New(ii_params)
  if ii_err != nil {
    fmt.Printf("Error creating Invoice Item: %s %s\n", description, ii_err)
    return
  }
  // fmt.Printf("Created Invoice Line Item: %s, %d\n", ii.Description, ii.Amount)

  params := &stripe.InvoiceParams{
    Customer: stripe.String(customerid),
    CollectionMethod: stripe.String("send_invoice"),
    DaysUntilDue: stripe.Int64(1),
    Description: stripe.String(description),
		AutoAdvance: stripe.Bool(true),
  }
	params.AddMetadata("offer_id", offerid)

  i, invoiceerr := StripeInvoice.New(params)

  if invoiceerr != nil {
    fmt.Printf("Error creating invoice for [%s] %s\n", email, invoiceerr)
    return
  }

  createdDate := time.Unix(i.Created,0).Format("2006-01-02 15:04")
  description1 := strings.Replace(i.Lines.Data[0].Description, ",", " -",-1)

  fmt.Printf("%s, %s, %s, %s, %s, offer_%s, %d, %s\n",
    i.ID,
    i.CustomerEmail,
    i.Customer.ID,
    createdDate,
    description1,
    i.Metadata["offer_id"],
    i.AmountDue,
    i.Status,
  )
}
