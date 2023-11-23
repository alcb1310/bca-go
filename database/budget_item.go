package database

import (
	"database/sql"

	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/google/uuid"
)

func (d *Database) GetAllBudgetItems(companyId uuid.UUID, pages *types.PaginationQuery, searchParam string) ([]types.BudgetItemType, *types.PaginationReturn, error) {
	var rows *sql.Rows
	var err error
	sqlQuery := "SELECT id, code, name, accumulates, level, parent_id, parent_code FROM budget_item_with_parents WHERE company_id = $1"

	if searchParam != "" {
		sqlQuery += " AND name ILIKE $2"
		sqlQuery += " ORDER BY code"
		if pages != nil && pages.Limit != 0 {
			sqlQuery += " LIMIT $3"
			sqlQuery += " OFFSET $4"
			rows, err = d.Query(sqlQuery, companyId, "%"+searchParam+"%", pages.Limit, (pages.Offset-1)*pages.Limit)
		} else {
			rows, err = d.Query(sqlQuery, companyId, "%"+searchParam+"%")
		}
	} else {
		sqlQuery += " ORDER BY code"
		if pages != nil {
			sqlQuery += " LIMIT $2"
			sqlQuery += " OFFSET $3"
			rows, err = d.Query(sqlQuery, companyId, pages.Limit, (pages.Offset-1)*pages.Limit)
		} else {
			rows, err = d.Query(sqlQuery, companyId)
		}
	}

	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var items []types.BudgetItemType
	for rows.Next() {
		var id, parent_id *uuid.UUID
		var code, name, accumulates, parent_code *string
		var level *uint
		if err := rows.Scan(&id, &code, &name, &accumulates, &level, &parent_id, &parent_code); err != nil {
			return nil, nil, err
		}

		items = append(items, types.BudgetItemType{
			ID:          *id,
			Code:        code,
			Name:        name,
			Accumulates: *accumulates == "true",
			Level:       level,
			ParentId:    parent_id,
			ParentCode:  parent_code,
		})
	}

	sqlQuery = "SELECT count(*) FROM budget_item_with_parents WHERE company_id = $1"
	if searchParam != "" {
		sqlQuery += " AND name LIKE $2"
	}

	pag, err := d.getPaginationStruct(sqlQuery, *pages, companyId, searchParam)
	if err != nil {
		return nil, nil, err
	}

	return items, &pag, nil
}

func (d *Database) AllBudgetItemsByAccumulates(companyId uuid.UUID, accumulates bool) ([]types.BudgetItemType, error) {
	sql := "SELECT id, code, name  FROM budget_item WHERE company_id = $1 AND accumulates = $2 ORDER BY code"
	rows, err := d.Query(sql, companyId, accumulates)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []types.BudgetItemType
	for rows.Next() {
		var id uuid.UUID
		var code, name *string
		if err := rows.Scan(&id, &code, &name); err != nil {
			return nil, err
		}

		items = append(items, types.BudgetItemType{
			ID:   id,
			Code: code,
			Name: name,
		})
	}

	return items, nil
}

func (d *Database) CreateBudgetItem(companyId uuid.UUID, budgetItem *types.BudgetItemType) error {
	sql := "INSERT INTO budget_item (company_id, code, name, accumulates, level, parent_id) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := d.Exec(sql, companyId, budgetItem.Code, budgetItem.Name, budgetItem.Accumulates, budgetItem.Level, budgetItem.ParentId)
	return err
}

func (d *Database) GetBudgetItemById(companyId uuid.UUID, budgetItemId uuid.UUID) (*types.BudgetItemType, error) {
	sql := "SELECT id, code, name, accumulates, level, parent_id, parent_code FROM budget_item_with_parents WHERE company_id = $1 and id = $2"
	var id, parent_id *uuid.UUID
	var level *uint
	var code, name, parent_code *string
	var accumulates bool
	err := d.QueryRow(sql, companyId, budgetItemId).Scan(&id, &code, &name, &accumulates, &level, &parent_id, &parent_code)
	if err != nil {
		return nil, err
	}
	return &types.BudgetItemType{
		ID:          *id,
		Code:        code,
		Name:        name,
		Accumulates: accumulates,
		Level:       level,
		ParentId:    parent_id,
		ParentCode:  parent_code,
	}, nil
}

func (d *Database) UpdateBudgetItem(companyId uuid.UUID, budgetItem *types.BudgetItemType) error {
	sql := "UPDATE budget_item SET code = $1, name = $2, accumulates = $3, level = $4, parent_id = $5 WHERE company_id = $6 AND id = $7"
	_, err := d.Exec(sql, budgetItem.Code, budgetItem.Name, budgetItem.Accumulates, budgetItem.Level, budgetItem.ParentId, companyId, budgetItem.ID)
	return err
}
