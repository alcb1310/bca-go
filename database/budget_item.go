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
		var id, parentId uuid.UUID
		var code, name, accumulates, level *string
		if err := rows.Scan(&id, &code, &name, &accumulates, &level, &parentId); err != nil {
			return nil, err
		}

		var accBool bool
		if *accumulates == "f" {
			accBool = false
		} else {
			accBool = true
		}
		levelInt, _ := strconv.Atoi(*level)

		items = append(items, types.BudgetItemType{
			ID:          id,
			Code:        code,
			Name:        name,
			Accumulates: accBool,
			Level:       uint(levelInt),
			ParentId:    &parentId,
		})
	}

	return items, nil
}
