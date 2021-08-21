package main

import (
	"crywealth/pkg/binance"
	"crywealth/storage"
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

	latestSnapshot, err := getLatestSnapshot(resp.SnapshotVos)
	if err != nil {
		log.Fatal(err)
	}

	tNow := time.Now()
	fmt.Println("Timestamp: ", tNow)
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

	grandTotalIDR := grandTotal * usdtbidrTickerFloat

	tbl = table.New("Grand Total (USDT)", "Grand Total (IDR)")
	tbl.WithHeaderFormatter(headerFmt)
	tbl.AddRow(fmt.Sprintf("%.2f", grandTotal), fmt.Sprintf("%.2f", grandTotalIDR))
	tbl.Print()

	fmt.Println()

	dataStorage := storage.New()

	if err := dataStorage.CreateIfNotExists(); err != nil {
		log.Fatal("err on create if not exists: ", err)
	}

	records, err := dataStorage.Read()
	if err != nil {
		log.Fatal("err on read: ", err)
	}

	if len(records) != 0 {
		latestRecords := records[len(records)-1]
		// file csv format
		// timestamp, totalUSDT, totalBIDR
		if len(latestRecords) == 3 {
			lastUSDT, err := strconv.ParseFloat(latestRecords[1], 64)
			if err != nil {
				log.Fatal(err)
			}

			lastIDR, err := strconv.ParseFloat(latestRecords[2], 64)
			if err != nil {
				log.Fatal(err)
			}

			diffUSDT := grandTotal - lastUSDT
			diffIDR := grandTotalIDR - lastIDR
			fmt.Println("difference from last check: ", latestRecords[0])
			tbl = table.New("Previous USDT", "Diff Grand Total (USDT)", "Previous IDR", "Diff Grand Total (IDR)")
			tbl.WithHeaderFormatter(headerFmt)
			tbl.AddRow(fmt.Sprintf("%.2f", lastUSDT), fmt.Sprintf("%.2f", diffUSDT), fmt.Sprintf("%.2f", diffIDR), fmt.Sprintf("%.2f", lastIDR))
			tbl.Print()

		}
	}

	// write to file latest record
	data := []string{tNow.String(), fmt.Sprintf("%f", grandTotal), fmt.Sprintf("%f", grandTotalIDR)}
	if err := dataStorage.Write(data); err != nil {
		log.Fatal("err on write: ", err)
	}
}

func getLatestSnapshot(snapshots []binance.AccountSnapshotVo) (binance.AccountSnapshotVo, error) {
	latest := binance.AccountSnapshotVo{}
	for _, snapshot := range snapshots {
		if latest.UpdateTime == 0 {
			latest = snapshot
			continue
		}

		//compare time
		latestTime := time.Unix(0, latest.UpdateTime*int64(time.Millisecond))

		currTime := time.Unix(0, snapshot.UpdateTime*int64(time.Millisecond))

		if latestTime.Before(currTime) {
			latest = snapshot
		}
	}
	return latest, nil
}
