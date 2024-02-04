package leetcode

// https://leetcode.cn/problems/valid-anagram/submissions/498034448/

func isAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}
	mp := map[int32]int{}
	for _, s1 := range s {
		mp[s1]++
	}
	for _, t1 := range t {
		if val, _ := mp[t1]; val == 0 {
			return false
		}
		mp[t1]--
	}
	return true
}
