package main

import (
  "flag"
	"fmt"
	"log"
  "io/ioutil"
	"encoding/csv"
	"strings"
	"strconv"
  "stripey"
)


func main() {
  keyArg := flag.String("key","test","Key to use: Live or Test")
  offerFileArg := flag.String("file","","Filename of offers to process")
  flag.Parse()

  if strings.Compare(*offerFileArg,"") == 0 {
    fmt.Printf("No file specified. Exiting.\n")
    return
  }

  stripey.SetKey(*keyArg)

  // Read File Contents
  content, err := ioutil.ReadFile(*offerFileArg)
  if err!= nil {
    log.Fatal(err)
  }

  lines := string(content)
	r := csv.NewReader(strings.NewReader(lines))
	r.Comma = ','
	r.Comment = '#'

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

  inputfilesize := len(records)
	fmt.Printf("Input file records found: %d\n", inputfilesize)

  // Iterate contents and create Invoices
  for i:=0 ; i<inputfilesize; i++ {
		if (records[i][0] == "S#") {
			continue
		}

    offerid := records[i][0]
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
