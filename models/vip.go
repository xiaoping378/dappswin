package models

type VIP struct {
	Level  uint8
	Amount uint64
	Rebate float32
}

var vipInfo = []VIP{
	{1, 0, 0},
	{2, 10, 0.001},
	{3, 100, 0.002},
	{4, 1e3, 0.003},
	{5, 1e4, 0.004},
	{6, 1e5, 0.005},
	{7, 1e6, 0.006},
	{8, 1e7, 0.007},
	{9, 1e8, 0.008},
}

func getLevel(amount uint64) uint8 {

	for _, a := range vipInfo {
		if amount < a.Amount {
			return a.Level - 1
		}
	}
	return 9
}
