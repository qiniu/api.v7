package auth

import (
	"context"
)

// MacContextKey 是用户的密钥信息
// context.Context中的键值不应该使用普通的字符串， 有可能导致命名冲突
type macContextKey struct{}

// WithCredentials 返回一个包含密钥信息的context
func WithCredentials(ctx context.Context, cred *Credentials) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, macContextKey{}, cred)
}

// CredentialsFromContext 从context获取密钥信息
func CredentialsFromContext(ctx context.Context) (cred *Credentials, ok bool) {
	cred, ok = ctx.Value(macContextKey{}).(*Credentials)
	return
}
