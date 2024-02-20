package leetcode

// https://leetcode.cn/problems/minimum-genetic-mutation/

//输入：start = "AAAAACCC", end = "AACCCCCC", bank = ["AAAACCCC","AAACCCCC","AACCCCCC"]
//输出：3

//start、end 和 bank[i] 仅由字符 ['A', 'C', 'G', 'T'] 组成

func genChangeGene(gene []byte, bankMap map[string]struct{}) [][]byte {
	var change = []byte{'A', 'C', 'G', 'T'}
	var result [][]byte
	for i, g := range gene {
		for _, c := range change {
			if c == g {
				continue
			}

			var src = make([]byte, len(gene))
			copy(src, gene)
			src[i] = c
			if _, ok := bankMap[string(src)]; ok {
				result = append(result, src)
				delete(bankMap, string(src))
			}
		}
	}
	return result
}

func minMutation(startGene string, endGene string, bank []string) int {
	if len(bank) == 0 {
		return -1
	}
	var bankMap = map[string]struct{}{}
	for _, b := range bank {
		bankMap[b] = struct{}{}
	}
	if _, ok := bankMap[endGene]; !ok {
		return -1
	}
	var queue [][]byte
	var min = -1
	queue = append(queue, []byte(startGene))
	for len(queue) > 0 {
		min++
		var newQueue [][]byte
		for _, q := range queue {
			if string(q) == endGene {
				return min
			}
			newQueue = append(newQueue, genChangeGene(q, bankMap)...)
		}
		queue = newQueue
	}
	return -1
}
