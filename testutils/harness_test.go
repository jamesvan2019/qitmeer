// Copyright (c) 2020 The qitmeer developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package testutils_test

import (
	"github.com/Qitmeer/qitmeer/params"
	. "github.com/Qitmeer/qitmeer/testutils"
	"sync"
	"testing"
	"time"
)

func TestHarness(t *testing.T) {
	h, err := NewHarness(t, params.PrivNetParam.Params)
	if err != nil {
		t.Errorf("create new test harness instance failed %v", err)
	}
	if err := h.Setup(); err != nil {
		t.Errorf("setup test harness instance failed %v", err)
	}

	h2, err := NewHarness(t, params.PrivNetParam.Params)
	defer func() {

		if err := h.Teardown(); err != nil {
			t.Errorf("tear down test harness instance failed %v", err)
		}
		numOfHarnessInstances := len(AllHarnesses())
		if numOfHarnessInstances != 10 {
			t.Errorf("harness num is wrong, expect %d , but got %d", 10, numOfHarnessInstances)
			for _, h := range AllHarnesses() {
				t.Errorf("%v\n", h.Id())
			}
		}

		if err := TearDownAll(); err != nil {
			t.Errorf("tear down all error %v", err)
		}
		numOfHarnessInstances = len(AllHarnesses())
		if numOfHarnessInstances != 0 {
			t.Errorf("harness num is wrong, expect %d , but got %d", 0, numOfHarnessInstances)
			for _, h := range AllHarnesses() {
				t.Errorf("%v\n", h.Id())
			}
		}

	}()
	numOfHarnessInstances := len(AllHarnesses())
	if numOfHarnessInstances != 2 {
		t.Errorf("harness num is wrong, expect %d , but got %d", 2, numOfHarnessInstances)
	}
	if err := h2.Teardown(); err != nil {
		t.Errorf("teardown h2 error:%v", err)
	}

	numOfHarnessInstances = len(AllHarnesses())
	if numOfHarnessInstances != 1 {
		t.Errorf("harness num is wrong, expect %d , but got %d", 1, numOfHarnessInstances)
	}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			NewHarness(t, params.PrivNetParam.Params)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestHarnessNodePorts(t *testing.T) {
	var setup, teardown sync.WaitGroup
	ha := make(map[int]*Harness, 10)
	for i := 0; i < 10; i++ {
		setup.Add(1)
		teardown.Add(1)
		// new and setup
		go func(index int) {
			defer setup.Done()
			h, err := NewHarness(t, params.PrivNetParam.Params)
			if err != nil {
				t.Errorf("new harness failed: %v", err)
			}
			ha[index] = h
			if err := ha[index].Setup(); err != nil {
				t.Errorf("setup harness failed: %v", err)
			}
			time.Sleep(500 * time.Millisecond)
		}(i)
		go func(index int) {
			defer teardown.Done()
			setup.Wait() //wait for all setup finished
			if err := ha[index].Teardown(); err != nil {
				t.Errorf("tear down harness failed: %v", err)
			}
		}(i)
	}
	teardown.Wait()
}
