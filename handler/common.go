package handler

import (
	"main/app"
	"strconv"
)

const pageSize = 10 // 每页显示的数据量

func offset(page int) int {
	if page == 0 {
		return 0
	}
	return (page - 1) * pageSize
}

// 成功时返回的数据，格式化失败时返回空字符串
func succ(data interface{}) app.SuccRes {
	res := app.SuccRes{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
	return res
}

// 失败时返回的数据，格式化失败时返回空字符串
func failed(msg string) app.ErrRes {
	res := app.ErrRes{
		Code: -1,
		Msg:  msg,
	}
	return res
}

// 根据数据总数、当前页页码、每页显示的数据量计算要显示的页码(当前页的前2后3)
func pages(total, cur, num int) []string {
	var res []string
	ps := total / num // 总页数
	if (total % num) != 0 {
		ps++
	}
	for i := cur - 2; i <= ps; i++ {
		if len(res) == 6 {
			break
		}
		if i > 0 {
			k := strconv.Itoa(i)
			if i == cur {
				k = "_"
			}
			res = append(res, k)
		}
	}
	return res
}
