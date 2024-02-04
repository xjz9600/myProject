package leetcode

//输入：s = "()[]{}"
//输出：true
//示例 3：
//
//输入：s = "(]"
//输出：false

// https://leetcode.cn/problems/valid-parentheses/description/

func isValid(s string) bool {
	var stack []byte
	bsMap := map[byte]byte{
		'}': '{',
		')': '(',
		']': '[',
	}
	for _, s1 := range s {
		bs := byte(s1)
		switch bs {
		case ')', ']', '}':
			if len(stack) == 0 || stack[len(stack)-1] != bsMap[bs] {
				return false
			}
			stack = stack[:len(stack)-1]
		default:
			stack = append(stack, bs)
		}
	}
	if len(stack) == 0 {
		return true
	}
	return false
}
