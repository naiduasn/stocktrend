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
	fmt.Println(delay, err)

	for {
		pricedata, err := paytm.FetchLivePrices(ptmjwt, livePriceUrl, data)
		fmt.Println(string(pricedata), err)
		time.Sleep(time.Duration(delay) * time.Minute)
	}
	//getSedgeData(url, delay)
}
