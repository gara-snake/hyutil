package hyudb_test

import (
	"hyutil"
	"hyutil/hyudb"
	"testing"

	"github.com/cheekybits/is"

	_ "github.com/go-sql-driver/mysql"
)

var connectionString = "root:root@tcp(localhost:3306)/asobism_affairs?parseTime=true&loc=Asia%2FTokyo"

type TestObj struct {
	ID      int64 `hyudb:"pk"`
	Name    string
	Age     int32
	Rate    float64
	Invalid bool
	InsDate hyutil.DateTime
	UpdDate hyutil.DateTime
}

func (t *TestObj) TableName() string {
	return "test"
}

func TestSelect(t *testing.T) {

	is := is.New(t)

	db := hyudb.MysqlNew(connectionString)

	tbl := db.SelectQuery(" SELECT * FROM hr_employee ")

	is.Equal(1, len(tbl.Rows))
	is.Equal("大志", tbl.Rows[0].Columns["first_name"])

}

func TestGet(t *testing.T) {

	is := is.New(t)

	db := hyudb.MysqlNew(connectionString)

	obj := &TestObj{
		ID: 1,
	}

	db.Debug = true
	db.Get(obj)

	is.Equal(1, obj.ID)
	is.Equal("テスト太郎", obj.Name)
	is.Equal(15, obj.Age)
	is.Equal(false, obj.Invalid)

	is.Equal(123.456, obj.Rate)

	is.Equal("2018-10-26 14:23:05", obj.InsDate.Format("2006-01-02 15:04:05"))
	is.Equal("2018-10-26 14:24:06", obj.UpdDate.Format("2006-01-02 15:04:05"))

}
