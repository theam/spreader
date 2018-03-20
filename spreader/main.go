package main

import (
	"log"
	"github.com/siddontang/go-mysql/canal"
	"github.com/alvaroloes/kinesis/rtc"
)

func main() {
	config := canal.NewDefaultConfig()
	config.Dump.Databases = []string{"spreader"}

	binlogScanner, err := canal.NewCanal(config)
	if err != nil {
		log.Fatal(err)
	}

	binlogScanner.SetEventHandler(&RowEventProcessor{bus: rtc.NewKinesis()})
	binlogScanner.Run()
}
