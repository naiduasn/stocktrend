package sedge

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"text/tabwriter"
)

type Security struct {
	SecurityID int     `json:"SecurityID"`
	Symbol     string  `json:"Symbol"`
	CZG        float64 `json:"CZG"`
	// Add other fields as needed
}

var increasedCounter map[string]int
var decreasedCounter map[string]int

func init() {
	increasedCounter = make(map[string]int)
	decreasedCounter = make(map[string]int)
}

func FetchJSONData(url string) ([]Security, error) {
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

func CompareJSON(oldJSON, newJSON []Security) []map[string]interface{} {
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
	//updateCounters(oldJSON, newJSON)

	return positionChanges
}

func CompareJSONByCZG(oldJSON, newJSON []Security) {

	var changes []Security

	for _, oldSec := range oldJSON {
		for _, newSec := range newJSON {
			if oldSec.Symbol == newSec.Symbol && math.Abs(oldSec.CZG-newSec.CZG) >= 1 {
				changes = append(changes, newSec)
				break // Exit inner loop once found
			}
		}
	}
	fmt.Println("========================================")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Name\tCZG\t")
	for _, obj := range changes {
		fmt.Printf("%v\t%v\n",
			obj.Symbol,
			obj.CZG,
			// Add other fields as needed...
		)
	}
	w.Flush()
	//fmt.Println(changes)

}

func updateCounters(oldJSON, newJSON []Security) {
	// Reset counters
	increasedCounter = make(map[string]int)
	decreasedCounter = make(map[string]int)

	// Update increasedCounter and decreasedCounter
	for _, oldSec := range oldJSON {
		decreasedCounter[oldSec.Symbol]++
	}
	for _, newSec := range newJSON {
		increasedCounter[newSec.Symbol]++
	}

	// Reduce counter for symbols in both increase and decrease
	for symbol := range increasedCounter {
		if count, found := decreasedCounter[symbol]; found {
			// Symbol present in both, reduce the counter
			diff := increasedCounter[symbol] - count
			if diff > 0 {
				increasedCounter[symbol] = diff
			} else {
				decreasedCounter[symbol] = -diff
				delete(increasedCounter, symbol)
			}
		}
	}
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

func sortByChangeDiff(changes []map[string]interface{}) {
	for i := range changes {
		for j := i + 1; j < len(changes); j++ {
			if changes[i]["CZG"].(int) < changes[j]["CZG"].(int) {
				changes[i], changes[j] = changes[j], changes[i]
			}
		}
	}
}

func UpdatePreviousData(oldJSON, newJSON []Security, previousData []Security) []Security {
	// Create a map to store names of stocks in oldJSON
	existingStocks := make(map[string]bool)
	for _, stock := range oldJSON {
		existingStocks[stock.Symbol] = true
	}

	// Iterate through newJSON to find and add new stocks to previousData
	for _, stock := range newJSON {
		if _, exists := existingStocks[stock.Symbol]; !exists {
			// If stock is not in existing stocks, add it to previousData
			previousData = append(previousData, stock)
		}
	}

	return previousData
}
