package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/linspacestrom/ChangeLogger/internal/domain"
)

type ProjectRepository interface {
	Create(ctx context.Context, projectCreate domain.ProjectCreateOrUpdate) (domain.Project, error)
	GetAll(ctx context.Context) ([]domain.Project, error)
	//Update(ctx context.Context, uuid string, projectUpdate domain.ProjectCreateOrUpdate) (domain.Project, error)
	//GetById(ctx context.Context, uuid string) (domain.ProjectDetail, error)
	//Delete(ctx context.Context, uuid string) error
}

type PoolProjectRepository struct {
	pool *pgxpool.Pool
}

func NewPoolProjectRepository(pool *pgxpool.Pool) *PoolProjectRepository {
	return &PoolProjectRepository{pool: pool}
}

func (r *PoolProjectRepository) Create(ctx context.Context, createProject domain.ProjectCreateOrUpdate) (domain.Project, error) {
	var project domain.Project

	row := r.pool.QueryRow(ctx, `INSERT INTO projects (name) VALUES ($1) RETURNING id, name`, createProject.Title)

	if err := row.Scan(&project.Id, &project.Title); err != nil {
		return domain.Project{}, err
	}

	return project, nil

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
