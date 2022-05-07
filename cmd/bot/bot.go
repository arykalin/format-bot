package bot

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/arykalin/format-bot/internal/telegram"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type bot struct {
}

type Bot interface {
	Start()
}

type Config struct {
	TeleToken  string `yaml:"telegram_token"`
	TeleChatID int64  `yaml:"telegram_chat_id"`
	SheetID    string `yaml:"sheet_id"`
}

func (r *bot) Start() {
	pathConfig := pflag.StringP("path", "c", "./config.yml", "path to config file")
	help := pflag.BoolP("help", "h", false, "show help")
	pflag.Parse()

	configFile, err := ioutil.ReadFile(*pathConfig)
	if err != nil {
		log.Fatalf("can't read file")
	}

	if *help {
		pflag.PrintDefaults()
		os.Exit(0)
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("can't unmarshal config: %s", err)
	}

	sLoggerConfig := zap.NewDevelopmentConfig()
	sLoggerConfig.DisableStacktrace = true
	sLoggerConfig.DisableCaller = true
	sLogger, err := sLoggerConfig.Build()
	if err != nil {
		panic(err)
	}
	logger := sLogger.Sugar()
	newTeleBot := telegram.NewBot(
		config.TeleChatID,
		config.TeleToken,
		config.SheetID,
		logger,
	)
	err = newTeleBot.Start()
	if err != nil {
		logger.Fatalw("starting bot failed", "reason", err)
	}
}

func NewBot() Bot {
	return &bot{}
}
