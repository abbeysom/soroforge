package stellar

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests verify SimulateCount reports the number of simulated
// transactions correctly, covering happy-path and error-path scenarios.
func TestMockClientSimulateCountHappyPath(t *testing.T) {
	client := NewMockClient()

	// Calling Simulate once should report count 1.
	result, err := client.Simulate(context.Background(), "tx1")
	require.NoError(t, err)
	assert.Equal(t, 1, client.SimulateCount())

	// Calling Simulate again should report count 2.
	result, err = client.Simulate(context.Background(), "tx2")
	require.NoError(t, err)
	assert.Equal(t, 2, client.SimulateCount())
}

// table-driven tests for SimulateCount
func TestMockClientSimulateCount(t *testing.T) {
	type testCase struct {
		name    string
		ops     func(*MockClient)
		want    int
		wantErr bool
	}

	tests := []testCase{
		{
			name: "zero simulations",
			ops: func(client *MockClient) {
				// no calls to Simulate
			},
			want: 0,
		},
		{
			name: "one simulation",
			ops: func(client *MockClient) {
				client.Simulate(context.Background(), "tx1")
			},
			want: 1,
		},
		{
			name: "multiple simulations",
			ops: func(client *MockClient) {
				client.Simulate(context.Background(), "tx1")
				client.Simulate(context.Background(), "tx2")
				client.Simulate(context.Background(), "tx3")
			},
			want: 3,
		},
		{
			name: "simulation after other operations",
			ops: func(client *MockClient) {
				// LoadAccount does not affect SimulateCount
				_ = client.LoadAccount(context.Background(), "GBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBX")
				client.Simulate(context.Background(), "tx1")
			},
			want: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := NewMockClient()
			tc.ops(client)
			got := client.SimulateCount()
			if tc.wantErr {
				assert.Error(t, got, "expected error but got count %d", got)
			} else {
				assert.Equal(t, tc.want, got, "SimulateCount %q", tc.name)
			}
		})
	}
}