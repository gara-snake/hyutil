package hyutil_test

import (
	"hyutil"
	"testing"

	"github.com/cheekybits/is"
)

func TestCamelToSnake(t *testing.T) {

	is := is.New(t)

	ret := hyutil.CamelToSnake("StringsTestNow!!")
	is.Equal(ret, "strings_test_now!!")

	ret = hyutil.CamelToSnake("SN")
	is.Equal(ret, "sn")

}

func TestSnakeToUcamel(t *testing.T) {

	is := is.New(t)

	ret := hyutil.SnakeToUcamel("strings_test_now!!")
	is.Equal(ret, "StringsTestNow!!")

	ret = hyutil.SnakeToUcamel("s_t_r_i_n_g_s_t_e_s_t_n_o_w!!")
	is.Equal(ret, "STRINGSTESTNOW!!")
}
