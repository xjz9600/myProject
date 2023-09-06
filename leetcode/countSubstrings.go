package leetcode

//给定一个字符串 s ，请计算这个字符串中有多少个回文子字符串。
//具有不同开始位置或结束位置的子串，即使是由相同的字符组成，也会被视作不同的子串。
//示例 1：
//输入：s = "abc"
//输出：3
//解释：三个回文子串: "a", "b", "c"

//示例 2：
//输入：s = "aaa"
//输出：6
//解释：6个回文子串: "a", "a", "a", "aa", "aa", "aaa"

//提示：
//1 <= s.length <= 1000
//s 由小写英文字母组成

var result int

func countSubstrings(s string) int {
	var data []byte = []byte(s)
	for i := 0; i < len(data); i++ {
		getString(data, i, i)
	}
	return result
}

func getString(str []byte, startIndex int, nowIndex int) {
	if nowIndex == len(str) {
		return
	}
	if str[startIndex] == str[nowIndex] {
		if StringIsOK(str[startIndex : nowIndex+1]) {
			result++
		}
	}
	getString(str, startIndex, nowIndex+1)
}

func StringIsOK(s []byte) bool {
	size := len(s) - 1
	i := 0
	for i <= size {
		if s[i] != s[size] {
			return false
		}
		i++
		size--
	}
	return true
}
