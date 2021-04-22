package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/justinas/alice"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trelore/todoapi/internal"
	"github.com/trelore/todoapi/internal/datastores/mem"
	"github.com/trelore/todoapi/internal/datastores/postgres"
	"github.com/trelore/todoapi/internal/datastores/redis"
	"github.com/trelore/todoapi/internal/middlewares"
	"go.uber.org/zap"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todoapi",
	Short: "A small todo API",
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(); err != nil {
			log.Fatal(err)
		}
	},
}

func run() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("new logger: %w", err)
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	// Create channel for shutdown signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	// go routine serving the swagger docs
	go func() {
		docsPort := os.Getenv("DOCS_PORT")
		if docsPort != "" {
			docsPort = ":8083"
		}
		sugar.Infof("serving docs on port: %s", docsPort)
		fs := http.FileServer(http.Dir("./docs"))
		http.Handle("/swaggerui/", http.StripPrefix("/swaggerui/", fs))
		err := http.ListenAndServe(docsPort, nil)
		if err != nil {
			sugar.Fatal(err)
		}
	}()

	// choose the datastore
	var db internal.Datastore
	switch strings.ToLower(os.Getenv("DATASTORE")) {
	case "redis":
		db = redis.New(sugar)
	case "postgres":
		var close func() error
		var err error
		db, close, err = postgres.New()
		if err != nil {
			return fmt.Errorf("new postgres: %w", err)
		}
		defer close()
	default:
		sugar.Warn("using in memory datastore")
		db = mem.New()
	}

	// go routine serving the todo app
	go func() {
		port := os.Getenv("PORT")
		if port != "" {
			port = ":8081"
		}
		s := internal.NewServer(db)
		sugar.Infof("running on address: %s", port)
		http.ListenAndServe(port, alice.New(
			middlewares.Recovery,
			middlewares.Logging(sugar),
		).Then(s))
	}()
	<-stop

	sugar.Infof("closing server")
	return nil
}

func init() {
	rootCmd.Flags().StringP("datastore", "d", viper.GetString("DATASTORE"), "Choose which datastore you want to use (set to `POSTGRES` to use postgres, etc..).")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
