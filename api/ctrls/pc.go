package ctrls

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/nomadit/antminer-api/api/models/db"
	"github.com/nomadit/antminer-api/api/util"
	"io/ioutil"
	"net/http"
	"strconv"
)

func NewPcCtrl() *PcCtrl {
	conn := db.GetDB()
	return &PcCtrl{
		conn:          conn,
		pcDAO:         db.NewPcDAO(conn),
		proxyDAO:      db.NewProxyDAO(conn),
		poolConfigDAO: db.NewPoolConfig(conn),
		commandDAO:    db.NewCommandDAO(conn),
	}
}

type PcCtrl struct {
	conn          *sqlx.DB
	pcDAO         db.IPcDAO
	proxyDAO      db.IProxyDAO
	poolConfigDAO db.IPoolConfigDAO
	commandDAO    db.ICommandDAO
}

func (cr *PcCtrl) PcListByNetIDIncludeDeleted(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("proxyID"), 10, 64)
	if util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	if pcList, err := cr.pcDAO.GetAllByProxyIDs([]int64{id}); err != nil {
		fmt.Println("ERROR " + err.Error())
		util.CheckErrorInHTTP(c, err)
	} else {
		c.JSON(http.StatusOK, pcList)
	}
}

func (cr *PcCtrl) UpsertPC(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("proxyID"), 10, 64)
	if util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	jsonBytes, err := ioutil.ReadAll(c.Request.Body)
	if util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	var req db.Pc
	if err = json.Unmarshal(jsonBytes, &req); util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	req.ProxyID = id

	net, err := cr.proxyDAO.GetByID(req.ProxyID)
	if util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	req.UserID = net.UserID
	item, err := cr.pcDAO.GetByMac(req.MacAddress)
	tx := cr.conn.MustBegin()
	cr.pcDAO.SetTx(tx)
	var ret *db.Pc
	if err != nil {
		if err == sql.ErrNoRows {
			req.Status = db.RunStatus
			if ret, err = cr.pcDAO.Insert(&req); err != nil {
				fmt.Println("ERROR " + err.Error())
				util.CheckErrorInHTTP(c, err)
				tx.Rollback()
				return
			}
		} else {
			fmt.Println("ERROR " + err.Error())
			util.CheckErrorInHTTP(c, err)
			tx.Rollback()
			return
		}
	} else {
		item.IP = req.IP
		item.ProxyID = req.ProxyID
		if ret, err = cr.pcDAO.Update(item); err != nil {
			fmt.Println("ERROR " + err.Error())
			util.CheckErrorInHTTP(c, err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	c.JSON(http.StatusOK, ret)
}

func (cr *PcCtrl) UpdatePcStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	jsonBytes, err := ioutil.ReadAll(c.Request.Body)
	if util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	var req db.Pc
	if err = json.Unmarshal(jsonBytes, &req); util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	item, err := cr.pcDAO.GetByID(id)
	if err != nil {
		util.CheckErrorInHTTP(c, err)
		return
	}
	var ret *db.Pc
	if item.Status != req.Status {
		item.Status = req.Status
		tx := cr.conn.MustBegin()
		cr.pcDAO.SetTx(tx)
		if ret, err = cr.pcDAO.Update(item); err != nil {
			fmt.Println("ERROR " + err.Error())
			util.CheckErrorInHTTP(c, err)
			tx.Rollback()
			return
		}
		tx.Commit()
	}
	c.JSON(http.StatusOK, ret)
}

func (cr *PcCtrl) UpdatePcName(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	jsonBytes, err := ioutil.ReadAll(c.Request.Body)
	if util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	var req db.Pc
	if err = json.Unmarshal(jsonBytes, &req); util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	item, err := cr.pcDAO.GetByID(id)
	if err != nil {
		util.CheckErrorInHTTP(c, err)
		return
	}
	var ret *db.Pc
	if item.Name == nil || *item.Name != *req.Name {
		item.Name = req.Name
		tx := cr.conn.MustBegin()
		cr.pcDAO.SetTx(tx)
		if ret, err = cr.pcDAO.Update(item); err != nil {
			fmt.Println("ERROR " + err.Error())
			util.CheckErrorInHTTP(c, err)
			tx.Rollback()
			return
		}
		tx.Commit()
	}
	c.JSON(http.StatusOK, ret)
}

func (cr *PcCtrl) UpdateInvalidByMac(c *gin.Context) {
	mac := c.Param("mac")
	item, err := cr.pcDAO.GetByMac(mac)
	if err != nil {
		util.CheckErrorInHTTP(c, err)
		return
	}
	if item.Status == "STOP" {
		c.JSON(http.StatusOK, item)
		return
	}
	item.Status = db.ErrorNoWorker
	tx := cr.conn.MustBegin()
	cr.pcDAO.SetTx(tx)
	var ret *db.Pc
	if ret, err = cr.pcDAO.Update(item); err != nil {
		fmt.Println("ERROR " + err.Error())
		util.CheckErrorInHTTP(c, err)
		tx.Rollback()
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, ret)
}

func (cr *PcCtrl) UpsertPcFrequencyConfig(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	jsonBytes, err := ioutil.ReadAll(c.Request.Body)
	if util.CheckErrorInHTTP(c, err) != nil {
		return
	}
	type request struct {
		Freq int                 `json:"freq,string"`
		Conf []map[string]string `json:"conf"`
	}
	var req request
	if err = json.Unmarshal(jsonBytes, &req); util.CheckErrorInHTTP(c, err) != nil {
		return
	}

	var pc *db.Pc
	if pc, err = cr.pcDAO.GetByID(id); err != nil {
		util.CheckErrorInHTTP(c, err)
		return
	}

	tx := cr.conn.MustBegin()
	cr.pcDAO.SetTx(tx)
	pc.Frequency = req.Freq
	if _, err := cr.pcDAO.Update(pc); err != nil {
		fmt.Println("ERROR pc.Update:" + err.Error())
		util.CheckErrorInHTTP(c, err)
		tx.Rollback()
		return
	}
	confList, err := cr.poolConfigDAO.GetListByPcIDs(&[]int64{pc.ID})
	if err != nil {
		fmt.Println("ERROR poolConfigDAO.GetListByPcIDs:" + err.Error())
		util.CheckErrorInHTTP(c, err)
		tx.Rollback()
		return
	}
	hashIDMap := map[string][]int64{}
	for _, conf := range *confList {
		if _, ok := hashIDMap[conf.Md5]; ok {
			hashIDMap[conf.Md5] = append(hashIDMap[conf.Md5], conf.ID)
		} else {
			hashIDMap[conf.Md5] = []int64{conf.ID}
		}
	}
	cr.poolConfigDAO.SetTx(tx)
	existMap := map[string]bool{}
	for _, raw := range req.Conf {
		password := raw["password"]
		hashBytes := md5.Sum([]byte(raw["url"] + raw["user"] + password))
		hash := fmt.Sprintf("%x", hashBytes)
		if _, ok := hashIDMap[hash]; !ok {
			if _, ok := existMap[hash]; !ok {
				conf := &db.PoolConfig{
					URL:      raw["url"],
					Wallet:   raw["user"],
					Password: &password,
					PcID:     pc.ID,
					UserID:   pc.UserID,
					Md5:      hash,
				}

				if _, err := cr.poolConfigDAO.Insert(conf); err != nil {
					fmt.Println("ERROR poolConfigDAO.Insert:" + err.Error())
					util.CheckErrorInHTTP(c, err)
					tx.Rollback()
					return
				}
			}
		}
		existMap[hash] = true
	}
	var deleteIDs []int64
	for key, idList := range hashIDMap {
		if _, ok := existMap[key]; !ok {
			deleteIDs = append(deleteIDs, idList[0])
		}
		if len(idList) > 1 {
			deleteIDs = append(deleteIDs, idList[0])
		}
	}
	if len(deleteIDs) > 0 {
		if err = cr.poolConfigDAO.DeleteByIDs(&deleteIDs); err != nil {
			fmt.Println("ERROR " + err.Error())
			util.CheckErrorInHTTP(c, err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	c.JSON(http.StatusOK, nil)
}

func (cr *PcCtrl) updatePcConfigList(userID int64, req *db.PcModify) error {
	pcList, err := cr.pcDAO.GetListByIDs(&req.IDs)
	if err != nil {
		return err
	}
	changedIDSet := map[int64]bool{}
	pcIDMap := map[int64]*db.Pc{}
	var pcIDs []int64
	for _, pc := range *pcList {
		if pc.Frequency != req.Freq {
			pc.Frequency = req.Freq
			changedIDSet[pc.ID] = true
			pcIDs = append(pcIDs, pc.ID)
		}
		pcIDMap[pc.ID] = &pc
	}
	if len(pcIDs) > 0 {
		if err := cr.pcDAO.UpdateFrequency(&pcIDs, req.Freq); err != nil {
			return err
		}
	}

	confList, err := cr.poolConfigDAO.GetListByPcIDs(&req.IDs)
	if err != nil {
		return err
	}
	pcIDConfsMap := map[int64][]*db.PoolConfig{}
	for i, conf := range *confList {
		if _, ok := pcIDConfsMap[conf.PcID]; ok {
			pcIDConfsMap[conf.PcID] = append(pcIDConfsMap[conf.PcID], &(*confList)[i])
		} else {
			pcIDConfsMap[conf.PcID] = []*db.PoolConfig{&(*confList)[i]}
		}
	}
	var reqList []*db.PoolConfig
	for idx, reqConf := range *req.Conf {
		hashBytes := md5.Sum([]byte(reqConf.URL + reqConf.Wallet + *reqConf.Password))
		hash := fmt.Sprintf("%x", hashBytes)
		(*req.Conf)[idx].Md5 = hash
		(*req.Conf)[idx].UserID = userID
		reqList = append(reqList, &(*req.Conf)[idx])
	}

	// TODO if the configures for the pool are changed, delete all the existed configures and
	// inert all the new configure. Because configures don't have ordered indexes. it don't know
	// which were changed configures.
	// So if you want to change the logic, add the ordered index in configure table,
	// retain the three configure for each pc
	confChangedIDSet := map[int64]bool{}
	for pcID, confs := range pcIDConfsMap {
		if len(reqList) != len(confs) {
			confChangedIDSet[pcID] = true
			changedIDSet[pcID] = true
		} else {
			var deleteIds []int64
			for idx, req := range reqList {
				if req.Md5 != confs[idx].Md5 {
					deleteIds = append(deleteIds, confs[idx].ID)
				}
			}
			if len(deleteIds) > 0 {
				confChangedIDSet[pcID] = true
				changedIDSet[pcID] = true
			}
		}
	}
	for pcID := range confChangedIDSet {
		if err := cr.poolConfigDAO.DeleteByPcIDs(&[]int64{pcID}); err != nil {
			return err
		}
		for _, conf := range *req.Conf {
			conf.PcID = pcID
			cr.poolConfigDAO.Insert(&conf)
		}
	}
	for pcID := range changedIDSet {
		param, err := json.Marshal(&db.CommandParamType{
			Freq: req.Freq,
			Conf: req.Conf,
		})
		if err != nil {
			return err
		}
		command := db.Command{
			PcID:   pcID,
			Status: db.RunStatus,
			Type:   db.CommandChangeConfig,
			Param:  string(param),
			UserID: userID,
		}
		if _, err := cr.commandDAO.Insert(&command); err != nil {
			return err
		}
	}
	return nil
}

