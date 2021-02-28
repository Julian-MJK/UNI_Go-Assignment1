package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

/*************\
|** STRUCTS **|
\*************/

/**
 * @struct Currency
 **/
type Currency struct {
	Code string `json:"code"`
}

/**
 * @struct Response - for REST countries
 **/
type Response struct {
	Country    string       `json:"country"`
	Currencies []Currency   `json:"currencies"`
	Borders    []string     `json:"borders"`
	ExchangeD  ExchangeData `json:"exchangedata"`
}

/**
 * @struct Data
 **/
type Data struct {
	Country  string `json:"country"`
	Code     string `json:"code"`
	Currency string `json:"currency"`
	Rate     Rate   `json:"rate"`
}

/**
 * @struct ExchangeData
 **/
type ExchangeData struct {
	Country string `json:"country"`
	Code    string `json:"code"`
	Date    string `json:"date"`
	Rates   Rate   `json:"rates"`
}

/**
 * @struct Rate
 **/
type Rate struct {
	AUD float64 `json:"AUD,omitempty"`
	BGN float64 `json:"BGN,omitempty"`
	BRL float64 `json:"BRL,omitempty"`
	CAD float64 `json:"CAD,omitempty"`
	CHF float64 `json:"CHF,omitempty"`
	CNY float64 `json:"CNY,omitempty"`
	CZK float64 `json:"CZK,omitempty"`
	DKK float64 `json:"DKK,omitempty"`
	EUR float64 `json:"EUR,omitempty"`
	GBP float64 `json:"GBP,omitempty"`
	HKD float64 `json:"HKD,omitempty"`
	HRK float64 `json:"HRK,omitempty"`
	HUF float64 `json:"HUF,omitempty"`
	IDR float64 `json:"IDR,omitempty"`
	ILS float64 `json:"ILS,omitempty"`
	INR float64 `json:"INR,omitempty"`
	ISK float64 `json:"ISK,omitempty"`
	JPY float64 `json:"JPY,omitempty"`
	KRW float64 `json:"KRW,omitempty"`
	MXN float64 `json:"MXN,omitempty"`
	MYR float64 `json:"MYR,omitempty"`
	NOK float64 `json:"NOK,omitempty"`
	NZD float64 `json:"NZD,omitempty"`
	PHP float64 `json:"PHP,omitempty"`
	PLN float64 `json:"PLN,omitempty"`
	RON float64 `json:"RON,omitempty"`
	RUB float64 `json:"RUB,omitempty"`
	SEK float64 `json:"SEK,omitempty"`
	SGD float64 `json:"SGD,omitempty"`
	THB float64 `json:"THB,omitempty"`
	TRY float64 `json:"TRY,omitempty"`
	USD float64 `json:"USD,omitempty"`
	ZAR float64 `json:"ZAR,omitempty"`
}

/**
 * @struct Diagnostic
 **/
type Diagnostic struct {
	Version       string `json:"version"`
	Uptime        string `json:"uptime"`
	ExchangeAPI   int    `json:"exchangeapi"`
	RestCountries int    `json:"restcountries"`
}

/**
 * @method	main
 * @desc	Main scope
 **/
func main() {
	handleRequests()
}

/**
 * @method	getPort
 * @desc	Returns port or fallback
 **/
func getport() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}

/**
 * @method	handleRequests
 * @desc	Routes URL requests to methods
 **/
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage).Methods("GET")
	myRouter.HandleFunc("/exchange/v1/exchangehistory/{country_name}", exchangeHistory).Methods("GET")
	myRouter.HandleFunc("/exchange/v1/exchangehistory/{country_name}/{begin_date-end_date}", exchangeHistoryDates).Methods("GET")
	myRouter.HandleFunc("/exchange/v1/exchangeborder/{country_name}", exchangeBorder).Methods("GET")
	myRouter.HandleFunc("/exchange/v1/diag", diagnostics).Methods("GET")
	log.Fatal(http.ListenAndServe(getport(), myRouter))
}

/**
 * @method 	init
 * @desc	Initializes startup time
 **/
func init() {
	startTime = time.Now()
}

/**
 * @method 	logIfErr
 * @desc	Logs given error if not nil.
 * @param	err - error
 **/
func logIfErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

/**
 * @method 	fatalIfErr
 * @desc	Logs given error and exits if not nil.
 * @param	err - error
 **/
func fatalIfErr(err error) {
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
}

/**
 * @method	homePage
 * @desc	Endpoint tied to the root address
 **/
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[Endpoint] homePage")
	fmt.Fprintf(w, "Welcome to the homepage! In order to access the following functions, append them to the current website root URL, replacing the items in curled braces with your desired information. \n\n`/exchange/v1/exchangehistory/{country_name}` \n`/exchange/v1/exchangehistory/{country_name}/{begin_date-end_date}` \n`/exchange/v1/exchangeborder/{country_name}` \n`/exchange/v1/diag`")
}

/**
 * @method	exchangeHistory
 * @desc	Receives and outputs exchange data on given country using restcountries.eu API
 **/
func exchangeHistory(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[Endpoint] ExchangeHistory")
	// Compile URL from parameters
	parameters := mux.Vars(r)
	url := "https://restcountries.eu/rest/v2/name/" + parameters["country_name"]

	// Get data and handle errors
	resp, err := http.Get(url)
	fatalIfErr(err)
	respData, err := ioutil.ReadAll(resp.Body)
	logIfErr(err)

	// Unmarshall response into response object
	var respObj []Response
	json.Unmarshal(respData, &respObj)

	url = "https://api.exchangeratesapi.io/latest?symbols=" + respObj[0].Currencies[0].Code
	if respObj[0].Currencies[0].Code == "EUR" {
		url += "&base=USD"
	}

	// Get data and handle errors
	resp, err = http.Get(url)
	logIfErr(err)

	defer resp.Body.Close()
	bodyByte, _ := ioutil.ReadAll(resp.Body)
	bodyStr := json.RawMessage(bodyByte)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bodyStr)
}

/**
 * @method	exchangeHistoryDates
 * @desc	Similair to exchangeHistory, but including begin and end dates.
 * @see		exchangeHistory
 **/
func exchangeHistoryDates(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[Endpoint] ExchangeHistoryDates")
	// Compile URL from parameters
	parameters := mux.Vars(r)
	url := "https://restcountries.eu/rest/v2/name/" + parameters["country_name"]

	// Get data and handle errors
	resp, err := http.Get(url)
	fatalIfErr(err)
	respData, err := ioutil.ReadAll(resp.Body)
	logIfErr(err)

	// Unmarshall response into response object
	var respObj []Response
	json.Unmarshal(respData, &respObj)
	input := parameters["begin_date-end_date"]

	inpRune := []rune(input)
	BegDate := string(inpRune[:10])
	EndDate := string(inpRune[11:])

	// Compile URL with exchange in mind
	url = "https://api.exchangeratesapi.io/history?start_at=" + BegDate + "&end_at=" + EndDate + "&symbols=" + respObj[0].Currencies[0].Code
	if respObj[0].Currencies[0].Code == "EUR" {
		url += "&base=USD"
	}

	// Get data and handle errors
	resp, err = http.Get(url)
	fatalIfErr(err)

	defer resp.Body.Close()
	bodyByte, _ := ioutil.ReadAll(resp.Body)
	bodyStr := json.RawMessage(bodyByte)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bodyStr)
}

/**
 * @method	exchangeBorder
 * @desc	Outputs the exchange rates from all neighboring countries
 **/
func exchangeBorder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[Endpoint] ExchangeBorder")
	// Compile URL from parameters
	parameters := mux.Vars(r)
	url := "https://restcountries.eu/rest/v2/name/" + parameters["country_name"]

	// Get data and handle errors
	resp, err := http.Get(url)
	fatalIfErr(err)
	respData, err := ioutil.ReadAll(resp.Body)
	logIfErr(err)

	// Unmarshall response into response object
	var respObj []Response
	json.Unmarshal(respData, &respObj)
	var data []Data
	var code string

	// Loop through bordering entities
	for i := 0; i < len(respObj[0].Borders); i++ {
		borderingCountries := respObj[0].Borders[i]

		// Request information
		url := "https://restcountries.eu/rest/v2/alpha?codes=" + borderingCountries
		// Get data and handle errors
		respCountry, err := http.Get(url)
		fatalIfErr(err)
		respDataCountry, err := ioutil.ReadAll(respCountry.Body)
		logIfErr(err)

		// Unmarshall response into response object
		var respObjCountry []Response
		json.Unmarshal(respDataCountry, &respObjCountry)

		// Get information and set up base code for currency conversion
		url = "https://api.exchangeratesapi.io/latest?symbols=" + respObjCountry[0].Currencies[0].Code
		if respObjCountry[0].Currencies[0].Code == respObj[0].Currencies[0].Code {
			if respObjCountry[0].Currencies[0].Code == "EUR" {
				url += "&base=USD"
				code = "USD"
			} else {
				url += "&base=EUR"
				code = "EUR"
			}
		} else {
			url += "&base=" + respObj[0].Currencies[0].Code
			code = respObj[0].Currencies[0].Code
		}

		// Get data and handle errors
		respExchange, err := http.Get(url)
		fatalIfErr(err)
		respDataExchange, err := ioutil.ReadAll(respExchange.Body)
		logIfErr(err)

		// Unmarshall response into response object
		var respObjExchange ExchangeData
		json.Unmarshal(respDataExchange, &respObjExchange)

		// Append country exchange information to data
		respObjCountry[0].ExchangeD = respObjExchange
		data = append(data, Data{
			Country:  respObjCountry[0].Country,
			Currency: respObjCountry[0].Currencies[0].Code,
			Rate:     respObjCountry[0].ExchangeD.Rates,
			Code:     code})
	}

	// Return data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

/**
 * @method 	uptime
 * @desc	Returns time since startup
 **/
var startTime time.Time

func uptime() time.Duration {
	return time.Since(startTime)
}

/**
 * @method	formatTime
 * @desc	Formats a string of time, shortening it.
 **/
func formatTime(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}

/**
 * @method	diagnostic
 * @desc	Endpoint that outputs status code from REST api's
 **/
func diagnostics(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[Endpoint] Diagnostics")
	// Get data and handle errors
	respEx, err := http.Get("https://api.exchangeratesapi.io")
	fatalIfErr(err)
	respCount, err := http.Get("https://api.exchangeratesapi.io")
	fatalIfErr(err)
	// Compile and return diagnostic information
	diagnostic := Diagnostic{ExchangeAPI: respEx.StatusCode, RestCountries: respCount.StatusCode, Version: "v1", Uptime: formatTime(uptime())}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(diagnostic)
}
