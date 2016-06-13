package main

import (
  "qiniupkg.com/api.v7/kodo"
  "qiniupkg.com/api.v7/conf"
  "qiniupkg.com/api.v7/kodocli"
  "fmt"
)

var (
    //设置上传到的空间
    bucket = "yourbucket"
    //设置上传文件的key
    key = "yourdefinekey"
    //设置转码参数
    fops = "avthumb/mp4/s/640x360/vb/1.25m"
    //设置转码用的队列
    pipeline = "yourpipeline"
)

//构造返回值字段
type PutRet struct {
    Hash    string `json:"hash"`
    Key     string `json:"key"`
    PersistentId  string `json:"persistentId"`
}

func main() {
    //初始化AK，SK
    conf.ACCESS_KEY = "ACCESS_KEY"
    conf.SECRET_KEY = "SECRET_KEY"

    //创建一个Client
    c := kodo.New(0, nil)

    //设置上传的策略
    policy := &kodo.PutPolicy{
        Scope:   bucket+ ":" + key,
        //设置Token过期时间
        Expires: 3600,
        InsertOnly: 1,
        PersistentOps: fops,
        PersistentPipeline: pipeline,
    }
    //生成一个上传token
    token := c.MakeUptoken(policy);

    //构建一个uploader
    zone := 0
    uploader := kodocli.NewUploader(zone, nil)

    var ret PutRet
    //设置上传文件的路径
    filepath := "/xxx/xxx/sample.flv"
    //调用PutFile方式上传，这里的key需要和上传指定的key一致
    res := uploader.PutFile(nil, &ret, token, key, filepath, nil)
    //打印返回的信息
    fmt.Println(ret)
    //打印出错信息
    if res != nil {
        fmt.Println("io.Put failed:", res)
        return
    }   

}