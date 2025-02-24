// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: user.sql

package sqlc

import (
	"context"
)

const activateUser = `-- name: ActivateUser :exec
UPDATE users SET
    is_active = TRUE
WHERE id = $1
`

func (q *Queries) ActivateUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, activateUser, id)
	return err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    first_name,
    last_name,
    address,
    phone,
    is_active,
    is_superuser
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, email, password_hash, first_name, last_name, address, phone, is_active, is_superuser, created_at, updated_at
`

type CreateUserParams struct {
	Email        *string `json:"email"`
	PasswordHash string  `json:"-"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Address      string  `json:"address"`
	Phone        string  `json:"phone"`
	IsActive     bool    `json:"is_active"`
	IsSuperuser  bool    `json:"is_superuser"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Email,
		arg.PasswordHash,
		arg.FirstName,
		arg.LastName,
		arg.Address,
		arg.Phone,
		arg.IsActive,
		arg.IsSuperuser,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PasswordHash,
		&i.FirstName,
		&i.LastName,
		&i.Address,
		&i.Phone,
		&i.IsActive,
		&i.IsSuperuser,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deactivateUser = `-- name: DeactivateUser :exec
UPDATE users SET
    is_active = FALSE
WHERE id = $1
`

func (q *Queries) DeactivateUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deactivateUser, id)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const demoteUser = `-- name: DemoteUser :exec
UPDATE users SET
    is_superuser = FALSE
WHERE id = $1
`

func (q *Queries) DemoteUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, demoteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, password_hash, first_name, last_name, address, phone, is_active, is_superuser, created_at, updated_at FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email *string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PasswordHash,
		&i.FirstName,
		&i.LastName,
		&i.Address,
		&i.Phone,
		&i.IsActive,
		&i.IsSuperuser,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, email, password_hash, first_name, last_name, address, phone, is_active, is_superuser, created_at, updated_at FROM users WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PasswordHash,
		&i.FirstName,
		&i.LastName,
		&i.Address,
		&i.Phone,
		&i.IsActive,
		&i.IsSuperuser,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listActiveUsers = `-- name: ListActiveUsers :many
SELECT id, email, password_hash, first_name, last_name, address, phone, is_active, is_superuser, created_at, updated_at FROM users WHERE is_active = TRUE ORDER BY id
`

func (q *Queries) ListActiveUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, listActiveUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.PasswordHash,
			&i.FirstName,
			&i.LastName,
			&i.Address,
			&i.Phone,
			&i.IsActive,
			&i.IsSuperuser,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsers = `-- name: ListUsers :many
SELECT id, email, password_hash, first_name, last_name, address, phone, is_active, is_superuser, created_at, updated_at FROM users ORDER BY id
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.PasswordHash,
			&i.FirstName,
			&i.LastName,
			&i.Address,
			&i.Phone,
			&i.IsActive,
			&i.IsSuperuser,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const promoteUser = `-- name: PromoteUser :exec
UPDATE users SET
    is_superuser = TRUE
WHERE id = $1
`

func (q *Queries) PromoteUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, promoteUser, id)
	return err
}

const searchUsersByName = `-- name: SearchUsersByName :many
SELECT id, email, password_hash, first_name, last_name, address, phone, is_active, is_superuser, created_at, updated_at FROM users
WHERE CONCAT(first_name, ' ', last_name) ILIKE '%' || $1::text || '%'
ORDER BY id
`

func (q *Queries) SearchUsersByName(ctx context.Context, name string) ([]User, error) {
	rows, err := q.db.Query(ctx, searchUsersByName, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.PasswordHash,
			&i.FirstName,
			&i.LastName,
			&i.Address,
			&i.Phone,
			&i.IsActive,
			&i.IsSuperuser,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE users SET
    first_name = $1,
    last_name = $2,
    email = $3,
    address = $4,
    phone = $5
WHERE id = $6
RETURNING id, email, password_hash, first_name, last_name, address, phone, is_active, is_superuser, created_at, updated_at
`

type UpdateUserParams struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     *string `json:"email"`
	Address   string  `json:"address"`
	Phone     string  `json:"phone"`
	ID        int64   `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Address,
		arg.Phone,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PasswordHash,
		&i.FirstName,
		&i.LastName,
		&i.Address,
		&i.Phone,
		&i.IsActive,
		&i.IsSuperuser,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
