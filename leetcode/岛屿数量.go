package leetcode

//https://leetcode.cn/problems/number-of-islands/description/
//输入：grid = [
//["1","1","1","1","0"],
//["1","1","0","1","0"],
//["1","1","0","0","0"],
//["0","0","0","0","0"]
//]
//输出：1

func numIslands(grid [][]byte) int {
	var num int
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == '1' {
				num++
				dfs(grid, i, j)
			}
		}
	}
	return num
}

func dfs(grid [][]byte, i, j int) {
	if i >= len(grid) || i < 0 {
		return
	}
	if j >= len(grid[i]) || j < 0 {
		return
	}
	if grid[i][j] == '0' {
		return
	}
	grid[i][j] = '0'
	dfs(grid, i, j+1)
	dfs(grid, i, j-1)
	dfs(grid, i+1, j)
	dfs(grid, i-1, j)
}
