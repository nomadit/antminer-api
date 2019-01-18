package db

import (
	"time"

	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Pool is the schema in DynamoDB
type PoolConfig struct {
	ID        int64          `json:"id" db:"id"`
	URL       string         `json:"url" db:"url"`
	Wallet    string         `json:"wallet" db:"wallet"`
	Password  *string        `json:"password" db:"password"`
	Md5       string         `json:"md5" db:"md5"`
	PcID      int64          `json:"pcID" db:"pc_id"`
	UserID    int64          `json:"user_id" db:"user_id"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	DeletedAt mysql.NullTime `json:"deleted_at" db:"deleted_at"`
}

type IPoolConfigDAO interface {
	GetListByPcIDs(pcIDs *[]int64) (*[]PoolConfig, error)
	DeleteByPcIDs(ids *[]int64) error
	DeleteByIDs(ids *[]int64) error
	Insert(m *PoolConfig) (*PoolConfig, error)
	Update(m *PoolConfig) (*PoolConfig, error)
	SetTx(tx *sqlx.Tx)
}

func NewPoolConfig(conn *sqlx.DB) IPoolConfigDAO {
	return &PoolConfigDAO{Conn: conn}
}

// PoolConfigDAO is dao for Wallet
type PoolConfigDAO struct {
	Conn *sqlx.DB
	tx   *sqlx.Tx
}

func (d *PoolConfigDAO) GetListByPcIDs(pcIDs *[]int64) (*[]PoolConfig, error) {
	if d.Conn == nil {
		return nil, errors.New("connection is null")
	}
	query, args, err := sqlx.In("SELECT * FROM pool_config WHERE pc_id IN (?) and deleted_at is null", *pcIDs)
	if err != nil {
		return nil, err
	}
	var list []PoolConfig
	err = d.Conn.Select(&list, query, args...)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (d *PoolConfigDAO) DeleteByPcIDs(ids *[]int64) error {
	if d.tx == nil {
		return errors.New("connection is null")
	}
	now := time.Now().UTC()
	query, args, err := sqlx.In(`UPDATE pool_config SET deleted_at=?
        WHERE pc_id IN (?) AND deleted_at IS NULL`, now, *ids)
	if err != nil {
		return err
	}
	_, err = d.tx.Exec(query, args...)
	return err
}

func (d *PoolConfigDAO) DeleteByIDs(ids *[]int64) error {
	if d.tx == nil {
		return errors.New("connection is null")
	}
	now := time.Now().UTC()
	query, args, err := sqlx.In(`UPDATE pool_config SET deleted_at=?
        WHERE id IN (?) AND deleted_at IS NULL`, now, *ids)
	if err != nil {
		return err
	}
	_, err = d.tx.Exec(query, args...)
	return err
}

func (d *PoolConfigDAO) Update(m *PoolConfig) (*PoolConfig, error) {
	if d.tx == nil {
		return nil, errors.New("tx is null")
	}
	m.UpdatedAt = time.Now().UTC()
	_, err := d.tx.NamedExec(`UPDATE pool_config SET
		url=:url,wallet=:wallet,password=:password,md5=:md5,pc_id=:pc_id,user_id=:user_id,updated_at=:updated_at`, m)
	return m, err
}

func (d *PoolConfigDAO) Insert(m *PoolConfig) (*PoolConfig, error) {
	if d.tx == nil {
		return nil, errors.New("tx is null")
	}
	now := time.Now().UTC()
	m.CreatedAt = now
	m.UpdatedAt = now
	res, err := d.tx.NamedExec(`INSERT INTO pool_config
		(url,wallet,password,md5,pc_id,user_id,updated_at,created_at)
		VALUES (:url,:wallet,:password,:md5,:pc_id,:user_id,:updated_at,:created_at)`, m)
	if err != nil {
		return nil, err
	}
	m.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (d *PoolConfigDAO) SetTx(tx *sqlx.Tx) {
	d.tx = tx
}
