package stripey

import (
  "fmt"
  "strings"
  "strconv"
  "time"
  "encoding/json"

  "github.com/stripe/stripe-go"
  "github.com/stripe/stripe-go/invoice"
  StripeInvoice "github.com/stripe/stripe-go/invoice"
  "github.com/stripe/stripe-go/invoiceitem"
)

func ListInvoices (startdateArg *string, statusArg *string) {
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

    count += 1
    fmt.Printf("%d, %s, %s, %s, %s, %s, %s, %d, %s, %s\n",
      count,
      i.ID,
      i.CustomerEmail,
      i.Customer.ID,
      createdDate,
      description,
      i.Metadata["offer_id"],
      i.AmountRemaining,
      i.Status,
      paidDate,
    )
  }
}

func CreateInvoice (email string, description string, amount int64, offerid string) string {

  customerid := GetCustomerId(email)
  if len(customerid)<1 {
    CreateCustomer(email)
    customerid = GetCustomerId(email)
  }

  description_clean := strings.Replace(description, ",", " -",-1)

  ii_params := &stripe.InvoiceItemParams{
    Customer: stripe.String(customerid),
    Amount: stripe.Int64(amount),
    Currency: stripe.String(string(stripe.CurrencyEUR)),
    Description: stripe.String(description_clean),
  }
  _, ii_err := invoiceitem.New(ii_params)
  if ii_err != nil {
    fmt.Printf("Error creating Invoice Item: %s %s\n", description, ii_err)
    return "err"
  }
  // fmt.Printf("Created Invoice Line Item: %s, %d\n", ii.Description, ii.Amount)

  params := &stripe.InvoiceParams{
    Customer: stripe.String(customerid),
    CollectionMethod: stripe.String("send_invoice"),
    DaysUntilDue: stripe.Int64(1),
    Description: stripe.String(description_clean),
    AutoAdvance: stripe.Bool(false),
  }
  params.AddMetadata("offer_id", offerid)

  i, invoiceerr := StripeInvoice.New(params)

  if invoiceerr != nil {
    fmt.Printf("Error creating invoice for [%s] %s\n", email, invoiceerr)
    return "err"
  }

  createdDate := time.Unix(i.Created,0).Format("2006-01-02 15:04")

  fmt.Printf("%s, %s, %s, %s, %s, %s, %d, %s\n",
    i.ID,
    i.CustomerEmail,
    i.Customer.ID,
    createdDate,
    i.Description,
    i.Metadata["offer_id"],
    i.AmountRemaining,
    i.Status,
  )

  return i.ID
}

func Send (invoiceArg *string) {
  invoice_id := *invoiceArg
  if len(invoice_id) < 1 {
    fmt.Printf("Invalid invoice id '%s'\n", *invoiceArg)
    return
  }
  in, _ := invoice.SendInvoice(
    invoice_id,
    nil,
  )
  fmt.Printf("Sent invoice '%s'\n", in.ID)
}

func Void (invoiceArg *string) {
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

func Delete (invoiceArg *string) {
  invoice_id := *invoiceArg
  if len(invoice_id) < 1 {
    fmt.Printf("Invalid invoice id '%s'\n", *invoiceArg)
    return
  }
  in, _ := invoice.Del(
    invoice_id,
    nil,
  )
  fmt.Printf("Deleted invoice '%s'\n", in.ID)
}

func DeleteAllDrafts () {

  listparams := &stripe.InvoiceListParams{}
  listparams.Filters.AddFilter("limit", "", "100")
  listparams.Filters.AddFilter("status", "", "draft")

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

    count += 1
    fmt.Printf("%d, %s, %s, %s, %s, %s, %s, %d, %s, %s\n",
      count,
      i.ID,
      i.CustomerEmail,
      i.Customer.ID,
      createdDate,
      description,
      i.Metadata["offer_id"],
      i.AmountRemaining,
      i.Status,
      paidDate,
    )

    in, _ := invoice.Del(
      i.ID,
      nil,
    )
    fmt.Printf("Deleted invoice '%s'\n", in.ID)
  }

}

func FinalizeInvoices (startdateArg *string, finalizeArg *string) {
  const shortFormDate = "2006-01-02"
  startdate, _ := time.Parse(shortFormDate,*startdateArg)
  finalize, _ := strconv.ParseBool(*finalizeArg)
  createdTimeUnix := strconv.FormatInt(startdate.Unix(),10)
  fmt.Printf("Invoices created after %s (Report date: %s)\n",
    startdate.Format("2006-01-02"), time.Now().Format("2006-01-02 15:04"))

  listparams := &stripe.InvoiceListParams{}
  listparams.Filters.AddFilter("limit", "", "100")
  listparams.Filters.AddFilter("created", "gt", createdTimeUnix)
  listparams.Filters.AddFilter("status", "", "draft")
  invoiceList := invoice.List(listparams)

  count := 0
  for invoiceList.Next() {
    i := invoiceList.Invoice()
    count = count + 1
    createdDate := time.Unix(i.Created,0).Format("2006-01-02 15:04")
    // Removing comma from description text - obsolete?
    description := strings.Replace(i.Lines.Data[0].Description, ",", " -",-1)

    fmt.Printf("%d, %s, %s, %s, %s, %s, %d, %s, %s",
      count,
      i.ID,
      i.CustomerEmail,
      i.Customer.ID,
      createdDate,
      description,
      i.AmountRemaining,
      i.Status,
      i.Metadata["offer_id"],
    )

    if (finalize) {
      finalizeInvoiceParams := &stripe.InvoiceFinalizeParams{
        AutoAdvance: stripe.Bool(false),
      }
      final_invoice, invoiceerr := invoice.FinalizeInvoice(i.ID, finalizeInvoiceParams)
      if invoiceerr != nil {
        fmt.Printf("Error finalizing invoice for [%i] %s\n", i.ID, invoiceerr)
      }
      fmt.Printf (" >>> %s\n", final_invoice.Status)

      invoice.SendInvoice(
        i.ID,
        nil,
      )

    } else {
      fmt.Printf ("\n")
    }
  }
}

func TestPayInvoice (invoiceArg *string, amountArg *string) {
  SetKey("test")

  if strings.Compare(*invoiceArg,"") == 0 {
    fmt.Printf("No *invoice* specified. Exiting.\n")
    return
  }

  amount,err1 := strconv.Atoi(*amountArg)
  if err1 != nil {
    fmt.Printf("Error converting *amount* [%s]. Exiting. \n", *amountArg)
    return
  }
  if amount <= 0 {
    fmt.Printf("Zero or negative or unspecified *amount*. Exiting.\n")
    return
  }

  i_invoiceid := *invoiceArg

  final_invoice, invoice_error := invoice.Pay(i_invoiceid, nil)
  if invoice_error != nil {
    var errormap map[string]interface{}
    error_json, _ := json.Marshal(invoice_error)
    json.Unmarshal([]byte(error_json), &errormap)
    for key, value := range errormap {
      switch key {
      case "code":
        fallthrough
      case "message":
        fallthrough
      case "type":
        fmt.Printf("    %s: %v\n", key, value)
      }
    }
    return
  }

  fmt.Printf("Paid Invoice: %s, %s, \"%s\", %d, %s\n",
    final_invoice.ID,
    final_invoice.CustomerEmail,
    final_invoice.Description,
    final_invoice.AmountPaid,
    final_invoice.Status,
  )
}
