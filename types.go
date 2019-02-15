package api

// Bool 获取一个bool类型的指针
func Bool(v bool) *bool {
	return &v
}

// String获取一个字符串指针
func String(v string) *string {
	return &v
}

// Int64 获取一个整形的指针
func Int64(v int64) *int64 {
	i := int64(v)
	return &i
}
