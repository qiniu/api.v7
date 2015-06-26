/*
包 qiniupkg.com/api.v7/kodocli 提供了在客户端（比如：Android/iOS 设备、Windows/Mac/Linux 桌面环境）调用七牛云存储部分服务的能力

客户端，严谨说是非可信环境，主要是指在用户端执行的环境，比如：Android/iOS 设备、Windows/Mac/Linux 桌面环境、也包括浏览器（如果浏览器能够执行 Go 语
言代码的话）。

注意，在这种场合下您不应该在任何地方配置 AccessKey/SecretKey。泄露 AccessKey/SecretKey 如同泄露您的用户名/密码一样十分危险，
会影响您的数据安全。

第一个问题是如何上传文件。因为是在非可信环境，所以我们首先是要授予它有上传文件的能力。答案是给它颁发上传凭证：

*/
package kodocli
