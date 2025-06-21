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
		return 0, fmt.Errorf("can't add new datas: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("can't find last insert id: %w", err)
	}
	return int(id), nil
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p

	// верните идентификатор последней добавленной записи
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow("SELECT Number,Client,Status,Address,created_at from parcel WHERE number = :number", sql.Named("number", number))

	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, fmt.Errorf("can't Scans data:%w", err)
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	row, err := s.db.Query("SELECT Number,Client,Status,Address,created_at FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, fmt.Errorf("can't find client:%d , %w", client, err)
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
			return nil, fmt.Errorf("can't scan data %w", err)
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
		return fmt.Errorf("can't update status:%s from number:%d ,%w", status, number, err)
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	row := s.db.QueryRow("SELECT Status FROM parcel WHERE number =:number", sql.Named("number", number))

	var currentStatus string
	err := row.Scan(&currentStatus)
	if err != nil {
		return fmt.Errorf("can't scan status for number:%d, %w", number, err)
	}
	if currentStatus != ParcelStatusRegistered {
		return fmt.Errorf("can't change Address for number:%d status is REGISTERED, %w", number, err)
	} else {
		_, err = s.db.Exec("UPDATE parcel SET Address =:address WHERE Number =:number",
			sql.Named("address", address),
			sql.Named("number", number))
		if err != nil {
			return fmt.Errorf("can't change Address for number:%d ,%w", number, err)
		}

	}
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	row := s.db.QueryRow("SELECT Status FROM parcel WHERE Number =:number", sql.Named("number", number))
	var currentStatus string
	err := row.Scan(&currentStatus)
	if err == sql.ErrNoRows {
		return fmt.Errorf("can't find parcel for number:%d, %w", number, err)
	}
	if currentStatus != ParcelStatusRegistered {
		return fmt.Errorf("can't delete row for number:%d status is %s", number, currentStatus)
	} else {
		_, err = s.db.Exec("DELETE FROM parcel WHERE Number =:number", sql.Named("number", number))
		if err != nil {
			return fmt.Errorf("can't delete row for number:%d ,%w", number, err)
		}

	}
	return nil
}
