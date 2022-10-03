package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"olympos.io/encoding/edn"
	"os"
	"time"
)

type Login struct {
	UserName string `edn:"mailaddress"`
	Password string `edn:"password"`
}
type IdToken struct {
	IdToken string `edn:"idToken"`
}

type JSONTime int64

/*
https://kenzo0107.github.io/2020/05/19/2020-05-20-go-json-time/
*/

// String converts the unix timestamp into a string
func (t JSONTime) String() string {
	tm := t.Time()
	return fmt.Sprintf("\"%s\"", tm.Format("2006-01-02"))
}

// Time returns a `time.Time` representation of this value.
func (t JSONTime) Time() time.Time {
	return time.Unix(int64(t), 0)
}

// UnmarshalJSON will unmarshal both string and int JSON values
func (t *JSONTime) UnmarshalJSON(buf []byte) error {
	s := bytes.Trim(buf, `"`)
	aa, _ := time.Parse("20060102", string(s))
	*t = JSONTime(aa.Unix())
	return nil
}

type Quote struct {
	Code             string   `json:"Code"`
	Close            float64  `json:"Close"`
	Date             JSONTime `json:"Date"`
	AdjustmentHigh   float64  `json:"AdjustmentHigh"`
	Volume           float64  `json:"Volume"`
	TurnoverValue    float64  `json:"TurnoverValue"`
	AdjustmentClose  float64  `json:"AdjustmentClose"`
	AdjustmentLow    float64  `json:"AdjustmentLow"`
	Low              float64  `json:"Low"`
	High             float64  `json:"High"`
	Open             float64  `json:"Open"`
	AdjustmentOpen   float64  `json:"AdjustmentOpen"`
	AdjustmentFactor float64  `json:"AdjustmentFactor"`
	AdjustmentVolume float64  `json:"AdjustmentVolume"`
}
type DailyQuotes struct {
	DailyQuotes []Quote `json:"daily_quotes"`
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func Daily(code string, date string, from string, to string) DailyQuotes {

	homeDir, err := os.UserHomeDir()
	configDir := homeDir + "/.config/jquants/"
	s, _ := os.ReadFile(configDir + "login.edn")
	//fmt.Printf("%s\n", s)
	var user Login
	edn.Unmarshal(s, &user)
	//fmt.Printf("%s\n", user)

	s2, _ := os.ReadFile(configDir + "id_token.edn")
	var idtoken IdToken
	edn.Unmarshal(s2, &idtoken)
	//fmt.Printf("%s\n", idtoken.IdToken)

	baseUrl := fmt.Sprintf("https://api.jpx-jquants.com/v1/prices/daily_quotes?code=%s", code)
	var url string
	if from != "" && to != "" {
		url = fmt.Sprintf("%s&from=%s&to=%s", baseUrl, from, to)
	} else {
		url = fmt.Sprintf("%s&date=%s", baseUrl, date)
	}
	fmt.Printf("URL: %s\n", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	Check(err)

	req.Header = http.Header{
		"Authorization": {"Bearer " + idtoken.IdToken},
	}

	client := http.Client{}
	res, _ := client.Do(req)

	var quotes DailyQuotes
	err_ := json.NewDecoder(res.Body).Decode(&quotes)
	Check(err_)
	return quotes
}
