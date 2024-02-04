package leetcode

//https://leetcode.cn/problems/top-k-frequent-elements/submissions/499112089/

//输入: nums = [1,1,1,2,2,3], k = 2
//输出: [1,2]

func topKFrequent(nums []int, k int) []int {
	var kNums = map[int]int{}
	var kvArr []kv
	for _, nu := range nums {
		kNums[nu]++
	}
	push := func(data kv) {
		kvArr = append(kvArr, data)
		i := len(kvArr) - 1
		j := (i - 1) / 2
		for j >= 0 && kvArr[i].key < kvArr[j].key {
			kvArr[i], kvArr[j] = kvArr[j], kvArr[i]
			i = j
			j = (i - 1) / 2
		}
	}
	pop := func() {
		kvArr[0] = kvArr[len(kvArr)-1]
		kvArr = kvArr[:len(kvArr)-1]
		i := 0
		len := len(kvArr) - 1
		for 2*i+1 <= len || 2*i+2 <= len {
			ni := 2*i + 1
			if 2*i+2 <= len {
				if kvArr[2*i+2].key < kvArr[2*i+1].key {
					ni = 2*i + 2
				}
			}
			if kvArr[ni].key >= kvArr[i].key {
				break
			}
			kvArr[ni], kvArr[i] = kvArr[i], kvArr[ni]
			i = ni
		}
	}
	i := 0
	for k1, v := range kNums {
		data := kv{v, k1}
		if i >= k {
			if kvArr[0].key < v {
				pop()
				push(data)
			}
		} else {
			push(data)
		}
		i++
	}
	var result []int
	for _, v := range kvArr {
		result = append(result, v.val)
	}
	return result
}

type kv struct {
	key int
	val int
}

func topKFrequentV1(nums []int, k int) []int {
	numMap := map[int]int{}
	var kvNums []kv
	for _, n := range nums {
		numMap[n]++
	}
	for k, v := range numMap {
		data := kv{v, k}
		kvNums = append(kvNums, data)
	}
	res := quick(kvNums, 0, len(kvNums)-1, k)
	return res
}

func quick(data []kv, start, end, k int) []int {
	if start >= end {
		return []int{data[start].val}
	}
	prev := start + (end-start)/2
	pivot := data[prev]
	data[start], pivot = pivot, data[start]
	i := start
	j := end
	for i < j {
		for i < j && data[j].key >= pivot.key {
			j--
		}
		data[i], data[j] = data[j], data[i]
		for i < j && data[i].key < pivot.key {
			i++
		}
		data[i], data[j] = data[j], data[i]
	}
	data[i] = pivot
	var res []int
	if end-i+1 == k {
		for i1 := i; i1 <= end; i1++ {
			res = append(res, data[i1].val)
		}
	} else if end-i+1 < k {
		for i1 := i; i1 <= end; i1++ {
			res = append(res, data[i1].val)
		}
		res = append(res, quick(data, start, i-1, k-(end-i)-1)...)
	} else {
		res = quick(data, i+1, end, k)
	}
	return res
}
