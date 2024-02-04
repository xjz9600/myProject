package leetcode

// https://leetcode.cn/problems/group-anagrams/submissions/498054328/

// 输入: strs = ["eat", "tea", "tan", "ate", "nat", "bat"]
// 输出: [["bat"],["nat","tan"],["ate","eat","tea"]]

func groupAnagrams(strs []string) [][]string {
	var resultMap = map[[26]int][]string{}
	for _, st := range strs {
		var re [26]int
		for _, s := range st {
			re[s-'a']++
		}
		if _, ok := resultMap[re]; !ok {
			resultMap[re] = []string{st}
			continue
		}
		resultMap[re] = append(resultMap[re], st)
	}
	var result [][]string
	for _, rp := range resultMap {
		result = append(result, rp)
	}
	return result
}
