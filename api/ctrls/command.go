package ctrls

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/nomadit/antminer-api/api/models/db"
	"github.com/nomadit/antminer-api/api/util"
	"io/ioutil"
	"net/http"
	"strconv"
)

func NewCommandCtrl() *CommandCtrl {
	conn := db.GetDB()
	return &CommandCtrl{
		conn:       conn,
		commandDAO: db.NewCommandDAO(conn),
	}
}

type CommandCtrl struct {
	conn       *sqlx.DB
	commandDAO db.ICommandDAO
}

func (cr *CommandCtrl) CommandListByPcIDs(c *gin.Context) {
	idsStr := c.QueryArray("id[]")
	if len(idsStr) == 0 {
		util.CheckErrorInHTTP(c, errors.New("query id list is 0"))
		return
	}
	idsList, err := util.SliceAtoi64(idsStr)
	if err != nil {
		util.CheckErrorInHTTP(c, err)
		return
	}
	if list, err := cr.commandDAO.GetCommandListByPcIDs(&idsList); err != nil {
		util.CheckErrorInHTTP(c, err)
		return
	} else {
		c.JSON(http.StatusOK, list)
	}
}

func (cr *CommandCtrl) UpdateCommandStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	jsonBytes, err := ioutil.ReadAll(c.Request.Body)
	if util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	var req map[string]string
	if err = json.Unmarshal(jsonBytes, &req); util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	tx := cr.conn.MustBegin()
	cr.commandDAO.SetTx(tx)
	if err := cr.commandDAO.UpdateStatusByID(id, req["status"]); err != nil {
		tx.Rollback()
		util.CheckErrorInHTTP(c, err)
	}
	tx.Commit()
	c.JSON(http.StatusOK, nil)
}
