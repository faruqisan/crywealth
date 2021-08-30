package repo

import (
	"crywealth/storage"
	"fmt"
	"strconv"
	"time"
)

type (
	WealthRepo struct {
		storage *storage.Engine
	}

	Record struct {
		CreatedAt      time.Time
		GrandTotalUSDT float64
		GrandTotalBIDR float64
	}
)

func New(storage *storage.Engine) *WealthRepo {
	return &WealthRepo{
		storage: storage,
	}
}

func (w *WealthRepo) GetLatestRecord() (Record, error) {
	latest := Record{}
	records, err := w.storage.Read()
	if err != nil {
		return latest, err
	}

	if len(records) != 0 {
		latestRec := records[len(records)-1]

		createdAt, err := time.Parse(time.RFC3339, latestRec[0])
		if err != nil {
			createdAt = time.Now()
		}

		totalUSDT, err := strconv.ParseFloat(latestRec[1], 64)
		if err != nil {
			return latest, err
		}

		totalBIDR, err := strconv.ParseFloat(latestRec[2], 64)
		if err != nil {
			return latest, err
		}

		latest.CreatedAt = createdAt
		latest.GrandTotalUSDT = totalUSDT
		latest.GrandTotalBIDR = totalBIDR
	}

	return latest, nil
}

func (w *WealthRepo) WriteNewRecord(record Record) error {
	data := []string{
		record.CreatedAt.Format(time.RFC3339),
		fmt.Sprintf("%f", record.GrandTotalUSDT),
		fmt.Sprintf("%f", record.GrandTotalBIDR),
	}
	return w.storage.Write(data)
}
