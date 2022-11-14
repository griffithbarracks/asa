package stripey

import (
  "fmt"
  "strings"
  "github.com/stripe/stripe-go"
  "github.com/stripe/stripe-go/customer"
  // "github.com/stripe/stripe-go/card"
  // "github.com/stripe/stripe-go/source"
)


func UpdateCustomerEmail(before_email string, after_email string) {
  clparams := &stripe.CustomerListParams{}
	clparams.Filters.AddFilter("limit", "", "5")
	clparams.Filters.AddFilter("email", "", before_email)
	customer_list := customer.List(clparams)
  found := 0
  fmt.Printf("Seaching for customer: %s\n", before_email)
	for customer_list.Next() {
    found += 1
		cust := customer_list.Customer()

    if found == 1 {
      fmt.Printf("    Customer: %s [%s]\n", cust.Email, cust.ID)
      customer_id := cust.ID
      params := &stripe.CustomerParams{}
      params.Email = &after_email
      c, _ := customer.Update(
        customer_id,
        params,
      )
      fmt.Printf("    ==> Updated Customer: %s [%s]\n", c.Email, c.ID)
    } else {
      fmt.Printf("Additional found customer: %s, %s\n", cust.Email, cust.ID)
    }
	}
	if found == 0 {
		fmt.Printf("Customer [%s] not found.\n", before_email)
	}

}

func GetCustomerId(email string) string {

  if strings.Compare(email,"") == 0 {
    return *stripe.String("")
  }

	clparams := &stripe.CustomerListParams{}
	clparams.Filters.AddFilter("limit", "", "5")
	clparams.Filters.AddFilter("email", "", email)
	customer_list := customer.List(clparams)
  customerid := stripe.String("")
  found := 0
	for customer_list.Next() {
    found += 1
		customer := customer_list.Customer()
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

  var cust stripe.Customer

  if strings.Compare(email,"") == 0 {
    return cust
  }

	clparams := &stripe.CustomerListParams{}
	clparams.Filters.AddFilter("limit", "", "5")
	clparams.Filters.AddFilter("email", "", email)
	customer_list := customer.List(clparams)
  found := 0
	for customer_list.Next() {
    found += 1
		cust := customer_list.Customer()
    if found == 1 {
      fmt.Printf("Customer: %s [%s]\n", cust.Email, cust.ID)
    } else {
      fmt.Printf("Additional found customer: %s, %s\n", cust.Email, cust.ID)
    }
	}

	if found == 0 {
		fmt.Printf("Customer: %s not found.\n", email)
    lowercase_email := strings.ToLower(email)
    if (lowercase_email != email) {
      return GetCustomer(lowercase_email)
    }
	}

	return cust
}

func ListCustomers() {
  clparams := &stripe.CustomerListParams{}
  clparams.Filters.AddFilter("limit", "", "200")

  customer_list := customer.List(clparams)
  found := 0
  for customer_list.Next() {
    found += 1
    customer := customer_list.Customer()
    fmt.Printf("%04d, %s, %s\n", found, customer.Email, customer.ID)
  }

  if found == 0 {
    fmt.Printf("Customers not found.\n")
  }

}

func CreateCustomer(email string) string {
  if strings.Compare(email,"") == 0 {
    return *stripe.String("")
  }

  lowercase_email := strings.ToLower(email)
  params := &stripe.CustomerParams{
    Email: &lowercase_email,
  }
  c, _ := customer.New(params)
	return c.ID
}
