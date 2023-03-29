package main

import (
	"fmt"
	"time"
)

type KeepAliveBroadcasterSource struct {
	isSourcePaused bool
	source         *BroadcasterSource
}

func (source *KeepAliveBroadcasterSource) Init() error {
	fmt.Printf("starting keep-alive for %s source ...\n", source.GetName())
	go func() {
		for {
			time.Sleep(2 * time.Second)
			if !source.isSourcePaused {
				continue
			}
			fmt.Printf("keep-alive heartbeat for %s\n", source.GetName())
			_, err := (*source.source).ReadFrame()
			if err != nil {
				fmt.Printf("seeing error on keep-alive for %s: %v\n", source.GetName(), err)
			}
		}
	}()
	return nil
}

func (source KeepAliveBroadcasterSource) ReadFrame() ([]byte, error) {
	return (*source.source).ReadFrame()
}

func (source KeepAliveBroadcasterSource) GetName() string {
	return fmt.Sprintf("keep-alive:%s", (*source.source).GetName())
}

func (source *KeepAliveBroadcasterSource) Pause() {
	source.isSourcePaused = true
}

func (source *KeepAliveBroadcasterSource) Unpause() {
	source.isSourcePaused = false
}
