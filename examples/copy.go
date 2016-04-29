package main

import (
  "github.com/qiniu/api.v7/kodo"
  "qiniupkg.com/api.v7/conf"
  "fmt"
)

var (
  //设置需要操作的空间
  bucket = "yourbucket"
  //设置需要操作的文件的key
  key = "yourkey"
  //设置复制后文件的文件名
  copykey = "yourcopykey"
)

func main() {

  conf.ACCESS_KEY = "ACCESS_KEY"
  conf.SECRET_KEY = "SECRET_KEY"

  //new一个Bucket管理对象
  c := kodo.New(0, nil)
    p := c.Bucket(bucket)

  //调用Copy方法移动文件
  res := p.Copy(nil, key, copykey)

  //打印返回值以及出错信息
  if res == nil {
    fmt.Println("Copy success")
  }else {
    fmt.Println("Copy failed:",res)
  }
}