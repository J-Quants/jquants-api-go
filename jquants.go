package main

import (
	"flag"
	"fmt"
)

func main() {
	code := flag.String("code", "86970", "Company Code")
	date := flag.String("date", "20220930", "Date of the quote")
	from := flag.String("from", "20220930", "Start Date for date range")
	to := flag.String("to", "20220930", "End Date for date range")
	flag.Parse()

	fmt.Printf("Code: %s and Date: %s\n", *code, *date)
	//if flag.NArg() != 0 {
	//	fmt.Fprintf(os.Stderr, "Usage of jquants:\n")
	//	flag.PrintDefaults()
	//	os.Exit(0)
	//}

	var quotes = jquants.Daily(*code, *date, *from, *to)
	fmt.Printf("[%d] Daily Quotes for %s \n", len(quotes.DailyQuotes), *code)
	for _, quote := range quotes.DailyQuotes {
		fmt.Printf("%s,%f\n", quote.Date, quote.Close)
	}

}
