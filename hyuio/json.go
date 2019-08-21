package hyuio

import "encoding/json"

// SaveJSON 指定のファイル名でJsonを保存する
func SaveJSON(v interface{}, fileName string) error {
	b, e := json.Marshal(v)
	if e != nil {
		return e
	}
	if e := SaveFile(b, fileName); e != nil {
		return e
	}

	return nil
}
