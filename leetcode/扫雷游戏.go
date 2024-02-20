package leetcode

//https://leetcode.cn/problems/minesweeper/description/

//输入：board = [["E","E","E","E","E"],["E","E","M","E","E"],["E","E","E","E","E"],["E","E","E","E","E"]], click = [3,0]
//输出：[["B","1","E","1","B"],["B","1","M","1","B"],["B","1","1","1","B"],["B","B","B","B","B"]]

func updateBoard(board [][]byte, click []int) [][]byte {
	x := click[0]
	y := click[1]
	if board[x][y] == 'M' {
		board[x][y] = 'X'
		return board
	}
	changeBoard(x, y, board)
	return board
}

func changeBoard(i, j int, board [][]byte) {
	if i < 0 || i > len(board)-1 {
		return
	}
	if j < 0 || j > len(board[i])-1 {
		return
	}
	if board[i][j] == 'M' {
		return
	}
	if board[i][j] != 'E' {
		return
	}
	var mineNum int
	mineNum += checkMine(i-1, j, board)
	mineNum += checkMine(i+1, j, board)
	mineNum += checkMine(i, j-1, board)
	mineNum += checkMine(i, j+1, board)
	mineNum += checkMine(i-1, j-1, board)
	mineNum += checkMine(i-1, j+1, board)
	mineNum += checkMine(i+1, j-1, board)
	mineNum += checkMine(i+1, j+1, board)
	if mineNum == 0 {
		board[i][j] = 'B'
		changeBoard(i-1, j, board)
		changeBoard(i+1, j, board)
		changeBoard(i, j-1, board)
		changeBoard(i, j+1, board)
		changeBoard(i-1, j-1, board)
		changeBoard(i-1, j+1, board)
		changeBoard(i+1, j-1, board)
		changeBoard(i+1, j+1, board)
	} else {
		board[i][j] = byte(mineNum + '0')
	}
}

func checkMine(i, j int, board [][]byte) int {
	if i < 0 || i > len(board)-1 {
		return 0
	}
	if j < 0 || j > len(board[i])-1 {
		return 0
	}
	if board[i][j] == 'M' {
		return 1
	}
	return 0
}
