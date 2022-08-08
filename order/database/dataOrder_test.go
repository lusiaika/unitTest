package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"order/entity"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockTvpConverter struct{}

func (converter *mockTvpConverter) ConvertValue(raw interface{}) (driver.Value, error) {

	// Since this function will take the place of every call of ConvertValue, we will inevitably
	// the fake string we return from this function so we need to check whether we've recieved
	// that or a TVP. More extensive logic may be required
	switch inner := raw.(type) {
	case string:
		return raw.(string), nil
	case mssql.TVP:

		// First, verify the type name
		if !strings.EqualFold(inner.TypeName, "dbo.ItemTableType") {
			return nil, fmt.Errorf("Invalid type")
		}
		// VERIFICATION LOGIC HERE

		// Finally, return a fake value that we can use when verifying the arguments
		return Tvp, nil
	case time.Time:
		return raw.(time.Time), nil

	case int:
		return raw.(int), nil

	}

	// We had an invalid type; return an error
	return nil, fmt.Errorf("Invalid type")
}

var DataItems = []entity.DataItem{
	{
		ItemCode:    "ITEM_001",
		Description: "Iphone 10X",
		Quantity:    1,
	},
	{
		ItemCode:    "ITEM_002",
		Description: "Samsung S21",
		Quantity:    3,
	},
	{
		ItemCode:    "ITEM_003",
		Description: "Iphone 12",
		Quantity:    5,
	},
	{
		ItemCode:    "ITEM_004",
		Description: "Samsung S20",
		Quantity:    8,
	},
}

var Tvp = mssql.TVP{
	TypeName: "dbo.ItemTableType",
	Value:    DataItems,
}

var OrderWithItems = entity.Orders{
	OrderID:      1,
	CustomerName: "Blacky",
	OrderedAt:    time.Now(),
	Items: []entity.Item{
		{
			ItemID:      1,
			ItemCode:    "ITEM_001",
			Description: "Iphone 10X",
			Quantity:    1,
		},
		{
			ItemID:      2,
			ItemCode:    "ITEM_002",
			Description: "Samsung S21",
			Quantity:    1,
		},
	},
}
var Items = []entity.Item{
	{
		ItemID:      1,
		ItemCode:    "ITEM_001",
		Description: "Iphone 10X",
		Quantity:    1,
	},
	{
		ItemID:      2,
		ItemCode:    "ITEM_002",
		Description: "Samsung S21",
		Quantity:    1,
	},
	{
		ItemID:      3,
		ItemCode:    "ITEM_003",
		Description: "Iphone 10X",
		Quantity:    2,
	},
	{
		ItemID:      4,
		ItemCode:    "ITEM_004",
		Description: "Samsung S20",
		Quantity:    2,
	},
}
var Order = entity.Orders{
	OrderID:      1,
	CustomerName: "Blacky",
	OrderedAt:    time.Now(),
}

var Orders = []entity.Orders{
	{
		OrderID:      1,
		CustomerName: "Blacky",
		OrderedAt:    time.Now(),
	},
	{
		OrderID:      2,
		CustomerName: "Bone",
		OrderedAt:    time.Now(),
	},
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func Test_GetOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}

	t.Run("database down", func(t *testing.T) {
		query := "select order_id,ordered_at,customer_name from orders"
		mock.ExpectQuery(query).
			WillReturnError(errors.New("db down"))
		got, err := dbtes.GetOrders(ctx)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("Test get orders success", func(t *testing.T) {
		query := "select order_id,ordered_at,customer_name from orders"
		rows := mock.NewRows([]string{"order_id", "ordered_at", "customer_name"}).
			AddRow(Orders[0].OrderID, Orders[0].OrderedAt, Orders[0].CustomerName).
			AddRow(Orders[1].OrderID, Orders[1].OrderedAt, Orders[1].CustomerName)
		mock.ExpectQuery(query).WillReturnRows(rows)
		got, err := dbtes.GetOrders(ctx)
		assert.NotNil(t, got)
		assert.NoError(t, err)
	})

}

func Test_GetItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	dbtes := Database{
		SqlDb: db,
	}

	t.Run("database down", func(t *testing.T) {
		query := "select item_id,item_code,description,quantity from items where order_id=@OrderID"
		mock.ExpectQuery(query).
			WillReturnError(errors.New("db down"))
		got, err := dbtes.GetItems(ctx, &Orders[0])
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.Equal(t, "db down", err.Error())
	})
	t.Run("Empty OrderID", func(t *testing.T) {
		query := "select item_id,item_code,description,quantity from items where order_id=@OrderID"
		mock.ExpectQuery(query).
			WillReturnError(errors.New("OrderID can not be empty"))
		got, err := dbtes.GetItems(ctx, &Orders[0])
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.Equal(t, "OrderID can not be empty", err.Error())
	})
	t.Run("Test get items success", func(t *testing.T) {
		query := "select item_id,item_code,description,quantity from items where order_id=@OrderID"
		rows := mock.NewRows([]string{"item_id", "item_code", "description", "quantity"}).
			AddRow(Items[0].ItemID, Items[0].ItemCode, Items[0].Description, Items[0].Quantity).
			AddRow(Items[1].ItemID, Items[1].ItemCode, Items[1].Description, Items[1].Quantity)
		mock.ExpectQuery(query).WithArgs(Orders[0].OrderID).WillReturnRows(rows)
		got, err := dbtes.GetItems(ctx, &Orders[0])
		assert.NotNil(t, got)
		assert.NoError(t, err)

	})
}

func Test_GetOrderByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	dbtes := Database{
		SqlDb: db,
	}

	t.Run("database down", func(t *testing.T) {
		query := "select order_id,ordered_at,customer_name from orders where order_id=@ID"
		mock.ExpectQuery(query).
			WillReturnError(errors.New("db down"))
		got, err := dbtes.GetOrderByID(ctx, OrderWithItems.OrderID)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.Equal(t, "db down", err.Error())
	})
	t.Run("Empty OrderID", func(t *testing.T) {
		query := "select order_id,ordered_at,customer_name from orders where order_id=@ID"
		mock.ExpectQuery(query).
			WillReturnError(errors.New("OrderID can not be empty"))
		got, err := dbtes.GetOrderByID(ctx, OrderWithItems.OrderID)
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.Equal(t, "OrderID can not be empty", err.Error())
	})

	t.Run("Test get order by ID success", func(t *testing.T) {
		query := "select order_id,ordered_at,customer_name from orders where order_id=@ID"
		rows := mock.NewRows([]string{"OrderID", "OrderedAt", "CustomerName"}).
			AddRow(OrderWithItems.OrderID, OrderWithItems.OrderedAt, OrderWithItems.CustomerName)

		mock.ExpectQuery(query).
			WithArgs(OrderWithItems.OrderID).
			WillReturnRows(rows)
		got, err := dbtes.GetOrderByID(ctx, OrderWithItems.OrderID)
		assert.NotNil(t, got)
		assert.NoError(t, err)
	})
}

func Test_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock, _ := sqlmock.New(sqlmock.ValueConverterOption(&mockTvpConverter{}))
	dbtes := Database{
		SqlDb: db,
	}

	t.Run("database down", func(t *testing.T) {
		mock.ExpectExec("spCreateOrder").
			WithArgs(OrderWithItems.OrderedAt, OrderWithItems.CustomerName, Tvp).
			WillReturnError(errors.New("db down"))
		got, err := dbtes.CreateOrder(ctx, OrderWithItems)
		assert.Error(t, err)
		assert.Equal(t, "", got)
		assert.Equal(t, "db down", err.Error())
	})
	t.Run("Empty CustomerName", func(t *testing.T) {
		OrderWithItems.CustomerName = ""
		mock.ExpectExec("spCreateOrder").
			WithArgs(OrderWithItems.OrderedAt, OrderWithItems.CustomerName, Tvp).
			WillReturnError(errors.New("CustomerName can not be empty"))
		got, err := dbtes.CreateOrder(ctx, OrderWithItems)
		assert.Error(t, err)
		assert.Equal(t, "", got)
		assert.Equal(t, "CustomerName can not be empty", err.Error())
	})
	t.Run("Test create orders success", func(t *testing.T) {
		mock.ExpectExec("spCreateOrder").
			WithArgs(OrderWithItems.OrderedAt, OrderWithItems.CustomerName, Tvp).
			WillReturnResult(sqlmock.NewResult(1, 1))
		got, err := dbtes.CreateOrder(ctx, OrderWithItems)
		assert.NotNil(t, got)
		assert.NoError(t, err)
	})
}

func Test_UpdateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	t.Run("database down", func(t *testing.T) {
		db, mock := NewMock()
		dbtes := Database{
			SqlDb: db,
		}
		query := "update orders set customer_name = @customer_name,ordered_at = @ordered_at where order_id = @id"
		mock.ExpectExec(query).
			WithArgs(OrderWithItems.CustomerName, OrderWithItems.OrderedAt, OrderWithItems.OrderID).
			WillReturnError(errors.New("db down"))

		for _, item := range OrderWithItems.Items {

			querys := "update items set item_code=@itemCode, description=@desc, quantity=@quantity where item_id=@itemID and order_id = @orderID"
			mock.ExpectExec(querys).
				WithArgs(item.ItemCode, item.Description, item.Quantity, item.ItemID, OrderWithItems.OrderID).
				WillReturnError(errors.New("db down"))
		}
		got, err := dbtes.UpdateOrder(ctx, OrderWithItems.OrderID, OrderWithItems)
		assert.Error(t, err)
		assert.Equal(t, "", got)
		assert.Equal(t, "db down", err.Error())
	})
	t.Run("Empty OrderID", func(t *testing.T) {
		OrderWithItems.OrderID = 0
		db, mock := NewMock()
		dbtes := Database{
			SqlDb: db,
		}
		query := "update orders set customer_name = @customer_name,ordered_at = @ordered_at where order_id = @id"
		mock.ExpectExec(query).
			WithArgs(OrderWithItems.CustomerName, OrderWithItems.OrderedAt, OrderWithItems.OrderID).
			WillReturnError(errors.New("OrderID can not be empty"))

		for _, item := range OrderWithItems.Items {

			querys := "update items set item_code=@itemCode, description=@desc, quantity=@quantity where item_id=@itemID and order_id = @orderID"
			mock.ExpectExec(querys).
				WithArgs(item.ItemCode, item.Description, item.Quantity, item.ItemID, OrderWithItems.OrderID).
				WillReturnError(errors.New("OrderID can not be empty"))
		}
		got, err := dbtes.UpdateOrder(ctx, OrderWithItems.OrderID, OrderWithItems)
		assert.Error(t, err)
		assert.Equal(t, "", got)
		assert.Equal(t, "OrderID can not be empty", err.Error())
	})
	t.Run("Test update orders success", func(t *testing.T) {
		db, mock := NewMock()
		dbtes := Database{
			SqlDb: db,
		}
		query := "update orders set customer_name = @customer_name,ordered_at = @ordered_at where order_id = @id"
		mock.ExpectExec(query).
			WithArgs(OrderWithItems.CustomerName, OrderWithItems.OrderedAt, OrderWithItems.OrderID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		for _, item := range OrderWithItems.Items {

			querys := "update items set item_code=@itemCode, description=@desc, quantity=@quantity where item_id=@itemID and order_id = @orderID"
			mock.ExpectExec(querys).
				WithArgs(item.ItemCode, item.Description, item.Quantity, item.ItemID, OrderWithItems.OrderID).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}
		got, err := dbtes.UpdateOrder(ctx, OrderWithItems.OrderID, OrderWithItems)
		assert.NotNil(t, got)
		assert.NoError(t, err)
	})
}

func Test_DeleteOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	t.Run("database down", func(t *testing.T) {
		db, mock := NewMock()
		dbtes := Database{
			SqlDb: db,
		}
		query := "delete from items where order_id=@id"
		mock.ExpectExec(query).
			WithArgs(Order.OrderID).
			WillReturnError(errors.New("db down"))

		querys := "delete from orders where order_id=@id"
		mock.ExpectExec(querys).
			WithArgs(Order.OrderID).
			WillReturnError(errors.New("db down"))

		got, err := dbtes.DeleteOrder(ctx, Order.OrderID)
		assert.Error(t, err)
		assert.Equal(t, "", got)
		assert.Equal(t, "db down", err.Error())
	})
	t.Run("Empty OrderID", func(t *testing.T) {
		db, mock := NewMock()
		dbtes := Database{
			SqlDb: db,
		}
		query := "delete from items where order_id=@id"
		mock.ExpectExec(query).
			WillReturnError(errors.New("OrderID can not be empty"))

		querys := "delete from orders where order_id=@id"
		mock.ExpectExec(querys).
			WillReturnError(errors.New("OrderID can not be empty"))

		got, err := dbtes.DeleteOrder(ctx, Order.OrderID)
		assert.Error(t, err)
		assert.Equal(t, "", got)
		assert.Equal(t, "OrderID can not be empty", err.Error())
	})
	t.Run("Test delete order success", func(t *testing.T) {
		db, mock := NewMock()
		dbtes := Database{
			SqlDb: db,
		}
		query := "delete from items where order_id=@id"
		mock.ExpectExec(query).
			WithArgs(Order.OrderID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		querys := "delete from orders where order_id=@id"
		mock.ExpectExec(querys).
			WithArgs(Order.OrderID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		got, err := dbtes.DeleteOrder(ctx, Order.OrderID)
		assert.NotNil(t, got)
		assert.NoError(t, err)
	})
}
