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
