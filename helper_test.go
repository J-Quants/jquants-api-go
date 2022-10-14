package jquants_api_go

import (
	"fmt"
	"os"
	"testing"
)

func TestPrepareLogin(t *testing.T) {
	PrepareLogin(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
}
func TestRefreshToken(t *testing.T) {
	token, _ := GetRefreshToken()
	fmt.Printf("%s\n", token)
}
func TestIdToken(t *testing.T) {
	token, _ := GetIdToken()
	fmt.Printf("%s\n", token)
}

func TestDaily(t *testing.T) {
	var quotes = Daily("86970", "", "20220929", "20221003")
	for _, quote := range quotes.DailyQuotes {
		fmt.Printf("%s,%f\n", quote.Date, quote.Close)
	}
}
