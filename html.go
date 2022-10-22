package libtools

import (
	"bytes"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/beego/beego/v2/core/logs"
)

func CKEditorFullHtml(content string) string {
	buf := bytes.NewBufferString(content)
	doc, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		logs.Warning("[CKEditorFullHtml] goquery can not parse input, content: %s", content)
		return ""
	}

	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		selection.RemoveAttr("href")
	})
	head := doc.Find("head")
	if head != nil {
		head.AppendHtml("<meta charset='UTF-8'>")
		head.AfterHtml(`<link type="text/css" rel="stylesheet" href="https://cdn.ckeditor.com/4.13.1/standard/contents.css">`)
		head.AfterHtml(`<link type="text/css" rel="stylesheet" href="https://cdn.ckeditor.com/4.13.1/standard/plugins/tableselection/styles/tableselection.css">`)
	}

	code, err := doc.Html()
	if err != nil {
		logs.Error("[CKEditorFullHtml] get full html exception, err: %v", err)
	}

	return code
}

func BuildFaceBookUserAvatar(userID string) string {
	return fmt.Sprintf(`https://graph.facebook.com/%s/picture?type=normal`, userID)
}
