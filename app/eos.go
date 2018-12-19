package app

import (
	"bytes"
	"dappswin/conf"
	"dappswin/database"
	"fmt"
	"os/exec"
	"strings"

	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

var eosConf *EosConf
var db *gorm.DB
var apiEndpoint []string

// InitEos 启动Eos Resolver
func Init() {
	eosConf = newEosConf()
	db = database.Db
	go gameRoutine()
	go resloveTXRoutine()
	go checkWinRoutine()
	if eosConf.EnableICO {
		go checkICORoutine()
	}
	go Huber.run()
	checkcleosExists()
	//
}

// EosConf :
type EosConf struct {
	RPCURL       string
	ChainID      string
	FetchIdleDur int    // 查询blk时间间隔
	FromBlkNum   uint32 // 从哪个blocknum开始查询
	GameAccount  string
	ICOAccount   string
	EnableICO    bool
	ICOStartTime int64
	TokenSymbol  string
	TokenAccount string
	EOS_CGG      uint64
	WalletURL    string
	WalletPW     string
}

func newEosConf() *EosConf {
	dur := conf.C.GetInt("eos.FetchIdleDur")
	num := conf.C.GetInt64("eos.FromBlkNum")
	return &EosConf{
		RPCURL:       conf.C.GetString("eos.RPCURL"),
		ChainID:      conf.C.GetString("eos.ChainID"),
		FetchIdleDur: dur,
		FromBlkNum:   uint32(num),
		GameAccount:  conf.C.GetString("eos.GameAccount"),
		ICOAccount:   conf.C.GetString("eos.ICOAccount"),
		EnableICO:    conf.C.GetBool("eos.EnableICO"),
		ICOStartTime: conf.C.GetInt64("eos.ICOStartTime"),
		EOS_CGG:      uint64(conf.C.GetInt("eos.EOS_CGG")),
		WalletURL:    conf.C.GetString("eos.WalletURL"),
		WalletPW:     conf.C.GetString("eos.WalletPW"),
		TokenSymbol:  conf.C.GetString("eos.TokenSymbol"),
		TokenAccount: conf.C.GetString("eos.TokenAccount"),
	}
}

func sendTokens(to string, quan string) error {

	cmd := exec.Command("cleos", "--wallet-url", eosConf.WalletURL, "--url", eosConf.RPCURL, "wallet", "unlock", "--password", eosConf.WalletPW)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		glog.Errorf("cmd.Run() failed with %s\nErr:\n%s", err, string(stderr.Bytes()))
	}
	defer exec.Command("cleos", "--wallet-url", eosConf.WalletURL, "--url", eosConf.RPCURL, "wallet", "lock").Run()

	// cleos push action eosio.token transfer '['xiaopingeos2', "xiaopingeos3", "2.0000 EOS", "转账EOS"]' -p xiaopingeos2@active
	// cleos push action xxptoken1234 transfer '['xiaopingeos2', "xiaopingeos3", "2.0000 CGG", "转账代币"]' -p xiaopingeos2@active
	var account string
	if strings.Contains(quan, "EOS") {
		account = "eosio.token"
	} else {
		account = eosConf.TokenAccount
	}

	actionData := fmt.Sprintf("[\"%s\", \"%s\", \"%s\", \"%s\"]", eosConf.ICOAccount, to, quan, " ")
	args := []string{"--wallet-url", eosConf.WalletURL, "--url", eosConf.RPCURL, "push", "action", account, "transfer", actionData, "-p", eosConf.ICOAccount}
	cmd = exec.Command("cleos", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		glog.Errorf("cmd.Run() failed with %s\nErr:\n%s", err, string(stderr.Bytes()))
		return err
	}

	return nil
}

func checkcleosExists() {
	path, err := exec.LookPath("cleos")
	if err != nil {
		glog.Fatalln("didn't find 'cleos' executable")
	} else {
		glog.Infof("'cleos' executable is in '%s'", path)
	}

	cmd := exec.Command("cleos", "--wallet-url", eosConf.WalletURL, "--url", eosConf.RPCURL, "get", "currency", "balance", "eosio.token", eosConf.ICOAccount)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		glog.Fatalf("cmd.Run() failed with %s\nErr:\n%s", err, string(stderr.Bytes()))
	}
	glog.Infof("cmd.Run() get balance of %s Out:%s", eosConf.ICOAccount, string(stdout.Bytes()))

}
