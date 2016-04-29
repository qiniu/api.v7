package main

import (
  "kodo"
  "conf"
  "fmt"
)

var (
  //设置需要操作的空间
  bucket = "yourbucket"
  //设置需要操作的文件的key
  key = "yourkey"
  //设置移动后文件的文件名
  movekey = "movekey"
)

func main() {

  conf.ACCESS_KEY = "ACCESS_KEY"
  conf.SECRET_KEY = "SECRET_KEY"

  //new一个Bucket管理对象
  c := kodo.New(0, nil)
    p := c.Bucket(bucket)

    //调用Move方法移动文件
  res := p.Move(nil, key, movekey)

  //打印返回值以及出错信息
  if res == nil {
    fmt.Println("Move success")
  }else {
    fmt.Println("Move failed:",res)
  }
}