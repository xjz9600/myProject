package leetcode

func LongestCommonPrefix(strs []string) string {
	var result string
	if len(strs) == 0 {
		return result
	}
	if len(strs) == 1 {
		return strs[0]
	}
	prefix := strs[0]
	for i := 1; i < len(strs); i++ {
		nextStr := strs[i]
		if len(prefix) > len(nextStr) {
			prefix = prefix[:len(nextStr)]
		}
		for j := 0; j < len(nextStr); j++ {
			if len(prefix) <= j || prefix[j] != nextStr[j] {
				prefix = prefix[:j]
				break
			}
		}
	}
	return prefix
}
