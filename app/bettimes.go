package app

import (
	"fmt"
	"math"
	"strings"
)

var starCount = 6
var lastStar = starCount - 1

/*
一个分号的情况有2种，
1：b/s,o/e;453456
2:b/s/o/e;12345
第1种如果中奖，返回的倍率有：2,4,10,12,14几种
flag表示是否猜中数字
b/s  o/e  flag  times
0    0   0       0
0    0   1      10
0    1   0       2
0    1   1       12
1    0   0       2
1    0   1      12
1    1   0       4
1    1   1      14
第2种如果中奖，返回的倍率有：2,10,12 三种
b/s/o/e    flag  times
0           0    0
0           1    10
1           0    2
1           1    12
*/
func HandleBsOeAndOneStar(str []string, gameNum []byte) (int, int) {
	var flag int8
	var betnum int
	for _, v := range str[1] { /*用户猜中数字了*/
		if byte(v) == gameNum[lastStar] {
			flag = 1
		}
	}
	betnum = len(str[1]) - 1 /*一星选取的数字个数*/
	if len(str[0]) == 4 {    /*str[0] maybe [b,o,] or[b,e,]or[s,o,] or[s,e,]*/
		var V1Bit, V2Bit int8
		switch str[0][0] {
		case 'b':
			if gameNum[lastStar] >= '5' && gameNum[lastStar] <= '9' {
				V1Bit = 1
			}
		case 's':
			if gameNum[lastStar] >= '0' && gameNum[lastStar] <= '4' {
				V1Bit = 1
			}
		}
		switch str[0][2] {
		case 'o':
			if gameNum[lastStar]%2 > 0 {
				V2Bit = 1
			}
		case 'e':
			if gameNum[lastStar]%2 == 0 {
				V2Bit = 1
			}
		}
		value := V1Bit<<2 | V2Bit<<1 | flag
		switch value {
		case 0:
			return 0, betnum + 2
		case 1:
			return 10, betnum + 2
		case 2:
			return 2, betnum + 2
		case 3:
			return 12, betnum + 2
		case 4:
			return 2, betnum + 2
		case 5:
			return 12, betnum + 2
		case 6:
			return 4, betnum + 2
		case 7:
			return 14, betnum + 2
		}
		return 0, betnum + 2
	} else { /* str[0] maybe [b,] or [s,] or[o,] or[e,]*/
		var bit int8
		switch str[0][0] {
		case 'b':
			if gameNum[lastStar] >= '5' && gameNum[lastStar] <= '9' {
				bit = 1
			}
		case 's':
			if gameNum[lastStar] >= '0' && gameNum[lastStar] <= '5' {
				bit = 1
			}
		case 'o':
			if gameNum[lastStar]%2 > 0 {
				bit = 1
			}
		case 'e':
			if gameNum[lastStar]%2 == 0 {
				bit = 1
			}
		}
		value := bit<<1 | flag
		switch value {
		case 0:
			return 0, betnum + 1
		case 1:
			return 10, betnum + 1
		case 2:
			return 2, betnum + 1
		case 3:
			return 12, betnum + 1
		}
		return 0, betnum + 1
	}
}

/*
This function will return the number of times for the bets
b
s
o
e
b,o
b,e
s,o
s,e   split(str,"[") 切分完之后，得到一个包含一个成员的数组。
------------------
b,[0~9]
s,[0~9]
o,[0~9]
e,[0~9]
b,o,[0~9]
b,e,[0~9]
s,o,[0~9]
s,e,[0~9]  split 返回包含2个元素的字符串数组
-------------------------
[0~9]
[0~9][0~9]
[0~9][0~9][0~9]
[0~9][0~9][0~9][0~9]
[0~9][0~9][0~9][0~9][0~9]  split 返回i+1个元素的字符串数组，第一个元素str[0]为空，从str[1]处理
-----------------------------------------
*/
func HandleBetInfo(betinfo string, gameNum []byte) (int, int) {
	str := strings.Split(betinfo, "[")
	strlen := len(str)
	switch strlen { /*以split 切分后返回的字符串数组长度做case 分支*/
	case 1: /*handle : b/s/o/e/b,o/b,e/s,o/s,e*/
		if len(str[0]) > 1 {
			if str[0][0] == 'b' && str[0][2] == 'o' {
				if (gameNum[lastStar] >= '5' && gameNum[lastStar] <= '9') && gameNum[lastStar]%2 > 0 {
					fmt.Println("the reward num is: 4")
					return 4, 2
				}
				if ((gameNum[lastStar] >= '5' && gameNum[lastStar] <= '9') && gameNum[lastStar]%2 == 0) || ((gameNum[lastStar] >= '0' && gameNum[lastStar] <= '4') && gameNum[lastStar]%2 > 0) {
					fmt.Println("the reward num is: 2")
					return 2, 2
				}
			}
			if str[0][0] == 'b' && str[0][2] == 'e' {
				if (gameNum[lastStar] >= '5' && gameNum[lastStar] <= '9') && gameNum[lastStar]%2 == 0 {
					fmt.Println("the reward num is: 4")
					return 4, 2
				}
				if ((gameNum[lastStar] >= '5' && gameNum[lastStar] <= '9') && gameNum[lastStar]%2 > 0) || ((gameNum[lastStar] >= '0' && gameNum[lastStar] <= '4') && gameNum[lastStar]%2 == 0) {
					fmt.Println("the reward num is: 2")
					return 2, 2
				}
			}
			if str[0][0] == 's' && str[0][2] == 'o' {
				if (gameNum[lastStar] >= '0' && gameNum[lastStar] <= '4') && gameNum[lastStar]%2 > 0 {
					fmt.Println("the reward num is: 4")
					return 4, 2
				}
				if ((gameNum[lastStar] >= '0' && gameNum[lastStar] <= '4') && gameNum[lastStar]%2 == 0) || ((gameNum[lastStar] >= '5' && gameNum[lastStar] <= '9') && gameNum[lastStar]%2 > 0) {
					fmt.Println("the reward num is: 2")
					return 2, 2
				}
			}
			if str[0][0] == 's' && str[0][2] == 'e' {
				if (gameNum[lastStar] >= '0' && gameNum[lastStar] <= '4') && gameNum[lastStar]%2 == 0 {
					fmt.Println("the reward num is: 4")
					return 4, 2
				}
				if ((gameNum[lastStar] >= '5' && gameNum[lastStar] <= '9') && gameNum[lastStar]%2 == 0) || ((gameNum[lastStar] >= '0' && gameNum[lastStar] <= '4') && gameNum[lastStar]%2 > 0) {
					fmt.Println("the reward num is: 2")
					return 2, 2
				}
			}
			return 0, 2
		} else {
			switch str[0][0] {
			case 'b':
				if gameNum[lastStar] >= '5' && gameNum[lastStar] <= '9' {
					fmt.Println("the reward num is: 2")
					return 2, 1
				}
			case 's':
				if gameNum[lastStar] >= '0' && gameNum[lastStar] <= '4' {
					fmt.Println("the reward num is: 2")
					return 2, 1
				}
			case 'o':
				if gameNum[lastStar]%2 > 0 {
					fmt.Println("the reward num is: 2")
					return 2, 1
				}
			case 'e':
				if gameNum[lastStar]%2 == 0 {
					fmt.Println("the reward num is: 2")
					return 2, 1
				}
			}
			return 0, 1
		}
	case 2:
		/*长度为2的时候，有2种情况：
		b,[0~9]
		s,[0~9]
		o,[0~9]
		e,[0~9]
		b,o,[0~9]
		b,e,[0~9]
		s,o,[0~9]
		s,e,[0~9]
		一种是：
		只有[0-9]*/
		var times, betnum int
		if len(str[0]) > 0 {
			times, betnum = HandleBsOeAndOneStar(str, gameNum)
		} else {
			times, betnum = HandleStarNum(str, gameNum, 1)
		}
		return times, betnum
	case 3:
		times, betnum := HandleStarNum(str, gameNum, 2)
		return times, betnum
	case 4:
		times, betnum := HandleStarNum(str, gameNum, 3)
		return times, betnum
	case 5:
		times, betnum := HandleStarNum(str, gameNum, 4)
		return times, betnum
	case 6:
		times, betnum := HandleStarNum(str, gameNum, 5)
		return times, betnum
	case 7:
		times, betnum := HandleStarNum(str, gameNum, 6)
		return times, betnum
	default:
		fmt.Println("the betinfo is not valid")
		return 0, 0
	}
}

/*
多星玩儿法；不涉及选择大小，单双
*/
func HandleStarNum(str []string, gameNum []byte, starnum int) (int, int) {
	flag := make([]int, starnum)
	var hitflag int
	var betnum int
	j := starCount - starnum
	for i := 1; i <= starnum; i++ {
		for _, v := range str[i] {
			if byte(v) == gameNum[j] {
				flag[i-1] = 1
			}
		}
		j++
	}
	betnum = len(str[1]) - 1 /*delete the last ']'*/
	if len(str) > 2 {
		for i := 2; i <= starnum; i++ {
			betnum *= (len(str[i]) - 1)
		}
	}
	hitflag = flag[0]
	for i := 1; i < starnum; i++ {
		hitflag = hitflag & flag[i]
	}
	if hitflag > 0 {
		return int(math.Pow10(starnum)), betnum
	} else {
		return 0, betnum
	}
}

/*
func main() {
	if len(os.Args) != 3 {
		fmt.Println("Please input the correct parameters: ./main \"b,e,[1223]/[123][456][789]\" \"45678\" ")
		return
	}
	betinfo := os.Args[1]
	gameNum := []byte(os.Args[2])
	wintimes, betnum := HandleBetInfo(betinfo, gameNum)
	fmt.Println("you should send", wintimes, "times to user", "the bet number is:", betnum)
}*/
