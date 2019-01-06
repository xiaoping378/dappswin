package app

import (
	"fmt"
	"testing"
)

func Test_getRefund(t *testing.T) {
	type args struct {
		Time  int64
		value float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test1", args: args{0, 10 * 1e4}, want: "200.0000"},
		{name: "test2", args: args{0, 100 * 1e4}, want: "2010.0000"},
		{name: "test3", args: args{0, 601.2 * 1e4}, want: "12144.2400"},
		{name: "test4", args: args{0, 1010 * 1e4}, want: "20503.0000"},
		{name: "test5", args: args{0, 1500 * 1e4}, want: "30600.0000"},
		{name: "test6", args: args{0, 1500.1234 * 1e4}, want: "30602.5174"},
		{name: "test7", args: args{0, 2001 * 1e4}, want: "41020.5000"},
		{name: "test8", args: args{2 + oneDayMills, 1500.1234 * 1e4}, want: "30002.4680"},
		{name: "test9", args: args{2 + oneDayMills, 600 * 1e4}, want: "12000.0000"},
	}

	// Init handle this.
	eosConf = &EosConf{EOS_CGG: 20, ICOStartTime: 1}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRefund(tt.args.Time, tt.args.value) / 1e4; fmt.Sprintf("%.4f", got) != tt.want {
				t.Errorf("getRefund() = %v, want %v", fmt.Sprintf("%.4f", got), tt.want)
			}
		})
	}
}
