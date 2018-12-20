package hyutil

import (
	"encoding/json"
	"log"
	"strconv"
	"time"
)

const parseTimeFormat1 = "2006-01-02T15:04:05-07:00"

const parseTimeFormat2 = "2006/01/02 15:04:05 -0700"

const parseTimeFormat3 = "2006-01-02T15:04:05Z"

const dateFormat = "2006-01-02"

//年齢計算用
const dateFormatOnlyNumber = "20060102"

//DateTimeZero ゼロ値
var DateTimeZero = DateTime{}

//DateTime はJsonで正しく変換できるよう書式を固定した日付型です
type DateTime struct {
	*time.Time
}

// NowDateTime 現在時間のDateTimeを作成する
func NowDateTime() DateTime {
	n := time.Now()
	dt := DateTime{}

	dt.Time = &n

	return dt
}

//DateTimeFormat はシステム内の標準日付、時間書式です
const DateTimeFormat = "2006/01/02 15:04:05 -0700"

//DatetimeParse 2006-01-02T15:04:05-07:00 形式の文字列をから時間を設定します
func DatetimeParse(s string) DateTime {
	if s == "" {
		return DateTimeZero
	}

	formats := []string{parseTimeFormat1, parseTimeFormat2, parseTimeFormat3, dateFormat}

	for _, f := range formats {

		t, err := time.Parse(f, s)

		if err == nil {
			if f == parseTimeFormat3 {
				//末尾Zの場合はUTCなのでローカル時間に変える
				if loc, e := time.LoadLocation("Asia/Tokyo"); e == nil {
					t = t.In(loc)
				}
			}

			return DateTime{&t}
		}
	}

	return DateTimeZero
}

//UnmarshalJSON はJson文字列から要素を取得する処理です
func (t *DateTime) UnmarshalJSON(data []byte) error {

	str := string(data)
	if str == "" || str == "\"\"" {
		*t = DateTimeZero
		return nil
	}

	time, err := time.Parse("\""+DateTimeFormat+"\"", str)
	if err != nil {
		log.Println("str : " + str)
	}
	*t = DateTime{&time}
	return err
}

//MarshalJSON はJson文字列化化する処理です
func (t DateTime) MarshalJSON() ([]byte, error) {
	if t == DateTimeZero {
		return json.Marshal("")
	}
	return json.Marshal(t.Format(DateTimeFormat))
}

//String は文字列化です
func (t *DateTime) String() string {
	if *t == DateTimeZero {
		return ""
	}
	return t.Format(DateTimeFormat)
}

//FirstDay 日にちを月初に設定します
func (t *DateTime) FirstDay() *DateTime {

	newTime := time.Date(t.Year(), t.Month(), 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)

	t.Time = &newTime

	return t
}

//LastDay 日にちを月末に設定します
func (t *DateTime) LastDay() *DateTime {

	newTime := time.Date(t.Year(), t.Month()+1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	newTime = newTime.AddDate(0, 0, -1)

	t.Time = &newTime

	return t
}

//SetHour 時を設定します
func (t *DateTime) SetHour(h int) *DateTime {

	if h < 0 && 23 < h {
		return t
	}

	newTime := time.Date(t.Year(), t.Month(), t.Day(), h, t.Minute(), t.Second(), t.Nanosecond(), time.Local)

	t.Time = &newTime

	return t
}

//SetMinutes 分を設定します
func (t *DateTime) SetMinutes(m int) *DateTime {

	if m < 0 && 59 < m {
		return t
	}

	newTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), m, t.Second(), t.Nanosecond(), time.Local)

	t.Time = &newTime

	return t
}

// FirstTime 時と分を00：00にします。
func (t *DateTime) FirstTime() *DateTime {
	return t.SetHour(0).SetMinutes(0)
}

// LastTime 時と分を23：59にします。
func (t *DateTime) LastTime() *DateTime {
	return t.SetHour(23).SetMinutes(59)
}

//CalcAgeNow 現在日時を基準に年齢相当の数値を返却する
func (t *DateTime) CalcAgeNow() int32 {

	now := time.Now().Format(dateFormatOnlyNumber)
	birthday := t.Time.Format(dateFormatOnlyNumber)

	nowInt, _ := strconv.Atoi(now)

	birthdayInt, _ := strconv.Atoi(birthday)

	return int32((nowInt - birthdayInt) / 10000)
}
