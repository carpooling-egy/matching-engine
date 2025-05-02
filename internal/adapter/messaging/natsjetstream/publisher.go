package natsjetstream

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
	re "matching-engine/internal/adapter/messaging"
	"matching-engine/internal/adapter/messaging/natsjetstream/mappers"
	"matching-engine/internal/model"
)

// NATSPublisher implements the Publisher interface using NATS JetStream
type NATSPublisher struct {
	nc     *nats.Conn
	js     jetstream.JetStream
	mapper mappers.Mapper
	config Config
}

// NewNATSPublisher creates a new publisher that uses NATS JetStream with default configuration
func NewNATSPublisher() (re.Publisher, error) {
	return NewNATSPublisherWithConfig(LoadConfig())
}

// NewNATSPublisherWithConfig creates a new publisher that uses NATS JetStream with the provided configuration
func NewNATSPublisherWithConfig(config Config) (re.Publisher, error) {
	// Connection options
	opts := []nats.Option{
		nats.Name(config.ConnectionName),
		//nats.RetryOnFailedConnect(true),
		//nats.MaxReconnects(config.MaxReconnects),
		//nats.ReconnectWait(config.ReconnectWait),
		nats.Timeout(config.ConnectTimeout),
		nats.UserInfo("publisher", "publisherpass"),

		// Connection event handlers for logging
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Error().Err(err).Msg("NATS connection disconnected")
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Info().Str("url", nc.ConnectedUrl()).Msg("NATS reconnected")
		}),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Error().Err(err).Msg("NATS error")
		}),
	}

	log.Info().Str("url", config.URL).Msg("Connecting to NATS")

	nc, err := nats.Connect(config.URL, opts...)
	if err != nil {
		log.Error().Err(err).Str("url", config.URL).Msg("Failed to connect to NATS")
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	log.Info().Str("url", nc.ConnectedUrl()).Msg("Connected to NATS")

	js, err := jetstream.New(nc)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create JetStream context")
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	log.Info().Str("url", nc.ConnectedUrl()).Msg("JetStream context created")

	publisher := &NATSPublisher{
		nc:     nc,
		js:     js,
		mapper: mappers.NewJsonMapper(),
		config: config,
	}

	return publisher, nil
}

// PublishMatchingResults publishes matching results to NATS JetStream
func (p *NATSPublisher) PublishMatchingResults(results []*model.MatchingResult) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.config.PublishTimeout)
	defer cancel()

	logCtx := log.With().Str("subject", p.config.Subject).Logger()

	if len(results) == 0 {
		logCtx.Warn().Msg("No matching results to publish")
		return nil
	}

	for i, result := range results {
		logCtx.Debug().
			Str("offerId", result.OfferID()).
			Int("index", i).
			Int("total", len(results)).
			Msg("Publishing matching result")

		data, err := p.mapper.Marshal(result)
		if err != nil {
			logCtx.Error().
				Err(err).
				Str("offerId", result.OfferID()).
				Int("index", i).
				Msg("Failed to marshal result")
			return fmt.Errorf("failed to marshal result: %w", err)
		}

		_, err = p.js.Publish(ctx, p.config.Subject, data)
		if err != nil {
			logCtx.Error().
				Err(err).
				Str("offerId", result.OfferID()).
				Int("index", i).
				Int("published", i).
				Int("remaining", len(results)-i).
				Msg("Failed to publish result; aborting batch")
			return fmt.Errorf("failed to publish result: %w", err)
		}

		logCtx.Debug().
			Str("offerId", result.OfferID()).
			Int("dataSize", len(data)).
			Msg("Successfully published result")
	}

	logCtx.Info().Int("count", len(results)).Msg("Successfully published all results")
	return nil
}

// Close releases resources used by the publisher
func (p *NATSPublisher) Close() error {
	if p.nc != nil {
		log.Info().Msg("Draining NATS connection")
		err := p.nc.Drain()
		if err != nil {
			log.Error().Err(err).Msg("Failed to drain NATS connection")
			return fmt.Errorf("failed to drain NATS connection: %w", err)
		}
		log.Info().Msg("NATS connection drained successfully")
	}
	log.Info().Msg("NATS publisher closed")
	return nil
}
