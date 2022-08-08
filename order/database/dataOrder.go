package database

import (
	"context"
	"database/sql"
	"order/entity"

	mssql "github.com/denisenkom/go-mssqldb"
)

func (s *Database) GetOrders(ctx context.Context) ([]entity.Orders, error) {
	var result []entity.Orders

	rows, err := s.SqlDb.QueryContext(ctx, "select order_id,ordered_at,customer_name from orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var row entity.Orders
		err := rows.Scan(
			&row.OrderID,
			&row.OrderedAt,
			&row.CustomerName,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}

func (s *Database) GetItems(ctx context.Context, order *entity.Orders) ([]entity.Item, error) {
	var result []entity.Item

	rows, err := s.SqlDb.QueryContext(ctx, "select item_id,item_code,description,quantity from items where order_id=@OrderID",
		sql.Named("OrderID", order.OrderID))

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var row entity.Item
		err := rows.Scan(
			&row.ItemID,
			&row.ItemCode,
			&row.Description,
			&row.Quantity,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}

func (s *Database) GetOrderByID(ctx context.Context, orderid int) (*entity.Orders, error) {
	result := &entity.Orders{}

	rows, err := s.SqlDb.QueryContext(ctx, "select order_id,ordered_at,customer_name from orders where order_id=@ID", sql.Named("ID", orderid))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&result.OrderID,
			&result.OrderedAt,
			&result.CustomerName,
		)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *Database) CreateOrder(ctx context.Context, order entity.Orders) (string, error) {
	var result string
	sp := "[dbo].[spCreateOrder]"
	var items []entity.DataItem

	for _, dataItem := range order.Items {
		newItem := entity.DataItem{
			ItemCode:    dataItem.ItemCode,
			Description: dataItem.Description,
			Quantity:    dataItem.Quantity,
		}
		items = append(items, newItem)

	}

	tvp := mssql.TVP{
		TypeName: "dbo.ItemTableType",
		Value:    items,
	}

	_, err := s.SqlDb.ExecContext(ctx, sp,
		sql.Named("orderedAt", order.OrderedAt),
		sql.Named("customerName", order.CustomerName),
		sql.Named("items", tvp))

	if err != nil {
		return "", err
	}

	result = "Inserted"

	return result, nil
}

func (s *Database) UpdateOrder(ctx context.Context, ID int, order entity.Orders) (string, error) {
	var result string
	_, err := s.SqlDb.ExecContext(ctx, "update orders set customer_name = @customer_name,ordered_at = @ordered_at where order_id = @id",
		sql.Named("customer_name", order.CustomerName),
		sql.Named("ordered_at", order.OrderedAt),
		sql.Named("id", ID))
	if err != nil {
		return "", err
	}

	for _, item := range order.Items {
		_, err = s.SqlDb.ExecContext(ctx, "update items set item_code=@itemCode, description=@desc, quantity=@quantity where item_id=@itemID and order_id = @orderID",
			sql.Named("itemCode", item.ItemCode),
			sql.Named("desc", item.Description),
			sql.Named("quantity", item.Quantity),
			sql.Named("itemID", item.ItemID),
			sql.Named("orderID", ID))
		if err != nil {
			return "", err
		}

	}

	result = "Updated"

	return result, nil
}

func (s *Database) DeleteOrder(ctx context.Context, ID int) (string, error) {
	var result string

	_, err := s.SqlDb.ExecContext(ctx, "delete from items where order_id=@id", sql.Named("id", ID))
	if err != nil {

		return "", err
	}

	_, err = s.SqlDb.ExecContext(ctx, "delete from orders where order_id=@id", sql.Named("id", ID))
	if err != nil {
		return "", err
	}

	result = "Deleted"

	return result, nil
}
