// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package adkrest

import (
	"strings"
	"testing"
	"time"
)

func TestNewServerRejectsInvalidSSETimeouts(t *testing.T) {
	tests := []struct {
		name string
		cfg  ServerConfig
		want string
	}{
		{
			name: "negative write timeout",
			cfg:  ServerConfig{SSEWriteTimeout: -time.Second},
			want: "SSEWriteTimeout must be non-negative",
		},
		{
			name: "negative heartbeat",
			cfg:  ServerConfig{SSEHeartbeatInterval: -time.Second},
			want: "SSEHeartbeatInterval must be non-negative",
		},
		{
			name: "heartbeat not below write timeout",
			cfg: ServerConfig{
				SSEWriteTimeout:      time.Second,
				SSEHeartbeatInterval: time.Second,
			},
			want: "SSEHeartbeatInterval must be less than SSEWriteTimeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewServer(tt.cfg)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("NewServer() error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestNewServerAcceptsZeroValueSSETimeoutDefaults(t *testing.T) {
	if _, err := NewServer(ServerConfig{}); err != nil {
		t.Fatalf("NewServer() with zero-value SSE timeouts error = %v", err)
	}
}
