// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

//go:build e2e

package e2e

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/elastic/inputrunner/input/assets/k8s"
	"github.com/elastic/inputrunner/input/testutil"
	stateless "github.com/elastic/inputrunner/input/v2/input-stateless"

	"github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/logp"
	v2 "github.com/elastic/inputrunner/input/v2"
	"github.com/stretchr/testify/assert"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestAssetsK8s_Run_startsAndStopsTheInput(t *testing.T) {
	publisher := testutil.NewInMemoryPublisher()

	ctx, cancel := context.WithCancel(context.Background())
	inputCtx := v2.Context{
		Logger:      logp.NewLogger("test"),
		Cancelation: ctx,
	}
	input, err := k8s.Plugin().Manager.(stateless.InputManager).Configure(config.NewConfig())
	assert.NoError(t, err)
	client := k8sfake.NewSimpleClientset()
	if err := k8s.SetClient(client, input); err != nil {
		t.Fatalf("Test failed: %s", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = input.Run(inputCtx, publisher)
		assert.NoError(t, err)
	}()

	time.Sleep(time.Millisecond)
	cancel()

	timeout := time.After(10 * time.Second)
	closeCh := make(chan struct{})
	go func() {
		defer close(closeCh)
		wg.Wait()
	}()

	select {
	case <-timeout:
		t.Fatal("Test timed out")
	case <-closeCh:
		// Waitgroup finished in time, nothing to do
	}
}