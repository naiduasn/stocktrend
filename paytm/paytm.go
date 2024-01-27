package paytm

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func FetchLivePrices(url string, data [][]string) ([]byte, error) {
	const maxConcatenations = 50
	var placeholders []string
	for i := 1; i < len(data); i++ {
		placeholder := fmt.Sprintf("NSE%%3A%s%%3AEQUITY%%2C", data[i][1])
		placeholders = append(placeholders, placeholder)
	}
	if len(placeholders) > maxConcatenations {
		// Split into chunks of 50 placeholders
		for i := 0; i < len(placeholders); i += maxConcatenations {
			end := i + maxConcatenations
			if end > len(placeholders) {
				end = len(placeholders)
			}
			chunk := placeholders[i:end]
			fmt.Println(chunk)
			url2 := url + strings.Join(chunk, "")
			fmt.Println(url)
			td, err := getLatestPrice(url2)
			fmt.Println(string(td), err)
		}
		return getLatestPrice(url)
	} else {
		url = url + strings.Join(placeholders, "")
		fmt.Println(url)
		return getLatestPrice(url)
	}
	//url := "https://developer.paytmmoney.com/data/v1/price/live?mode=LTP&pref=NSE%3A13061%3AEQUITY"
	//fmt.Println(string(body))
	//return getLatestPrice(url)
}

func getLatestPrice(url string) ([]byte, error) {
	//url := "https://developer.paytmmoney.com/data/v1/price/live?mode=LTP&pref=NSE%3A13061%3AEQUITY"

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("x-jwt-token", "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJtZXJjaGFudCIsImlzcyI6InBheXRtbW9uZXkiLCJpZCI6ODA1NjY2LCJleHAiOjE3MDYzODAxOTl9.vgq4F37y-Tk1Bq-r1dQpwUVsukXFiLRx_PxVTe_1jsI")

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
