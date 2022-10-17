package jquants_api_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"olympos.io/encoding/edn"
	"os"
	"time"
)

const BASE_URL = "https://api.jpx-jquants.com/v1"
const REFRESH_TOKEN_FILE = "refresh_token.edn"
const ID_TOKEN_FILE = "id_token.edn"

type Login struct {
	UserName string `edn:"mailaddress" json:"mailaddress"`
	Password string `edn:"password" json:"password"`
}
type RefreshToken struct {
	RefreshToken string `edn:"refreshToken" json:"refreshToken"`
}
type IdToken struct {
	IdToken string `edn:"idToken" json:"idToken"`
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

func sendRequest(url string, idToken string) *http.Response {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	Check(err)

	req.Header = http.Header{
		"Authorization": {"Bearer " + idToken},
	}

	client := http.Client{}
	res, _ := client.Do(req)
	return res
}

/**
PRIVATE METHODS TO READ / WRITE CONFIG FILES
*/
func getConfigDir() string {
	homeDir, _ := os.UserHomeDir()
	configDir := homeDir + "/.config/jquants/"
	os.MkdirAll(configDir, os.ModePerm)
	return configDirg
}
func readConfigFile(file string) []byte {
	s, _ := os.ReadFile(getConfigDir() + file)
	return s
}
func writeConfigFile(file string, content []byte) {
	os.WriteFile(getConfigFile(file), content, 0664)
}
func getConfigFile(file string) string {
	return fmt.Sprintf("%s/%s", getConfigDir(), file)
}

func GetUser() Login {
	s, _ := os.ReadFile(getConfigFile("login.edn"))
	var user Login
	edn.Unmarshal(s, &user)
	return user
}
func ReadRefreshToken() RefreshToken {
	var refreshToken RefreshToken
	s := readConfigFile(REFRESH_TOKEN_FILE)
	edn.Unmarshal(s, &refreshToken)
	return refreshToken
}

func ReadIdToken() IdToken {
	var idToken IdToken
	s := readConfigFile(ID_TOKEN_FILE)
	edn.Unmarshal(s, &idToken)
	return idToken
}

func PrepareLogin(username string, password string) {
	var user = Login{username, password}
	encoded, _ := edn.Marshal(&user)
	writeConfigFile("login.edn", encoded)
}

func GetRefreshToken() (RefreshToken, error) {
	var user = GetUser()
	url := fmt.Sprintf("%s/token/auth_user", BASE_URL)
	// fmt.Printf("%s", url)

	data, err := json.Marshal(user)
	// https://golang.cafe/blog/golang-convert-byte-slice-to-io-reader.html

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

	// fmt.Println(req)
	client := http.Client{}
	res, err := client.Do(req)

	var rt RefreshToken
	json.NewDecoder(res.Body).Decode(&rt)

	encoded, err := edn.Marshal(&rt)
	writeConfigFile(REFRESH_TOKEN_FILE, encoded)

	return rt, err
}

func GetIdToken() (IdToken, error) {
	var token = ReadRefreshToken()

	url := fmt.Sprintf("%s/token/auth_refresh?refreshtoken=%s", BASE_URL, token.RefreshToken)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	client := http.Client{}
	res, err := client.Do(req)

	var rt IdToken
	json.NewDecoder(res.Body).Decode(&rt)

	encoded, err := edn.Marshal(&rt)
	writeConfigFile(ID_TOKEN_FILE, encoded)

	return rt, err
}

func Daily(code string, date string, from string, to string) DailyQuotes {
	idtoken := ReadIdToken()

	baseUrl := fmt.Sprintf("%s/prices/daily_quotes?code=%s", BASE_URL, code)
	var url string
	if from != "" && to != "" {
		url = fmt.Sprintf("%s&from=%s&to=%s", baseUrl, from, to)
	} else {
		url = fmt.Sprintf("%s&date=%s", baseUrl, date)
	}
	res := sendRequest(url, idtoken.IdToken)

	var quotes DailyQuotes
	err_ := json.NewDecoder(res.Body).Decode(&quotes)
	Check(err_)
	return quotes
}
