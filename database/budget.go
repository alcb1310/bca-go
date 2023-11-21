package database

import (
	"errors"
	"fmt"

	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/google/uuid"
)

func (d *Database) GetBudgets(companyId uuid.UUID, pages *types.PaginationQuery, searchParam string) ([]types.BudgetType, *types.PaginationReturn, error) {
	var ret []types.BudgetType
	sqlQuery := "SELECT id, project, code, budget_item_name, initial_quantity, initial_cost, initial_total, spent_quantity, spent_total, to_spend_quantity, to_spend_cost, to_spend_total, updated_budget FROM budget_description WHERE company_id = $1"
	sqlQuery += " ORDER BY project, code"

	rows, err := d.Query(sqlQuery, companyId)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id *uuid.UUID
		var project, code, budgetItemName *string
		var initialQuantity, initialCost, initialTotal, spentQuantity, spendTotal, toSpendQuantity, toSpendCost, toSpendTotal, updatedBudget *float64

		if err := rows.Scan(&id, &project, &code, &budgetItemName, &initialQuantity, &initialCost, &initialTotal, &spentQuantity, &spendTotal, &toSpendQuantity, &toSpendCost, &toSpendTotal, &updatedBudget); err != nil {
			return nil, nil, err
		}

		ret = append(ret, types.BudgetType{
			ID:              id,
			ProjectName:     project,
			BudgetItemCode:  code,
			BudgetItemName:  budgetItemName,
			InitialQuantity: initialQuantity,
			InitialCost:     initialCost,
			InitialTotal:    initialTotal,
			SpentQuantity:   spentQuantity,
			SpentTotal:      spendTotal,
			ToSpendQuantity: toSpendQuantity,
			ToSpendCost:     toSpendCost,
			ToSpendTotal:    toSpendTotal,
			UpddatedBudget:  updatedBudget,
		})
	}

	sqlQuery = "SELECT count(*) FROM budget_description WHERE company_id = $1"
	pag, err := d.getPaginationStruct(sqlQuery, *pages, companyId, searchParam)
	if err != nil {
		return nil, nil, err
	}

	return ret, &pag, nil
}

func (d *Database) CreateBudget(budget *types.BudgetCreate, companyId uuid.UUID) error {
	return d.SaveBudget(budget, budget.BudgetItemId, companyId)
}

func (d *Database) SaveBudget(budget *types.BudgetCreate, budgetItem *uuid.UUID, companyId uuid.UUID) error {
	if budgetItem == nil {
		return nil
	}

	row := d.QueryRow("SELECT id, code, name, accumulates, level, parent_id FROM budget_item WHERE company_id = $1 AND id = $2", companyId, *budgetItem)
	bi := &types.BudgetItemType{}
	err := row.Scan(&bi.ID, &bi.Code, &bi.Name, &bi.Accumulates, &bi.Level, &bi.ParentId)
	if err != nil {
		return err
	}
	if bi.ID == uuid.Nil {
		return errors.New("budget item not found")
	}

	sqlQuery := "SELECT count(*) FROM budget WHERE company_id = $1 and budget_item_id = $2 and project_id = $3"
	row = d.QueryRow(sqlQuery, companyId, bi.ID, budget.ProjectId)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 && bi.Accumulates {
		sqlQuery = "INSERT INTO budget (company_id, budget_item_id, project_id, initial_total, spent_total, to_spend_total, updated_budget) VALUES ($1, $2, $3, $4, $5, $6, $7)"
		_, err = d.Exec(sqlQuery, companyId, bi.ID, budget.ProjectId, budget.Total, 0, budget.Total, budget.Total)
		if err != nil {
			return err
		}
	} else if count == 0 && !bi.Accumulates {
		sqlQuery = "INSERT INTO budget (company_id, budget_item_id, project_id, initial_quantity, initial_cost, initial_total, spent_quantity, spent_total, to_spend_quantity, to_spend_cost, to_spend_total, updated_budget) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"
		_, err = d.Exec(sqlQuery, companyId, bi.ID, budget.ProjectId, budget.Quantity, budget.Cost, budget.Total, 0, 0, budget.Quantity, budget.Cost, budget.Total, budget.Total)
		if err != nil {
			return err
		}
	} else if bi.Accumulates {
		sqlQuery = "UPDATE budget SET initial_total = initial_total + $1, to_spend_total = to_spend_total + $1, updated_budget = updated_budget + $1 WHERE company_id = $2 AND budget_item_id = $3 AND project_id = $4"
		_, err = d.Exec(sqlQuery, budget.Total, companyId, bi.ID, budget.ProjectId)
		if err != nil {
			return err
		}
	} else {
		return errors.New(fmt.Sprintf("Este presupuesto ya existe, %s", *bi.Code))
	}

	return d.SaveBudget(budget, bi.ParentId, companyId)
}
