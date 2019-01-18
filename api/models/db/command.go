package db

import (
	"time"

	"errors"
	"github.com/jmoiron/sqlx"
)

type Command struct {
	ID        int64     `json:"id" db:"id"`
	PcID      int64     `json:"pcID" db:"pc_id"`
	Status    string    `json:"status" db:"status"`
	Type      string    `json:"type" db:"type"`
	Param     string    `json:"param" db:"param"`
	UserID    int64     `json:"userID" db:"user_id"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CommandParamType struct {
	Freq int `json:"freq"`
	Host string `json:"hostname"`
	Conf *[]PoolConfig `json:"conf"`
}

type ICommandDAO interface {
	Insert(m *Command) (*Command, error)
	SetTx(tx *sqlx.Tx)
	GetCommandListByPcIDs(pcIDs *[]int64) (*[]Command, error)
	UpdateStatusByID(id int64, status string) error
}

func NewCommandDAO(conn *sqlx.DB) ICommandDAO {
	return &CommandDAO{Conn: conn}
}

type CommandDAO struct {
	Conn *sqlx.DB
	tx   *sqlx.Tx
}

func (d *CommandDAO) UpdateStatusByID(id int64, status string) error {
	if d.tx == nil {
		return errors.New("connection is null")
	}
	query, args, err := sqlx.In(`UPDATE command SET status=?, updated_at=?
        WHERE id = ?`, status, time.Now().UTC(), id)
	if err != nil {
		return err
	}
	_, err = d.tx.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (d *CommandDAO) GetCommandListByPcIDs(pcIDs *[]int64) (*[]Command, error) {
	query, args, err := sqlx.In("SELECT * FROM command WHERE pc_id IN (?) AND status = 'RUN'", *pcIDs)
	if err != nil {
		return nil, err
	}
	var l []Command
	err = d.Conn.Select(&l, query, args...)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (d *CommandDAO) Insert(m *Command) (*Command, error) {
	if d.tx == nil {
		return nil, errors.New("tx is null")
	}
	now := time.Now().UTC()
	m.CreatedAt = now
	m.UpdatedAt = now
	res, err := d.tx.NamedExec(`INSERT INTO command
		(pc_id, status, type, param, user_id, created_at, updated_at) VALUES 
		(:pc_id, :status, :type, :param, :user_id, :created_at, :updated_at)`, m)
	if err != nil {
		return nil, err
	}
	m.ID, err = res.LastInsertId()
	return m, nil
}

func (d *CommandDAO) SetTx(tx *sqlx.Tx) {
	d.tx = tx
}
