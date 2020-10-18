package utils

import "stu-manager/logger"

// 获取最小的可插入值
// 如果中间不可以插入就默认返回max($slice)+$minGap
func GetMinInsertValueWhenGap(slice []int64, minGap int64) int64 {
	if slice == nil {
		return minGap
	}
	var result int64
	//先对slice进行排序
	slice = BubbleAsort(slice)
	logger.Info.Printf("sorted slice:%v", slice)
	var index int
	for _, value := range slice {
		result = value + minGap
		if index == len(slice)-1 { //已经到了最后一项
			return result
		}
		nextValue := slice[index+1]
		if result < nextValue { //可以插入
			return result
		}
		index++
	}
	return result
}

// 冒泡升序排序，直接照搬:https://studygolang.com/articles/6127
func BubbleAsort(values []int64) []int64 {
	if values == nil {
		return nil
	}
	for i := 0; i < len(values)-1; i++ {
		for j := i + 1; j < len(values); j++ {
			if values[i] > values[j] {
				values[i], values[j] = values[j], values[i] //go的元素交换
			}
		}
	}
	return values
}
