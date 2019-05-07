package ut

/**
求公共前缀
返回两个字符串公共的序号
-1 代表无重复部分
*/
func CommonPrefix(a, b string) int {
	s := Min(len(a), len(b))
	for i := 0; i < s; i++ {
		if a[i] != b[i] {
			return i - 1
		}
	}
	return s - 1
}

func Min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
