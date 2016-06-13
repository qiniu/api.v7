package main

import (
  "qiniupkg.com/api.v7/kodo"
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

    //调用Delete方法删除文件
  res := p.Delete(nil, key)
  //打印返回值以及出错信息
  if res == nil {
    fmt.Println("Delete success")
  }else {
    fmt.Println(res)
  }
}