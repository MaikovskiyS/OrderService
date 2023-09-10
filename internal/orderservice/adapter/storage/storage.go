package storage

import (
	"context"
	"fmt"
	"orderservice/internal/orderservice/model"
	"orderservice/internal/orderservice/service"
	"orderservice/pkg/postgres"
	"time"

	"github.com/sirupsen/logrus"
)

const timeout = 5 * time.Second

type storage struct {
	log  *logrus.Logger
	col  tableColumns
	conn *postgres.Client
}

// Constructor
func New(conn *postgres.Client, l *logrus.Logger) service.Storage {
	columns := NewTables()
	return &storage{
		log:  l,
		col:  columns,
		conn: conn,
	}
}

// Create Order in db
func (s *storage) Create(or *model.Order) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	//create order
	o := s.col.Order
	var orderId uint64
	sql, args, err := s.conn.Builder.
		Insert("orders").
		Columns(o.OrderUid, o.TrackNumber, o.Entry, o.Locale, o.InternalSignature, o.CustomerId, o.DeliveryService, o.Shardkey, o.SmId, o.DateCreated, o.OofShard).
		Suffix("RETURNING \"id\"").
		Values(or.OrderUid, or.TrackNumber, or.Entry, or.Locale, or.InternalSignature, or.CustomerId, or.DeliveryService, or.Shardkey, or.SmId, or.DateCreated, or.OofShard).
		ToSql()
	if err != nil {
		s.log.Errorf("OrderStorage -Create - s.OrderInsert.Builder: %s", err)
		return 0, err
	}
	row := s.conn.Pool.QueryRow(ctx, sql, args...)
	err = row.Scan(&orderId)
	if err != nil {
		return 0, fmt.Errorf("OrderStorage -Create - s.OrderInsert.row.Scan: %w", err)
	}
	//create delivery
	d := s.col.Delivery
	sql, args, err = s.conn.Builder.
		Insert("delivery").
		Columns(d.OrderId, d.Name, d.Phone, d.Zip, d.City, d.Address, d.Region, d.Email).
		Values(orderId, or.Delivery.Name, or.Delivery.Phone, or.Delivery.Zip, or.Delivery.City, or.Delivery.Address, or.Delivery.Region, or.Delivery.Email).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("OrderStorage -Create - s.DeliveryInsert.Builder: %w", err)
	}
	dtag, err := s.conn.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return 0, fmt.Errorf("OrderStorage -Create - s.DeliveryInsert.Exec: %w", err)
	}
	dtag.Insert()
	//create payment
	p := s.col.Payment
	sql, args, err = s.conn.Builder.
		Insert("payments").
		Columns(p.OrderId, p.Transaction, p.Request_id, p.Currency, p.Provider, p.Amount, p.Payment_dt, p.Bank, p.Delivery_cost, p.Goods_total, p.Custom_fee).
		Values(orderId, or.Payment.Transaction, or.Payment.Request_id, or.Payment.Currency, or.Payment.Provider, or.Payment.Amount, or.Payment.Payment_dt,
			or.Payment.Bank, or.Payment.Delivery_cost, or.Payment.Goods_total, or.Payment.Custom_fee).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("OrderStorage -Create - s.PaymentInsert.Builder: %w", err)
	}
	ptag, err := s.conn.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return 0, fmt.Errorf("OrderStorage -Create - s.PaymentInsert.Exec: %w", err)
	}
	ptag.Insert()
	//create items
	i := s.col.Items
	for _, it := range or.Items {
		sql, args, err = s.conn.Builder.
			Insert("items").
			Columns(i.OrderId, i.Chrt_id, i.Track_number, i.Price, i.Rid, i.Name, i.Sale, i.Size, i.Total_price, i.Nm_id, i.Brand, i.Status).
			Values(orderId, it.Chrt_id, it.Track_number, it.Price, it.Rid, it.Name, it.Sale, it.Size, it.Total_price, it.Nm_id, it.Brand, it.Status).
			ToSql()
		if err != nil {
			return 0, fmt.Errorf("OrderStorage -Create - s.ItemInsert.Builder: %w", err)
		}
		itag, err := s.conn.Pool.Exec(ctx, sql, args...)
		if err != nil {
			return 0, fmt.Errorf("OrderStorage -Create - s.ItemInsert.Exec: %w", err)
		}
		itag.Insert()
	}
	return orderId, nil
}

// Get all Orders from db
func (s *storage) GetAll(ctx context.Context) ([]*model.Order, error) {
	orders := make([]*model.Order, 0)
	//select order,delivery,payment
	sql, args, err := s.conn.Builder.Select("*").From("orders").Join("delivery ON (id=order_id)").Join("payments ON (delivery.order_id=payments.order_id)").ToSql()
	if err != nil {
		s.log.Info("OrderStorage -GetAll - s.conn.Builder.Select.Err: %w", err)
		return nil, err
	}
	rows, err := s.conn.Pool.Query(ctx, sql, args...)
	if err != nil {
		s.log.Info("OrderStorage -GetAll - s.conn.Pool.Query.Err: %w", err)
		return nil, err
	}
	for rows.Next() {
		items := make([]*model.Item, 0)
		dorderid := 0
		p := 0
		order := &model.Order{Delivery: &model.Delivery{}, Payment: &model.Payment{}, Items: items}
		err = rows.Scan(
			&order.OrderId, &order.OrderUid, &order.TrackNumber, &order.Entry,
			&order.Locale, &order.InternalSignature, &order.CustomerId, &order.DeliveryService, &order.Shardkey,
			&order.SmId, &order.DateCreated, &order.OofShard,

			&dorderid, &order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address,
			&order.Delivery.Region, &order.Delivery.Email,

			&p, &order.Payment.Transaction, &order.Payment.Request_id, &order.Payment.Currency, &order.Payment.Provider,
			&order.Payment.Amount, &order.Payment.Payment_dt, &order.Payment.Bank, &order.Payment.Delivery_cost, &order.Payment.Goods_total, &order.Payment.Custom_fee)
		if err != nil {
			s.log.Info("OrderStorage -GetAll - SelectScan.Err: %w", err)
			return nil, err
		}
		//select items for each order
		sql, args, err := s.conn.Builder.Select("*").From("items").Where("order_id=?", dorderid).ToSql()
		if err != nil {
			s.log.Info("OrderStorage -GetAll - s.conn.Builder.Select.Err: %w", err)
			return nil, err
		}
		rows, err := s.conn.Pool.Query(ctx, sql, args...)
		if err != nil {
			s.log.Info("OrderStorage -GetAll - s.conn.Pool.Query.Err: %w", err)
			return nil, err
		}
		for rows.Next() {
			i := &model.Item{}
			err = rows.Scan(&dorderid, &i.Chrt_id, &i.Track_number, &i.Price, &i.Rid, &i.Name, &i.Sale, &i.Size,
				&i.Total_price, &i.Nm_id, &i.Brand, &i.Status)
			if err != nil {
				s.log.Info("OrderStorage -GetAll - s.Item.Scan.Err: %w", err)
				return nil, err
			}
			items = append(items, i)
		}
		order.Items = items
		orders = append(orders, order)

	}
	return orders, nil
}

// Get Order by Id
func (s *storage) GetById(ctx context.Context, id uint64) (*model.Order, error) {
	//select orders
	order := &model.Order{}
	sql, args, err := s.conn.Builder.Select("*").From("orders").Where("id=?", id).ToSql()
	if err != nil {
		s.log.Info("OrderStorage -GetById - s.Order.SelectBuilder: %w", err)
		return nil, err
	}
	err = s.conn.Pool.QueryRow(ctx, sql, args...).Scan(&order.OrderId, &order.OrderUid, &order.TrackNumber, &order.Entry,
		&order.Locale, &order.InternalSignature, &order.CustomerId, &order.DeliveryService, &order.Shardkey,
		&order.SmId, &order.DateCreated, &order.OofShard)
	if err != nil {
		s.log.Info("OrderStorage -GetById - s.Order.Pool.QueryRow: %w", err)
		return nil, err
	}
	//select delivery
	dorderId := 0
	d := &model.Delivery{}
	sql, args, err = s.conn.Builder.Select("*").From("delivery").Where("order_id=?", id).ToSql()
	if err != nil {
		s.log.Info("OrderStorage -GetById - s.Delivery.SelectBuilder: %w", err)
		return nil, err
	}
	err = s.conn.Pool.QueryRow(ctx, sql, args...).Scan(&dorderId, &d.Name, &d.Phone, &d.Zip, &d.City, &d.Address,
		&d.Region, &d.Email)
	if err != nil {
		s.log.Info("OrderStorage -GetById - s.DeliveryRowsScanErr: %w", err)
		return nil, err
	}

	//select payments
	porderId := 0
	p := &model.Payment{}
	sql, args, err = s.conn.Builder.Select("*").From("payments").Where("order_id=?", id).ToSql()
	if err != nil {
		s.log.Info("OrderStorage -GetById - s.Payments.SelectBuilder: %w", err)
		return nil, err
	}
	err = s.conn.Pool.QueryRow(ctx, sql, args...).Scan(&porderId, &p.Transaction, &p.Request_id, &p.Currency,
		&p.Provider, &p.Amount, &p.Payment_dt, &p.Bank, &p.Delivery_cost, &p.Goods_total, &p.Custom_fee)
	if err != nil {
		s.log.Info("OrderStorage -GetById - s.PaymentsRowsScanErr: %w", err)
		return nil, err
	}

	//select items
	iorderId := 0
	items := make([]*model.Item, 0)

	sql, args, err = s.conn.Builder.Select("*").From("items").Where("order_id=?", id).ToSql()
	if err != nil {
		s.log.Info("OrderStorage -GetById - s.Items.SelectBuilder: %w", err)
		return nil, err
	}
	rows, err := s.conn.Pool.Query(ctx, sql, args...)
	if err != nil {
		s.log.Info("OrderStorage -GetById - s.ItemsQueryErr: %w", err)
		return nil, err
	}
	for rows.Next() {
		i := &model.Item{}
		err = rows.Scan(&iorderId, &i.Chrt_id, &i.Track_number, &i.Price, &i.Rid, &i.Name, &i.Sale, &i.Size,
			&i.Total_price, &i.Nm_id, &i.Brand, &i.Status)
		if err != nil {
			s.log.Info("OrderStorage -GetById - s.ItemsScanErr: %w", err)
			return nil, err
		}
		items = append(items, i)
	}
	order.Delivery = d
	order.Payment = p
	order.Items = items

	return order, nil
}
