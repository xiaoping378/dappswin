package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"dappswin/models"

	"github.com/golang/glog"
)

const (
	blockType = iota
	txType
	gameType
	winType
)

var cachedHeadNum uint32

type InfoRsp struct {
	Num uint32 `json:"head_block_num"`
}

type BlockRsp struct {
	Hash string                   `json:"id"`
	Num  uint32                   `json:"block_num"`
	Time string                   `json:"timestamp"`
	Txs  []models.TransactionResp `json:"transactions"`
}

func getBlockByNum(num uint32) (*BlockRsp, error) {
	params := fmt.Sprintf(`{"block_num_or_id": %d}`, num)
	url := eosConf.RPCURL + "/v1/chain/get_block"
	resp, err := http.Post(url, "application/json", strings.NewReader(params))
	if nil != err {
		glog.Errorf("getBlockByNum - http.Post(%s) with params %s failed : %v", url, params, err)
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		glog.Errorf("getBlockByNum - ioutil.ReadAll failed : %v", err)
		return nil, err
	}
	blk := &BlockRsp{}
	if err = json.Unmarshal(buf, blk); nil != err {
		glog.Errorf("getBlockByNum - json.Unmarshall failed : %v", err)
		return nil, err
	}

	return blk, nil
}

func resolveBlock(num uint32) {
	glog.V(7).Infof("resolving Block num: %d", num)
	blkRsp, err := getBlockByNum(num)
	if err != nil {
		// TODO, notify chan to recheck it.
		glog.Errorf("resolve block num %d : err %v", num, err)
	}

	tm, _ := time.Parse("2006-01-02T15:04:05.999999999", blkRsp.Time)
	timemills := tm.UnixNano() / 1e6
	// glog.Info("timemills is %d, %s", timemills, blkRsp.Time)
	txschan <- &models.Message{Type: txType, BlockNum: num, Hash: "", Time: timemills, Data: blkRsp.Txs}

	blk := models.Block{blkRsp.Hash, blkRsp.Num, timemills}

	// 游戏轮数需要统计
	glog.Infof("Pushing Game needed block... %#v", blk)
	gameChan <- &blk

	// 广播区块信息
	msg := blk.Message()
	Huber.broadcast <- msg

	if err := models.AddBlock(&blk); err != nil {
		glog.Error("save Block error", err)
	}
}

func getHeadNum() uint32 {
	url := eosConf.RPCURL + "/v1/chain/get_info"
	resp, err := http.Post(url, "application/json", nil)
	if nil != err {
		glog.Errorf("getBlockByNum - http.Post(%s) failed : %v", url, err)
		return 0
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		glog.Errorf("getBlockByNum - ioutil.ReadAll failed : %v", err)
		return 0
	}
	info := &InfoRsp{}
	if err = json.Unmarshal(buf, info); nil != err {
		glog.Errorf("getBlockByNum - json.Unmarshall failed : %v", err)
		return 0
	}

	return info.Num
}

func ResolveRoutine() {
	ticker := time.NewTicker(time.Duration(eosConf.FetchIdleDur) * time.Millisecond)
	defer func() {
		ticker.Stop()
	}()

	blk, err := models.GetLastestBlock()
	glog.Warningf("get Latest block is %v", blk)
	if err != nil {
		glog.Errorf("get Latest Block error %s", err.Error())
	}
	if blk != nil && blk.Num > eosConf.FromBlkNum {
		cachedHeadNum = blk.Num
	} else {
		cachedHeadNum = eosConf.FromBlkNum
	}

	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			head := getHeadNum()
			glog.Infof("head num is %d, cached is %d", head, cachedHeadNum)

			for i := cachedHeadNum + 1; i <= head; i++ {
				resolveBlock(i)
				cachedHeadNum = i
			}

			ticker = time.NewTicker(time.Duration(eosConf.FetchIdleDur) * time.Millisecond)
		}
	}
}
