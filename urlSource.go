package main

import (
	"sync"

	"github.com/mattn/go-mjpeg"
)

type MJpegUrlBroadcasterSource struct {
	m       sync.Mutex
	decoder *mjpeg.Decoder
}

func (source *MJpegUrlBroadcasterSource) Init() error {
	return nil
}

func (source *MJpegUrlBroadcasterSource) ReadFrame() ([]byte, error) {
	source.m.Lock()
	defer source.m.Unlock()

	b, err := source.decoder.DecodeRaw()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (source MJpegUrlBroadcasterSource) GetName() string {
	return "urlsource"
}

func (source MJpegUrlBroadcasterSource) Pause() {}

func (source MJpegUrlBroadcasterSource) Unpause() {}
