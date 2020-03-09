package resolver

import (
	"errors"
	"log"

	"github.com/c-rainbow/simplechatbot/api/helix"
	helix_api "github.com/nicklaw5/helix"
)

/*
Resolver that handles StreamAPIType variables. It is responsible for
(1) Calling Twitch API and returning result or error
(2) TODO: Caching results. Can be useful for $(uptime) or $(uptime_at)
(3) TODO: Aggregate call requests from multiple channels and make one API call together.
	When the bot is heavily loaded, a short wait of 500ms is good enough to reduce total number of API calls
*/

var (
	// ErrChannelNotFound Channel not found
	ErrChannelNotFound = errors.New("Channel is not found")
	defaultResolver    StreamsAPIResolverT
)

// StreamsAPIResolverT interface for Streams API resolver
type StreamsAPIResolverT interface {
	Resolve(channel string) (*helix_api.Stream, error)
}

// StreamsAPIResolver resolver struct
type StreamsAPIResolver struct {
	client helix.HelixClientT
}

// DefaultStreamsAPIResolver default Streams API resolver
func DefaultStreamsAPIResolver() StreamsAPIResolverT {
	if defaultResolver == nil {
		defaultResolver = &StreamsAPIResolver{client: helix.DefaultHelixClient()}
	}
	return defaultResolver
}

// Resolve resolve results
func (resolver *StreamsAPIResolver) Resolve(channel string) (*helix_api.Stream, error) {
	streams, err := resolver.client.GetStreams([]string{}, []string{channel})
	if err != nil {
		log.Println("Could not get stream from Helix Streams API. ", err)
		return nil, err
	}
	if len(streams) == 0 {
		log.Println("Stream could not found for channel ", channel)
		return nil, ErrChannelNotFound
	}
	return &streams[0], nil

}
