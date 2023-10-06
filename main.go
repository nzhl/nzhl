package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type CMCResp struct {
	Status status            `json:"status"`
	Data   map[string][]data `json:"data"`
}

type status struct {
	Timestamp    time.Time `json:"timestamp"`
	ErrorCode    int       `json:"error_code"`
	ErrorMessage *string   `json:"error_message"`
	Elapsed      int       `json:"elapsed"`
	CreditCount  int       `json:"credit_count"`
	Notice       *string   `json:"notice"`
}

type data struct {
	Quote map[string]quote `json:"quote"`
}

type quote struct {
	Price                 float64   `json:"price"`
	Volume24H             float64   `json:"volume_24h"`
	VolumeChange24H       float64   `json:"volume_change_24h"`
	PercentChange1H       float64   `json:"percent_change_1h"`
	PercentChange24H      float64   `json:"percent_change_24h"`
	PercentChange7D       float64   `json:"percent_change_7d"`
	PercentChange30D      float64   `json:"percent_change_30d"`
	PercentChange60D      float64   `json:"percent_change_60d"`
	PercentChange90D      float64   `json:"percent_change_90d"`
	MarketCap             float64   `json:"market_cap"`
	MarketCapDominance    float64   `json:"market_cap_dominance"`
	FullyDilutedMarketCap float64   `json:"fully_diluted_market_cap"`
	TVL                   *float64  `json:"tvl"`
	LastUpdated           time.Time `json:"last_updated"`
}

func main() {
	cmcApiKey := os.Getenv("CMC_PRO_API_KEY")
	watchList := []string{"BTC", "ETH", "BNB", "CRV", "ARB", "OKB", "DOGE"}
	url := "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest?symbol=" + strings.Join(watchList, ",")

	log.Default().Println("CMC API URL", url)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-CMC_PRO_API_KEY", cmcApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request to server: %s\n", err)
		return
	}
	defer resp.Body.Close()

	// Read response body
	respBody, _ := io.ReadAll(resp.Body)

	// Parse response body
	var cmcResp CMCResp
	if err := json.Unmarshal(respBody, &cmcResp); err != nil {
		fmt.Printf("Error parsing response body: %s\n", err)
		return
	}

	file, err := os.OpenFile("README.md", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
		return
	}
	defer file.Close()

	location, _ := time.LoadLocation("Asia/Hong_Kong")
	timeInUTC8 := cmcResp.Status.Timestamp.In(location)

	fmt.Fprintf(file, "## 📈 Crypto Prices\n\n")
	fmt.Fprintf(file, "| Coin | Price |\n")
	fmt.Fprintf(file, "| ---- | ----- |\n")
	for _, symbol := range watchList {
		fmt.Fprintf(file, "| %s | $%.4f |\n", symbol, cmcResp.Data[symbol][0].Quote["USD"].Price)
	}
	fmt.Fprintf(file, "\n_Last Updated: %s_", timeInUTC8)
}
