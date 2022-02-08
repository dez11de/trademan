package main

import (
	"fmt"

	"github.com/BoRuDar/configuration/v3"
	"github.com/dez11de/cryptodb"
	"github.com/dez11de/exchange"
)

// TODO: provide well documented .example configuration file

// TODO: set defaults and descriptions
type RESTServerConfig struct {
	Host string `default:""`
	Port string `default:"8888"`
}

type trademanConfig struct {
	// TODO: allow setting name of configfile as flag
	Database   cryptodb.DatabaseConfig
	RESTServer RESTServerConfig
	Exchange   exchange.ExchangeConfig
}

func readConfig(cfg *trademanConfig) error {
	fileProvider, err := configuration.NewFileProvider("trademan.yml")
	if err != nil {
		return err
	}

	configurator, err := configuration.New(
		cfg,
		configuration.NewEnvProvider(),
		fileProvider,
		configuration.NewFlagProvider(cfg),
		configuration.NewDefaultProvider(),
	)
	if err != nil {
		return err
	}

	configurator.EnableLogging(true)
	configurator.InitValues()

    fmt.Printf("Configuration:\n%+v", cfg)
	return nil
}
