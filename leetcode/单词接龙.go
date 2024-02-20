package leetcode

//https://leetcode.cn/problems/word-ladder/

//输入：beginWord = "hit", endWord = "cog", wordList = ["hot","dot","dog","lot","log","cog"]
//输出：5
//解释：一个最短转换序列是 "hit" -> "hot" -> "dot" -> "dog" -> "cog", 返回它的长度 5。

func ladderLength(beginWord string, endWord string, wordList []string) int {
	if len(wordList) == 0 {
		return 0
	}
	var endIndex = -1
	for i, e := range wordList {
		if e == endWord {
			endIndex = i
		}
	}
	if endIndex == -1 {
		return 0
	}
	var grid = make([][]bool, len(wordList))
	for i := range grid {
		grid[i] = make([]bool, len(wordList))
	}
	var dist = make([]int, len(wordList))
	var queue []int
	for i := 0; i < len(wordList); i++ {
		if isOneStep(beginWord, wordList[i]) {
			dist[i] = 2
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
			for i, g := range grid[q] {
				if g == true && dist[i] == 0 {
					dist[i] = dist[q] + 1
					newQueue = append(newQueue, i)
				}
			}
		}
		queue = newQueue
	}
	return dist[endIndex]
}

//func isOneStep(a, b string) bool {
//	var step int
//	for i := 0; i < len(a); i++ {
//		if a[i] != b[i] {
//			step++
//		}
//	}
//	return step == 1
//}
