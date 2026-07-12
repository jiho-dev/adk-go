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

package helper

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
)

type deadlineRefreshingRecorder struct {
	*httptest.ResponseRecorder
	refreshErr   error
	refreshCalls int
}

func (r *deadlineRefreshingRecorder) RefreshWriteDeadline() error {
	r.refreshCalls++
	return r.refreshErr
}

func TestEmitJSONRefreshesWriteDeadline(t *testing.T) {
	rw := &deadlineRefreshingRecorder{ResponseRecorder: httptest.NewRecorder()}

	if err := EmitJSON(rw, map[string]string{"event": "ok"}); err != nil {
		t.Fatalf("EmitJSON() failed: %v", err)
	}

	if got := rw.refreshCalls; got != 1 {
		t.Errorf("RefreshWriteDeadline calls = %d, want 1", got)
	}
	if got := rw.Body.String(); !strings.Contains(got, `"event":"ok"`) {
		t.Errorf("EmitJSON() body = %q, want event payload", got)
	}
}

func TestEmitJSONWriteDeadlineRefreshError(t *testing.T) {
	wantErr := errors.New("deadline failed")
	rw := &deadlineRefreshingRecorder{
		ResponseRecorder: httptest.NewRecorder(),
		refreshErr:       wantErr,
	}

	err := EmitJSON(rw, map[string]string{"event": "not written"})
	if !errors.Is(err, wantErr) {
		t.Fatalf("EmitJSON() error = %v, want %v", err, wantErr)
	}
	if got := rw.Body.Len(); got != 0 {
		t.Errorf("EmitJSON() wrote %d bytes after deadline refresh failed, want 0", got)
	}
}
