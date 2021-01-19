package main

import (
	"sync"
	"testing"
	"time"
)

func TestParseCfg(t *testing.T) {
	parseCfg()
}

func TestSaveCfg(t *testing.T) {
	resumeForUpdates = append(resumeForUpdates, "2")
	saveCfg()
	resumeForUpdates = []string{"1"}
	saveCfg()
}

func TestGetUsrHomeDir(t *testing.T) {
	homeDir := getUsrHomeDir()
	if homeDir != `C:\Users\dimon` {
		t.Error("wrong homedir")
	}
}

func TestPrepareChrome1(t *testing.T) {
	ctx, cancel := prepareChrome(true)
	defer cancel()
	if ctx == nil || cancel == nil {
		t.Error("error ctx or cancel func")
	}
}

func TestPrepareChrome2(t *testing.T) {
	ctx, cancel := prepareChrome(false)
	if ctx == nil || cancel == nil {
		t.Error("error ctx or cancel func")
	}
	updateResume(ctx, "1")
	cancel()
	ctx, cancel = prepareChrome(false)
	getResumeList(ctx, cancel)
	cancel()
	getResumeList(nil, cancel)
	cancel()
	ctx, cancel = prepareChrome(false)
	firstRunChrome(ctx, cancel)
	cancel()
}

func TestGoUpdateMonitor(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go goUpdateMonitor(false)
	time.Sleep(2 * time.Second)
	wg.Done()
}