package leetcode

// https://leetcode.cn/problems/walking-robot-simulation/description/

//输入：commands = [4,-1,4,-2,4], obstacles = [[2,4]]
//输出：65
//解释：机器人开始位于 (0, 0)：
//1. 向北移动 4 个单位，到达 (0, 4)
//2. 右转
//3. 向东移动 1 个单位，然后被位于 (2, 4) 的障碍物阻挡，机器人停在 (1, 4)
//4. 左转
//5. 向北走 4 个单位，到达 (1, 8)
//距离原点最远的是 (1, 8) ，距离为 12 + 82 = 65

func robotSim(commands []int, obstacles [][]int) int {
	var obstMap = map[[2]int]struct{}{}
	for _, val := range obstacles {
		pos := [2]int{}
		pos[0] = val[0]
		pos[1] = val[1]
		obstMap[pos] = struct{}{}
	}
	var x, y int
	// 0为上，1为右，2为下，3为左
	remove := [][2]int{
		{0, 1},
		{1, 0},
		{0, -1},
		{-1, 0},
	}
	var dire int
	var ans int
	for _, c := range commands {
		if c == -1 {
			dire = (dire + 1) % 4
			continue
		}
		if c == -2 {
			dire = (dire + 3) % 4
			continue
		}
		for i := 0; i < c; i++ {
			checkPost := [2]int{}
			checkPost[0] = x + remove[dire][0]
			checkPost[1] = y + remove[dire][1]
			if _, ok := obstMap[checkPost]; ok {
				break
			}
			x = x + remove[dire][0]
			y = y + remove[dire][1]
			ans = max(ans, x*x+y*y)
		}
	}
	return ans
}

//func max(a, b int) int {
//	if a >= b {
//		return a
//	}
//	return b
//}
