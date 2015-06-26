package kodo

import (
	"encoding/base64"
	"io"
	"net/url"
	"strconv"

	. "golang.org/x/net/context"
)

// ----------------------------------------------------------

func (p *Client) Batch(ctx Context, ret interface{}, op []string) (err error) {

	return p.CallWithForm(ctx, ret, "POST", p.RSHost+"/batch", map[string][]string{"op": op})
}

// ----------------------------------------------------------

type Bucket struct {
	Conn *Client
	Name string
}

func (p *Client) Bucket(name string) Bucket {

	return Bucket{p, name}
}

type Entry struct {
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize"`
	PutTime  int64  `json:"putTime"`
	MimeType string `json:"mimeType"`
	EndUser  string `json:"endUser"`
}

func (p Bucket) Stat(ctx Context, key string) (entry Entry, err error) {
	err = p.Conn.Call(ctx, &entry, "POST", p.Conn.RSHost+URIStat(p.Name, key))
	return
}

func (p Bucket) Delete(ctx Context, key string) (err error) {
	return p.Conn.Call(ctx, nil, "POST", p.Conn.RSHost+URIDelete(p.Name, key))
}

func (p Bucket) Move(ctx Context, keySrc, keyDest string) (err error) {
	return p.Conn.Call(ctx, nil, "POST", p.Conn.RSHost+URIMove(p.Name, keySrc, p.Name, keyDest))
}

func (p Bucket) Copy(ctx Context, keySrc, keyDest string) (err error) {
	return p.Conn.Call(ctx, nil, "POST", p.Conn.RSHost+URICopy(p.Name, keySrc, p.Name, keyDest))
}

func (p Bucket) ChangeMime(ctx Context, key, mime string) (err error) {
	return p.Conn.Call(ctx, nil, "POST", p.Conn.RSHost+URIChangeMime(p.Name, key, mime))
}

// ----------------------------------------------------------

type ListItem struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize"`
	PutTime  int64  `json:"putTime"`
	MimeType string `json:"mimeType"`
	EndUser  string `json:"endUser"`
}

// 1. 首次请求 marker = ""
// 2. 无论 err 值如何，均应该先看 entries 是否有内容
// 3. 如果后续没有更多数据，err 返回 EOF，markerOut 返回 ""（但不通过该特征来判断是否结束）
//
func (p Bucket) List(
	ctx Context, prefix, marker string, limit int) (entries []ListItem, markerOut string, err error) {

	listUrl := p.makeListURL(prefix, marker, limit)

	var listRet struct {
		Marker string     `json:"marker"`
		Items  []ListItem `json:"items"`
	}
	err = p.Conn.Call(ctx, &listRet, "POST", listUrl)
	if err != nil {
		return
	}
	if listRet.Marker == "" {
		return listRet.Items, "", io.EOF
	}
	return listRet.Items, listRet.Marker, nil
}

func (p Bucket) makeListURL(prefix, marker string, limit int) string {

	query := make(url.Values)
	query.Add("bucket", p.Name)
	if prefix != "" {
		query.Add("prefix", prefix)
	}
	if marker != "" {
		query.Add("marker", marker)
	}
	if limit > 0 {
		query.Add("limit", strconv.FormatInt(int64(limit), 10))
	}
	return p.Conn.RSFHost + "/list?" + query.Encode()
}

// ----------------------------------------------------------

type BatchStatItemRet struct {
	Data  Entry  `json:"data"`
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func (p Bucket) BatchStat(ctx Context, keys []string) (ret []BatchStatItemRet, err error) {

	b := make([]string, len(keys))
	for i, key := range keys {
		b[i] = URIStat(p.Name, key)
	}
	err = p.Conn.Batch(ctx, &ret, b)
	return
}

type BatchItemRet struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func (p Bucket) BatchDelete(ctx Context, keys []string) (ret []BatchItemRet, err error) {

	b := make([]string, len(keys))
	for i, key := range keys {
		b[i] = URIDelete(p.Name, key)
	}
	err = p.Conn.Batch(ctx, &ret, b)
	return
}

type KeyPair struct {
	Src  string
	Dest string
}

func (p Bucket) BatchMove(ctx Context, entries []KeyPair) (ret []BatchItemRet, err error) {

	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = URIMove(p.Name, e.Src, p.Name, e.Dest)
	}
	err = p.Conn.Batch(ctx, &ret, b)
	return
}

func (p Bucket) BatchCopy(ctx Context, entries []KeyPair) (ret []BatchItemRet, err error) {

	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = URICopy(p.Name, e.Src, p.Name, e.Dest)
	}
	err = p.Conn.Batch(ctx, &ret, b)
	return
}

// ----------------------------------------------------------

func encodeURI(uri string) string {
	return base64.URLEncoding.EncodeToString([]byte(uri))
}

func URIDelete(bucket, key string) string {
	return "/delete/" + encodeURI(bucket+":"+key)
}

func URIStat(bucket, key string) string {
	return "/stat/" + encodeURI(bucket+":"+key)
}

func URICopy(bucketSrc, keySrc, bucketDest, keyDest string) string {
	return "/copy/" + encodeURI(bucketSrc+":"+keySrc) + "/" + encodeURI(bucketDest+":"+keyDest)
}

func URIMove(bucketSrc, keySrc, bucketDest, keyDest string) string {
	return "/move/" + encodeURI(bucketSrc+":"+keySrc) + "/" + encodeURI(bucketDest+":"+keyDest)
}

func URIChangeMime(bucket, key, mime string) string {
	return "/chgm/" + encodeURI(bucket+":"+key) + "/mime/" + encodeURI(mime)
}

// ----------------------------------------------------------

