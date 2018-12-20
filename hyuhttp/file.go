package hyuhttp

import (
	"fmt"
	"net/http"
	"strings"
)

// FileData 扱うファイルの情報
type FileData struct {
	Name string
	Bin  []byte
}

// GetContentType 拡張子からContentTypeを判定する
func GetContentType(f *FileData) string {

	str := strings.Split(f.Name, ".")

	if len(str) <= 1 {
		return "application/octet-stream"
	}

	switch str[len(str)-1] {
	case "txt":
		return "text/plain"
	case "csv":
		return "text/csv"
	case "html", "htm":
		return "text/html"
	case "png":
		return "image/png"
	case "pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}

}

// DownloadFile Fileをダウンロード"させる"
func DownloadFile(f *FileData, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Disposition", "attachment; filename="+f.Name)
	w.Header().Set("Content-Type", GetContentType(f))
	w.Header().Set("Content-Length", fmt.Sprint(len(f.Bin)))

	w.Write(f.Bin)

}
