package stripey

import (
  "github.com/joho/godotenv"
  "log"
  "os"
  "strings"
  "fmt"

  "github.com/stripe/stripe-go"
  "github.com/stripe/stripe-go/customer"
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
