package main

import (
	"crywealth/pkg/binance"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/rodaine/table"
)

func main() {

	godotenv.Load()

	client := binance.NewClient("https://api.binance.com", os.Getenv("BINANCE_API_KEY"), os.Getenv("BINANCE_API_SECRET"))
	resp, err := client.AccountSnapshot("SPOT")
	if err != nil {
		log.Fatal(err)
	}

	tickers, err := client.GetAllTicker()
	if err != nil {
		log.Fatal(err)
	}

	tickerMap := map[string]string{}

	for _, ticker := range tickers {
		tickerMap[ticker.Symbol] = ticker.Price
	}

	usdtbidrTicker := tickerMap["USDTBIDR"]
	usdtbidrTickerFloat, err := strconv.ParseFloat(usdtbidrTicker, 64)
	if err != nil {
		log.Fatal(err)
	}

	latestSnapshot := resp.SnapshotVos[len(resp.SnapshotVos)-1]

	fmt.Println("Timestamp: ", time.Now())
	fmt.Println()

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Symbol", "Price (USDT)", "Amount", "Total (USDT)")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	grandTotal := float64(0)
	for _, balance := range latestSnapshot.Data.Balances {
		if balance.Free == "0" && balance.Locked == "0" {
			continue
		}

		symUSDT := balance.Asset + "USDT"
		price, tickerExists := tickerMap[symUSDT]
		if !tickerExists {
			if balance.Asset == "USDT" {
				price = "1"
			} else {
				price = "0"
			}
		}

		amountFloat, err := strconv.ParseFloat(balance.Free, 64)
		if err != nil {
			log.Fatal(err)
		}

		priceFloat, err := strconv.ParseFloat(price, 64)
		if err != nil {
			log.Fatal(err)
		}

		if balance.Asset == "BIDR" {
			priceFloat = 1 / usdtbidrTickerFloat
		}

		total := amountFloat * priceFloat
		grandTotal += total

		tbl.AddRow(balance.Asset, fmt.Sprintf("%.6f", priceFloat), balance.Free, fmt.Sprintf("%.6f", total))
	}

	tbl.Print()

	fmt.Println()

	tbl = table.New("Grand Total (USDT)", "Grand Total (IDR)")
	tbl.WithHeaderFormatter(headerFmt)
	tbl.AddRow(fmt.Sprintf("%.2f", grandTotal), fmt.Sprintf("%.2f", grandTotal*usdtbidrTickerFloat))
	tbl.Print()

}
