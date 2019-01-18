package db

import (
	"time"

	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Proxy struct {
	ID         int64          `json:"id" db:"id"`
	Name       string         `json:"name" db:"name"`
	SerialKey  string         `json:"serialKey" db:"serial_key"`
	Network    *string        `json:"network" db:"network"`
	MacAddress *string        `json:"macAddress" db:"mac_address"`
	UserID     int64          `json:"userID" db:"user_id"`
	UpdatedAt  time.Time      `json:"updatedAt" db:"updated_at"`
	CreatedAt  time.Time      `json:"createdAt" db:"created_at"`
	DeletedAt  mysql.NullTime `json:"deletedAt" db:"deleted_at"`
}

type IProxyDAO interface {
	GetByMac(serial string) (*Proxy, error)
	GetByID(id int64) (*Proxy, error)
	Insert(m *Proxy) (*Proxy, error)
	Update(m *Proxy) (*Proxy, error)
	SetTx(tx *sqlx.Tx)
}

func NewProxyDAO(conn *sqlx.DB) IProxyDAO {
	return &ProxyDAO{Conn: conn}
}

type ProxyDAO struct {
	Conn *sqlx.DB
	tx   *sqlx.Tx
}

func (d *ProxyDAO) GetByID(id int64) (*Proxy, error) {
	if d.Conn == nil {
		return nil, errors.New("connection is null")
	}
	net := Proxy{}
	err := d.Conn.Get(&net, `SELECT * FROM proxy WHERE id=? AND deleted_at IS NULL`, id)
	if err != nil {
		return nil, err
	}
	return &net, nil

}

func (d *ProxyDAO) GetByMac(mac string) (*Proxy, error) {
	if d.Conn == nil {
		return nil, errors.New("connection is null")
	}
	net := Proxy{}
	err := d.Conn.Get(&net, `SELECT * FROM proxy WHERE mac_address=? AND deleted_at IS NULL`, mac)
	if err != nil {
		return nil, err
	}
	return &net, nil
}

func (d *ProxyDAO) Update(m *Proxy) (*Proxy, error) {
	if d.tx == nil {
		return nil, errors.New("tx is null")
	}
	m.UpdatedAt = time.Now().UTC()
	_, err := d.tx.NamedExec(`UPDATE proxy SET
		name=:name,serial_key=:serial_key,network=:network,user_id=:user_id,updated_at=:updated_at 
		where id = :id`, m)
	return m, err
}

func (d *ProxyDAO) Insert(m *Proxy) (*Proxy, error) {
	if d.tx == nil {
		return nil, errors.New("tx is null")
	}
	now := time.Now().UTC()
	m.CreatedAt = now
	m.UpdatedAt = now
	res, err := d.tx.NamedExec(`INSERT INTO proxy
		(name,serial_key,network,mac_address,user_id,updated_at,created_at)
		VALUES (:name,:serial_key,:network,:mac_address,:user_id,:updated_at,:created_at)`, m)
	if err != nil {
		return nil, err
	}
	m.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (d *ProxyDAO) SetTx(tx *sqlx.Tx) {
	d.tx = tx
}
