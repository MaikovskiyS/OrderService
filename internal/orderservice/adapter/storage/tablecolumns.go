package storage

type tableColumns struct {
	Order    OrderTable
	Payment  PaymentTable
	Delivery DeliveryTable
	Items    ItemsTable
}

type OrderTable struct {
	OrderUid          string
	TrackNumber       string
	Entry             string
	Locale            string
	InternalSignature string
	CustomerId        string
	DeliveryService   string
	Shardkey          string
	SmId              string
	DateCreated       string
	OofShard          string
}
type OrderColumns struct {
}
type DeliveryTable struct {
	OrderId string
	Name    string
	Phone   string
	Zip     string
	City    string
	Address string
	Region  string
	Email   string
}
type PaymentTable struct {
	OrderId       string
	Transaction   string
	Request_id    string
	Currency      string
	Provider      string
	Amount        string
	Payment_dt    string
	Bank          string
	Delivery_cost string
	Goods_total   string
	Custom_fee    string
}

type ItemsTable struct {
	OrderId      string
	Chrt_id      string
	Track_number string
	Price        string
	Rid          string
	Name         string
	Sale         string
	Size         string
	Total_price  string
	Nm_id        string
	Brand        string
	Status       string
}

func NewTables() tableColumns {
	return tableColumns{
		Order: OrderTable{
			OrderUid:          "order_uid",
			TrackNumber:       "track_number",
			Entry:             "entry",
			Locale:            "locate",
			InternalSignature: "internalsignature",
			CustomerId:        "customer_id",
			DeliveryService:   "delivery_service",
			Shardkey:          "shardkey",
			SmId:              "sm_id",
			DateCreated:       "date_created",
			OofShard:          "oof_shard",
		},
		Delivery: DeliveryTable{
			OrderId: "order_id",
			Name:    "name",
			Phone:   "phone",
			Zip:     "zip",
			City:    "city",
			Address: "address",
			Region:  "region",
			Email:   "email",
		},
		Payment: PaymentTable{
			OrderId:       "order_id",
			Transaction:   "transaction",
			Request_id:    "request_id",
			Currency:      "currency",
			Provider:      "provider",
			Amount:        "amount",
			Payment_dt:    "payment_dt",
			Bank:          "bank",
			Delivery_cost: "delivery_cost",
			Goods_total:   "goods_total",
			Custom_fee:    "custom_fee",
		},
		Items: ItemsTable{
			OrderId:      "order_id",
			Chrt_id:      "chrt_id",
			Track_number: "track_number",
			Price:        "price",
			Rid:          "rid",
			Name:         "name",
			Sale:         "sale",
			Size:         "size",
			Total_price:  "total_price",
			Nm_id:        "nm_id",
			Brand:        "brand",
			Status:       "status",
		},
	}
}
