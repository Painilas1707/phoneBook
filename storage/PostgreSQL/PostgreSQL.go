package PostgreSQL

import (
	"Study/Demo/internal/StructUser"
	"Study/Demo/storage"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dbPath string) (*Storage, error) {
	const fn = "Storage.postgreSQL.NewStorage"

	db, err := sql.Open("pgx", dbPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	query := (`
	CREATE TABLE IF NOT EXISTS phonebook(
	id SERIAL PRIMARY KEY,
	contact_fio VARCHAR(50) NOT NULL,
	birth_date DATE NOT NULL,
	phone_number VARCHAR(20) NOT NULL,
	email VARCHAR(50) NOT NULL UNIQUE,
	time_create TIMESTAMP DEFAULT NOW()
	)`)

	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return &Storage{db: db}, nil
}
func (s *Storage) SaveContact(contact StructUser.UserPhoneBook) (int64, error) {
	const fn = "Storage.postgreSQL.saveContact"
	var id int64
	err := s.db.QueryRow("INSERT INTO phonebook(contact_fio, birth_date, phone_number, email) VALUES($1, $2, $3, $4)RETURNING id, time_create",
		contact.ContactFIO, contact.BirthDate, contact.PhoneNumber, contact.Email).Scan(&id, &contact.TimeCreate)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, fmt.Errorf("%s: %w", fn, storage.ErrUserAlreadyExists)
			}
		}

		return 0, fmt.Errorf("%s: %w", fn, err)
	}
	return id, nil
}
func (s *Storage) GetContact(searchQuery string) (StructUser.UserPhoneBook, error) {
	const fn = "Storage.postgreSQL.GetContact"
	var contact StructUser.UserPhoneBook
	err := s.db.QueryRow("SELECT id, contact_fio, birth_date, phone_number, email, time_create FROM phonebook WHERE contact_fio = $1 OR phone_number = $1 OR email = $1", searchQuery).Scan(&contact.ID, &contact.ContactFIO, &contact.BirthDate, &contact.PhoneNumber, &contact.Email, &contact.TimeCreate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return contact, storage.ErrUserNotFound
		}
		return contact, fmt.Errorf("%s: %w", fn, err)
	}
	return contact, nil
}

func (s *Storage) DeleteContact(id int64) error {
	const fn = "Storage.postgreSQL.DeleteContact"
	res, err := s.db.Exec("DELETE FROM phonebook where id =$1", id)
	if err != nil {
		return fmt.Errorf("%s %w", fn, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s %w", fn, err)
	}
	if rowsAffected == 0 {
		return storage.ErrUserNotFound
	}
	return nil
}
