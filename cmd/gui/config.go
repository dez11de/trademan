package main

import (
	"github.com/BoRuDar/configuration/v3"
)

// TODO: provide well documented .example configuration file

// TODO: set defaults and descriptions
type RESTServerConfig struct {
	Host string `default:""`
    Port string `default:"8888"`
}

type trademanConfig struct {
	// TODO: allow setting name of configfile as flag
	RESTServer RESTServerConfig
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

	configurator.InitValues()

	return nil
}
