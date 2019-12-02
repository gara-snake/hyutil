package hyudb

//mainに以下が必要
//import _ "github.com/go-sql-driver/mysql"

import (
	"database/sql"
	"errors"
	"fmt"
	"hyutil"
	"log"
	"reflect"
	"strings"
)

// NoID はInt型プライマリーキーの新規値です
const NoID = 0

const dbTimeFormat = "2006-01-02T15:04:05-07:00"

//DB への参照です
type DB struct {
	IsOpen      bool
	Debug       bool
	connection  *sql.DB
	transaction *sql.Tx
	hasErr      bool
}

//Row カラム名ごとに文字列型で値を代入したMap
type Row struct {
	Columns map[string]string
}

//Table 行の集合体
type Table struct {
	Rows []Row
}

//Modeler モデルインターフェイス
type Modeler interface {
	TableName() string
}

// DBID id型 NOT NULL 成約のために必要
type DBID int64

//DbBool はデータベース上でのBool値の表現文字列を返します。
func DbBool(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

//DbEsc は文字列をエスケープ処理します
func DbEsc(s string) string {
	if s == "" {
		return "NULL"
	}
	return "'" + strings.Replace(s, "'", "\\'", -1) + "'"

}

//DbNum は数値の文字列表現を返します
func DbNum(num interface{}) string {
	return fmt.Sprint(num)
}

const dbDatetimeFormat = "2006-01-02 15:04:05"

//DbDt はデータベース上でのDatetime値の表現文字列を返します。
func DbDt(t *hyutil.DateTime) string {

	if t == nil {
		return " NULL "
	}

	return "'" + t.Format(dbDatetimeFormat) + "'"

}

// Nval 空文字列を空白に変換します
func Nval(s string) string {
	if s == "" {
		return " "
	}
	return s
}

//ColEsc はカラム名のエスケープです
func ColEsc(s string) string {
	return "`" + s + "`"
}

//New データベースへの新規接続を開始します
func New(dbType string, connectionstr string) *DB {

	db, err := sql.Open(dbType, connectionstr)

	if err != nil {
		log.Fatalln(err)
		return &DB{IsOpen: false}
	}

	return &DB{
		IsOpen:     true,
		Debug:      false,
		connection: db,
		hasErr:     false,
	}

}

//MysqlNew 任意のMysqlサーバへの接続を開始します
func MysqlNew(connectionstr string) *DB {
	return New("mysql", connectionstr)
}

//BeginTx トランザクションを開始します
func (db *DB) BeginTx() {

	tx, err := db.connection.Begin()

	if err != nil {
		log.Fatalln(err)
		return
	}

	db.transaction = tx

}

//Rollback トランザクションが開始されている場合、ロールバックします
func (db *DB) Rollback() {

	if db.transaction != nil {

		db.transaction.Rollback()
		db.transaction = nil
	}

}

//Close DBへの接続を閉じます。未完了のトランザクションは”コミット”されます
func (db *DB) Close() {

	if err := recover(); err != nil {
		db.Rollback()
	}

	if db.connection != nil {

		if db.transaction != nil {

			if db.hasErr {
				db.Rollback()
			} else {
				//オートコミット
				db.transaction.Commit()
			}
		}

		db.connection.Close()
		db.IsOpen = false
		db.connection = nil
	}

}

//Exec INSERT、UPDATE、DELETEを実行します RowsAffected LastInsertId
func (db *DB) Exec(query string) (int64, int64) {

	if db.Debug {
		log.Println("EXEC QUERY : " + query)
	}

	var result sql.Result
	var err error

	if db.transaction == nil {
		result, err = db.connection.Exec(query)
	} else {
		result, err = db.transaction.Exec(query)
	}

	if err != nil {
		log.Fatalln(err)
		db.hasErr = true
	}

	ret1, err := result.RowsAffected()
	ret2, err := result.LastInsertId()

	if err != nil {
		log.Println(err)
		ret1 = 0
		ret2 = NoID
		db.hasErr = true
		return -1, -1
	}

	return ret1, ret2

}

// SelectExists queryで行が取得できたかどうかを返却します
func (db *DB) SelectExists(query string) bool {

	if db.Debug {
		log.Println("SELECT QUERY : " + query)
	}

	var rows *sql.Rows
	var err error

	if db.transaction == nil {
		rows, err = db.connection.Query(query)
	} else {
		rows, err = db.transaction.Query(query)
	}

	if err != nil {
		log.Fatalln(err)
		db.hasErr = true
		return false
	}

	defer rows.Close()

	return rows.Next()

}

// SelectTop queryを実行し、先頭の要素をDBFillします
func (db *DB) SelectTop(query string, model interface{}) error {

	tbl := db.SelectQuery(query)
	for _, r := range tbl.Rows {

		DBFill(model, &r)
		return nil

	}

	return errors.New("レコードが取得できませんでした。")
}

//SelectQuery SELECTを実行します
func (db *DB) SelectQuery(query string) *Table {

	if db.Debug {
		log.Println("SELECT QUERY : " + query)
	}

	var rows *sql.Rows
	var err error

	if db.transaction == nil {
		rows, err = db.connection.Query(query)
	} else {
		rows, err = db.transaction.Query(query)
	}

	defer rows.Close()

	if err != nil {
		log.Fatalln(err)
		db.hasErr = true
		return nil
	}

	columns, err := rows.Columns()
	if err != nil {
		log.Fatalln(err)
		db.hasErr = true
		return nil
	}

	var ret = Table{
		Rows: make([]Row, 0),
	}

	values := make([]sql.NullString, len(columns))
	scanArgs := make([]interface{}, len(columns))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {

		err = rows.Scan(scanArgs...)

		if err != nil {
			log.Fatalln(err)
			db.hasErr = true
			return nil
		}

		cols := make(map[string]string, len(columns))

		for i, col := range values {

			cols[columns[i]] = col.String

		}

		var row = Row{
			Columns: cols,
		}

		ret.Rows = append(ret.Rows, row)

	}

	return &ret

}

// Get でmodelのプライマリーキーでデータを取得します。プライマリーが未指定の場合はデータが登録されません。
func (db *DB) Get(model interface{}) error {

	query, err := createSelectQuery(model)

	if err != nil {
		return err
	}

	tbl := db.SelectQuery(query)

	if len(tbl.Rows) <= 0 {
		return errors.New("DB.Get:該当するレコードがありません")
	}

	DBFill(model, &tbl.Rows[0])

	return nil
}

// DBFill はすでに存在するモデルにRowを展開します。プライマリーキーは考慮（再検索）されません。
func DBFill(model interface{}, row *Row) {

	hyutil.ObjFill(model, row.Columns, true)

}

func createSelectQuery(model interface{}) (string, error) {

	val := reflect.ValueOf(model)
	tp := val.Type()

	if tp.Kind() == reflect.Ptr {
		val = val.Elem()
		tp = tp.Elem()
	}

	if tp.Kind() != reflect.Struct {
		return "", errors.New("引数が構造体ではありません")
	}

	columns := make([]string, 0)
	pk := ""

	tableName := hyutil.CamelToSnake(tp.Name())

	if m, ok := model.(Modeler); ok {
		tableName = m.TableName()
	}

	var pkFieldName string

	for i := 0; i < tp.NumField(); i++ {

		field := tp.Field(i)

		key := field.Tag.Get("hyudb")

		if key == "non" {
			continue
		}

		// カラム名作成
		col := field.Tag.Get("hyudb_col")
		alias := ""

		if col == "" {
			col = field.Tag.Get("json")
		} else {
			// DBFillでhyudb_colが優先されたのでエイリアスは不要に
			// alias = field.Tag.Get("json")
		}

		if col == "" {
			col = strings.ToLower(hyutil.CamelToSnake(field.Name))
		}

		if key == "pk" {
			pk = col
			pkFieldName = field.Name
		}

		if alias != "" {
			col = ColEsc(col) + " AS " + alias
		} else {
			col = ColEsc(col)
		}

		//DB予約文字エスケープ
		columns = append(columns, col)

	}

	if pk == "" {
		return "", errors.New("PrimaryKeyが指定されていません")
	}

	keyprm := ""
	i := val.FieldByName(pkFieldName)

	switch v := i.Interface().(type) {
	case string:
		keyprm = DbEsc(v)
	case DBID:
		keyprm = fmt.Sprint(v)
	default:
		keyprm = fmt.Sprint(i.Interface())
	}

	query :=
		" SELECT " + strings.Join(columns, ",") +
			" FROM " + tableName +
			" WHERE " + pk + " = " + keyprm

	return query, nil
}

// Save 要素を作成または更新します
func (db *DB) Save(model interface{}) error {

	val := reflect.ValueOf(model)
	tp := val.Type()

	if tp.Kind() == reflect.Ptr {
		val = val.Elem()
		tp = tp.Elem()
	}

	if tp.Kind() != reflect.Struct {
		return errors.New("引数が構造体ではありません")
	}

	var pkVal *reflect.Value
	var pkName string

	isNew := false

	for i := 0; i < tp.NumField(); i++ {

		field := tp.Field(i)
		key := field.Tag.Get("hyudb")

		if key == "non" {
			continue
		}

		if key == "pk" {
			v := val.FieldByName(field.Name)
			pkVal = &v
			pkName = field.Name
			break
		}

	}

	if pkVal == nil {
		return errors.New("プライマリーキーの指定がありません")
	}

	switch v := pkVal.Interface().(type) {
	case int, int32, int64, DBID:
		isNew = (fmt.Sprint(v) == fmt.Sprint(NoID))
	default:
		return errors.New("プライマーキーの型が不明です。")
	}

	if isNew {
		query := createInsertQuery(model, val, tp)

		_, id := db.Exec(query)
		if id != NoID {
			pkVal.SetInt(id)
		}

	} else {
		query := createUpdateQuery(model, pkName, fmt.Sprint(pkVal.Interface()), val, tp)

		db.Exec(query)
	}

	return nil
}

const (
	mapIns int = iota
	mapUpd
)

func createInsertQuery(model interface{}, val reflect.Value, tp reflect.Type) string {

	columns := make([]string, 0)
	vals := make([]string, 0)

	tableName := hyutil.CamelToSnake(tp.Name())

	if m, ok := model.(Modeler); ok {
		tableName = m.TableName()
	}

	colVal := createColValMap(val, tp, mapIns)

	// カラム名と値文字列の順番を揃える
	for k, v := range colVal {
		columns = append(columns, k)
		vals = append(vals, v)
	}

	query :=
		" INSERT INTO " + tableName + " (" +
			strings.Join(columns, ",") +
			" ) VALUES ( " +
			strings.Join(vals, ",") +
			" ) "

	return query

}

func createUpdateQuery(model interface{}, pkName string, pkVal string, val reflect.Value, tp reflect.Type) string {

	sets := make([]string, 0)

	tableName := hyutil.CamelToSnake(tp.Name())

	if m, ok := model.(Modeler); ok {
		tableName = m.TableName()
	}

	colVal := createColValMap(val, tp, mapUpd)

	// カラム名と値文字列の順番を揃える
	for k, v := range colVal {
		sets = append(sets, k+" = "+v)
	}

	query :=
		" UPDATE " + tableName + " SET " +
			strings.Join(sets, ",") +
			" WHERE " + ColEsc(pkName) + " = " + pkVal

	return query
}

func createColValMap(val reflect.Value, tp reflect.Type, mode int) map[string]string {

	ret := make(map[string]string)

	for i := 0; i < tp.NumField(); i++ {

		field := tp.Field(i)

		col := field.Tag.Get("hyudb_col")

		if col == "" {
			col = field.Tag.Get("json")
		}

		key := field.Tag.Get("hyudb")

		if key == "pk" {
			continue
		}

		if key == "non" {
			continue
		}

		if col == "" {
			col = strings.ToLower(hyutil.CamelToSnake(field.Name))
		}

		//DB予約文字エスケープ
		col = ColEsc(col)

		switch v := val.FieldByName(field.Name).Interface().(type) {
		case string:
			ret[col] = DbEsc(v)
		case hyutil.DateTime:
			if v == hyutil.DateTimeZero {
				ret[col] = "NULL"
			} else {
				ret[col] = DbDt(&v)
			}
		case DBID:
			if v <= 0 {
				ret[col] = "NULL"
			} else {
				ret[col] = fmt.Sprint(v)
			}
		case bool:
			if v {
				ret[col] = "1"
			} else {
				ret[col] = "0"
			}
		case nil:
			ret[col] = "NULL"
		default:
			ret[col] = fmt.Sprint(v)
		}

	}

	return ret
}

// Del 要素を論理削除します TODO
func (db *DB) Del(model interface{}) error {

	return nil
}

func createDeleteQuery(model interface{}) (string, error) {

	return "", nil
}

func createDeleteForeverQuery(model interface{}) (string, error) {

	return "", nil
}
