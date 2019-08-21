package hyunet_test

import (
	"hyutil/hyunet"
	"log"
	"testing"
)

func TestIPV4(t *testing.T) {

	ip := hyunet.IPV4("192.168.")
	log.Println(ip)

}
