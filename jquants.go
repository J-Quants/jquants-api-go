package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hellonico/helper"
	"net/http"
	"olympos.io/encoding/edn"
	"os"
)

func main() {
	code := flag.String("code", "86970", "Company Code")
	date := flag.String("date", "20220930", "Date of the quote")
	flag.Parse()

	fmt.Printf("Code: %s and Date: %s\n", *code, *date)
	//if flag.NArg() != 0 {
	//	fmt.Fprintf(os.Stderr, "Usage of jquants:\n")
	//	flag.PrintDefaults()
	//	os.Exit(0)
	//}

	homeDir, err := os.UserHomeDir()
	configDir := homeDir + "/.config/jquants/"
	s, _ := os.ReadFile(configDir + "login.edn")
	//fmt.Printf("%s\n", s)
	var user helper.Login
	edn.Unmarshal(s, &user)
	//fmt.Printf("%s\n", user)

	s2, _ := os.ReadFile(configDir + "id_token.edn")
	var idtoken helper.IdToken
	edn.Unmarshal(s2, &idtoken)
	//fmt.Printf("%s\n", idtoken.IdToken)

	url := fmt.Sprintf("https://api.jpx-jquants.com/v1/prices/daily_quotes?code=%s&date=%s", *code, *date)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	helper.Check(err)

	req.Header = http.Header{
		"Authorization": {"Bearer " + idtoken.IdToken},
	}

	client := http.Client{}
	res, _ := client.Do(req)

	var quotes helper.DailyQuotes
	err_ := json.NewDecoder(res.Body).Decode(&quotes)
	helper.Check(err_)

	hello := quotes.DailyQuotes[0]
	fmt.Printf("Daily: %s > %f\n", hello.Date.String(), hello.Close)
}
