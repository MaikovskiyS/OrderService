package model

type Order struct {
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

type BrokerOutputData struct {
	Order *Order
	Err   error
}
