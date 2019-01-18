package ctrls

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/nomadit/antminer-api/api/models/db"
	"github.com/nomadit/antminer-api/api/util"
	"io/ioutil"
	"log"
	"net/http"
)

func NewProxyCtrl() *ProxyCtrl {
	conn := db.GetDB()
	return &ProxyCtrl{
		conn:     conn,
		pcDAO:    db.NewPcDAO(conn),
		proxyDAO: db.NewProxyDAO(conn),
	}
}

type ProxyCtrl struct {
	conn     *sqlx.DB
	pcDAO    db.IPcDAO
	proxyDAO db.IProxyDAO
}

func (cr *ProxyCtrl) ProxyByMac(c *gin.Context) {
	mac := c.Param("mac")
	item, err := cr.proxyDAO.GetByMac(mac)
	if err != nil {
		log.Println("ERROR " + err.Error())
		util.CheckErrorInHTTP(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (cr *ProxyCtrl) UpdateProxy(c *gin.Context) {
	jsonBytes, err := ioutil.ReadAll(c.Request.Body)
	if util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	var req db.Proxy
	if err = json.Unmarshal(jsonBytes, &req); util.CheckErrorInHTTP(c, err) != nil {
		return
	}

	if net, err := cr.proxyDAO.GetByID(req.ID); err != nil {
		log.Println("ERROR " + err.Error())
		util.CheckErrorInHTTP(c, err)
	} else {
		net.Name = req.Name
		net.Network = req.Network
		tx := cr.conn.MustBegin()
		cr.proxyDAO.SetTx(tx)
		ret, err := cr.proxyDAO.Update(net)
		if err != nil {
			tx.Rollback()
			util.CheckErrorInHTTP(c, err)
			return
		}
		tx.Commit()
		c.JSON(http.StatusOK, ret)
	}
}

func (cr *ProxyCtrl) CreateProxy(c *gin.Context) {
	jsonBytes, err := ioutil.ReadAll(c.Request.Body)
	if util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	var req db.Proxy
	if err = json.Unmarshal(jsonBytes, &req); util.CheckErrorInHTTP(c, err) != nil {
		return
	}

	tx := cr.conn.MustBegin()
	cr.proxyDAO.SetTx(tx)
	ret, err := cr.proxyDAO.Insert(&req)
	if err != nil {
		tx.Rollback()
		util.CheckErrorInHTTP(c, err)
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, ret)
}

