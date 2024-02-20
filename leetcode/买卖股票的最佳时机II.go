package leetcode

//https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-ii/description/

//示例 1：
//
//输入：prices = [7,1,5,3,6,4]
//输出：7
//解释：在第 2 天（股票价格 = 1）的时候买入，在第 3 天（股票价格 = 5）的时候卖出, 这笔交易所能获得利润 = 5 - 1 = 4 。
//随后，在第 4 天（股票价格 = 3）的时候买入，在第 5 天（股票价格 = 6）的时候卖出, 这笔交易所能获得利润 = 6 - 3 = 3 。
//总利润为 4 + 3 = 7 。

func maxProfit(prices []int) int {
	var income int
	for i, p := range prices {
		if i > 0 {
			if p-prices[i-1] > 0 {
				income += p - prices[i-1]
			}
		}
	}
	return income
}
