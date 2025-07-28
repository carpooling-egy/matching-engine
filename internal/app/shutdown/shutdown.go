package shutdown

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

// TODO: Add a proper graceful shutdown implementation, for example, we can output some intermediate matched results this is just a placeholder

// Setup configures graceful shutdown handling
func Setup(cancel func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Info().
			Str("signal", sig.String()).
			Msg("Received signal. Initiating graceful shutdown...")

		// Signal to cancel the context
		cancel()

		// If we receive a second signal, force exit
		select {
		case <-sigChan:
			log.Warn().Msg("Received second signal, forcing exit")
			os.Exit(1)
		}
	}()
}
