package stripey

import (
  "github.com/joho/godotenv"
  "log"
  "os"
  "fmt"
  "strings"
  "time"

  "github.com/stripe/stripe-go"
  "github.com/stripe/stripe-go/customer"
  StripeInvoice "github.com/stripe/stripe-go/invoice"
  "github.com/stripe/stripe-go/invoiceitem"
)

func SetKey(key string) {
  stripe.LogLevel = int(1)

  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  livekey := os.Getenv("livekey")
  testkey := os.Getenv("testkey")
  stripe.Key = testkey
  if key == "live" {
    stripe.Key = livekey
  }
}

func GetCustomerId(email string) string {

  if strings.Compare(email,"") == 0 {
    return *stripe.String("")
  }

	clparams := &stripe.CustomerListParams{}
	clparams.Filters.AddFilter("limit", "", "5")
	clparams.Filters.AddFilter("email", "", email)
	i := customer.List(clparams)

	found := 0
	customerid := stripe.String("")

	for i.Next() {
    found = found + 1
		customer := i.Customer()
		customerid = stripe.String(customer.ID)
    if found > 1 {
      fmt.Printf("Additional found customer: %s, %s\n", customer.Email, customer.ID)
    }
	}

	if found == 0 {
		fmt.Printf("Customer [%s] not found.\n", email)
	}
	return *customerid
}

func CreateInvoice (email string, description string, amount int64, offerid string) string {
  customerid := GetCustomerId(email)

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
    AutoAdvance: stripe.Bool(true),
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
    i.AmountDue,
    i.Status,
  )

  return i.ID

}
