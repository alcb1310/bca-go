package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/google/uuid"
)

func (d *Database) getPaginationStruct(sqlQuery string, pagQuery types.PaginationQuery, companyId uuid.UUID, searchParam string) (types.PaginationReturn, error) {
	var err error
	var rows *sql.Rows
	pag := types.PaginationReturn{}

	if searchParam == "" {
		rows, err = d.Query(sqlQuery, companyId)
	} else {
		searchParam = "%" + searchParam + "%"
		rows, err = d.Query(sqlQuery, companyId, searchParam)
	}
	if err != nil {
		log.Println(fmt.Sprintf("Query used: %s", sqlQuery))
		log.Println(fmt.Sprintf("Company Id: %s", companyId))
		log.Println(fmt.Sprintf("Pagination Error: %s", err.Error()))
		return pag, err
	}
	defer rows.Close()

	for rows.Next() {
		var total uint

		if err := rows.Scan(&total); err != nil {
			log.Println(fmt.Sprintf("Query used: %s", sqlQuery))
			log.Println(fmt.Sprintf("Error scanning row: %s", err.Error()))
			return pag, err
		}
		pag.TotalResults = total
		pag.CurrentPage = pagQuery.Offset

		if pagQuery.Limit != 0 {
			pag.Last = total / pagQuery.Limit
			if (total % pagQuery.Limit) != 0 {
				pag.Last++
			}
		} else {
			pag.Last = 1
		}

		if pag.CurrentPage < pag.Last {
			pag.Next = pag.CurrentPage + 1
		} else {
			pag.Next = pag.Last
		}

		if pag.CurrentPage > 1 {
			pag.Prev = pag.CurrentPage - 1
		} else {
			pag.Prev = 1
		}
	}

	return pag, nil
}
