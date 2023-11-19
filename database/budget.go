package database

import (
	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/google/uuid"
)

func (d *Database) GetBudgets(companyId uuid.UUID, pages *types.PaginationQuery, searchParam string) ([]types.BudgetType, *types.PaginationReturn, error) {
	var ret []types.BudgetType
	sqlQuery := "SELECT id, project, code, budget_item_name, initial_quantity, initial_cost, initial_total, spent_quantity, spent_total, to_spend_quantity, to_spend_cost, to_spend_total, updated_budget FROM budget_description WHERE company_id = $1"

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
