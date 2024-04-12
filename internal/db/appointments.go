package db

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Appointment struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	Phone     string    `db:"phone"`
	Date      time.Time `db:"date"`
	CreatedAt time.Time `db:"created_at"`
	Attended  bool      `db:"attended"`
}

func FindAppointmentById(id string) (*Appointment, error) {
	conn, err := GetPool()
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var appointment Appointment

	row, err := conn.Query(
		ctx,
		"SELECT * FROM appointments WHERE id = $1 LIMIT 1",
		id,
	)

	if err != nil {
		return nil, err
	}

	appointment, err = pgx.CollectOneRow[Appointment](row, pgx.RowToStructByName[Appointment])

	if err != nil {
		return nil, err
	}

	return &appointment, nil
}

func CreateAppointment(data *Appointment) (id uuid.UUID, err error) {
	conn, err := GetPool()
	if err != nil {
		return
	}
	defer conn.Release()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	uid, err := uuid.NewV7()

	if err != nil {
		return
	}

	tag, err := conn.Exec(
		ctx,
		"INSERT INTO appointments (id, name, phone, date) VALUES ($1, $2, $3, $4)",
		id,
		data.Name,
		data.Phone,
		data.Date,
	)

	if err != nil {
		return
	}

	if tag.RowsAffected() == 0 {
		err = errors.New("No se pudo crear la cita")
		return
	}

	id = uid
	return
}

func DeleteAppointmentById(id string) error {
	conn, err := GetPool()
	if err != nil {
		return err
	}
	defer conn.Release()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := conn.Exec(
		ctx,
		"DELETE FROM appointments WHERE id = $1",
		id,
	)

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("No se encontro cita con el id especificado. ID:" + id)
	}

	return nil
}
