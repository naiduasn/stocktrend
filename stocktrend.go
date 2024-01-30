package main

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/naiduasn/stocktrend/paytm"
	sedge "github.com/naiduasn/stocktrend/sedge"
	"github.com/naiduasn/stocktrend/utils"

	"github.com/joho/godotenv"
)

func getEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}
	url := os.Getenv("URL")
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
	fmt.Println("Fetching data from", url)
	//getSedgeData(url, delay)
}

func getSedgeData(url string, delay int) {
	previousData, err := sedge.FetchJSONData(url)
	if err != nil {
		fmt.Println("Failed to fetch data:", err)
		return
	}

	for {
		time.Sleep(time.Duration(delay) * time.Minute)

		newData, err := sedge.FetchJSONData(url)
		if err != nil {
			fmt.Println("Failed to fetch new data:", err)
			continue
		}

		// Compare new data with previous data
		changes := sedge.CompareJSON(previousData, newData)
		sedge.CompareJSONByCZG(previousData, newData)
		// Print ranked changes
		//fmt.Println("Ranked changes based on PositionDiff (highest rank):")
		printTable(changes, previousData)
		//fmt.Println(increasedCounter)
		//fmt.Println(decreasedCounter)
		// Update previous data for the next comparison
		//previousData = newData
		previousData = sedge.UpdatePreviousData(previousData, newData, previousData)
	}
}

func printTable(changes []map[string]interface{}, oldJSON []sedge.Security) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// Print table headers
	fmt.Println("=======================================================")
	fmt.Fprintln(w, "Name\tCZG\tOld Position\tNew Position\tPositionDiff")

	// Print table rows
	for _, change := range changes {
		oldPos := int(change["OldPosition"].(int))
		newPos := int(change["NewPosition"].(int))
		positionDiff := int(change["PositionDiff"].(int))
		czg := change["CZG"].(float64)
		name := oldJSON[oldPos].Symbol

		fmt.Fprintf(w, "%s\t%.2f\t%d\t%d\t%d\n",
			name, czg, oldPos, newPos, positionDiff)
	}

	w.Flush()
}
