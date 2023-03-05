package cart

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Item struct {
	productID uuid.UUID
	quantity  int64
	price     decimal.Decimal
}

func NewItem(productID uuid.UUID, quantity int64, price decimal.Decimal) Item {
	return Item{
		productID: productID,
		quantity:  quantity,
		price:     price,
	}
}

func (i Item) ProductID() uuid.UUID {
	return i.productID
}

func (i Item) Quantity() int64 {

	return i.quantity
}

func (i Item) Price() decimal.Decimal {
	return i.price
}

type Cart struct {
	id uuid.UUID

	customerID uuid.UUID
	items      []Item

	totalPrice decimal.Decimal
}

func NewCart(id, customerID uuid.UUID, items []Item) Cart {
	totalPrice := decimal.NewFromInt(0)

	for _, item := range items {
		totalPrice = totalPrice.Add(item.price.Mul(decimal.NewFromInt(item.quantity)))
	}

	return Cart{
		id:         id,
		customerID: customerID,
		items:      items,
		totalPrice: totalPrice,
	}
}

func (c *Cart) ID() uuid.UUID {
	return c.id
}

func (c *Cart) CustomerID() uuid.UUID {
	return c.customerID
}

func (c *Cart) Items() []Item {
	return c.items
}

func (c *Cart) TotalPrice() decimal.Decimal {
	return c.totalPrice
}

func (c *Cart) AddItem(item Item) {
	c.items = append(c.items, item)
	c.totalPrice = c.totalPrice.Add(item.price.Mul(decimal.NewFromInt(item.quantity)))
}
