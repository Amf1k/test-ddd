package command

import (
	"context"
	"ddd-cart/internal/domain/cart"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AddItem struct {
	CustomerID uuid.UUID
	ProductID  uuid.UUID
	Quantity   int64
	Price      decimal.Decimal
}

type AddItemHandler struct {
	repository cart.Repository
}

func NewAddItemHandler(repository cart.Repository) AddItemHandler {
	return AddItemHandler{
		repository: repository,
	}
}

func (h AddItemHandler) Handle(ctx context.Context, cmd AddItem) error {
	item := cart.NewItem(cmd.ProductID, cmd.Quantity, cmd.Price)
	_, err := h.repository.FindByCustomerID(ctx, cmd.CustomerID)
	switch {
	case errors.Is(err, cart.NotFoundError{}):
		return h.createNewCart(ctx, cmd.CustomerID, item)
	case err != nil:
		return fmt.Errorf("find cart: %w", err)
	}

	return h.updateCart(ctx, cmd.CustomerID, item)
}

func (h AddItemHandler) createNewCart(ctx context.Context, customerID uuid.UUID, item cart.Item) error {
	c := cart.NewCart(uuid.New(), customerID, []cart.Item{item})
	return h.repository.Create(ctx, c)
}

func (h AddItemHandler) updateCart(ctx context.Context, customerID uuid.UUID, item cart.Item) error {
	return h.repository.Update(ctx, customerID, func(ctx context.Context, c *cart.Cart) error {
		c.AddItem(item)
		return nil
	})
}
