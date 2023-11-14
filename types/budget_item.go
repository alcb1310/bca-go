package types

import "github.com/google/uuid"

type BudgetItemType struct {
	ID          uuid.UUID
	Code        *string
	Name        *string
	Accumulates bool
	Level       *uint
	ParentId    *uuid.UUID
	ParentCode  *string
}
