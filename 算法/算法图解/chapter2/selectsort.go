//选择排序

package main

import "fmt"

func main() {
	l := []int{4, 1, 10, 3, 2, 5, 0}
	sort(l)
	fmt.Printf("%v\n", l)
}

//每次都选择数组中最小的元素,O(n^2)
func sort(l []int) {
	for i := 0; i < len(l)-1; i++ {
		min := i
		for j := i + 1; j < len(l); j++ {
			if l[j] < l[min] {
				min = j
			}
		}
		l[i], l[min] = l[min], l[i]
	}
}
