module stripe.go

go 1.17

replace github.com/griffithbarracks/asa/stripey v0.0.1 => ./stripey

require (
	github.com/griffithbarracks/asa/stripey v0.0.1
	github.com/joho/godotenv v1.4.0
	github.com/stripe/stripe-go v70.15.0+incompatible
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
