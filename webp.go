package libtools

//
//import (
//	"bytes"
//	"fmt"
//	"image"
//	"image/gif"
//	"image/jpeg"
//	"image/png"
//	"net/http"
//	"strings"
//
//	"github.com/beego/beego/v2/core/logs"
//	"github.com/chai2010/webp"
//	"golang.org/x/image/bmp"
//)
//
//func GetFileContentType(buffer []byte) string {
//	// Use the net/http package's handy DectectContentType function. Always returns a valid
//	// content-type by returning "application/octet-stream" if no others seemed to match.
//	contentType := http.DetectContentType(buffer)
//	return contentType
//}
//
//func WebpEncoder(origin []byte, quality float32) (*bytes.Buffer, error) {
//	// if convert fails, return error; success nil
//	var (
//		err error
//		img image.Image
//		buf = new(bytes.Buffer)
//	)
//
//	contentType := GetFileContentType(origin[:512])
//	if strings.Contains(contentType, "jpeg") {
//		img, _ = jpeg.Decode(bytes.NewReader(origin))
//	} else if strings.Contains(contentType, "png") {
//		img, _ = png.Decode(bytes.NewReader(origin))
//	} else if strings.Contains(contentType, "bmp") {
//		img, _ = bmp.Decode(bytes.NewReader(origin))
//	} else if strings.Contains(contentType, "gif") {
//		// TODO: need to support animated webp
//		logs.Warn("gif support is not perfect!")
//		img, _ = gif.Decode(bytes.NewReader(origin))
//	}
//
//	if img == nil {
//		msg := "image file is corrupted or not supported"
//		logs.Warning(msg)
//		err = fmt.Errorf(msg)
//		return buf, err
//	}
//
//	if err = webp.Encode(buf, img, &webp.Options{Lossless: false, Quality: quality}); err != nil {
//		logs.Error(err)
//		return buf, err
//	}
//
//	return buf, nil
//}
