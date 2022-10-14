package main

import (
	"flag"
	"fmt"
	jquants "github.com/hellonico/jquants-api-go"
)

func main() {

	refreshToken := flag.Bool("refresh", false, "refresh RefreshToken")
	refreshId := flag.Bool("id", false, "refresh IdToken")

	code := flag.String("code", "86970", "Company Code")
	date := flag.String("date", "20220930", "Date of the quote")
	from := flag.String("from", "", "Start Date for date range")
	to := flag.String("to", "", "End Date for date range")

	flag.Parse()

	fmt.Printf("Code: %s and Date: %s [From: %s To: %s]\n", *code, *date, *from, *to)
	//if flag.NArg() != 0 {
	//	fmt.Fprintf(os.Stderr, "Usage of jquants:\n")
	//	flag.PrintDefaults()
	//	os.Exit(0)
	//}
	if *refreshToken {
		jquants.GetRefreshToken()
	}
	if *refreshId {
		jquants.GetIdToken()
	}

	var quotes = jquants.Daily(*code, *date, *from, *to)
	fmt.Printf("[%d] Daily Quotes for %s \n", len(quotes.DailyQuotes), *code)
	for _, quote := range quotes.DailyQuotes {
		fmt.Printf("%s,%f\n", quote.Date, quote.Close)
	}

}
