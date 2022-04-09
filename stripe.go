package main

import (
  "os"
  "fmt"
  "strings"
  "strconv"
  "flag"
  "io/ioutil"
  "encoding/csv"
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

  offersCmd := flag.NewFlagSet("offers", flag.ExitOnError)
  offers_key := offersCmd.String("key","test","Key to use: Live or Test")
  offers_file := offersCmd.String("file","","Filename of offers to process")

  testPayCmd := flag.NewFlagSet("testpay", flag.ExitOnError)
  testPay_invoice := testPayCmd.String("invoice","","invoice id: e.g. 'in_1Klz3u2eZvKYlo2CYU1wdKoW'")
  testPay_amount := testPayCmd.String("amount","0","Amount to invoice in cent")

  testcardCmd := flag.NewFlagSet("testcard", flag.ExitOnError)
  testcard_email := testcardCmd.String("email","","Email of customer")
  testcard_token := testcardCmd.String("token","tok_1234","Card Token")

  getcustomerCmd := flag.NewFlagSet("getcustomer", flag.ExitOnError)
  getcustomer_key := getcustomerCmd.String("key","test","Key to use: Live or Test")
  getcustomer_email := getcustomerCmd.String("email","","Email of customer")


  flag.Parse()

  if len(os.Args) < 2 {
      fmt.Println("Missing subcommand: e.g. 'ls', 'finalize', 'offers', 'getcustomer'")
      os.Exit(1)
  }

  switch os.Args[1] {
  case "ls":
      lsCmd.Parse(os.Args[2:])
      stripey.SetKey(*key)
      stripey.ListInvoices (startdate, status)

  case "invoice":
      invoiceCmd.Parse(os.Args[2:])
      Invoice (invoice_key, invoice_email, invoice_amount, invoice_desc, invoice_offer)

  case "finalize":
      finalizeCmd.Parse(os.Args[2:])
      stripey.SetKey(*finalize_key)
      stripey.FinalizeInvoices(finalize_startdate, finalize_perform)

  case "void":
      voidCmd.Parse(os.Args[2:])
      stripey.SetKey(*void_key)
      stripey.Void (void_invoiceid)

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

  case "testcard":
      testcardCmd.Parse(os.Args[2:])
      stripey.SetKey("test")
      stripey.CustomerAddCard(testcard_email, testcard_token)

  case "getcustomer":
      getcustomerCmd.Parse(os.Args[2:])
      stripey.SetKey(*getcustomer_key)
      stripey.GetCustomer (*getcustomer_email)

  default:
      fmt.Println("Expected subcommands: 'ls', 'invoice', 'finalize', 'void', 'charges', 'offers', 'testpay', 'testcard', 'getcustomer'")
      return
  }
}

func Invoice (keyArg *string, emailArg *string, amountArg *string, descArg *string, offerArg *string) {
  stripey.SetKey(*keyArg)
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

    offerid := "offer_"+records[i][0]
		childname := records[i][2]
		email := records[i][5]
		asa := records[i][9]
		amounteuro, _ := strconv.ParseFloat(records[i][10],32)
		amount := int64(amounteuro * 100)
    donationstring := records[i][11]
		donation, _ := strconv.ParseFloat(donationstring,32)

		description := asa + " for " + childname
		if (donation > 0) {
			description = description + ", incl. donation of €" + donationstring;
		}

    stripey.CreateInvoice(email, description, int64(amount), offerid)
	}
}
