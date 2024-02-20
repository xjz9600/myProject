package leetcode

import "math/bits"

// https://leetcode.cn/problems/n-queens/submissions/500520575/
// 输入：n = 4
// 输出：[[".Q..","...Q","Q...","..Q."],["..Q.","Q...","...Q",".Q.."]]
// 解释：如上图所示，4 皇后问题存在两个不同的解法。

var queuePositions [][]string

func solveNQueensV1(n int) [][]string {
	queuePositions = [][]string{}
	getDataV1([]int{}, map[int]struct{}{}, map[int]struct{}{}, map[int]struct{}{}, 0, n)
	return queuePositions
}

func getDataV1(re []int, col, left, right map[int]struct{}, row, n int) {
	if row == n {
		queuePosition := genQueue(re)
		queuePositions = append(queuePositions, queuePosition)
		return
	}
	for i := 0; i < n; i++ {
		if _, ok := col[i]; ok {
			continue
		}
		l := row + i
		if _, ok := left[l]; ok {
			continue
		}
		r := row - i
		if _, ok := right[r]; ok {
			continue
		}
		re = append(re, i)
		col[i] = struct{}{}
		left[l] = struct{}{}
		right[r] = struct{}{}
		getDataV1(re, col, left, right, row+1, n)
		re = re[:len(re)-1]
	}
	return
}

func genQueue(queue []int) []string {
	var result []string
	for i := 0; i < len(queue); i++ {
		var res string
		for j := 0; j < len(queue); j++ {
			if queue[i] == j {
				res += "Q"
				continue
			}
			res += "."
		}
		result = append(result, res)
	}
	return result
}

func solveNQueens(n int) [][]string {
	queuePositions = [][]string{}
	getData([]int{}, 0, 0, 0, 0, n)
	return queuePositions
}

func getData(re []int, col, left, right, row, n int) {
	if row == n {
		queuePosition := genQueue(re)
		queuePositions = append(queuePositions, queuePosition)
		return
	}
	availablePositions := (1<<n - 1) & ^(col | left | right)
	for availablePositions != 0 {
		position := availablePositions & (-availablePositions)
		availablePositions = availablePositions & (availablePositions - 1)
		selPosition := bits.OnesCount(uint(position - 1))
		re = append(re, selPosition)
		getData(re, col|position, (left|position)<<1, (right|position)>>1, row+1, n)
		re = re[:len(re)-1]
	}
}
