package repository

import (
	"film-rental/db"
	"film-rental/model"
	"fmt"
)

const queryColumns = "staff_id, first_name, last_name, address_id, email, store_id, active, username, last_update, picture"

func scanStaffRow(scanner interface {
	Scan(dest ...any) error
}) (*model.Staff, error) {
	var f model.Staff
	err := scanner.Scan(
		&f.StaffId, &f.FirstName, &f.LastName, &f.AddressId,
		&f.Email, &f.StoreId, &f.Active,
		&f.Username, &f.LastUpdate, &f.Picture,
	)
	return &f, err
}
func GetAllStaff(page int, limit int) ([]*model.Staff, int, error) {
	queryStr := fmt.Sprintf(`SELECT %s FROM staff LIMIT %d OFFSET %d`, queryColumns, limit, (page-1)*limit)

	rows, err := db.DB.Query(queryStr)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	rowCount := db.DB.QueryRow("SELECT COUNT (*) FROM staff")

	var totalCount int
	if err := rowCount.Scan(&totalCount); err != nil {

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
			first_name, last_name, address_id, email, store_id, active, username, password, picture
		) VALUES  ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	    RETURNING staff_id
	`

	var lastID int64
	err := DB.QueryRow(query,
		staff.FirstName,
		staff.LastName,
		staff.AddressId,
		staff.Email,
		staff.StoreId,
		staff.Active,
		staff.Username,
		staff.Password,
		staff.Picture,
	).Scan(&lastID)

	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func GetStaff(username string) (*model.Staff, error) {
	query := `SELECT username, password FROM staff WHERE username = $1`
	row := db.DB.QueryRow(query, username)
	var user model.Staff
	if err := row.Scan(&user.Username, &user.Password); err != nil {
		return nil, err
	}
	return &user, nil
}
