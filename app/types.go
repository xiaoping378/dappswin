package app

import "github.com/shopspring/decimal"

// ws msg types
const (
	block          = 0
	lotteryEOSBuy  = 101
	lotteryEOSWin  = 102
	lotteryGame    = 103
	lotteryCGGBuy  = 111
	lotteryCGGWin  = 112
	luckyNumEosBuy = 201
	luckyNumEosWin = 211
	luckyNumGame   = 213
)

var wsTypes = map[string]int{
	"lotteryEOSBuy":        101,
	"lotteryEOSWin":        102,
	"lotteryCGGBuy":        111,
	"lotteryCGGWin":        112,
	"lotteryCGGTotalVoted": 121,
	"lotteryEOSTotalVoted": 122,
}

var totalVotedEOS decimal.Decimal
var totalVotedCGG decimal.Decimal

const (
	eos = iota
	cgg
)

var coinIDs = map[string]int{
	"EOS": eos,
	"CGG": cgg,
}

var coinNames = map[int]string{
	eos: "EOS",
	cgg: "CGG",
}
