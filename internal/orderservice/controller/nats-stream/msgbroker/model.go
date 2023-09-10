package msgbroker

import (
	"errors"
	"orderservice/internal/orderservice/model"
)

type OrderDTO struct {
	OrderId           uint64
	OrderUid          string
	TrackNumber       string
	Entry             string
	Delivery          *Delivery
	Payment           *Payment
	Items             []*Item
	Locale            string
	InternalSignature string
	CustomerId        string
	DeliveryService   string
	Shardkey          string
	SmId              string
	DateCreated       string
	OofShard          string
}
type Delivery struct {
	Name    string
	Phone   string
	Zip     string
	City    string
	Address string
	Region  string
	Email   string
}
type Payment struct {
	Transaction   string
	Request_id    string
	Currency      string
	Provider      string
	Amount        uint64
	Payment_dt    int
	Bank          string
	Delivery_cost uint64
	Goods_total   uint64
	Custom_fee    uint64
}
type Item struct {
	Chrt_id      uint64
	Track_number string
	Price        float64
	Rid          string
	Name         string
	Sale         uint8
	Size         string
	Total_price  float64
	Nm_id        uint64
	Brand        string
	Status       uint64
}

func (o *OrderDTO) Validate() error {
	if o.Delivery == nil {
		return errors.New("delivery required")
	}
	if o.Payment == nil {
		return errors.New("payment required")
	}
	if o.Items == nil {
		return errors.New("atleast one item required")
	}
	if o.OrderUid == "" {
		return errors.New("order_uid required")
	}
	return nil
}
func (o *OrderDTO) toModel() *model.Order {
	items := make([]*model.Item, 0)
	for _, v := range o.Items {
		item := &model.Item{
			Chrt_id:      v.Chrt_id,
			Track_number: v.Track_number,
			Price:        v.Price,
			Rid:          v.Rid,
			Name:         v.Name,
			Sale:         v.Sale,
			Size:         v.Size,
			Total_price:  v.Total_price,
			Nm_id:        v.Nm_id,
			Brand:        v.Brand,
			Status:       v.Status,
		}
		items = append(items, item)
	}
	return &model.Order{
		OrderUid:    o.OrderUid,
		TrackNumber: o.TrackNumber,
		Entry:       o.Entry,
		Delivery: &model.Delivery{
			Name:    o.Delivery.Name,
			Phone:   o.Delivery.Phone,
			Zip:     o.Delivery.Zip,
			City:    o.Delivery.City,
			Address: o.Delivery.Address,
			Region:  o.Delivery.Region,
			Email:   o.Delivery.Email,
		},
		Payment: &model.Payment{
			Transaction:   o.Payment.Transaction,
			Request_id:    o.Payment.Request_id,
			Currency:      o.Payment.Currency,
			Provider:      o.Payment.Provider,
			Amount:        o.Payment.Amount,
			Payment_dt:    o.Payment.Payment_dt,
			Bank:          o.Payment.Bank,
			Delivery_cost: o.Payment.Delivery_cost,
			Goods_total:   o.Payment.Goods_total,
			Custom_fee:    o.Payment.Custom_fee,
		},
		Items:             items,
		Locale:            o.Locale,
		InternalSignature: o.InternalSignature,
		CustomerId:        o.CustomerId,
		DeliveryService:   o.DeliveryService,
		Shardkey:          o.Shardkey,
		SmId:              o.SmId,
		DateCreated:       o.DateCreated,
		OofShard:          o.OofShard,
	}
}
