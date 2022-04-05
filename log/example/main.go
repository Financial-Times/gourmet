package main

import (
	"fmt"

	"github.com/gotha/gourmet/log"
)

func main() {

	logger := log.NewStructuredLogger(log.InfoLevel)

	logger.Info("this is my test log",
		log.WithField("transaction_id", "tid_gosho1"),
		log.WithField("name", "gosho"),
	)

	l2 := log.NewStructuredLogger(log.DebugLevel,
		log.WithServiceName("myService"),
		log.WithField("appSystemCode", "systemCode"),
	)

	l2.Debug("my debug message",
		log.WithField("key", "val"),
	)
	l2.Error("my error log",
		log.WithError(fmt.Errorf("my error")),
	)
}
