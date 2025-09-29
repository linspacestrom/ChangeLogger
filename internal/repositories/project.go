package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/linspacestrom/ChangeLogger/internal/domain"
)

type ProjectRepository interface {
	GetAll(ctx context.Context) ([]domain.Project, error)
	GetById(ctx context.Context, id string) (domain.ProjectDetail, error)
	Create(ctx context.Context, projectCreate domain.ProjectCreateOrUpdate) (domain.Project, error)
	Update(ctx context.Context, id string, projectUpdate domain.ProjectCreateOrUpdate) (domain.Project, error)
	Delete(ctx context.Context, id string) error
}

type PoolProjectRepository struct {
	pool *pgxpool.Pool
}

func NewPoolProjectRepository(pool *pgxpool.Pool) *PoolProjectRepository {
	return &PoolProjectRepository{pool: pool}
}

func (r *PoolProjectRepository) GetAll(ctx context.Context) ([]domain.Project, error) {
	rows, err := r.pool.Query(ctx, `SELECT * FROM projects`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]domain.Project, 0)

	for rows.Next() {
		var project domain.Project

		err = rows.Scan(&project.Id, &project.Title)
		if err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *PoolProjectRepository) GetById(ctx context.Context, id string) (domain.ProjectDetail, error) {
	var project domain.ProjectDetail

	row := r.pool.QueryRow(ctx, `SELECT * FROM projects WHERE projects.id = $1`, id)

	if err := row.Scan(&project.Id, &project.Title); err != nil {
		return domain.ProjectDetail{}, err
	}

	return project, nil
}

func (r *PoolProjectRepository) Create(ctx context.Context, createProject domain.ProjectCreateOrUpdate) (domain.Project, error) {
	var project domain.Project

	row := r.pool.QueryRow(ctx, `INSERT INTO projects (name) VALUES ($1) RETURNING id, name`, createProject.Title)

	if err := row.Scan(&project.Id, &project.Title); err != nil {
		return project, err
	}

	return project, nil

}

func (r *PoolProjectRepository) Update(ctx context.Context, id string, updateProject domain.ProjectCreateOrUpdate) (domain.Project, error) {
	var project domain.Project

	row := r.pool.QueryRow(ctx, `
		UPDATE projects
		SET name = $1
		WHERE id = $2
		RETURNING id, name
	`, updateProject.Title, id)

	if err := row.Scan(&project.Id, &project.Title); err != nil {
		return project, err
	}

	return project, nil
}

func (r *PoolProjectRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM projects WHERE projects.id = ($1)`, id)
	if err != nil {
		return err
	}
	return nil
}
