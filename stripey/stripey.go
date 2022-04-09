package stripey

import (
  "github.com/joho/godotenv"
  "log"
  "os"

  "github.com/stripe/stripe-go"
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
