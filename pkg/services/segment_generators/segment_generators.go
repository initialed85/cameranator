package segment_generators

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/initialed85/cameranator/pkg/liveness"

	"github.com/initialed85/glue/pkg/network"

	"github.com/initialed85/cameranator/pkg/segments/segment_generator"
)

type SegmentGenerators struct {
	mu                     sync.Mutex
	feeds                  []segment_generator.Feed
	host                   string
	port                   int64
	sender                 *network.Sender
	segmentGeneratorByFeed map[segment_generator.Feed]*segment_generator.SegmentGenerator
	livenessAgent          *liveness.Agent
}

func NewSegmentGenerators(feeds []segment_generator.Feed, host string, port int64) *SegmentGenerators {
	s := SegmentGenerators{
		feeds:                  feeds,
		host:                   host,
		port:                   port,
		segmentGeneratorByFeed: make(map[segment_generator.Feed]*segment_generator.SegmentGenerator),
	}

	return &s
}

func (s *SegmentGenerators) completeFn(event segment_generator.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, err := json.Marshal(event)
	if err != nil {
		log.Printf("err: failed to marshal %#+v because %v", event, err)
	}

	err = s.sender.Send(b)
	if err != nil {
		log.Printf("err: failed to send %#+v to %v:%v because %v", event, s.host, s.port, err)
	}
}

func (s *SegmentGenerators) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sender = network.NewSender(fmt.Sprintf("%v:%v", s.host, s.port))
	err := s.sender.Open()
	if err != nil {
		return err
	}

	s.livenessAgent, err = liveness.Open(
		[]liveness.HasLiveness{s},
		8080, // TODO
	)
	if err != nil {
		return err
	}

	for _, feed := range s.feeds {
		s.segmentGeneratorByFeed[feed] = segment_generator.NewSegmentGenerator(
			feed,
			s.completeFn,
		)
		err = s.segmentGeneratorByFeed[feed].Start()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SegmentGenerators) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, feed := range s.feeds {
		s.segmentGeneratorByFeed[feed].Stop()
	}

	s.sender.Close()
	s.livenessAgent.Close()
}

func (s *SegmentGenerators) IsLive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, segmentGenerator := range s.segmentGeneratorByFeed {
		if !segmentGenerator.IsLive() {
			return false
		}
	}

	return true
}
