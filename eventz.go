package eventz

import (
	"fmt"
	"sync"
)

var (
	events chan Event
	subs   map[string]*destSubs
	rw     sync.RWMutex
)

type destSubs struct {
	subs map[int]subInfo
	rw   sync.RWMutex
}

type subInfo struct {
	id int
	ch chan<- Event
}

func init() {
	events = make(chan Event, 128)
	subs = make(map[string]*destSubs)

	go mainLoop()
}

// Event is a struct used to distribute events. Has a destination, a payload, and
// an ID.
type Event struct {
	ID          int
	Destination string
	Payload     interface{}
}

// Subscription is returned when a consumer subscribes to a destination. The Events
// channel is used to propagate the Events themselves. The Unsubscribe function should
// be used when the consumer no longer wants to receive events.
//
// Calling Unsubscribe closes the Events channel and removes the subscription from all
// of the registered destinations
type Subscription struct {
	ids         map[string]int
	Events      <-chan Event
	Unsubscribe func()
}

func mainLoop() {
	for e := range events {
		fmt.Println(e)
	}
}

// Subscribe subscribes the consumer to the specified endpoints. Returns an error
// if no endpoints have been specified
func Subscribe(endpoints ...string) (Subscription, error) {
	if len(endpoints) == 0 {
		return Subscription{}, fmt.Errorf("no endpoint specified")
	}

	// Check for empty strings
	for ix, endpoint := range endpoints {
		if endpoint == "" {
			return Subscription{}, fmt.Errorf("endpoint at index %d is empty string", ix)
		}
	}

	eventQ := make(chan Event)

	ids := make(map[string]int)
	for _, endpoint := range endpoints {
		subID := sub(endpoint, eventQ)
		ids[endpoint] = subID
	}

	unsubscribe := func() {
		for endpoint, subID := range ids {
			unsub(endpoint, subID)
		}
		close(eventQ)
	}

	return Subscription{
		ids:         ids,
		Events:      eventQ,
		Unsubscribe: unsubscribe,
	}, nil
}

func sub(endpoint string, eventQ chan Event) int {
	rw.RLock()
	endpointSub, ok := subs[endpoint]
	rw.RUnlock()
	if !ok {
		endpointSub = &destSubs{
			subs: make(map[int]subInfo),
		}
	}

	rw.Lock()
	defer rw.Unlock()

	subID := ID()

	sub := subInfo{
		id: subID,
		ch: eventQ,
	}

	endpointSub.subs[subID] = sub

	subs[endpoint] = endpointSub

	return subID
}

func unsub(endpoint string, subID int) {
	// TODO implement
}

/*
ch := eventz.Subscribe("endpoint1", "endpoint2", ...)

eventz.SubscribeFunc("endpoint", func(){...})



func Subscribe(endpoints ...string) (Subscription)
func SubscribeFunc(f func(endpoint string, payload interface{}))
func Publish(endpoint string, payload interface{})

*/
