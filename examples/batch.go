package main

import (
	"qiniupkg.com/api.v7/kodo"
	"qiniupkg.com/api.v7/conf"
	"fmt"
)

func main() {

	//new一个数组，需要批量操作的数组
	entryPairs := []kodo.KeyPair {
		kodo.KeyPair{
			Src: "xxx.jpg",
			Dest: "xxxx.jpg",
		},kodo.KeyPair{
			Src: "xxxxx.jpg",
			Dest: "xxxxx.jpg",
		},
	}

	conf.ACCESS_KEY = "xxxx"
	conf.SECRET_KEY = "xxxx"

	//new一个Bucket对象
	c := kodo.New(0, nil)
    p := c.Bucket("xxxx")

    //调用BatchCopy方法
	batchCopyRets, err := p.BatchCopy(nil, entryPairs...)

	if err == nil {
		fmt.Println("Move success")
	}else {
		fmt.Println("Move failed:",err)
	}
	
	//循环遍历每个操作的返回结果
	for _, item := range batchCopyRets {
	    fmt.Println(item.Code, item.Error)
	}

}








