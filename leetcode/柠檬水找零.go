package leetcode

//https://leetcode.cn/problems/lemonade-change/description/

//输入：bills = [5,5,5,10,20]
//输出：true
//解释：
//前 3 位顾客那里，我们按顺序收取 3 张 5 美元的钞票。
//第 4 位顾客那里，我们收取一张 10 美元的钞票，并返还 5 美元。
//第 5 位顾客那里，我们找还一张 10 美元的钞票和一张 5 美元的钞票。
//由于所有客户都得到了正确的找零，所以我们输出 true。

func lemonadeChange(bills []int) bool {
	var five int
	var ten int
	for _, b := range bills {
		if b == 5 {
			five++
		}
		if b == 10 {
			if five < 1 {
				return false
			}
			five--
			ten++
		}
		if b == 20 {
			if ten < 1 {
				if five < 3 {
					return false
				}
				five -= 3
			} else {
				if five < 1 {
					return false
				}
				ten--
				five--
			}
		}
	}
	return true
}
