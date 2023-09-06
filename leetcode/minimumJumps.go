package leetcode

import "sort"

func MaxSize(a, b int, forbidden []int, x int) int {
	if a >= b {
		return b + x
	}
	var maxForbidden int
	for _, f := range forbidden {
		if maxForbidden <= f {
			maxForbidden = f
		}
	}
	return Max(maxForbidden+a+b, x)
}

func Max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func minimumJumps(forbidden []int, a int, b int, x int) int {
	sort.Ints(forbidden)
	var minNum int
	index := sort.SearchInts(forbidden, a)
	if index != len(forbidden) && forbidden[index] == a {
		return -1
	}
	if x == 0 {
		return 0
	}
	maxSize := MaxSize(a, b, forbidden, x)
	var startPosition = []positionData{{tpy: "add", position: a}}
	var poistioned = map[int]struct{}{}
	for len(startPosition) > 0 {
		minNum++
		var nextPosition []positionData
		for i := 0; i < len(startPosition); i++ {
			positionInfo := startPosition[i]
			if positionInfo.position > maxSize {
				continue
			}
			var direction int = 1
			if positionInfo.tpy == "reduce" {
				direction = -1
			}
			if _, ok := poistioned[positionInfo.position*direction]; ok {
				continue
			} else {
				poistioned[positionInfo.position*direction] = struct{}{}
			}
			if positionInfo.position == x {
				return minNum
			}
			index := sort.SearchInts(forbidden, positionInfo.position+a)
			if index == len(forbidden) || forbidden[index] != positionInfo.position+a {
				nextPosition = append(nextPosition, positionData{tpy: "add", position: positionInfo.position + a})
			}
			if positionInfo.tpy == "add" {
				if positionInfo.position-b > 0 {
					index = sort.SearchInts(forbidden, positionInfo.position-b)
					if index == len(forbidden) || forbidden[index] != positionInfo.position-b {
						nextPosition = append(nextPosition, positionData{tpy: "reduce", position: positionInfo.position - b})
					}
				}
			}

		}
		startPosition = nextPosition
	}
	return -1
}

type positionData struct {
	tpy      string
	position int
}
