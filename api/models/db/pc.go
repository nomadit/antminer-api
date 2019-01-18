package db

import (
	"log"
	"time"

	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Pc is the schema in DynamoDB
type Pc struct {
	ID         int64          `json:"id" db:"id"`
	MacAddress string         `json:"macAddress" db:"mac_address"`
	IP         string         `json:"ip" db:"ip"`
	Name       *string        `json:"name" db:"name"`
	Desc       *string        `json:"desc" db:"desc"`
	Frequency  int            `json:"frequency" db:"frequency"`
	Status     string         `json:"status" db:"status"`
	ProxyID    int64          `json:"proxyID" db:"proxy_id"`
	UserID     int64          `json:"userID" db:"user_id"`
	CreatedAt  time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time      `json:"updatedAt" db:"updated_at"`
	DeletedAt  mysql.NullTime `json:"deletedAt,omitempty" db:"deleted_at"`
}

type PcModify struct {
	IDs []int64 `json:"ids"`
	Name *string `json:"name"`
	CommandParamType
}

type IPcDAO interface {
	GetAllByProxyIDs(netIDs []int64) (*[]Pc, error)
	GetByMac(mac string) (*Pc, error)

	DeleteByIDs(ids *[]int64) error
	Insert(m *Pc) (*Pc, error)
	Update(m *Pc) (*Pc, error)
	SetTx(tx *sqlx.Tx)
	GetByID(id int64) (*Pc, error)
	GetListByIDs(ids *[]int64) (*[]Pc, error)
	UpdateFrequency(ids *[]int64, freq int) error
}

func NewPcDAO(conn *sqlx.DB) IPcDAO {
	return &PcDAO{Conn: conn}
}

// PcDAO is dao for Miner
type PcDAO struct {
	Conn *sqlx.DB
	tx   *sqlx.Tx
}

func (d *PcDAO) UpdateFrequency(ids *[]int64, freq int) error {
	if d.tx == nil {
		return errors.New("tx is null")
	}
	query, args, err := sqlx.In(`UPDATE pc SET frequency=?, updated_at=?
        WHERE id IN (?) AND deleted_at IS NULL`, freq, time.Now().UTC(), *ids)
	if err != nil {
		return err
	}
	_, err = d.tx.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (d *PcDAO) GetListByIDs(ids *[]int64) (*[]Pc, error) {
	if d.Conn == nil {
		return nil, errors.New("connection is null")
	}
	query, args, err := sqlx.In("SELECT * FROM pc WHERE id IN (?) AND deleted_at IS NULL", *ids)
	if err != nil {
		return nil, err
	}
	l := make([]Pc, 0)
	err = d.Conn.Select(&l, query, args...)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (d *PcDAO) GetByID(id int64) (*Pc, error) {
	if d.Conn == nil {
		return nil, errors.New("connection is null")
	}
	m := Pc{}
	if err := d.Conn.Get(&m, "SELECT * FROM pc WHERE id=? AND deleted_at IS NULL", id); err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *PcDAO) GetAllByProxyIDs(netIDs []int64) (*[]Pc, error) {
	if d.Conn == nil {
		return nil, errors.New("connection is null")
	}
	query, args, err := sqlx.In("SELECT * FROM pc WHERE proxy_id IN (?)", netIDs)
	if err != nil {
		return nil, err
	}
	l := make([]Pc, 0)
	err = d.Conn.Select(&l, query, args...)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (d *PcDAO) GetByMac(mac string) (*Pc, error) {
	if d.Conn == nil {
		return nil, errors.New("connection is null")
	}
	m := Pc{}
	if err := d.Conn.Get(&m, "SELECT * FROM pc WHERE mac_address=? AND deleted_at IS NULL", mac); err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *PcDAO) DeleteByIDs(ids *[]int64) error {
	if d.tx == nil {
		return errors.New("tx is null")
	}
	query, args, err := sqlx.In(`UPDATE pc SET deleted_at =?
        WHERE id IN (?) AND deleted_at IS NULL`, time.Now().UTC(), *ids)
	if err != nil {
		return err
	}
	_, err = d.tx.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (d *PcDAO) Insert(m *Pc) (*Pc, error) {
	if d.tx == nil {
		return nil, errors.New("tx is null")
	}
	now := time.Now().UTC()
	m.CreatedAt = now
	m.UpdatedAt = now
	log.Println(*m)
	res, err := d.tx.NamedExec(`INSERT INTO pc
		(mac_address, ip, name, `+"`desc`"+`, frequency, status, proxy_id, user_id, created_at, updated_at) VALUES 
		(:mac_address,:ip,:name,:desc,:frequency,:status,:proxy_id,:user_id,:created_at,:updated_at)`, m)
	if err != nil {
		return nil, err
	}
	m.ID, err = res.LastInsertId()
	return m, nil
}

func (d *PcDAO) Update(m *Pc) (*Pc, error) {
	if d.tx == nil {
		return nil, errors.New("tx is null")
	}
	now := time.Now().UTC()
	m.UpdatedAt = now
	_, err := d.tx.NamedExec(
		`UPDATE pc SET
        mac_address=:mac_address,ip=:ip,name=:name,`+"`desc`"+`=:desc,frequency=:frequency,status=:status,
        proxy_id=:proxy_id,user_id=:user_id,updated_at=:updated_at
        WHERE id=:id AND deleted_at IS NULL`, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (d *PcDAO) SetTx(tx *sqlx.Tx) {
	d.tx = tx
}
