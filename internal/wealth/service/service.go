package service

import (
	"crywealth/internal/wealth/repo"
	"crywealth/pkg/binance"
	"log"
)

type (
	WealthService struct {
		binanceClient *binance.Client
		repo          *repo.WealthRepo
	}

	TickersMap map[string]string
)

func New(binanceClient *binance.Client, repo *repo.WealthRepo) *WealthService {
	return &WealthService{
		binanceClient: binanceClient,
		repo:          repo,
	}
}

func (w *WealthService) Tickers() (TickersMap, error) {
	tickers, err := w.binanceClient.GetAllTicker()
	if err != nil {
		return nil, err
	}

	tickerMap := TickersMap{}

	for _, ticker := range tickers {
		tickerMap[ticker.Symbol] = ticker.Price
	}

	return tickerMap, nil
}

func (w *WealthService) LatestSnapshot() (binance.AccountSnapshotVo, error) {
	resp, err := w.binanceClient.AccountSnapshot("SPOT")
	if err != nil {
		return binance.AccountSnapshotVo{}, err
	}
	return getLatestSnapshot(resp.SnapshotVos)
}

func getLatestSnapshot(snapshots []binance.AccountSnapshotVo) (binance.AccountSnapshotVo, error) {
	//
	latest := binance.AccountSnapshotVo{}
	for _, snapshot := range snapshots {
		if latest.UpdateTime == 0 {
			latest = snapshot
			continue
		}

		log.Println("comparing: ",snapshot.UpdateTime, " to: ", latest.UpdateTime)

		if snapshot.UpdateTime > latest.UpdateTime {
			log.Println(snapshot.UpdateTime, " is bigger than: ", latest.UpdateTime)
			latest = snapshot
		}

		//compare time
		// latestTime := time.Unix(0, latest.UpdateTime*int64(time.Millisecond))

		// currTime := time.Unix(0, snapshot.UpdateTime*int64(time.Millisecond))

		// if latestTime.Before(currTime) {
		// 	latest = snapshot
		// }
	}
	return latest, nil
}

func (w *WealthService) GetLatestRecord() (repo.Record, error) {
	return w.repo.GetLatestRecord()
}

func (w *WealthService) WriteNewRecord(record repo.Record) error {
	// TODO: validation
	return w.repo.WriteNewRecord(record)
}
