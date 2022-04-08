package main

import (
  "os"
  "fmt"
  "strings"
  "strconv"
  "flag"
  "time"
  "github.com/stripe/stripe-go"
  "github.com/stripe/stripe-go/invoice"
  "github.com/stripe/stripe-go/charge"
  "github.com/griffithbarracks/utils/stripey"
)

func main() {
  lsCmd := flag.NewFlagSet("ls", flag.ExitOnError)
  key := lsCmd.String("key","test","Key to use: Live or Test")
  startdate := lsCmd.String("startdate","2020-01-01","Earliest date for invoice retrieval yyyy-mmm-dd")
  status := lsCmd.String("status","","Invoice Status: draft, open, paid, uncollectible, or void")

  invoiceCmd := flag.NewFlagSet("invoice", flag.ExitOnError)
  invoice_key := invoiceCmd.String("key","test","Key to use: Live or Test")
  invoice_email := invoiceCmd.String("email","","Email to send invoice")
  invoice_amount := invoiceCmd.String("amount","0","Amount to invoice in cent")
  invoice_desc := invoiceCmd.String("desc","GBMDS ASA","Description in invoice")
  invoice_offer := invoiceCmd.String("offerid","","Offer Id for Tracking")

  finalizeCmd := flag.NewFlagSet("finalize", flag.ExitOnError)
  finalize_key := finalizeCmd.String("key","test","Key to use: Live or Test")
  finalize_startdate := finalizeCmd.String("startdate","2020-01-02","Earliest date for invoice retrieval yyyy-mmm-dd")
  finalize_perform := finalizeCmd.String("finalize","false","Finalize if set to true; false will list only")

  voidCmd := flag.NewFlagSet("void", flag.ExitOnError)
  void_key := voidCmd.String("key","test","Key to use: Live or Test")
  void_invoiceid := voidCmd.String("invoice","","invoice id: e.g. 'in_1Klz3u2eZvKYlo2CYU1wdKoW'")

  chargesCmd := flag.NewFlagSet("charges", flag.ExitOnError)
  charges_key := chargesCmd.String("key","test","Key to use: Live or Test")
  charges_email := chargesCmd.String("email","","Email of customer")
  charges_startdate := chargesCmd.String("startdate","2022-01-01","Earliest date for charge retrieval")

  flag.Parse()

  if len(os.Args) < 2 {
      fmt.Println("Expected subcommand: 'ls', 'invoice', 'finalize'")
      os.Exit(1)
  }

  switch os.Args[1] {
  case "ls":
      lsCmd.Parse(os.Args[2:])
      // fmt.Println("subcommand 'ls'")
      // fmt.Println("  key:", *key)
      // fmt.Println("  startdate:", *startdate)
      // fmt.Println("  status:", *status)
      ListInvoices (key, startdate, status)

  case "invoice":
      invoiceCmd.Parse(os.Args[2:])
      // fmt.Println("subcommand 'invoice'")
      // fmt.Println("  key:", *invoice_key)
      // fmt.Println("  email:", *invoice_email)
      // fmt.Println("  amount:", *invoice_amount)
      // fmt.Println("  description:", *invoice_desc)
      Invoice (invoice_key, invoice_email, invoice_amount, invoice_desc, invoice_offer)

  case "finalize":
      finalizeCmd.Parse(os.Args[2:])
      // fmt.Println("subcommand 'finalize'")
      // fmt.Println("  key:", *finalize_key)
      // fmt.Println("  startdate:", *finalize_startdate)
      // fmt.Println("  finalize:", *finalize_perform)
      FinalizeInvoices (finalize_key, finalize_startdate, finalize_perform)

  case "void":
      voidCmd.Parse(os.Args[2:])
      // fmt.Println("subcommand 'void'")
      // fmt.Println("  key:", *void_key)
      // fmt.Println("  invoice:", *void_invoiceid)
      Void (void_key, void_invoiceid)

  case "charges":
      chargesCmd.Parse(os.Args[2:])
      // fmt.Println("subcommand 'charges'")
      // fmt.Println("  key:", *charges_key)
      // fmt.Println("  email:", *charges_email)
      // fmt.Println("  startdate:", *charges_startdate)
      Charges (charges_key, charges_email, charges_startdate)

  default:
      // fmt.Printf("subcommand '%s'\n", os.Args[1])
      fmt.Println("Expected subcommand: 'ls', 'invoice', 'finalize', 'void'")
      return
  }
}

func ListInvoices (keyArg *string, startdateArg *string, statusArg *string) {
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
  fmt.Printf("#, invoice_id, customer_email, customer_id, date_created, description, asa_offer_id, amount, status, date_paid\n")

  for invoiceList.Next() {
    i := invoiceList.Invoice()
    createdDate := time.Unix(i.Created,0).Format("2006-01-02 15:04")
    paidDate := ""
    if (i.StatusTransitions.PaidAt>0) {
      paidDate = time.Unix(i.StatusTransitions.PaidAt,0).Format("2006-01-02 15:04")
    }
    description := strings.Replace(i.Lines.Data[0].Description, ",", " -",-1)

    count = count + 1
    fmt.Printf("%d, %s, %s, %s, %s, %s, %s, %d, %s, %s\n",
      count,
      i.ID,
      i.CustomerEmail,
      i.Customer.ID,
      createdDate,
      description,
      i.Metadata["offer_id"],
      i.AmountDue,
      i.Status,
      paidDate,
    )
  }
}

func Invoice (keyArg *string, emailArg *string, amountArg *string,descArg *string, offerArg *string) {
  email := *emailArg

  amount,err1 := strconv.Atoi(*amountArg)
  if err1 != nil {
    fmt.Printf("Error converting amount [%s]. Exiting. \n", *amountArg)
    return
  }
  description := *descArg
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

func FinalizeInvoices (keyArg *string, startdateArg *string, finalizeArg *string) {

  stripey.SetKey(*keyArg)

  const shortFormDate = "2006-01-02"
  startdate, _ := time.Parse(shortFormDate,*startdateArg)
  // fmt.Println("Parsed Date:",startdate,*startdateArg)

  finalize, _ := strconv.ParseBool(*finalizeArg)

  createdTimeUnix := strconv.FormatInt(startdate.Unix(),10)

  fmt.Printf("Invoices created after %s (Report date: %s)\n",
    startdate.Format("2006-01-02"), time.Now().Format("2006-01-02 15:04"))

  count := 0

  listparams := &stripe.InvoiceListParams{}
  listparams.Filters.AddFilter("limit", "", "100")
  listparams.Filters.AddFilter("created", "gt", createdTimeUnix)
  listparams.Filters.AddFilter("status", "", "draft")

  invoiceList := invoice.List(listparams)
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
      }

      fmt.Printf (" >>> %s\n", final_invoice.Status)
    } else {
      fmt.Printf ("\n")
    }
  }
}

func Void (keyArg *string, invoiceArg *string) {
  stripey.SetKey(*keyArg)
  invoice_id := *invoiceArg
  if len(invoice_id) < 1 {
    fmt.Printf("Invalid invoice id '%s'\n", *invoiceArg)
    return
  }
  in, _ := invoice.VoidInvoice(
    invoice_id,
    nil,
  )
  fmt.Printf("Voided invoice '%s'\n", in.ID)
}

func Charges (keyArg *string, emailArg *string, startdateArg *string) {
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

    fmt.Println("CustID, Email, ID, Last4, Amt, Refunded, Status")
    for j.Next() {
      c := j.Charge()
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

    fmt.Println("")

  }
}
