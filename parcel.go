package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address,created_at) VALUES (:client, :status,:address,:created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p

	// верните идентификатор последней добавленной записи
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow("SELECT number,client,status,address,created_at FROM parcel WHERE number = :number", sql.Named("number", number))

	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	row, err := s.db.Query("SELECT number,client,status,address,created_at FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer row.Close()

	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for row.Next() {
		var par Parcel
		err := row.Scan(&par.Number, &par.Client, &par.Status, &par.Address, &par.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, par)
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status=:status WHERE number =:number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	row := s.db.QueryRow("SELECT status FROM parcel WHERE number =:number", sql.Named("number", number))

	var currentStatus string
	err := row.Scan(&currentStatus)
	if err != nil {
		return err
	}
	if currentStatus != ParcelStatusRegistered {
		return err
	} else {
		_, err = s.db.Exec("UPDATE parcel SET address =:address WHERE number =:number",
			sql.Named("address", address),
			sql.Named("number", number))
		if err != nil {
			return err
		}

	}
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	row := s.db.QueryRow("SELECT status FROM parcel WHERE number =:number", sql.Named("number", number))
	var currentStatus string
	err := row.Scan(&currentStatus)
	if err == sql.ErrNoRows {
		return err
	}
	if currentStatus != ParcelStatusRegistered {
		return fmt.Errorf("can't delete row for number:%d status is %s", number, currentStatus)
	} else {
		_, err = s.db.Exec("DELETE FROM parcel WHERE Number =:number", sql.Named("number", number))
		if err != nil {
			return err
		}

	}
	return nil
}
