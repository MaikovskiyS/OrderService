package ordergenerator

import (
	"context"
	"encoding/json"
	"orderservice/internal/orderservice/model"
	"orderservice/pkg/nats/jet"

	"github.com/sirupsen/logrus"
)

type Generator interface {
	SendOrder() error
}
type generator struct {
	log    *logrus.Logger
	stream *jet.Stream
}

func New(stream *jet.Stream, l *logrus.Logger) Generator {
	return &generator{
		log:    l,
		stream: stream,
	}
}

// Send Order message to stream by stream subject
func (g *generator) SendOrder() error {

	order, err := g.GenerateOrder()
	if err != nil {
		g.log.Info("Generator- cant generate order")
		return err
	}
	orderbytes, err := json.Marshal(order)
	if err != nil {
		g.log.Info("Generator- cant marshall order")
		return err
	}
	_, err = g.stream.Jet.Publish(context.Background(), g.stream.Cfg.StreamSubject()[0], orderbytes)
	if err != nil {
		g.log.Info("Generator- cant publish order")
		return err
	}

	return nil
}

// Generate Order model
func (g *generator) GenerateOrder() (*model.Order, error) {
	items := g.GenerateItems()
	return &model.Order{
		OrderUid:    "hfjahdksjdks",
		TrackNumber: "247343h2n3j2",
		Entry:       "",
		Delivery: &model.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmId:              "99",
		DateCreated:       "2021-11-26T06:22:19Z",
		OofShard:          "1",
		Payment: &model.Payment{
			Transaction:   "b563feb7b2b84b6test",
			Request_id:    "",
			Currency:      "USD",
			Provider:      "wbpay",
			Amount:        1817,
			Payment_dt:    1637907727,
			Bank:          "alpha",
			Delivery_cost: 1500,
			Goods_total:   317,
			Custom_fee:    0,
		},
		Items: items,
	}, nil
}

// Generate array of Item models
func (g *generator) GenerateItems() []*model.Item {
	items := make([]*model.Item, 0)
	for i := 0; i < 5; i++ {
		item := &model.Item{
			Chrt_id:      9934930,
			Track_number: "WBILMTESTTRACK",
			Price:        453,
			Rid:          "ab4219087a764ae0btest",
			Name:         "Mascaras",
			Sale:         30,
			Size:         "0",
			Total_price:  317,
			Nm_id:        2389212,
			Brand:        "Vivienne Sabo",
			Status:       202,
		}
		items = append(items, item)
	}
	return items
}
