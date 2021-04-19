package cmd

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/justinas/alice"
	"github.com/spf13/cobra"
	"github.com/trelore/todoapi/internal"
	"github.com/trelore/todoapi/internal/datastores/mem"
	"github.com/trelore/todoapi/internal/middlewares"
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
	// Create channel for shutdown signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	// go routine serving the swagger docs
	go func() {
		docsPort := ":8083"
		log.Printf("serving docs on port: %s", docsPort)
		fs := http.FileServer(http.Dir("."))
		http.Handle("/swaggerui/", http.StripPrefix("/swaggerui/", fs))
		err := http.ListenAndServe(docsPort, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// go routine serving the todo app
	go func() {
		s := internal.NewServer(mem.New())
		port := ":8081"
		log.Printf("running on address: %s", port)
		http.ListenAndServe(port, alice.New(
			middlewares.Recovery,
			middlewares.Logging,
		).Then(s))
	}()
	<-stop

	log.Println("closing server")
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
