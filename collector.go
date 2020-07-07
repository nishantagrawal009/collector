package main

import (
	"log"
)


type Collector struct {
	logger *log.Logger
	sw     Writer
}


func NewCollector(logger *log.Logger, sw Writer) *Collector {
	return &Collector{
		logger: logger,
		sw:     sw,
	}
}