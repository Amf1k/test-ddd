package app

import "ddd-cart/internal/app/command"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AddItem command.AddItemHandler
}

type Queries struct {
}
