package coindcx

import (
	"context"
	"encoding/json"
	"github.com/dev4fun007/autobot-common"
	log "github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	BrokerCoinDcxProducerTag = "BrokerCoinDcxProducerUsdtInr"
)

type TickerItem struct {
	Market       string `json:"market"`
	Change24Hour string `json:"change_24_hour"`
	High         string `json:"high"`
	Low          string `json:"low"`
	Volume       string `json:"volume"`
	LastPrice    string `json:"last_price"`
	//Bid          string `json:"bid"`
	//Ask          int    `json:"ask"`
	Timestamp int `json:"timestamp"`
}

type BrokerCoinDcxProducer struct {
	client              *http.Client
	request             *http.Request
	ticker              *time.Ticker
	tickerDataPublisher common.TickerPublisher
}

func NewCoinDcxProducer(tickerPublisher common.TickerPublisher) *BrokerCoinDcxProducer {
	url := PublicBaseUrl + TickerEndpoint
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error().Str(common.LogComponent, BrokerCoinDcxProducerTag).Err(err).Msg("error creating ticker request")
		return nil
	}
	return &BrokerCoinDcxProducer{
		tickerDataPublisher: tickerPublisher,
		ticker:              time.NewTicker(PollingTimeInSeconds * time.Second),
		client: &http.Client{
			Timeout: HttpTimeoutInSeconds * time.Second,
		},
		request: req,
	}
}

func (p BrokerCoinDcxProducer) StartMarketDataFetcher(ctx context.Context) {
	log.Info().Str(common.LogComponent, BrokerCoinDcxProducerTag).
		Msg("starting coin dcx api calls")
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Str(common.LogComponent, BrokerCoinDcxProducerTag).
					Msg("closing producer data emission")
				return
			case _ = <-p.ticker.C:
				res, errReq := p.client.Do(p.request)
				if errReq == nil {
					body, err := ioutil.ReadAll(res.Body)
					if err != nil {
						log.Error().Str(common.LogComponent, BrokerCoinDcxProducerTag).
							Err(err).Msg("error reading request body")
					}
					var tickers []common.TickerData
					err = json.Unmarshal(body, &tickers)
					if err != nil {
						log.Error().Str(common.LogComponent, BrokerCoinDcxProducerTag).
							Err(err).Msg("error marshalling request body")
					}
					p.tickerDataPublisher.Publish(tickers)
					_ = res.Body.Close()
				} else {
					log.Error().Str(common.LogComponent, BrokerCoinDcxProducerTag).
						Err(errReq).Msg("error making get request to read ticker data")
				}
			}
		}
	}()
}
