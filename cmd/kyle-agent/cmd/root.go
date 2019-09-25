package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mee6aas/kyle/internal/pkg/agent"
)

var (
	optDebug bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kyle-agent",
	Short: "kyle-agent is the agent delegate of the Mee6aaS",

	RunE: func(cmd *cobra.Command, args []string) (e error) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		log.Info("Setting up the agent")
		if e = agent.Setup(agent.Config{
			Host: "localhost",
			Port: 5122,
		}); e != nil {
			e = errors.Wrap(e, "Failed to setup the agent")
			log.Error(e)
			return
		}
		log.Info("Agent setup")

		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt, syscall.SIGINT)

		go func() {
			<-sig
			log.Warn("SIGINT received, stopping the agent")
			cancel()
		}()

		log.Info("Starting the agent")
		if e = agent.Start(ctx); e != nil && e != context.Canceled {
			e = errors.Wrap(e, "Failed to start the agent")
			log.Error(e)
			return
		}
		log.Info("Agent stopped")

		// listen siganl
		// print "Destroying agent"
		// print "Press ctrl+c to force shutdown"
		// agent.destory(ctx)

		return
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&optDebug, "debug", false, "print debug messages")

	rootCmd.PersistentPreRun = func(_ *cobra.Command, _ []string) {
		if optDebug {
			log.SetLevel(log.DebugLevel)
		}
	}
}
