package hyucsv

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

	"github.com/gara-snake/hyutil"
	"github.com/gara-snake/hyutil/hyudb"
)

const (
	extension      = ".csv"
	fieldSeparater = ":"
)

// Csv ファイル
type Csv struct {
	Name   string
	Fields []CsvField
	Rows   []CsvRow
}

// CsvField Csvファイルのカラム名
type CsvField struct {
	Label string
	Key   string
}

// String 文字列化
func (cf *CsvField) String() string {
	return cf.Key + fieldSeparater + cf.Label
}

// CsvRow 行情報
type CsvRow map[string]string

// Decode CSVを解析して構造体に設定する
func (row *CsvRow) Decode(obj interface{}) {

	buf := make(map[string]string)

	for k, v := range *row {
		buf[k] = v
	}

	hyutil.ObjFill(obj, buf, false)

}

func addEx(name string) string {
	if strings.HasSuffix(name, extension) {
		return name
	}
	return name + extension
}

// Create Csvデータを作成する
func Create(name string, fields []CsvField, data []interface{}) *Csv {

	csv := &Csv{
		Name:   addEx(name),
		Fields: fields,
	}

	for _, d := range data {

		val := reflect.ValueOf(d)
		tp := val.Type()

		if tp.Kind() == reflect.Ptr {
			val = val.Elem()
			tp = tp.Elem()
		}

		row := createCsvRow(val, tp)

		csv.Rows = append(csv.Rows, row)

	}

	return csv
}

// CreateType CSV作成種別
type CreateType int

// CreateTypeKeyLabel "キー値:ラベル"で格納する
var CreateTypeKeyLabel CreateType = 0

// CreateTypeLabel label値のみで格納する
var CreateTypeLabel CreateType = 1

// CreateFromFile Csvデータを作成する
func CreateFromFile(name string, fields []CsvField, r io.Reader) *Csv {
	return CreateFromFileEx(name, fields, r, CreateTypeKeyLabel)
}

// CreateFromFileEx Csvデータを作成する
func CreateFromFileEx(name string, fields []CsvField, r io.Reader, opt CreateType) *Csv {

	csv := &Csv{
		Name:   addEx(name),
		Fields: fields,
	}

	reader := newCsvReader(r)
	records, e := reader.ReadAll()

	if e != nil {
		log.Println(e)
		return csv
	}

	if len(records) < 2 {
		// 行なし
		return csv
	}

	header := records[0]
	detail := records[1:]

	for _, row := range detail {

		csvRow := CsvRow{}

		for _, f := range fields {
			for i, h := range header {
				h = strings.Trim(h, " ")
				h = strings.Trim(h, "　")

				var headerText string

				switch opt {
				case CreateTypeLabel:
					headerText = f.Label
				default:
					headerText = f.String()
				}

				if headerText == h {
					val := strings.Trim(row[i], " ")
					val = strings.Trim(row[i], "　")
					csvRow[f.Key] = val

					break
				}
			}
		}

		csv.Rows = append(csv.Rows, csvRow)
	}

	return csv
}

func createCsvRow(val reflect.Value, tp reflect.Type) CsvRow {

	ret := make(CsvRow)

	for i := 0; i < tp.NumField(); i++ {

		field := tp.Field(i)

		key := field.Tag.Get("json")

		if key == "" {
			key = strings.ToLower(hyutil.CamelToSnake(field.Name))
		}

		switch v := val.FieldByName(field.Name).Interface().(type) {
		case string:
			ret[key] = v
		case hyutil.DateTime:
			ret[key] = v.String()
		case hyudb.DBID:
			if v <= 0 {
				ret[key] = "NULL"
			} else {
				ret[key] = fmt.Sprint(v)
			}
		case bool:
			if v {
				ret[key] = "1"
			} else {
				ret[key] = "0"
			}
		case nil:
			ret[key] = "NULL"
		default:
			ret[key] = fmt.Sprint(v)
		}

	}

	return ret
}

// Buffer 文字列化表現を内包した bytes.Buffer
func (csv *Csv) Buffer() *bytes.Buffer {

	var buf bytes.Buffer

	w := newCsvWriter(io.Writer(&buf), true)

	fields := make([]string, 0)
	for _, f := range csv.Fields {
		fields = append(fields, f.String())
	}

	w.Write(fields)

	for _, row := range csv.Rows {
		cols := make([]string, 0)

		for _, f := range csv.Fields {

			cols = append(cols, row[f.Key])

		}
		w.Write(cols)
	}

	w.Flush()

	return &buf
}

// String 文字列化（BOMは考えない）
func (csv *Csv) String() string {
	return csv.Buffer().String()
}

// BOMつき読み取り
func newCsvReader(r io.Reader) *csv.Reader {
	br := bufio.NewReader(r)
	bs, err := br.Peek(3)
	if err != nil {
		return csv.NewReader(br)
	}
	if bs[0] == 0xEF && bs[1] == 0xBB && bs[2] == 0xBF {
		br.Discard(3)
	}
	return csv.NewReader(br)
}

// BOMつき書き込み
func newCsvWriter(w io.Writer, bom bool) *csv.Writer {
	bw := bufio.NewWriter(w)
	if bom {
		bw.Write([]byte{0xEF, 0xBB, 0xBF})
	}

	writer := csv.NewWriter(bw)
	writer.UseCRLF = true

	return writer
}
