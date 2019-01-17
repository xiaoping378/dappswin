package app

import (
	"bytes"
	"dappswin/conf"
	"dappswin/database"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

var eosConf *EosConf
var db *gorm.DB
var apiEndpoint []string

// InitEos 启动Eos Resolver
func Init() {
	eosConf = newEosConf()
	db = database.Db.Debug()
	go gameRoutine()
	go resloveTXRoutine()
	go votedRoutine()
	go checkWinRoutine()
	if eosConf.EnableICO {
		go checkICORoutine()
	}
	go Huber.run()
	checkcleosExists()
	go Forever(updateICOBalance, time.Second*10)
	go Forever(updateTotalLocked, 1*time.Minute)
	go Forever(updateCirculat, 1*time.Minute)
	// cancel paltform
	// go checkNotifyRoutine()
}

// EosConf :
type EosConf struct {
	RPCURL              string
	ChainID             string
	FetchIdleDur        int    // 查询blk时间间隔
	FromBlkNum          uint32 // 从哪个blocknum开始查询
	GameAccount         string
	ICOAccount          string
	EnableICO           bool
	ICOStartTime        int64
	TokenSymbol         string
	TokenAccount        string
	EOS_CGG             float64
	WalletURL           string
	WalletPW            string
	TotalAmount         float64
	LockAccount         string
	OfficialLockAccount string
	TotalCGGAmoount     float64
}

func newEosConf() *EosConf {
	dur := conf.C.GetInt("eos.FetchIdleDur")
	num := conf.C.GetInt64("eos.FromBlkNum")
	return &EosConf{
		RPCURL:              conf.C.GetString("eos.RPCURL"),
		ChainID:             conf.C.GetString("eos.ChainID"),
		FetchIdleDur:        dur,
		FromBlkNum:          uint32(num),
		GameAccount:         conf.C.GetString("eos.GameAccount"),
		ICOAccount:          conf.C.GetString("eos.ICOAccount"),
		EnableICO:           conf.C.GetBool("eos.EnableICO"),
		ICOStartTime:        conf.C.GetInt64("eos.ICOStartTime"),
		EOS_CGG:             conf.C.GetFloat64("eos.EOS_CGG"),
		WalletURL:           conf.C.GetString("eos.WalletURL"),
		WalletPW:            conf.C.GetString("eos.WalletPW"),
		TokenSymbol:         conf.C.GetString("eos.TokenSymbol"),
		TokenAccount:        conf.C.GetString("eos.TokenAccount"),
		TotalAmount:         conf.C.GetFloat64("eos.ICOTotalAmount"),
		LockAccount:         conf.C.GetString("eos.LockAccount"),
		OfficialLockAccount: conf.C.GetString("eos.OfficialLockAccount"),
		TotalCGGAmoount:     conf.C.GetFloat64("eos.TotalCGGAmoount"),
	}
}

func sendTokens(from, to string, quan string, memo string) (hash string, err error) {

	cmd := exec.Command("cleos", "--wallet-url", eosConf.WalletURL, "--url", eosConf.RPCURL, "wallet", "unlock", "--password", eosConf.WalletPW)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		glog.Warningf("cleos unlock failed with %s\nErr:\n%s", err, string(stderr.Bytes()))
	}

	// defer exec.Command("cleos", "--wallet-url", eosConf.WalletURL, "--url", eosConf.RPCURL, "wallet", "lock").Run()

	// cleos push action eosio.token transfer '['xiaopingeos2', "xiaopingeos3", "2.0000 EOS", "转账EOS"]' -p xiaopingeos2@active
	// cleos push action xxptoken1234 transfer '['xiaopingeos2', "xiaopingeos3", "2.0000 CGG", "转账代币"]' -p xiaopingeos2@active
	var account string
	if strings.Contains(quan, "EOS") {
		account = "eosio.token"
	} else {
		account = eosConf.TokenAccount
	}

	// var sender string
	// if eosConf.EnableICO {
	// 	sender = eosConf.ICOAccount
	// } else if strings.Contains(memo, "unstaked") {
	// 	sender = eosConf.LockAccount
	// } else {
	// 	sender = eosConf.GameAccount
	// }

	actionData := fmt.Sprintf("[\"%s\", \"%s\", \"%s\", \"%s\"]", from, to, quan, memo)
	args := []string{"--wallet-url", eosConf.WalletURL, "--url", eosConf.RPCURL, "push", "action", account, "transfer", actionData, "-p", from + "@active"}
	cmd = exec.Command("cleos", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		glog.Errorf("push transfer failed with %s\nErr:\n%s", err, string(stderr.Bytes()))
		return "", err
	}

	output := string(stderr.Bytes())
	glog.V(7).Infof("output is : \n%s\n", output)
	hash1 := strings.SplitN(output, "executed transaction: ", 2)
	if len(hash1) != 2 {
		return "", errors.New("reslove hash error")

	}
	hash2 := strings.SplitN(hash1[1], " ", 2)
	if len(hash2) != 2 {
		return "", errors.New("reslove hash error")

	}

	return hash2[0], nil
}

func checkcleosExists() {
	path, err := exec.LookPath("cleos")
	if err != nil {
		glog.Fatalln("didn't find 'cleos' executable")
	} else {
		glog.Infof("'cleos' executable is in '%s'", path)
	}

	cmd := exec.Command("cleos", "--wallet-url", eosConf.WalletURL, "--url", eosConf.RPCURL, "get", "currency", "balance", "eosio.token", eosConf.GameAccount)
	glog.Info(cmd.Args)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		glog.Fatalf("cmd.Run() failed with %s\nErr:\n%s", err, string(stderr.Bytes()))
	}
	glog.Infof("cmd.Run() get balance of %s Out:%s", eosConf.ICOAccount, string(stdout.Bytes()))

}

// // EosRegister 注册balance相关路由
// func EosRegister(router *gin.RouterGroup) {
// 	router.POST("/chain/get_currency_balance", getCurrencyBalance)
// }

type balancePost struct {
	Code    string `json:"code" binding:"required,max=12"`
	Account string `json:"account" binding:"required,len=12"`
	Symbol  string `json:"symbol" binding:"required,len=3"`
}

var percentBalance json.Number = "0.00"
var cacheLock sync.RWMutex

func getPercent() json.Number {
	cacheLock.RLock()
	cached := percentBalance
	cacheLock.RUnlock()

	return cached
}

func setPercent(percent string) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	percentBalance = json.Number(percent)
}

func getCurrencyBalance(c *gin.Context) {

	post := balancePost{}
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(200, gin.H{
			"status":  -1,
			"message": "post参数错误！",
			"data":    nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  0,
		"message": "",
		"data": map[string]json.Number{
			"result": getPercent(),
		}})
}

func updateICOBalance() {
	result := getBalance(eosConf.ICOAccount, "eosio.token", "EOS")
	balance, _ := decimal.NewFromString(result)
	balance = balance.Add(decimal.NewFromFloat(conf.C.GetFloat64("eos.ICOFakeAmount")))

	percent := balance.Div(decimal.NewFromFloat(eosConf.TotalAmount)).Mul(decimal.NewFromFloat(100))
	setPercent(percent.StringFixed(2))
}

func getBalance(account string, code string, symbol string) string {
	// ICOTotalAmount = 60000
	url := eosConf.RPCURL + "/v1/chain/get_currency_balance"

	payload := strings.NewReader("{\"code\":\"" + code + "\", \"account\":\"" + account + "\",\"symbol\":\"" + symbol + "\"}")

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		glog.Error(err)
		return ""
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		glog.Error(err)
		return ""
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	results := []string{}
	if err := json.Unmarshal(body, &results); err != nil {
		glog.Error("unmarshal error", err)
		return ""
	}
	if len(results) != 1 {
		glog.Warningf("%s balance is %d, %v", account, 0, results)
		return ""
	}
	result := strings.Split(results[0], " ")
	if len(result) != 2 {
		glog.Error("result 格式有问题")
		return ""
	}
	return result[0]
}
