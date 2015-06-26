/*
包 qiniupkg.com/api.v7/kodo 提供了在您的业务服务器（服务端）调用七牛云存储服务的能力

<b>首先，我们要创建一个 Client 对象：</b>

	zone := 0 // 您空间(Bucket)所在的区域
	c := kodo.New(zone, nil) // 使用默认配置创建 Client 对下

有了 Client，你就可以操作您的空间(Bucket)了，比如我们要上传一个文件：

	import "golang.org/x/net/context"

	ctx := context.Background()
	localFile := "/your/local/image/file.jpg"
	err := c.Bucket("your-bucket-name").PutFile(ctx, nil, "foo/bar.jpg", localFile, nil)
	if err != nil {
		... // 上传文件失败处理
		return
	}
	... // 上传文件成功，这时登陆七牛的 Portal，在 your-bucket-name 这个空间里面，你就可以看到一个 foo/bar.jpg 的文件了
*/
package kodo
