package app

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go-auth/utils/config"
	zaplog "go-auth/utils/zap"
	"go.uber.org/zap"
)

func parse() (*zap.Logger, *config.Config, func() error) {
	var (
		logPath    = flag.String("log", "log.txt", "log file path (default 'log.txt')")
		configPath = flag.String("config", "config.toml", "config path (default 'config.toml')")

		env = flag.Bool("env", false, "load from environment variables")
	)
	flag.Parse()

	// // //

	file, err := os.OpenFile(*logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalln(fmt.Errorf("create/open log file (%s): %w", *logPath, err))
	}

	info, err := file.Stat()
	if err != nil {
		log.Fatalln(fmt.Errorf("get file stats: %w", err))
	}

	if info.Size() > 0 {
		_, err = file.Write([]byte("\n\n"))
		if err != nil {
			log.Fatalln(fmt.Errorf("write to file: %w", err))
		}
	}

	logger := zaplog.NewLogger(file)

	conf, err := config.ReadConfig(*configPath)
	if err != nil {
		logger.Fatal("read config file", zap.String("config", *configPath), zap.Error(err))
	}

	if *env {
		err = conf.LoadEnv()
		if err != nil {
			logger.Fatal("load environments", zap.Error(err))
		}
	}

	err = conf.Verify()
	if err != nil {
		logger.Fatal("verify config", zap.Error(err))
	}

	return logger, conf, file.Close
}
