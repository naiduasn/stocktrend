package paytm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/naiduasn/stocktrend/store"
)

func FetchLivePrices(token string, url string, data [][]string, symbolMap map[float64]string) ([]byte, error) {
	const maxConcatenations = 50
	var placeholders []string
	for i := 1; i < len(data); i++ {
		placeholder := fmt.Sprintf("NSE%%3A%s%%3AEQUITY%%2C", data[i][1])
		placeholders = append(placeholders, placeholder)
	}
	var results []map[string][]map[string]interface{}

	if len(placeholders) > maxConcatenations {
		// Split into chunks of 50 placeholders
		for i := 0; i < len(placeholders); i += maxConcatenations {
			end := i + maxConcatenations
			if end > len(placeholders) {
				end = len(placeholders)
			}
			chunk := placeholders[i:end]
			//fmt.Println(chunk)
			url2 := url + strings.Join(chunk, "")
			//fmt.Println(url)
			td, err := getLatestPrice(token, url2)
			if err != nil {
				// Handle the error, e.g. log it or return an error
				fmt.Println(err)
				return nil, err

			}
			var result map[string][]map[string]interface{}
			if err := json.Unmarshal([]byte(td), &result); err != nil {
				fmt.Println("Error unmarshalling JSON:", err)
			}
			results = append(results, result)
			//fmt.Println(len(results))
		}
		var filteredArray []map[string]interface{}
		for _, result := range results {
			for _, item := range result["data"] {
				tradable, tradableOK := item["tradable"].(bool)
				found, foundOK := item["found"].(bool)
				if tradableOK && foundOK && tradable && found {
					filteredArray = append(filteredArray, item)
				}
			}
		}

		store.InsertData(filteredArray, symbolMap)

		// Create a new map with the filtered array
		filteredResult := map[string][]map[string]interface{}{"data": filteredArray}
		//fmt.Println(len(filteredArray))

		jsonData, err := json.Marshal(filteredResult)
		if err != nil {
			// Handle the error, e.g. log it or return an error
			fmt.Println(err)
			return nil, err
		}
		fmt.Println(string(jsonData))
		// Use the jsonData as needed
		return jsonData, nil
	} else {
		url = url + strings.Join(placeholders, "")
		//fmt.Println(url)
		return getLatestPrice(token, url)
	}
	//url := "https://developer.paytmmoney.com/data/v1/price/live?mode=LTP&pref=NSE%3A13061%3AEQUITY"
	//fmt.Println(string(body))
	//return getLatestPrice(url)
}

func getLatestPrice(token string, url string) ([]byte, error) {
	//url := "https://developer.paytmmoney.com/data/v1/price/live?mode=LTP&pref=NSE%3A13061%3AEQUITY"

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("x-jwt-token", token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return body, nil
}
