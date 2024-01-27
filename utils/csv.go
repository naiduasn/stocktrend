package utils

import (
	"encoding/csv"
	"net/http"
	"os"
)

func WriteCSV(filename string, data [][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	err = w.WriteAll(data)
	if err != nil {
		return err
	}
	return nil
}

func ReadCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	data, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetCSVData(filename string) ([][]string, error) {
	data, err := ReadCSV(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func WriteCSVData(filename string, data [][]string) error {
	err := WriteCSV(filename, data)
	if err != nil {
		return err
	}
	return nil
}

func ReadCSVData(filename string) ([][]string, error) {
	data, err := ReadCSV(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetCSVDataFromURL(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := csv.NewReader(resp.Body).ReadAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func WriteCSVDataFromURL(filename string, url string) error {
	data, err := GetCSVDataFromURL(url)
	if err != nil {
		return err
	}
	return WriteCSVData(filename, data)
}
