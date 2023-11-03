package database

import (
	"github.com/alcb1310/bca-go-w-test/types"
	"github.com/google/uuid"
)

func (d *Database) GetAllProjects(company_id uuid.UUID) ([]types.Project, error) {
	sql := "SELECT id, name, is_active  FROM project WHERE company_id = $1"
	rows, err := d.Query(sql, company_id)
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
