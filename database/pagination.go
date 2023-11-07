package database

import (
	"fmt"
	"log"

	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/google/uuid"
)

func (d *Database) getPaginationStruct(sql string, pagQuery types.PaginationQuery, companyId uuid.UUID) (types.PaginationReturn, error) {
	pag := types.PaginationReturn{}

	rows, err := d.Query(sql, companyId)
	if err != nil {
		log.Println(fmt.Sprintf("Query used: %s", sql))
		log.Println(fmt.Sprintf("Company Id: %s", companyId))
		log.Println(fmt.Sprintf("Pagination Error: %s", err.Error()))
		return pag, err
	}
	defer rows.Close()

	for rows.Next() {
		var total uint

		if err := rows.Scan(&total); err != nil {
			return pag, err
		}
		pag.TotalResults = total
		pag.CurrentPage = pagQuery.Offset

		if pagQuery.Limit != 0 {
			pag.TotalPages = total / pagQuery.Limit
			if (total % pagQuery.Limit) != 0 {
				pag.TotalPages++
			}
		} else {
			pag.TotalPages = 1
		}

		if pag.CurrentPage < pag.TotalPages {
			pag.NextPage = pag.CurrentPage + 1
		} else {
			pag.NextPage = pag.TotalPages
		}
	}

	return pag, nil
}
