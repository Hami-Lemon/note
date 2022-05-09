//二分查找

package main

import "fmt"

func main() {
	l := []int{1, 2, 3, 4, 5, 6}
	fmt.Printf("%d\n", binarySearch(l, 3))
	fmt.Printf("%d\n", binarySearch(l, 1))
	fmt.Printf("%d\n", binarySearch(l, 6))
	fmt.Printf("%d\n", binarySearch(l, 7))
}

func binarySearch(arr []int, target int) int {
	i, j := 0, len(arr)-1
	for i <= j {
		p := i + (j-i)/2 //防止数据过大溢出
		if arr[p] == target {
			return p
		} else if arr[p] < target {
			i = p + 1
		} else {
			j = p - 1
		}
	}
	return -1
}
