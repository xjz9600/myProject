package leetcode

//https://leetcode.cn/problems/serialize-and-deserialize-binary-tree/submissions/499933195/

import (
	"strconv"
	"strings"
)

type Codec struct {
}

func ConstructorV2() Codec {
	return Codec{}
}

func (this *Codec) serialize(root *TreeNode) string {
	return genStr(root)
}

func genStr(node *TreeNode) string {
	var result string
	if node == nil {
		return "1001,"
	} else {
		str := strconv.Itoa(node.Val)
		result += str + ","
		result += genStr(node.Left)
		result += genStr(node.Right)
	}
	return result
}

// Deserializes your encoded data to tree.
func (this *Codec) deserialize(data string) *TreeNode {
	dataArr := strings.Split(data, ",")
	node, _ := genNode(dataArr, 0)
	return node
}

func genNode(dataArr []string, i int) (*TreeNode, int) {
	if dataArr[i] == "1001" {
		return nil, i + 1
	}
	var node = &TreeNode{}
	intVal, _ := strconv.Atoi(dataArr[i])
	node.Val = intVal
	node.Left, i = genNode(dataArr, i+1)
	node.Right, i = genNode(dataArr, i)
	return node, i
}
