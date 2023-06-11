/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/inarithefox/partsy/server/api"
	"github.com/inarithefox/partsy/server/app"
	"github.com/inarithefox/partsy/server/config"
	"github.com/inarithefox/partsy/server/public/logger"
	"github.com/inarithefox/partsy/server/web"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the Partsy server",
	RunE:  serveCmdF,
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.RunE = serveCmdF
}

func serveCmdF(command *cobra.Command, args []string) error {
	interrupt := make(chan os.Signal, 1)

	configStore, err := config.NewStoreFromDSN(getConfigDSN(command, config.GetEnvironment()), false, nil, true)
	if err != nil {
		return errors.Wrap(err, "failed to load configuration.")
	}

	defer configStore.Close()

	return runServer(configStore, interrupt)
}

func runServer(configStore *config.Store, interrupt chan os.Signal) error {
	options := []app.Option{
		app.ConfigStore(configStore),
	}

	server, err := app.NewServer(options...)
	if err != nil {
		logger.Error(err, "unable to create server")
		return err
	}

	defer server.Shutdown()

	_, err = api.Init(server)
	if err != nil {
		logger.Error(err, "unable to initialize API")
		return err
	}

	web.New(server)

	err = server.Start()
	if err != nil {
		logger.Error(err, "unable to start server")
		return err
	}

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	return nil
}
