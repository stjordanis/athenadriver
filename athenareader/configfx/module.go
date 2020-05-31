package configfx

import (
	secret "github.com/uber/athenadriver/examples/constants"
	drv "github.com/uber/athenadriver/go"
	"go.uber.org/config"
	"go.uber.org/fx"
)

var Module = fx.Provide(new)

// Params defines the dependencies or inputs
type Params struct {
	fx.In
}

type ReaderOutputConfig struct {
	Render    string `yaml:"render"`
	Page      int    `yaml:"pagesize"`
	Style     string `yaml:"style"`
	Rowonly   bool   `yaml:"rowonly"`
	Moneywise bool   `yaml:"moneywise"`
}

type ReaderInputConfig struct {
	Bucket   string `yaml:"bucket"`
	Region   string `yaml:"region"`
	Database string `yaml:"database"`
	Admin    bool   `yaml:"admin"`
}

// Result defines output
type Result struct {
	fx.Out

	OConfig   ReaderOutputConfig
	DrvConfig *drv.Config
}

func new(p Params) (Result, error) {
	var readerOutputConfig ReaderOutputConfig
	var readerInputConfig ReaderInputConfig
	provider, err := config.NewYAML(config.File("athenareader/config.yaml"))
	if err != nil {
		return Result{}, err
	}
	provider.Get("athenareader.output").Populate(&readerOutputConfig)
	provider.Get("athenareader.input").Populate(&readerInputConfig)

	conf, err := drv.NewDefaultConfig(readerInputConfig.Bucket, readerInputConfig.Region, secret.AccessID, secret.SecretAccessKey)
	if err != nil {
		return Result{}, err
	}
	if readerOutputConfig.Moneywise {
		conf.SetMoneyWise(true)
	}
	conf.SetDB(readerInputConfig.Database)
	if !readerInputConfig.Admin {
		conf.SetReadOnly(true)
	}
	if err != nil {
		return Result{}, err
	}

	return Result{
		OConfig:   readerOutputConfig,
		DrvConfig: conf,
	}, nil
}
