package types

import "github.com/google/uuid"

type BudgetType struct {
	ID              *uuid.UUID
	BudgetItemCode  *string
	BudgetItemName  *string
	ProjectName     *string
	InitialQuantity *float64
	InitialCost     *float64
	InitialTotal    *float64
	SpentQuantity   *float64
	SpentTotal      *float64
	ToSpendQuantity *float64
	ToSpendCost     *float64
	ToSpendTotal    *float64
	UpddatedBudget  *float64
}

type BudgetCreate struct {
	ID           *uuid.UUID
	BudgetItemId *uuid.UUID
	ProjectId    *uuid.UUID
	Quantity     *float64
	Cost         *float64
	Total        *float64
}
