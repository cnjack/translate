# translate
百度翻译API接口

```
ts := NewBaiduTranslate("xxxx", "xxxxxxx", BaiduHttpsApiGateway)
dst, err := ts.Translate("你好啊", Option{
  From: BZhLang,
  To:   BEnLang,
})
// How do you do
```
