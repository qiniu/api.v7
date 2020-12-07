# Changelog

## 7.8.0 (2020-12-07)
对象存储：
* 新增[v2版分片上传](https://developer.qiniu.com/kodo/api/6364/multipartupload-interface)支持

## 7.7.0 (2020-11-27)
QVS视频监控：
* 新增录制管理相关api
* 新增数据统计管理相关api
* 流管理api增加按需截图，删除截图
* 设备管理api增加通道管理，查询本地录像回放列表
* 空间增加ondemandPull字段，模板修改recordFileFormat字段说明，增加mp4，flv录制格式字段

## 7.6.0（2020-09-04）
* QVS视频监控：增加设备管理模块功能Device；增加流管理停止推流功能StopStream API；修改历史单测问题

## 7.5.0（2020-05-29）
* 增加了QVS视频监控产品功能

## 7.4.2 (2020-04-26)
* 为 Go SDK 添加了 MIT License

## 7.4.1 (2019-12-09)
* 修正 FormUploader.put 方法 extra 参数为 nil 时会导致空指针问题
* 增加多个 Bucket 管理相关功能

## 7.4.0 (2019-10-30)
* 添加 linking 模块
* PutPolicy 添加 fsizeMin  forceSaveKey 属性

## 7.3.1 (2019-10-18)
* 支持 go.mod
* 要求 go 1.10 及以上

## 7.2.4 (2018-03-01)
* 增加新加坡机房，新机房上线
* 增加异步fetch的功能
* 修复构建访问外链时兼容域名带/后缀
* 默认开启crc32校验功能，表单+分片上传
* 使用go内置的context功能
* 修复qiniu rpc并发上传共用token的bug
* 增加七牛云rtc服务端功能

## 7.2.3 (2017-09-25)
* 增加Qiniu的鉴权方式
* 删除prefop域名检测功能
* 暴露分片上传的接口以支持复杂的自定义业务逻辑

## 7.2.2 (2017-09-19)
* 为表单上传和分片上传增加代理支持
* 优化表单上传的crc32计算方式，减少内存消耗
* 增加网页图片的Base64上传方式

## 7.2.1 (2017-08-20)
* 设置FormUpload默认支持crc32校验
* ResumeUpload从API层面即支持crc32校验

## 7.2.0 (2017-07-28)
* 重构了v7 SDK的所有代码

## 7.1.0 (2016-6-22)

### 增加
* 增加多机房相关功能

## 7.0.5 (2015-11-20)

### 增加
* add delimiter support to Bucket.List
* 增加回调校验

## 7.0.4 (2015-09-03)

### 增加
* 上传返回参数PutRet增加PersistentId，用于获取上传对应的fop操作的id

### 修复
* token 覆盖问题

## 7.0.3 (2015-07-11)

### 增加
* support NestedObject

## 7.0.2 (2015-07-7-10)

### 增加
* 增加跨空间移动文件(Bucket.MoveEx)

## 7.0.1 (2015-07-7-10)

### 增加
* 完善 PutPolicy：支持 MimeLimit、CallbackHost、CallbackFetchKey、 CallbackBodyType、 Checksum

## 7.0.0 (2016-06-29)

* 重构，初始版本
