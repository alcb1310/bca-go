package types

import "github.com/google/uuid"

type SupplierType struct {
	ID           uuid.UUID
	Ruc          *string
	Name         *string
	ContactName  *string
	ContactEmail *string
	ContactPhone *string
	CompanyId    uuid.UUID
}
