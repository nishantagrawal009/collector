package config

import (
	"collector/log"
	"flag"
	"fmt"
	"strings"
	"time"
	storageBadger "collector/storage/badger"
)

const (
	defaultAddr        = ":10100"
	defaultExitTimeout = 5 * time.Second

	defaultStorageType = "auto"
	StorageTypeBadger  = "badger"

)

var storageTypes = []string{StorageTypeBadger}

type Config struct {
	Addr        string
	ExitTimeout time.Duration
	Logger      log.Config
	storageType string
	Badger      storageBadger.Config
}

func (conf *Config) RegisterFlags(f *flag.FlagSet) {
	f.StringVar(&conf.Addr, "addr", defaultAddr, "address to listen")
	f.DurationVar(&conf.ExitTimeout, "exit-timeout", defaultExitTimeout, "server shutdown timeout")

	conf.Logger.RegisterFlags(f)

	f.StringVar(&conf.storageType, "storage-type", defaultStorageType, fmt.Sprintf("storage type: %s", strings.Join(storageTypes, ", ")))

	conf.Badger.RegisterFlags(f)

}

func (conf *Config) StorageType() ([]string, error) {
	return []string{StorageTypeBadger}, nil
}

