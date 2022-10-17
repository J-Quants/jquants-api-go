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
	if token.RefreshToken == "" {
		t.Errorf("Could not retrieve refresh token")
	}
}

func TestIdToken(t *testing.T) {
	token, _ := GetIdToken()
	fmt.Printf("%s\n", token)
	if token.IdToken == "" {
		t.Errorf("Could not retrieve id token")
	}
}

func TestDaily(t *testing.T) {
	var quotes = Daily("86970", "", "20220929", "20221003")
	for _, quote := range quotes.DailyQuotes {
		fmt.Printf("%s,%f\n", quote.Date, quote.Close)
	}
	if quotes.DailyQuotes == nil {
		t.Errorf("Could not retrieve quotes")
	}
}
