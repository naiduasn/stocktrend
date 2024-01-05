package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/joho/godotenv"
)

type Security struct {
	SecurityID int     `json:"SecurityID"`
	Symbol     string  `json:"Symbol"`
	CZG        float64 `json:"CZG"`
	// Add other fields as needed
}

func fetchJSONData(url string) ([]Security, error) {
	var securities []Security

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &securities)
	if err != nil {
		return nil, err
	}

	return securities, nil
}

func compareJSON(oldJSON, newJSON []Security) []map[string]interface{} {
	// Create maps to store CZG values and positions
	oldCZGPositions := make(map[float64]int)
	newCZGPositions := make(map[float64]int)

	// Populate oldCZGPositions
	for i, item := range oldJSON {
		oldCZGPositions[item.CZG] = i
	}

	// Populate newCZGPositions
	for i, item := range newJSON {
		newCZGPositions[item.CZG] = i
	}

	// Find position changes
	var positionChanges []map[string]interface{}
	for czg, oldPos := range oldCZGPositions {
		if newPos, exists := newCZGPositions[czg]; exists {
			if oldPos != newPos {
				change := map[string]interface{}{
					"CZG":          czg,
					"OldPosition":  oldPos,
					"NewPosition":  newPos,
					"PositionDiff": abs(oldPos - newPos), // Ensure PositionDiff is an int
				}

				positionChanges = append(positionChanges, change)
			}
		}
	}

	// Sort positionChanges by PositionDiff (highest rank)
	sortByPositionDiff(positionChanges)

	return positionChanges
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func sortByPositionDiff(changes []map[string]interface{}) {
	for i := range changes {
		for j := i + 1; j < len(changes); j++ {
			if changes[i]["PositionDiff"].(int) < changes[j]["PositionDiff"].(int) {
				changes[i], changes[j] = changes[j], changes[i]
			}
		}
	}
}

func printTable(changes []map[string]interface{}, oldJSON []Security) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// Print table headers
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

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}
	url := os.Getenv("URL")
	delayStr := os.Getenv("DELAY")
	fmt.Println(delayStr)
	delay, err := strconv.Atoi(delayStr)
	if err != nil {
		fmt.Println("Invalid timeout value:", err)
		return
	}

	previousData, err := fetchJSONData(url)
	if err != nil {
		fmt.Println("Failed to fetch data:", err)
		return
	}

	for {
		time.Sleep(time.Duration(delay) * time.Minute)

		newData, err := fetchJSONData(url)
		if err != nil {
			fmt.Println("Failed to fetch new data:", err)
			continue
		}

		// Compare new data with previous data
		changes := compareJSON(previousData, newData)
		// Print ranked changes
		fmt.Println("Ranked changes based on PositionDiff (highest rank):")
		printTable(changes, previousData)

		// Update previous data for the next comparison
		previousData = newData
	}
}
