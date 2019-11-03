package hyuio

import (
	"log"
	"os"
	"os/signal"
)

// WaitQuit Ctrl + C での終了を待機する
func WaitQuit() {

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		select {
		case <-quit:
			log.Println("Good luck !")
			signal.Stop(quit)
			os.Exit(0)
		}
	}()

}
