package cli

import (
	"crywealth/internal/wealth/repo"
	"crywealth/internal/wealth/service"
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type (
	WealthCLIHandler struct {
		wealthSVC *service.WealthService // TODO: define interface instead of struct
	}

	assetBalance struct {
		symbol    string
		priceUSDT float64
		priceBIDR float64
		amount    float64
		totalUSDT float64
		totalBIDR float64
	}

	assetOverview struct {
		balances       []assetBalance
		grandTotalUSDT float64
		grandTotalBIDR float64
	}
)

var (
	headerFmt = color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt = color.New(color.FgYellow).SprintfFunc()
)

const (
	usdtbidr = "USDTBIDR"
)

func New(wealthSVC *service.WealthService) *WealthCLIHandler {
	return &WealthCLIHandler{
		wealthSVC: wealthSVC,
	}
}

func (h *WealthCLIHandler) Print() error {
	overview, err := h.assetOverview()
	if err != nil {
		return err
	}
	record, err := h.wealthSVC.GetLatestRecord()
	if err != nil {
		return err
	}

	printAssetsTable(overview)
	printGrandTotalTable(overview)
	printDifferenceTable(record, overview)

	newRecord := repo.Record{
		CreatedAt:      time.Now(),
		GrandTotalUSDT: overview.grandTotalUSDT,
		GrandTotalBIDR: overview.grandTotalBIDR,
	}

	return h.wealthSVC.WriteNewRecord(newRecord)
}

func (h *WealthCLIHandler) assetOverview() (assetOverview, error) {

	var (
		overview       assetOverview
		grandTotalUSDT float64
		grandTotalBIDR float64
	)

	tickerMap, err := h.wealthSVC.Tickers()
	if err != nil {
		return overview, err
	}

	usdtbidrTicker := tickerMap[usdtbidr]
	usdtbidrTickerFloat, err := strconv.ParseFloat(usdtbidrTicker, 64)
	if err != nil {
		return overview, err
	}

	latestSnapshot, err := h.wealthSVC.LatestSnapshot()
	if err != nil {
		return overview, err
	}

	for _, balance := range latestSnapshot.Data.Balances {
		if balance.Free == "0" && balance.Locked == "0" {
			continue
		}

		symUSDT := balance.Asset + "USDT"
		priceUSDT, tickerExists := tickerMap[symUSDT]
		if !tickerExists {
			if balance.Asset == "USDT" {
				priceUSDT = "1"
			} else {
				priceUSDT = "0"
			}
		}

		amountFloat, err := strconv.ParseFloat(balance.Free, 64)
		if err != nil {
			return overview, err
		}

		priceUSDTFloat, err := strconv.ParseFloat(priceUSDT, 64)
		if err != nil {
			return overview, err
		}

		priceBIDR := priceUSDTFloat * usdtbidrTickerFloat

		if balance.Asset == "BIDR" {
			priceUSDTFloat = 1 / usdtbidrTickerFloat
			priceBIDR = 1
		}

		totalUSDT := amountFloat * priceUSDTFloat
		totalBIDR := amountFloat * priceBIDR
		grandTotalUSDT += totalUSDT
		grandTotalBIDR += totalBIDR

		overview.balances = append(overview.balances, assetBalance{
			symbol:    balance.Asset,
			priceUSDT: priceUSDTFloat,
			priceBIDR: priceBIDR,
			amount:    amountFloat,
			totalUSDT: totalUSDT,
			totalBIDR: totalBIDR,
		})
	}

	overview.grandTotalBIDR = grandTotalBIDR
	overview.grandTotalUSDT = grandTotalUSDT

	return overview, nil
}

func printAssetsTable(assetOverview assetOverview) {

	tNow := time.Now()
	fmt.Println("Timestamp: ", tNow)
	fmt.Println()

	tbl := table.New("Symbol", "Price (USDT)", "Price (BIDR)", "Amount", "Total (USDT)", "Total (BIDR)")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, balance := range assetOverview.balances {
		tbl.AddRow(balance.symbol, fmt.Sprintf("%.6f", balance.priceUSDT), fmt.Sprintf("%.6f", balance.priceBIDR), balance.amount, fmt.Sprintf("%.6f", balance.totalUSDT), fmt.Sprintf("%.6f", balance.totalBIDR))
	}

	tbl.Print()
	fmt.Println()

}

func printGrandTotalTable(assetOverview assetOverview) {
	tbl := table.New("Grand Total (USDT)", "Grand Total (IDR)")
	tbl.WithHeaderFormatter(headerFmt)
	tbl.AddRow(fmt.Sprintf("%.2f", assetOverview.grandTotalUSDT), fmt.Sprintf("%.2f", assetOverview.grandTotalBIDR))
	tbl.Print()
	fmt.Println()
}

func printDifferenceTable(record repo.Record, assetOverview assetOverview) {

	diffUSDT := assetOverview.grandTotalUSDT - record.GrandTotalUSDT
	diffIDR := assetOverview.grandTotalBIDR - record.GrandTotalBIDR

	fmt.Println("difference from last check: ", record.CreatedAt)
	tbl := table.New("Previous USDT", "Diff Grand Total (USDT)", "Previous IDR", "Diff Grand Total (IDR)")
	tbl.WithHeaderFormatter(headerFmt)
	tbl.AddRow(fmt.Sprintf("%.2f", record.GrandTotalUSDT), fmt.Sprintf("%.2f", diffUSDT), fmt.Sprintf("%.2f", diffIDR), fmt.Sprintf("%.2f", record.GrandTotalBIDR))
	tbl.Print()
}
