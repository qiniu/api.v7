package main

import (
	"qiniupkg.com/api.v7/kodo"
	"qiniupkg.com/api.v7/conf"
	"fmt"
)

var (
	//指定需要抓取到的空间
	bucket = "xxxx"
	//指定需要抓取的文件的url，必须是公网上面可以访问到的
	target_url = "xxxx"
	//指定抓取保存到空间的文件的key指
	key = "test.jpg"
)

func main() {

	conf.ACCESS_KEY = "xxxx"
	conf.SECRET_KEY = "xxxx"

	//new一个Bucket对象
	c := kodo.New(0, nil)
    p := c.Bucket(bucket)

    //调用Fetch方法
	err := p.Fetch(nil, key, target_url)
	if err != nil {
		fmt.Println("bucket.Fetch failed:", err)
	}else {
		fmt.Println("fetch success")
	}
}