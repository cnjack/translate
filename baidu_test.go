package translate

import "testing"

func TestBaiduTranslate_Translate(t *testing.T) {
	ts := NewBaiduTranslate("xxxx", "xxxxxxx", BaiduHttpsApiGateway)
	dst, err := ts.Translate("你好啊", Option{
		From: BZhLang,
		To:   BEnLang,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("translate 你好啊 => %s", dst)
}
