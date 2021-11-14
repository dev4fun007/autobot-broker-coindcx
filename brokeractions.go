package coindcx

import (
	"encoding/json"
	"fmt"
	common "github.com/dev4fun007/autobot-common"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	NameOfBroker    = "TestBroker"
	BrokerActionTag = "CoinDcxBrokerAction"
)

type BrokerActionCoinDcx struct {
}

func NewBrokerActionCoinDcx() BrokerActionCoinDcx {
	return BrokerActionCoinDcx{}
}

func (receiver BrokerActionCoinDcx) GetBrokerName() string {
	return NameOfBroker
}

func (receiver BrokerActionCoinDcx) ExecuteMarketOrder(order common.RequestMarketOrder) (common.Order, error) {
	log.Info().Str(common.LogComponent, BrokerActionTag).
		Float64("quantity", order.Quantity).
		Str("action", string(order.ActionType)).
		Str("order-type", string(order.OrderType)).
		Msg("broker action completed")

	res := `{"orders":[{"id":"ead19992-43fd-11e8-b027-bb815bcb14ed","market":"%s","order_type":"%s","side":"%s","status":"open","fee_amount":%f,"fee":0.1,"total_quantity":%f,"remaining_quantity":%f,"avg_price":0,"price_per_unit":%f,"created_at":"%s","updated_at":"%s"}]}`
	feeAmount := 0.001 * order.Quantity * order.LastPrice
	res = fmt.Sprintf(res, order.Market, string(order.OrderType), string(order.ActionType), feeAmount, order.Quantity, order.Quantity, order.LastPrice, time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339))
	var orderResponse common.OrderResponse
	err := json.Unmarshal([]byte(res), &orderResponse)
	if err != nil {
		return common.Order{}, err
	}
	if len(orderResponse.Orders) == 1 {
		return orderResponse.Orders[0], nil
	}
	return common.Order{}, nil
}

func (receiver BrokerActionCoinDcx) ExecuteLimitOrder(order common.RequestLimitOrder) (common.Order, error) {
	log.Info().Str(common.LogComponent, BrokerActionTag).
		Float64("quantity", order.Quantity).
		Str("action", string(order.ActionType)).
		Str("order-type", string(order.OrderType)).
		Msg("broker action completed")

	res := `{"orders":[{"id":"ead19992-43fd-11e8-b027-bb815bcb14ed","market":"%s","order_type":"%s","side":"%s","status":"open","fee_amount":%f,"fee":0.1,"total_quantity":%f,"remaining_quantity":%f,"avg_price":0,"price_per_unit":%f,"created_at":"%s","updated_at":"%s"}]}`
	feeAmount := 0.001 * order.Quantity * order.PricePerUnit
	res = fmt.Sprintf(res, order.Market, string(order.OrderType), string(order.ActionType), feeAmount, order.Quantity, order.Quantity, order.PricePerUnit, time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339))
	var orderResponse common.OrderResponse
	err := json.Unmarshal([]byte(res), &orderResponse)
	if err != nil {
		return common.Order{}, err
	}
	if len(orderResponse.Orders) == 1 {
		return orderResponse.Orders[0], nil
	}
	return common.Order{}, nil
}
