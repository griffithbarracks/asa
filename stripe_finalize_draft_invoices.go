package main

import (
    "fmt"
    "flag"
    "github.com/stripe/stripe-go"
    "github.com/stripe/stripe-go/invoice"
    "log"
    "strconv"
    "strings"
    "time"
    "stripey"
)

func main() {
  // Arg handling
  keyArg := flag.String("key","test","Key to use: Live or Test")
  startdateArg := flag.String("startdate","2020-01-02","Earliest date for invoice retrieval yyyy-mmm-dd")
  finalizeArg := flag.String("finalize","false","Finalize if set to true; false will list only")
  flag.Parse()

  stripey.SetKey(*keyArg)

  const shortFormDate = "2006-01-02"
  startdate, _ := time.Parse(shortFormDate,*startdateArg)
  // fmt.Println("Parsed Date:",startdate,*startdateArg)

  finalize, _ := strconv.ParseBool(*finalizeArg)


  createdTimeUnix := strconv.FormatInt(startdate.Unix(),10)

  listparams := &stripe.InvoiceListParams{}
  listparams.Filters.AddFilter("limit", "", "100")
  listparams.Filters.AddFilter("created", "gt", createdTimeUnix)
  listparams.Filters.AddFilter("status", "", "draft")

  invoiceList := invoice.List(listparams)
  fmt.Printf("Invoices created after %s. Report date: %s,,,,,\n",
    startdate.Format("2006-01-02"), time.Now().Format("2006-01-02 15:04"))

  count := 0

  for invoiceList.Next() {
    i := invoiceList.Invoice()

    count = count + 1
    createdDate := time.Unix(i.Created,0).Format("2006-01-02 15:04")
    // Removing comma in description for csv format
    description := strings.Replace(i.Lines.Data[0].Description, ",", " -",-1)

    fmt.Printf("%d, %s, %s, %s, %s, %s, %d, %s, %s",
      count,
      i.ID,
      i.CustomerEmail,
      i.Customer.ID,
      createdDate,
      description,
      i.AmountDue,
      i.Status,
      i.Metadata["offer_id"],
    )

    if (finalize) {
      finalizeInvoiceParams := &stripe.InvoiceFinalizeParams{
        AutoAdvance: stripe.Bool(true),
      }
      final_invoice, invoiceerr := invoice.FinalizeInvoice(i.ID, finalizeInvoiceParams)

      if invoiceerr != nil {
        fmt.Printf("Error finalizing invoice for [%i] %s\n", i.ID, invoiceerr)
        log.Fatal(invoiceerr)
      }

      fmt.Printf (
        " >>> %s\n",
        final_invoice.Status,
      )
    } else {
      fmt.Printf ("\n")
    }
  }
}
