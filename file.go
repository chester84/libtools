package libtools

import (
	"archive/zip"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/h2non/filetype"
)

const fileChunk = 8192 // we settle for 8KB

// BuildFileHashName 创建本地文件的hash名
func BuildFileHashName(localFile string) (hashDir, hashName, fileMd5 string, err error) {
	file, err := os.Open(localFile)
	if err != nil {
		return
	}
	defer file.Close()

	// calculate the file size
	info, _ := file.Stat()
	filesize := info.Size()
	blocks := uint64(math.Ceil(float64(filesize) / float64(fileChunk)))
	hash := md5.New()
	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(fileChunk, float64(filesize-int64(i*fileChunk))))
		buf := make([]byte, blocksize)

		_, _ = file.Read(buf)
		_, _ = io.WriteString(hash, string(buf)) // append into the hash
	}

	fileMd5 = fmt.Sprintf("%x", hash.Sum(nil)) // 文件md5值
	//fileSuffix := path.Ext(localFile)          //获取文件后缀
	fileSuffix := GetFileExt(localFile) //获取文件后缀

	hashDir, hashName = BuildHashName(fileMd5, fileSuffix)

	return
}

// 为app端上传文件生成hash文件名
func BuildUploadFileHashName(buf []byte, suffix string) (hashDir, hashName, fileMd5 string) {
	fileMd5 = Md5Bytes(buf)
	hashDir, hashName = BuildHashName(fileMd5, suffix)

	return
}

func GetS3Key(fileName string) string {
	fileBytes, _ := ioutil.ReadFile(fileName)
	fileMd5 := Md5Bytes(fileBytes)
	var extension string
	index := strings.LastIndex(fileName, ".")
	extension = fileName[index+1:]
	_, s3Key := BuildHashName(fileMd5, extension)
	return s3Key
}

func BuildHashName(fileMd5, suffix string) (hashDir, hashName string) {
	hashDir = fmt.Sprintf("%s/%s/%s", GetCurrentEnv(), SubString(fileMd5, 0, 2), SubString(fileMd5, 2, 4))
	// [dev|pro]/XX/YYYY/fileMd5.后缀
	hashName = fmt.Sprintf("%s/%s.%s", hashDir, fileMd5, suffix)

	return
}

// 生成本地hash目录
func LocalHashDir(hashDir string) string {
	return fmt.Sprintf("%s/%s", GetLocalUploadPrefix(), hashDir)
}

func DetectFileType(filename string) (string, string, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return "unknown", "", err
	}

	return DetectFileByteType(buf)
}

func DetectFileByteType(buf []byte) (extension, mime string, err error) {
	kind, unknown := filetype.Match(buf)
	if unknown != nil {
		extension = "unknown"
		err = unknown
		return
	}

	extension = kind.Extension
	mime = kind.MIME.Value

	return
}

// 简易版取文件名后缀,path.Ext()方法会带着个`.`
func GetFileExt(filename string) (suffix string) {
	exp := strings.Split(filename, ".")
	expLen := len(exp)
	if expLen > 1 {
		suffix = exp[expLen-1]
	}

	return
}

// 安全删除文件
func Remove(filename string) (err error) {
	_, err = os.Stat(filename)
	if err != nil {
		logs.Warning("file does not exist: ", filename)
		return
	}

	err = os.Remove(filename)

	return
}

var gitRevParseHead string = ""

func GitRevParseHead() string {
	if len(gitRevParseHead) <= 0 {
		filename := "conf/git-rev-hash"

		_, err := os.Stat(filename)
		if err != nil {
			logs.Error("[GitRevParseHead] file does not exist, filename:", filename)
			gitRevParseHead = "-1"
			return gitRevParseHead
		}

		const bufferSize = 32
		file, err := os.Open(filename)
		if err != nil {
			logs.Error("[GitRevparseHead] can NOT open filename:", filename)
			gitRevParseHead = "-2"
			return gitRevParseHead
		}

		defer file.Close()

		buffer := make([]byte, bufferSize)
		bytesRead, err := file.Read(buffer)
		if err != nil {
			logs.Error("[GitRevparseHead] can read hash data from file:", filename)
			gitRevParseHead = "-3"
			return gitRevParseHead
		}

		gitRevParseHead = string(buffer[:bytesRead])
	}

	return gitRevParseHead
}

func FileDownload(fileName, url string) (realFileName string, err error) {
	realFileName = fmt.Sprintf("/tmp/%s", fileName)
	res, err := http.Get(url)
	if err != nil {
		logs.Error("[FileDownload] Get file failed, err:", err)
		return
	}

	defer res.Body.Close()

	f, err := os.Create(realFileName)
	if err != nil {
		logs.Error("[FileDownload] Create file failed, err:", err)
		return
	}
	defer f.Close()

	_, _ = io.Copy(f, res.Body)

	return
}

func GetFileContentType(out multipart.File) (string, error) {
	// 只需要前 512 个字节就可以了
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func GetFileType(out multipart.File) (string, error) {
	// 只需要前 512 个字节就可以了
	buf := make([]byte, 512)

	_, err := out.Read(buf)
	if err != nil {
		return "", err
	}

	t, err := filetype.Get(buf)

	return t.MIME.Value, err
}

func GetFileExtension(f multipart.File, h *multipart.FileHeader) (string, error) {
	contentType := ""
	fileContentType, err := GetFileContentType(f)
	switch fileContentType {
	case "image/jpeg":
		contentType = "jpeg"
	case "image/png":
		contentType = "png"
	case "image/gif":
		contentType = "gif"
	case "application/octet-stream":
		sArr := strings.Split(h.Filename, ".")
		contentType = sArr[len(sArr)-1]
	case "application/zip":
		sArr := strings.Split(h.Filename, ".")
		contentType = sArr[len(sArr)-1]
	case "application/pdf":
		contentType = "pdf"
	}

	return contentType, err
}

// ZipDirectory 将整个目录压缩成一个 zip 文件
func ZipDirectory(sourceDir, zipFileName string) error {
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("创建 ZIP 文件失败: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(filepath.Dir(sourceDir), path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// **使用 CreateHeader 以支持压缩**
		header := &zip.FileHeader{
			Name:   relPath,
			Method: zip.Deflate, // 启用压缩
		}
		header.SetModTime(info.ModTime()) // 保留原文件时间戳
		header.SetMode(info.Mode())       // 保留原文件权限

		zipFileWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipFileWriter, file)
		return err
	})

	return err
}
