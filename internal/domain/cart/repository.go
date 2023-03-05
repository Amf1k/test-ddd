package cart

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

type NotFoundError struct {
	CustomerID uuid.UUID
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("cart for customer - '%s' not found", e.CustomerID)
}

type Repository interface {
	Create(ctx context.Context, cart Cart) error
	FindByCustomerID(ctx context.Context, customerID uuid.UUID) (Cart, error)
	Update(ctx context.Context, customerID uuid.UUID, updateFn func(ctx context.Context, cart *Cart) error) error
}
