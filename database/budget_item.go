package database

import (
	"strconv"

	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/google/uuid"
)

func (d *Database) GetAllBudgetItems(companyId uuid.UUID) ([]types.BudgetItemType, error) {
	sql := "SELECT id, code, name, accumulates, level, parent_id FROM budget_item WHERE company_id = $1"
	rows, err := d.Query(sql, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []types.BudgetItemType
	for rows.Next() {
		var id, parentId *uuid.UUID
		var code, name, accumulates, level *string
		if err := rows.Scan(&id, &code, &name, &accumulates, &level, &parentId); err != nil {
			return nil, err
		}

		accBool := *accumulates == "t"
		levelInt, _ := strconv.Atoi(*level)
		levelUint := uint(levelInt)

		items = append(items, types.BudgetItemType{
			ID:          *id,
			Code:        code,
			Name:        name,
			Accumulates: accBool,
			Level:       &levelUint,
			ParentId:    parentId,
		})
	}

	return items, nil
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
