package hyutil

import (
	"encoding/json"
	"time"
)

//DateTimeZero ゼロ値
var DateTimeZero = DateTime{}

//DateTime はJsonで正しく変換できるよう書式を固定した日付型です
type DateTime struct {
	*time.Time
}

//DateTimeFormat はシステム内の標準日付、時間書式です
const DateTimeFormat = "2006/01/02 15:04:05 -0700"

//UnmarshalJSON はJson文字列から要素を取得する処理です
func (t *DateTime) UnmarshalJSON(data []byte) error {

	time, err := time.Parse("\""+DateTimeFormat+"\"", string(data))
	*t = DateTime{&time}
	return err
}

//MarshalJSON はJson文字列化化する処理です
func (t DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format(DateTimeFormat))
}

//String は文字列化です
func (t *DateTime) String() string {
	return t.Format(DateTimeFormat)
}
