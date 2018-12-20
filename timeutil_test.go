package hyutil_test

import (
	"hyutil"
	"testing"

	"github.com/cheekybits/is"
)

func TestDatetimeParse(t *testing.T) {

	is := is.New(t)

	d1 := hyutil.DatetimeParse("2018-12-09T15:00:00Z")

	d2 := hyutil.DatetimeParse("2018-12-10T00:00:00+09:00")

	s1 := d1.Format("2006-01-02 15:04")
	s2 := d2.Format("2006-01-02 15:04")

	is.Equal(s1, s2)

}
