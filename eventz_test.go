package eventz

import (
	"fmt"
	"testing"
)

func Test_sub(t *testing.T) {
	type args struct {
		endpoint string
		eventQ   chan Event
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Valid", args{"testEndpoint", make(chan Event)}, 0},
		{"Valid2", args{"testEndpoint", make(chan Event)}, 0},
		{"Valid3", args{"testEndpoint", make(chan Event)}, 0},
	}

	ID = func() int {
		return 0
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sub(tt.args.endpoint, tt.args.eventQ); got != tt.want {
				t.Errorf("sub() = %v, want %v", got, tt.want)
				return
			}

			if err := checkSubscriptionExists(tt.args.endpoint, tt.want); err != nil {
				t.Errorf("sub failed: %v", err)
			}
		})
	}
}

func TestSubscribe(t *testing.T) {
	type args struct {
		endpoints []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int
		wantErr bool
	}{
		{"Valid", args{[]string{"someEndpoint"}}, map[string]int{"someEndpoint": 0}, false},
		{"Valid multi args", args{[]string{"someEndpoint", "someEndpointTwo"}}, map[string]int{"someEndpoint": 0, "someEndpointTwo": 0}, false},

		{"Missing endpoints", args{[]string{""}}, nil, true},
	}

	ID = func() int {
		return 0
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Subscribe(tt.args.endpoints...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Subscribe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			for _, endpoint := range tt.args.endpoints {
				id, ok := got.ids[endpoint]
				if !ok {
					t.Errorf("Subscription does not contain endpoint %q", endpoint)
					return
				}

				if id != tt.want[endpoint] {
					t.Errorf("Saved ID mismatch. Got: %d, expected: %d", id, tt.want[endpoint])
					return
				}

				if err := checkSubscriptionExists(endpoint, id); err != nil {
					t.Errorf("subscription does not exist: %v", err)
				}
			}

			if got.Events == nil {
				t.Errorf("Events channel is nil")
				return
			}
		})
	}
}

func checkSubscriptionExists(endpoint string, subID int) error {
	endpointSub, ok := subs[endpoint]
	if !ok {
		return fmt.Errorf("endpointSub not created for endpoint: %s", endpoint)
	}

	if endpointSub == nil {
		return fmt.Errorf("endpointSub is nil")
	}

	subInfo, ok := endpointSub.subs[subID]
	if !ok {
		return fmt.Errorf("subscription not registered at endpoint: %s", endpoint)
	}

	if subInfo.id != subID {
		return fmt.Errorf("subscription id mismatch. Expected: %d, got: %d", subID, subInfo.id)
	}

	if subInfo.ch == nil {
		return fmt.Errorf("subscription channel is nil")
	}

	return nil
}
