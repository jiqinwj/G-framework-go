package main

import "fmt"

func main() {

	/*
		   demo1
		输入array = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10], part = 2
		输出 【[1,2,3,4,5], [6,7,8,9,10]】
	*/
	var arr1 []int
	arr1 = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	res1 := SplitArray(arr1, 2)
	fmt.Println("res1:", res1)
	/*
			demo2
		输入array = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10], part = 3
		输出 【[1,2,3], [4,5,6], [7,8,9,10]】 或【[1,2,3,10], [4,5,6], [7,8,9]】等
	*/
	var arr2 []int
	arr2 = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	res2 := SplitArray(arr2, 3)
	fmt.Println("res2:", res2)
	/*
				  最终打印
			res1: [[1 2 3 4 5] [6 7 8 9 10]]
		    res2: [[1 2 3] [4 5 6] [7 8 9 10]]
	*/
}

func SplitArray(collection []int, size int) [][]int {
	if size <= 0 {
		panic("Second parameter must be greater than 0")
	}
	//每一份取多少数据
	chunksNum := len(collection) / size
	result := make([][]int, 0, size)
	//循环逻辑处理
	for i := 0; i < size; i++ {
		last := (i + 1) * chunksNum //切割的最后位置
		//判断是否大于总长度 如果大于就让它等于总长度 以免越界
		if last >= len(collection) {
			last = len(collection)
		}
		lastLen := last + 1 //判断是否到达最后一次切割。满足demo2 最终呈现的结果：[7 8 9 10]
		if lastLen == len(collection) {
			result = append(result, collection[i*chunksNum:lastLen])
		} else {
			result = append(result, collection[i*chunksNum:last])
		}

	}
	return result
}
