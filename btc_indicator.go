package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	//replace with a free API key from https://free.currencyconverterapi.com/
	ApiKey = "0d7805bde1a33096c74d"
)

// Binance
type Binance struct {
	Mins  int    `json:"mins" bson:"mins"`
	Price string `json:"price" bson:"price"`
}

// Converter
type Converter struct {
	USDINR float64 `json:"USD_INR" bson:"USD_INR"`
}

// Ticker
type Ticker struct {
	Buy  string `json:"buy" bson:"buy"`
	High string `json:"high" bson:"high"`
	Last string `json:"last" bson:"last"`
	Low  string `json:"low" bson:"low"`
	Sell string `json:"sell" bson:"sell"`
	Vol  string `json:"vol" bson:"vol"`
}

// Wazirx
type Wazirx struct {
	At     int     `json:"at" bson:"at"`
	Ticker *Ticker `json:"ticker" bson:"ticker"`
}

func main() {
	wazirXPrice := getWazirXPrice()
	binancePrice := getBinancePrice() * getUSDTPrice()
	delta := binancePrice - wazirXPrice
	fmt.Println("Price delta = ", delta)
	if delta > 0 {
		fmt.Println("Buy bitcoin!")
	} else {
		fmt.Println("Sell bitcoin!")
	}
}

func getWazirXPrice() float64 {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.wazirx.com/api/v2/tickers/btcinr", nil)

	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var responseObject Wazirx
	json.Unmarshal(bodyBytes, &responseObject)

	if s, err := strconv.ParseFloat(responseObject.Ticker.Last, 64); err == nil {
		return s
	}
	return 0

}

func getBinancePrice() float64 {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.binance.com/api/v3/avgPrice?symbol=BTCUSDT", nil)

	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var responseObject Binance
	json.Unmarshal(bodyBytes, &responseObject)

	if s, err := strconv.ParseFloat(responseObject.Price, 64); err == nil {
		return s
	}
	return 0

}

func getUSDTPrice() float64 {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://free.currconv.com/api/v7/convert?q=USD_INR&compact=ultra&apiKey="+ApiKey, nil)

	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var responseObject Converter
	json.Unmarshal(bodyBytes, &responseObject)

	return responseObject.USDINR

}
