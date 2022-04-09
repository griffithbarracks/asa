package stripey

import (
  "fmt"
  "strings"
  "github.com/stripe/stripe-go"
  "github.com/stripe/stripe-go/customer"
  // "github.com/stripe/stripe-go/card"
  // "github.com/stripe/stripe-go/source"
)


func GetCustomerId(email string) string {

  if strings.Compare(email,"") == 0 {
    return *stripe.String("")
  }

	clparams := &stripe.CustomerListParams{}
	clparams.Filters.AddFilter("limit", "", "5")
	clparams.Filters.AddFilter("email", "", email)
	i := customer.List(clparams)
  customerid := stripe.String("")
  found := 0
	for i.Next() {
    found += 1
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

func GetCustomer(email string) stripe.Customer  {

  customerid := GetCustomerId(email)
  c, _ := customer.Get(customerid, nil)

  fmt.Printf("    Customer: %s [%s]\n", c.Email, c.ID)
  fmt.Printf("      Default: %s\n", c.DefaultSource)
	return *c
}


func CreateCustomer(email string) string {

  if strings.Compare(email,"") == 0 {
    return *stripe.String("")
  }

  params := &stripe.CustomerParams{
    Email: &email,
  }
  c, _ := customer.New(params)
	return c.ID
}

func CustomerAddCard (emailArg *string, tokenArg *string) {

  if strings.Compare(*emailArg,"") == 0 {
    fmt.Printf("No email specified. Exiting.\n")
    return
  }

  if ! (strings.Compare(*tokenArg,"tok_visa") == 0 || strings.Compare(*tokenArg,"tok_mastercard") == 0) {
    fmt.Printf("Invalid token. Exiting.\n")
    return
  }

  customerid := GetCustomerId (*emailArg)
  // customer := GetCustomer(*emailArg)

  // card_params := &stripe.CardParams{
  //   Customer: stripe.String(customerid),
  //   Token: stripe.String(*tokenArg),
  // }
  //
  // c, err := card.New(card_params)
  // if err!= nil {
  //   fmt.Println(err)
  //   return
  // }
  //
  source_params := &stripe.SourceParams{
    Token: stripe.String(*tokenArg),
  }
  // s, err := source.New(source_params)

  customer_params := &stripe.CustomerParams{
    Source: source_params,
  }

  updated_cust, update_err := customer.Update (
    customerid,
    customer_params,
  )
  if update_err != nil {
    fmt.Printf("Error updating customer\n")
  }

  fmt.Printf("Card source %s added for %s [%s]\n", updated_cust.Sources, *emailArg, customerid)
}


func CustomerGetCards (emailArg *string) {

  if strings.Compare(*emailArg,"") == 0 {
    fmt.Printf("No email specified. Exiting.\n")
    return
  }

  // GetCustomer (*emailArg)
  // customerid := GetCustomerId (*emailArg)
  // card_params := &stripe.CardParams{
  //   Customer: stripe.String(customerid),
  // }
  // c, _ := card.Get(card_params)
  //
  // fmt.Printf("Card %s added for %s\n", c.ID, *emailArg)
}
