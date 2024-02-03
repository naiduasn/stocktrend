package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/naiduasn/stocktrend/paytm"
	"github.com/naiduasn/stocktrend/utils"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}
	//url := os.Getenv("URL")
	delayStr := os.Getenv("DELAY")
	niftyurl := os.Getenv("NIFTYURL")
	livePriceUrl := os.Getenv("PAYTMLIVEPRICEURL")
	ptmjwt := os.Getenv("PTMJWT")
	fmt.Println(delayStr)
	delay, err := strconv.Atoi(delayStr)
	if err != nil {
		fmt.Println("Invalid timeout value:", err)
		return
	}

	data, err := utils.GetCSVDataFromURL(niftyurl)
	SymbolMap := make(map[float64]string)
	for i, record := range data {
		if i == 0 {
			continue
		}
		symbol := record[0]
		securityID, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			panic(err)
		}
		SymbolMap[securityID] = symbol
	}
	fmt.Println(delay, err)

	for {
		pricedata, err := paytm.FetchLivePrices(ptmjwt, livePriceUrl, data, SymbolMap)
		fmt.Println(string(pricedata), err)
		time.Sleep(time.Duration(delay) * time.Minute)
	}
	//getSedgeData(url, delay)
}
