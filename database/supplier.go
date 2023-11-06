package database

import (
	"database/sql"
	"errors"

	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/google/uuid"
)

func (d *Database) CreateSupplier(s *types.SupplierType) error {
	sql := "INSERT INTO supplier (name, supplier_id, contact_name, contact_email, contact_phone, company_id) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := d.Query(sql, s.Name, s.Ruc, s.ContactName, s.ContactEmail, s.ContactPhone, s.CompanyId)
	return err
}

func (d *Database) UpdateSupplier(s *types.SupplierType) error {
	sql := "UPDATE supplier SET name = $1, supplier_id = $2, contact_name = $3, contact_email = $4, contact_phone = $5 WHERE company_id = $6 AND id = $7"
	_, err := d.Query(sql, s.Name, s.Ruc, s.ContactName, s.ContactEmail, s.ContactPhone, s.CompanyId, s.ID)
	return err
}

func (d *Database) GetSingleSupplier(supplierId, companyId uuid.UUID) (*types.SupplierType, error) {
	sup := &types.SupplierType{}
	sql := "SELECT id, name, supplier_id, contact_name, contact_email, contact_phone FROM supplier WHERE company_id = $2 AND id = $1"

	rows, err := d.Query(sql, supplierId, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id uuid.UUID
		var name, ruc, contactName, contactEmail, contactPhone *string
		if err := rows.Scan(&id, &name, &ruc, &contactName, &contactEmail, &contactPhone); err != nil {
			return nil, err
		}
		sup.ID = id
		sup.Name = name
		sup.Ruc = ruc
		sup.ContactName = contactName
		sup.ContactEmail = contactEmail
		sup.ContactPhone = contactPhone
		sup.CompanyId = companyId
	}

	if sup.Name == nil {
		err := errors.New("Proveedor no encontrado")
		return nil, err
	}

	return sup, nil
}

func (d *Database) GetAllSuppliers(companyId uuid.UUID) ([]types.SupplierType, error) {
	var rows *sql.Rows
	var err error
	sql := "SELECT id, name, supplier_id, contact_name, contact_email, contact_phone FROM supplier WHERE company_id = $1"

	rows, err = d.Query(sql, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []types.SupplierType
	for rows.Next() {
		var id uuid.UUID
		var name, ruc, contactName, contactEmail, contactPhone *string

		if err := rows.Scan(&id, &name, &ruc, &contactName, &contactEmail, &contactPhone); err != nil {
			return nil, err
		}

		suppliers = append(suppliers, types.SupplierType{
			ID:           id,
			Name:         name,
			Ruc:          ruc,
			ContactName:  contactName,
			ContactEmail: contactEmail,
			ContactPhone: contactPhone,
			CompanyId:    companyId,
		})
	}

	return suppliers, nil
}
