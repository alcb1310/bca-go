package database

import (
	"database/sql"

	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/google/uuid"
)

func (d *Database) GetAllProjects(company_id uuid.UUID, query string) ([]types.Project, error) {
	var rows *sql.Rows
	var err error
	sql := "SELECT id, name, is_active  FROM project WHERE company_id = $1"
	if query == "" {
		rows, err = d.Query(sql, company_id)
	} else {
		sql += " AND name like $2"
		query = "%" + query + "%"
		rows, err = d.Query(sql, company_id, query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var projects []types.Project
	for rows.Next() {
		var name string
		var id uuid.UUID
		var active bool
		if err := rows.Scan(&id, &name, &active); err != nil {
			return nil, err
		}

		projects = append(projects, types.Project{
			ID:        id,
			Name:      &name,
			IsActive:  active,
			CompanyId: company_id,
		})
	}
	return projects, nil
}
func (d *Database) CreateProject(project *types.Project) error {
	sql := "INSERT INTO project (name, is_active, company_id) VALUES ($1, $2, $3)"
	_, err := d.Exec(sql, project.Name, project.IsActive, project.CompanyId)
	return err
}

func (d *Database) GetSingleProject(id, company_id uuid.UUID) (*types.Project, error) {
	sql := "SELECT id, name, is_active  FROM project WHERE id = $1 AND company_id = $2"
	rows, err := d.Query(sql, id, company_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	project := &types.Project{}
	for rows.Next() {
		var name string
		var id uuid.UUID
		var active bool
		if err := rows.Scan(&id, &name, &active); err != nil {
			return nil, err
		}

		project.ID = id
		project.CompanyId = company_id
		project.Name = &name
		project.IsActive = active
	}
	return project, nil
}

func (d *Database) UpdateProject(project *types.Project) error {
	sql := "UPDATE project SET name = $1, is_active = $2 WHERE id = $3 AND company_id = $4"
	_, err := d.Exec(sql, project.Name, project.IsActive, project.ID, project.CompanyId)
	return err
}
