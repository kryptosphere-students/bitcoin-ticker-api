package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

func main() {

	// URLs
	binanceURL := "https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT"
	krakenURL := "https://api.kraken.com/0/public/Ticker?pair=XXBTZUSD"
	blockchainURL := "https://blockchain.info/tobtc?currency=USD&value=1"
	bitfinexURL := "https://api-pub.bitfinex.com/v2/ticker/tBTCUSD"

	// Binance Ticker
	binanceTicker := request(binanceURL)
	binancePrice := binanceExtractPrice(binanceTicker)
	fmt.Println("Binance:", binancePrice)

	//Kraken Ticker
	krakenTicker := request(krakenURL)
	krakenPrice := krakenExtractPrice(krakenTicker)
	fmt.Println("Kraken:", krakenPrice)

	//Blockchain(.)com Ticker
	blockchainTicker := request(blockchainURL)
	blockchainPrice := blockcainExtractPrice(blockchainTicker)
	fmt.Println("Blockchain(.)com:", blockchainPrice)

	//Bitfinex Ticker
	bitfinexTicker := request(bitfinexURL)
	bitfinexPrice := bitfinexExtractPrice(bitfinexTicker)
	fmt.Println("Bitfinex:", bitfinexPrice)

	// Weighted average
	weights := map[string]float64{
		"Binance":          0.25,
		"Kraken":           0.25,
		"Blockchain(.)com": 0.25,
		"Bitfinex":         0.25,
	}

	weightedAverage := weights["Binance"]*binancePrice + weights["Kraken"]*krakenPrice + weights["Blockchain(.)com"]*blockchainPrice + weights["Bitfinex"]*bitfinexPrice
	log.Println("Weighted Average:", weightedAverage)

}

func request(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

func binanceExtractPrice(body []byte) float64 {
	bodyString := string(body)
	price := gjson.Get(bodyString, "price")
	priceResult := price.Str
	priceFloat, _ := strconv.ParseFloat(priceResult, 64)

	return priceFloat
}

func krakenExtractPrice(body []byte) float64 {
	bodyString := string(body)
	last := gjson.Get(bodyString, "result.XXBTZUSD.c")
	lastArray := last.Array()
	price := lastArray[0]
	priceResult := price.Str
	priceFloat, _ := strconv.ParseFloat(priceResult, 64)

	return priceFloat

}

func blockcainExtractPrice(body []byte) float64 {
	bodyString := string(body)
	invFloat, _ := strconv.ParseFloat(bodyString, 64)
	priceFloat := 1 / invFloat
	return priceFloat
}

func bitfinexExtractPrice(body []byte) float64 {
	bodyString := string(body)
	bodySplit := strings.Split(bodyString, ",")
	price := bodySplit[6]
	priceFloat, _ := strconv.ParseFloat(price, 64)

	return priceFloat
}
