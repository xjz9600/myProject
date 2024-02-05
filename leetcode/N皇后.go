package leetcode

import "math/bits"

// https://leetcode.cn/problems/n-queens/submissions/500520575/
// 输入：n = 4
// 输出：[[".Q..","...Q","Q...","..Q."],["..Q.","Q...","...Q",".Q.."]]
// 解释：如上图所示，4 皇后问题存在两个不同的解法。
var queueStr [][]string

func solveNQueens(n int) [][]string {
	queueStr = [][]string{}
	getData([]int{}, map[int]struct{}{}, map[int]struct{}{}, map[int]struct{}{}, 0, n)
	return queueStr
}

func getData(re []int, col, left, right map[int]struct{}, row, n int) {
	if n == row {
		q := genQueue(re)
		queueStr = append(queueStr, q)
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
		getData(re, col, left, right, row+1, n)
		re = re[:len(re)-1]
		delete(col, i)
		delete(left, l)
		delete(right, r)
	}
}

func genQueue(res []int) []string {
	var result []string
	for i := 0; i < len(res); i++ {
		var queue string
		for j := 0; j < len(res); j++ {
			if res[i] == j {
				queue += "Q"
			} else {
				queue += "."
			}
		}
		result = append(result, queue)
	}
	return result
}

func solveNQueensV1(n int) [][]string {
	queueStr = [][]string{}
	getData1([]int{}, 0, 0, 0, 0, n)
	return queueStr
}

func getData1(re []int, col, left, right, row, n int) {
	if n == row {
		q := genQueue(re)
		queueStr = append(queueStr, q)
		return
	}
	availablePosition := ((1 << n) - 1) &^ (col | left | right)
	for availablePosition != 0 {
		position := availablePosition & (-availablePosition)
		availablePosition = availablePosition & (availablePosition - 1)
		selPosition := bits.OnesCount(uint(position - 1))
		re = append(re, selPosition)
		getData1(re, col|position, left|position<<1, right|position>>1, row+1, n)
		re = re[:len(re)-1]
	}
}
