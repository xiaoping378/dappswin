package app

// 1 1000 0.01%
// 2 5000 0.02%
// 3 10000 0.03%
// 4 50000 0.04%
// 5 100000 0.05%
// 6 500000 0.06%
// 7 1000000 0.08%
// 8 5000000 0.12%
// 9 10000000 0.16%
// 10 50000000 0.20%
type VIP struct {
	Level  uint8
	Amount float64
	Rebate float64
}

var vip = map[float64]float64{
	1e3: 0.01,
	5e3: 0.02,
	1e4: 0.03,
	5e4: 0.04,
	1e5: 0.05,
	5e5: 0.06,
	1e6: 0.08,
	5e6: 0.12,
	1e7: 0.16,
	5e7: 0.20,
}

var vipInfo = []VIP{
	{1, 1e3, 0.01},
	{2, 5e3, 0.02},
	{3, 1e4, 0.03},
	{4, 5e4, 0.04},
	{5, 1e5, 0.05},
	{6, 5e5, 0.06},
	{7, 1e6, 0.08},
	{8, 5e6, 0.12},
	{9, 1e7, 0.16},
	{10, 5e7, 0.20},
}

func getVIPLevel(amount float64) uint8 {

	for _, a := range vipInfo {
		if amount < a.Amount {
			return a.Level - 1
		}
	}
	return 10
}

func getNewVIP(name string, newTotal float64) uint8 {
	return 0
}
