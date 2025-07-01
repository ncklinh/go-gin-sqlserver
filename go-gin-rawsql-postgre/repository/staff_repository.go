package repository

import (
	"film-rental/db"
	"film-rental/model"
)

const queryColumns = "staff_id, first_name, last_name, address_id, email, store_id, active, username, role, last_update, picture"

func scanStaffRow(scanner interface {
	Scan(dest ...any) error
}) (*model.Staff, error) {
	var f model.Staff
	err := scanner.Scan(
		&f.StaffId, &f.FirstName, &f.LastName, &f.AddressId,
		&f.Email, &f.StoreId, &f.Active,
		&f.Username, &f.Role, &f.LastUpdate, &f.Picture,
	)
	return &f, err
}

func GetAllStaff(page int, limit int) ([]*model.Staff, int, error) {
	queryStr := `SELECT ` + queryColumns + ` FROM staff LIMIT $1 OFFSET $2`

	rows, err := db.DB.Query(queryStr, limit, (page-1)*limit)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	rowCount := db.DB.QueryRow("SELECT COUNT (*) FROM staff")

	var totalCount int
	if err := rowCount.Scan(&totalCount); err != nil {
		return nil, 0, err
	}

	var staffs []*model.Staff
	for rows.Next() {
		if f, err := scanStaffRow(rows); err != nil {
			continue
		} else {
			staffs = append(staffs, f)

		}
	}

	return staffs, totalCount, nil
}

func InsertStaff(staff model.Staff) (int64, error) {
	query := `
		INSERT INTO staff (
			first_name, last_name, address_id, email, store_id, active, username, password, role, picture
		) VALUES  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	    RETURNING staff_id
	`

	var lastID int64
	err := db.DB.QueryRow(query,
		staff.FirstName,
		staff.LastName,
		staff.AddressId,
		staff.Email,
		staff.StoreId,
		staff.Active,
		staff.Username,
		staff.Password,
		staff.Role,
		staff.Picture,
	).Scan(&lastID)

	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func GetStaff(username string) (*model.Staff, error) {
	query := `SELECT username, password, role FROM staff WHERE username = $1`
	row := db.DB.QueryRow(query, username)
	var user model.Staff
	if err := row.Scan(&user.Username, &user.Password, &user.Role); err != nil {
		return nil, err
	}
	return &user, nil
}

func IsUsernameExists(username string) (bool, error) {
	query := `SELECT COUNT(*) FROM staff WHERE username = $1`
	var count int
	err := db.DB.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
