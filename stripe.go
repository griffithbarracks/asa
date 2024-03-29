package main

import (
  "os"
  "fmt"
  "strings"
  "strconv"
  "flag"
  "io/ioutil"
  "encoding/csv"
  "github.com/griffithbarracks/asa/stripey"
)

func main() {

  finalizeCmd := flag.NewFlagSet("finalize", flag.ExitOnError)
  finalize_key := finalizeCmd.String("key","test","Key to use: Live or Test")
  finalize_startdate := finalizeCmd.String("startdate","2022-05-26","Earliest date for charge retrieval")
  finalize_perform := finalizeCmd.String("finalize","false","Finalize if set to true; false will list only")

  sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
  send_key := sendCmd.String("key","test","Key to use: Live or Test")
  send_invoiceid := sendCmd.String("invoice","","invoice id: e.g. 'in_1Klz3u2eZvKYlo2CYU1wdKoW'")

  mailOpenInvoicesCmd := flag.NewFlagSet("mailopeninvoices", flag.ExitOnError)
  mailOpenInvoices_key := mailOpenInvoicesCmd.String("key","test","Key to use: Live or Test")
  mailOpenInvoices_startdate := mailOpenInvoicesCmd.String("startdate","2022-05-26","Earliest date for charge retrieval")
  mailOpenInvoices_send := mailOpenInvoicesCmd.String("send","false","Send email if true, false will list only")

  voidCmd := flag.NewFlagSet("void", flag.ExitOnError)
  void_key := voidCmd.String("key","test","Key to use: Live or Test")
  void_invoiceid := voidCmd.String("invoice","","invoice id: e.g. 'in_1Klz3u2eZvKYlo2CYU1wdKoW'")

  voidOpenInvoicesCmd := flag.NewFlagSet("voidopeninvoices", flag.ExitOnError)
  voidOpenInvoices_key := voidOpenInvoicesCmd.String("key","test","Key to use: Live or Test")
  voidOpenInvoices_perform := voidOpenInvoicesCmd.String("perform","false","Void invoices if true, false will list only")
  voidOpenInvoices_startdate := voidOpenInvoicesCmd.String("startdate","2020-01-01","Earliest date of invoice creation")
  voidOpenInvoices_enddate := voidOpenInvoicesCmd.String("enddate","2022-05-26","Latest date of invoice creation")

  deleteCmd := flag.NewFlagSet("deletedraft", flag.ExitOnError)
  delete_key := deleteCmd.String("key","test","Key to use: Live or Test")
  delete_invoiceid := deleteCmd.String("invoice","","invoice id: e.g. 'in_1Klz3u2eZvKYlo2CYU1wdKoW'")

  delAllDraftsCmd := flag.NewFlagSet("deletealldrafts", flag.ExitOnError)
  delAllDrafts_key := delAllDraftsCmd.String("key","test","Key to use: Live or Test")

  chargesCmd := flag.NewFlagSet("charges", flag.ExitOnError)
  charges_key := chargesCmd.String("key","test","Key to use: Live or Test")
  charges_email := chargesCmd.String("email","","Email of customer")
  charges_startdate := chargesCmd.String("startdate","2022-01-01","Earliest date for charge retrieval")

  offersCmd := flag.NewFlagSet("offers", flag.ExitOnError)
  offers_key := offersCmd.String("key","test","Key to use: Live or Test")
  offers_file := offersCmd.String("file","","Filename of offers to process")

  testPayCmd := flag.NewFlagSet("testpay", flag.ExitOnError)
  testPay_invoice := testPayCmd.String("invoice","","invoice id: e.g. 'in_1Klz3u2eZvKYlo2CYU1wdKoW'")
  testPay_amount := testPayCmd.String("amount","0","Amount to invoice in cent")

  // testcardCmd := flag.NewFlagSet("testcard", flag.ExitOnError)
  // testcard_email := testcardCmd.String("email","","Email of customer")
  // testcard_token := testcardCmd.String("token","tok_1234","Card Token")

  getcustomerCmd := flag.NewFlagSet("getcustomer", flag.ExitOnError)
  getcustomer_key := getcustomerCmd.String("key","test","Key to use: Live or Test")
  getcustomer_email := getcustomerCmd.String("email","","Email of customer")

  updatecustomerCmd := flag.NewFlagSet("updatecustomer", flag.ExitOnError)
  updatecustomer_key := updatecustomerCmd.String("key","test","Key to use: Live or Test")
  updatecustomer_email1 := updatecustomerCmd.String("email1","","Current email of customer")
  updatecustomer_email2 := updatecustomerCmd.String("email2","","New email of customer")

  lsCustomersCmd := flag.NewFlagSet("lscustomers", flag.ExitOnError)
  lsCustomers_key := lsCustomersCmd.String("key","test","Key to use: Live or Test")


  if len(os.Args) < 2 {
      fmt.Println("Missing subcommand: e.g. 'ls', 'finalize', 'offers', 'getcustomer'")
      os.Exit(1)
  }

  // flag.Parse()

  switch os.Args[1] {
  case "ls":
      lsCmd := flag.NewFlagSet("ls", flag.ExitOnError)
      key := lsCmd.String("key","test","Key to use: Live or Test")
      startdate := lsCmd.String("startdate","2020-01-01","Earliest date for invoice retrieval yyyy-mmm-dd")
      status := lsCmd.String("status","","Invoice Status: draft, open, paid, uncollectible, or void")
      lsCmd.Parse(os.Args[2:])
      stripey.SetKey(*key)
      stripey.ListInvoices (startdate, status)

  case "invoice":
      invoiceCmd := flag.NewFlagSet("invoice", flag.ExitOnError)
      invoice_key := invoiceCmd.String("key","test","Key to use: Live or Test")
      invoice_email := invoiceCmd.String("email","","Email to send invoice")
      invoice_amount := invoiceCmd.String("amount","0","Amount to invoice in cent")
      invoice_desc := invoiceCmd.String("desc","GBMDS ASA","Description in invoice")
      invoice_offer := invoiceCmd.String("offerid","","Offer Id for Tracking")
      invoiceCmd.Parse(os.Args[2:])
      stripey.SetKey(*invoice_key)
      Invoice (invoice_email, invoice_amount, invoice_desc, invoice_offer)

  case "finalize":
      finalizeCmd.Parse(os.Args[2:])
      stripey.SetKey(*finalize_key)
      stripey.FinalizeInvoices(finalize_startdate, finalize_perform)

  case "void":
      voidCmd.Parse(os.Args[2:])
      stripey.SetKey(*void_key)
      stripey.Void (void_invoiceid)

  case "voidopeninvoices":
      voidOpenInvoicesCmd.Parse(os.Args[2:])
      stripey.SetKey(*voidOpenInvoices_key)
      stripey.VoidOpenInvoices (voidOpenInvoices_startdate, voidOpenInvoices_enddate, voidOpenInvoices_perform)

  case "send":
      sendCmd.Parse(os.Args[2:])
      stripey.SetKey(*send_key)
      stripey.Send (send_invoiceid)

  case "mailopeninvoices":
      mailOpenInvoicesCmd.Parse(os.Args[2:])
      stripey.SetKey(*mailOpenInvoices_key)
      stripey.MailOpenInvoices(mailOpenInvoices_startdate, mailOpenInvoices_send)

  case "delete":
      deleteCmd.Parse(os.Args[2:])
      stripey.SetKey(*delete_key)
      stripey.Delete (delete_invoiceid)

  case "delalldrafts":
      delAllDraftsCmd.Parse(os.Args[2:])
      stripey.SetKey(*delAllDrafts_key)
      stripey.DeleteAllDrafts()

  case "charges":
      chargesCmd.Parse(os.Args[2:])
      stripey.SetKey(*charges_key)
      stripey.Charges(charges_email, charges_startdate)

  case "offers":
      offersCmd.Parse(os.Args[2:])
      Offers (offers_key, offers_file)

  case "testpay":
      testPayCmd.Parse(os.Args[2:])
      stripey.SetKey("test")
      stripey.TestPayInvoice(testPay_invoice, testPay_amount)

  // case "testcard":
  //     testcardCmd.Parse(os.Args[2:])
  //     stripey.SetKey("test")
      // stripey.CustomerAddCard(testcard_email, testcard_token)

  case "getcustomer":
      getcustomerCmd.Parse(os.Args[2:])
      stripey.SetKey(*getcustomer_key)
      stripey.GetCustomer (*getcustomer_email)

  case "updatecustomer":
      updatecustomerCmd.Parse(os.Args[2:])
      stripey.SetKey(*updatecustomer_key)
      stripey.UpdateCustomerEmail (*updatecustomer_email1, *updatecustomer_email2)

  case "lscustomers":
      lsCustomersCmd.Parse(os.Args[2:])
      stripey.SetKey(*lsCustomers_key)
      stripey.ListCustomers()

  default:
      fmt.Println("Invoicing subcommands: 'ls', 'invoice', 'finalize', 'send', 'void', 'delete', 'delalldrafts'")
      fmt.Println("Other subcommands: 'charges', 'offers', 'testpay', 'testcard', 'getcustomer', 'updatecustomer', 'lscustomers'")
      return
  }
}

func Invoice (emailArg *string, amountArg *string, descArg *string, offerArg *string) {
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

  stripey.CreateInvoice(email, description, int64(amount), offerid)
}


func Offers (keyArg *string, offerFileArg *string) {
  stripey.SetKey(*keyArg)

  // Read File Contents
  content, err := ioutil.ReadFile(*offerFileArg)
  if err!= nil {
    fmt.Println(err)
    return
  }
  lines := string(content)
	r := csv.NewReader(strings.NewReader(lines))
	r.Comma = ','
	r.Comment = '#'
	records, err := r.ReadAll()
	if err != nil {
    fmt.Println(err)
    return
	}
  inputfilesize := len(records)
	fmt.Printf("- %s contains %d lines (including header)\n", *offerFileArg, inputfilesize)

  // Iterate contents and create Invoices
  for i:=0 ; i<inputfilesize; i++ {
		if (records[i][0] == "S#") {
			continue
		}

		email := records[i][0]
    description :=  records[i][1]
    amountstring := strings.Replace(records[i][2], "€", "",-1)
    amounteuro, _ := strconv.ParseFloat(amountstring,32)
    amount := int64(amounteuro * 100)
    offerid := "offer_"+records[i][3]

    stripey.CreateInvoice(email, description, int64(amount), offerid)
	}
}
