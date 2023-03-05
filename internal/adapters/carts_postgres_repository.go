package adapters

import (
	"context"
	"database/sql"
	"ddd-cart/internal/domain/cart"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var _ cart.Repository = (*CartsPostgresRepository)(nil)

type CartsPostgresRepository struct {
	db *sql.DB
}

func (c CartsPostgresRepository) Create(ctx context.Context, cart cart.Cart) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "INSERT INTO carts (id, customer_id, total_price) VALUES ($1, $2, $3)", cart.ID(), cart.CustomerID(), cart.TotalPrice())
	if err != nil {
		return fmt.Errorf("insert cart: %w", err)
	}

	for _, item := range cart.Items() {
		_, err = tx.ExecContext(ctx, "INSERT INTO cart_items (cart_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)", cart.ID(), item.ProductID(), item.Quantity(), item.Price())
		if err != nil {
			return fmt.Errorf("insert cart item: %w", err)
		}
	}

	_ = tx.Commit()
	return nil
}

func (c CartsPostgresRepository) FindByCustomerID(ctx context.Context, customerID uuid.UUID) (cart.Cart, error) {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return cart.Cart{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var id uuid.UUID
	var totalPrice int64

	err = tx.QueryRowContext(ctx, "SELECT id, total_price FROM carts WHERE customer_id = $1", customerID).Scan(&id, &totalPrice)
	if err != nil {
		if err == sql.ErrNoRows {
			return cart.Cart{}, cart.NotFoundError{CustomerID: customerID}
		}
		return cart.Cart{}, fmt.Errorf("select cart: %w", err)
	}

	rows, err := tx.QueryContext(ctx, "SELECT product_id, quantity, price FROM cart_items WHERE cart_id = $1", id)
	if err != nil {
		return cart.Cart{}, fmt.Errorf("select cart items: %w", err)
	}
	defer rows.Close()

	var items []cart.Item
	for rows.Next() {
		var productID uuid.UUID
		var quantity int64
		var price decimal.Decimal

		err := rows.Scan(&productID, &quantity, &price)
		if err != nil {
			return cart.Cart{}, fmt.Errorf("scan cart item: %w", err)
		}

		items = append(items, cart.NewItem(productID, quantity, price))
	}

	_ = tx.Commit()

	return cart.NewCart(id, customerID, items), nil
}

func (c CartsPostgresRepository) Update(ctx context.Context, customerID uuid.UUID, updateFn func(ctx context.Context, cart *cart.Cart) error) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Find cart
	// .....

	err = updateFn(ctx, &cart)

	// Save cart

	_ = tx.Commit()

	return nil
}

func NewCartsPostgresRepository(db *sql.DB) *CartsPostgresRepository {
	return &CartsPostgresRepository{db: db}
}
