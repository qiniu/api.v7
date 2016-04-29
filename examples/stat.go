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
)

func main() {

  conf.ACCESS_KEY = "ACCESS_KEY"
  conf.SECRET_KEY = "SECRET_KEY"

  //new一个Bucket管理对象
  c := kodo.New(0, nil)
    p := c.Bucket(bucket)

    //调用Stat方法获取文件的信息
  entry, err := p.Stat(nil, key)
  //打印列取的信息
  fmt.Println(entry)
  //打印出错时返回的信息
  if err != nil {
    fmt.Println(err)
  }