package repositories

import (
	"context"
	"database/sql"
	"strings"

	"github.com/alexm/fuzzy-builder/internal/models"
)

type SQLiteUsers struct{ db *sql.DB }

func NewSQLiteUsers(db *sql.DB) *SQLiteUsers { return &SQLiteUsers{db: db} }

func (r *SQLiteUsers) Create(ctx context.Context, user *models.User) error {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO users (email, password_hash, role, full_name) VALUES (?,?,?,?)`,
		user.Email, user.PasswordHash, user.Role, user.FullName,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return r.db.QueryRowContext(ctx, `SELECT created_at FROM users WHERE id=?`, id).Scan(&user.CreatedAt)
}

func (r *SQLiteUsers) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, role, full_name, created_at FROM users WHERE email=?`, email,
	)
	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.FullName, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *SQLiteUsers) GetByID(ctx context.Context, id int64) (*models.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, role, full_name, created_at FROM users WHERE id=?`, id,
	)
	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.FullName, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

type SQLiteProjects struct{ db *sql.DB }

func NewSQLiteProjects(db *sql.DB) *SQLiteProjects { return &SQLiteProjects{db: db} }

func (r *SQLiteProjects) Create(ctx context.Context, p *models.Project) error {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO projects (name, description, created_by) VALUES (?,?,?)`,
		p.Name, p.Description, p.CreatedBy,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = id
	return r.db.QueryRowContext(ctx, `SELECT created_at FROM projects WHERE id=?`, id).Scan(&p.CreatedAt)
}

func (r *SQLiteProjects) Update(ctx context.Context, p *models.Project) error {
	_, err := r.db.ExecContext(ctx, `UPDATE projects SET name=?, description=? WHERE id=?`, p.Name, p.Description, p.ID)
	return err
}

func (r *SQLiteProjects) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM projects WHERE id=?`, id)
	return err
}

func (r *SQLiteProjects) List(ctx context.Context, limit, offset int32) ([]models.Project, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, description, created_by, created_at FROM projects ORDER BY id DESC LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedBy, &p.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, rows.Err()
}

type SQLiteDefects struct{ db *sql.DB }

func NewSQLiteDefects(db *sql.DB) *SQLiteDefects { return &SQLiteDefects{db: db} }

func (r *SQLiteDefects) Create(ctx context.Context, d *models.Defect) error {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO defects (title, description, project_id, assigned_to, status, priority, due_date, created_by)
		 VALUES (?,?,?,?,?,?,?,?)`,
		d.Title, d.Description, d.ProjectID, d.AssignedTo, d.Status, d.Priority, d.DueDate, d.CreatedBy,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	d.ID = id
	return r.db.QueryRowContext(ctx, `SELECT created_at FROM defects WHERE id=?`, id).Scan(&d.CreatedAt)
}

func (r *SQLiteDefects) GetByID(ctx context.Context, id int64) (*models.Defect, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, title, description, project_id, assigned_to, status, priority, due_date, created_by, created_at FROM defects WHERE id=?`, id,
	)
	var d models.Defect
	if err := row.Scan(&d.ID, &d.Title, &d.Description, &d.ProjectID, &d.AssignedTo, &d.Status, &d.Priority, &d.DueDate, &d.CreatedBy, &d.CreatedAt); err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *SQLiteDefects) UpdateStatus(ctx context.Context, id int64, status models.DefectStatus) error {
	_, err := r.db.ExecContext(ctx, `UPDATE defects SET status=? WHERE id=?`, status, id)
	return err
}

func (r *SQLiteDefects) List(ctx context.Context, f DefectFilter) ([]models.Defect, error) {
	where := []string{"1=1"}
	args := []any{}
	if f.ProjectID != nil {
		where = append(where, "project_id=?")
		args = append(args, *f.ProjectID)
	}
	if f.AssignedTo != nil {
		where = append(where, "assigned_to=?")
		args = append(args, *f.AssignedTo)
	}
	if f.Status != nil {
		where = append(where, "status=?")
		args = append(args, *f.Status)
	}
	if f.Priority != nil {
		where = append(where, "priority=?")
		args = append(args, *f.Priority)
	}
	query := "SELECT id, title, description, project_id, assigned_to, status, priority, due_date, created_by, created_at FROM defects WHERE " + strings.Join(where, " AND ") + " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, f.Limit, f.Offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Defect
	for rows.Next() {
		var d models.Defect
		if err := rows.Scan(&d.ID, &d.Title, &d.Description, &d.ProjectID, &d.AssignedTo, &d.Status, &d.Priority, &d.DueDate, &d.CreatedBy, &d.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *SQLiteDefects) CountByStatus(ctx context.Context) (map[string]int64, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT status, COUNT(*) FROM defects GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make(map[string]int64)
	for rows.Next() {
		var s string
		var c int64
		if err := rows.Scan(&s, &c); err != nil {
			return nil, err
		}
		out[s] = c
	}
	return out, rows.Err()
}

func (r *SQLiteDefects) CountByProject(ctx context.Context) (map[int64]int64, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT project_id, COUNT(*) FROM defects GROUP BY project_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make(map[int64]int64)
	for rows.Next() {
		var pid int64
		var c int64
		if err := rows.Scan(&pid, &c); err != nil {
			return nil, err
		}
		out[pid] = c
	}
	return out, rows.Err()
}

type SQLiteAttachments struct{ db *sql.DB }

func NewSQLiteAttachments(db *sql.DB) *SQLiteAttachments { return &SQLiteAttachments{db: db} }

func (r *SQLiteAttachments) Create(ctx context.Context, a *models.Attachment) error {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO attachments (defect_id, filename, filepath, uploaded_by) VALUES (?,?,?,?)`,
		a.DefectID, a.Filename, a.Filepath, a.UploadedBy,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	a.ID = id
	return r.db.QueryRowContext(ctx, `SELECT created_at FROM attachments WHERE id=?`, id).Scan(&a.CreatedAt)
}

func (r *SQLiteAttachments) ListByDefect(ctx context.Context, defectID int64, limit, offset int32) ([]models.Attachment, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, defect_id, filename, filepath, uploaded_by, created_at FROM attachments WHERE defect_id=? ORDER BY id DESC LIMIT ? OFFSET ?`,
		defectID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Attachment
	for rows.Next() {
		var a models.Attachment
		if err := rows.Scan(&a.ID, &a.DefectID, &a.Filename, &a.Filepath, &a.UploadedBy, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *SQLiteAttachments) GetByID(ctx context.Context, id int64) (*models.Attachment, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, defect_id, filename, filepath, uploaded_by, created_at FROM attachments WHERE id=?`, id,
	)
	var a models.Attachment
	if err := row.Scan(&a.ID, &a.DefectID, &a.Filename, &a.Filepath, &a.UploadedBy, &a.CreatedAt); err != nil {
		return nil, err
	}
	return &a, nil
}
