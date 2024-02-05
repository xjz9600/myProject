package leetcode

import "fmt"

//https://leetcode.cn/problems/letter-combinations-of-a-phone-number/submissions/500491778/

func letterCombinations(digits string) (ans []string) {
	if len(digits) == 0 {
		return
	}
	var phone = map[byte][]byte{}
	phone['2'] = []byte{'a', 'b', 'c'}
	phone['3'] = []byte{'d', 'e', 'f'}
	phone['4'] = []byte{'g', 'h', 'i'}
	phone['5'] = []byte{'j', 'k', 'l'}
	phone['6'] = []byte{'m', 'n', 'o'}
	phone['7'] = []byte{'p', 'q', 'r', 's'}
	phone['8'] = []byte{'t', 'u', 'v'}
	phone['9'] = []byte{'w', 'x', 'y', 'z'}
	var dfs func(re string, i int, digits []byte)
	dfs = func(re string, i int, digits []byte) {
		if i == len(digits) {
			ans = append(ans, re)
			return
		}
		arr := phone[digits[i]]
		for _, a := range arr {
			re += string(a)
			fmt.Println(re)
			dfs(re, i+1, digits)
			re = re[:len(re)-1]
		}
	}
	dfs("", 0, []byte(digits))
	return
}
