package file

import (
	"sort"
)

type int64array []int64
func (a int64array) Len() int { return len(a) }
func (a int64array) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a int64array) Less(i, j int) bool { return a[i] < a[j] }
func (a int64array) Sort() { sort.Sort(a) }
