package leetcode

//https://leetcode.cn/problems/word-ladder-ii/description/

//输入：beginWord = "hit", endWord = "cog", wordList = ["hot","dot","dog","lot","log","cog"]
//输出：[["hit","hot","dot","dog","cog"],["hit","hot","lot","log","cog"]]
//解释：存在 2 种最短的转换序列：
//"hit" -> "hot" -> "dot" -> "dog" -> "cog"
//"hit" -> "hot" -> "lot" -> "log" -> "cog"

func findLadders(beginWord string, endWord string, wordList []string) [][]string {
	var result [][]string
	var endIndex = -1
	for i, w := range wordList {
		if w == endWord {
			endIndex = i
			break
		}
	}
	if endIndex == -1 {
		return result
	}
	var grid = make([][]bool, len(wordList))
	for i := range grid {
		grid[i] = make([]bool, len(wordList))
	}
	var dist = make([]int, len(wordList))
	var queue []int
	for i := 0; i < len(wordList); i++ {
		if isOneStep(beginWord, wordList[i]) {
			dist[i] = 1
			queue = append(queue, i)
		}
		for j := i + 1; j < len(wordList); j++ {
			if isOneStep(wordList[i], wordList[j]) {
				grid[i][j] = true
				grid[j][i] = true
			}
		}
	}
	for len(queue) > 0 {
		var newQueue []int
		for _, q := range queue {
			for i, ok := range grid[q] {
				if dist[i] == 0 && ok {
					dist[i] = dist[q] + 1
					newQueue = append(newQueue, i)
				}
			}
		}
		queue = newQueue
	}
	var dfs func(int)
	var res []string
	dfs = func(i int) {
		if isOneStep(wordList[i], beginWord) {
			resArr := make([]string, len(res))
			copy(resArr, res)
			resArr = append(resArr, wordList[i], beginWord)
			reversal(resArr)
			result = append(result, resArr)
			return
		}
		for j, ok := range grid[i] {
			if ok && dist[j] == dist[i]-1 {
				res = append(res, wordList[i])
				dfs(j)
				res = res[:len(res)-1]
			}
		}
	}
	dfs(endIndex)
	return result
}

func reversal(arr []string) {
	for i := 0; i < len(arr)/2; i++ {
		arr[i], arr[len(arr)-1-i] = arr[len(arr)-1-i], arr[i]
	}
}

func isOneStep(a, b string) bool {
	var step int
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			step++
		}
	}
	return step == 1
}
