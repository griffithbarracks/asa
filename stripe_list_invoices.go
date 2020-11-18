package main

import (
  "flag"
  "fmt"
  "time"
  "github.com/stripe/stripe-go"
  "github.com/stripe/stripe-go/invoice"
  "strconv"
  "strings"
  "stripey"
)

func main() {
  keyArg := flag.String("key","test","Key to use: Live or Test")
  startdateArg := flag.String("startdate","2020-01-01","Earliest date for invoice retrieval yyyy-mmm-dd")
  statusArg := flag.String("status","","Invoice Status")
  flag.Parse()
  stripey.SetKey(*keyArg)

  const shortFormDate = "2006-01-02"
  startdate, _ := time.Parse(shortFormDate,*startdateArg)

  createdTimeUnix := strconv.FormatInt(startdate.Unix(),10)

  listparams := &stripe.InvoiceListParams{}
  listparams.Filters.AddFilter("limit", "", "100")
  listparams.Filters.AddFilter("created", "gt", createdTimeUnix)
  if strings.Compare(*statusArg,"") != 0 {
    listparams.Filters.AddFilter("status", "", *statusArg)
  }

  invoiceList := invoice.List(listparams)

  count := 0
  fmt.Printf("#, invoice_id, customer_email, customer_id, date_created, description, amount, status, asa_offer_id\n")

  for invoiceList.Next() {
    i := invoiceList.Invoice()
    count = count + 1
    createdDate := time.Unix(i.Created,0).Format("2006-01-02 15:04")
    description := strings.Replace(i.Lines.Data[0].Description, ",", " -",-1)

    fmt.Printf("%d, %s, %s, %s, %s, %s, offer_%s, %d, %s\n",
      count,
      i.ID,
      i.CustomerEmail,
      i.Customer.ID,
      createdDate,
      description,
      i.Metadata["offer_id"],
      i.AmountDue,
      i.Status,
    )
  }
}
