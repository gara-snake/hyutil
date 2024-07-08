package hyunet_test

import (
	"log"
	"testing"

	"github.com/gara-snake/hyutil/hyunet"
)

func TestIPV4(t *testing.T) {

	ip := hyunet.IPV4("192.168.")
	log.Println(ip)

}
