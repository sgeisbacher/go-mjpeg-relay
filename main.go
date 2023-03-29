package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/mattn/go-mjpeg"
)

type BroadcasterSource interface {
	Init() error
	GetName() string
	ReadFrame() ([]byte, error)
	Pause()
	Unpause()
}

type BroadcasterMode int

const (
	NORMAL BroadcasterMode = iota
	FALLBACK
)

type Broadcaster struct {
	mode           BroadcasterMode
	source         BroadcasterSource
	fallbackSource BroadcasterSource
}

func (broadcaster *Broadcaster) Broadcast(wg *sync.WaitGroup, stream *mjpeg.Stream) {
	defer wg.Done()

	for {
		source := broadcaster.source
		if broadcaster.mode == FALLBACK {
			source = broadcaster.fallbackSource
		}
		b, err := source.ReadFrame()
		if err != nil {
			fmt.Printf("error while reading frame: %v\n", err)
			fmt.Printf("switching to fallback-source %s ...\n", broadcaster.fallbackSource.GetName())
			broadcaster.SwitchToFallback()
			continue
		}
		stream.Update(b)
	}
}

func (broadcaster *Broadcaster) SwitchSource(newSource BroadcasterSource) {
	if broadcaster.source != nil {
		broadcaster.source.Pause()
	}
	newSource.Unpause()
	broadcaster.source = newSource // TODO mutex?
}

func (broadcaster *Broadcaster) SwitchToFallback() {
	broadcaster.mode = FALLBACK
}

func createMJpegUrlSource(url string) BroadcasterSource {
	decoder, err := mjpeg.NewDecoderFromURL(url)
	if err != nil {
		log.Fatal(err)
	}

	source := &MJpegUrlBroadcasterSource{decoder: decoder}
	err = source.Init()
	if err != nil {
		log.Fatal(err)
	}

	return source
}

func createTextSource(text string) BroadcasterSource {
	source := &TextBroadcasterSource{text}
	err := source.Init()
	if err != nil {
		log.Fatal(err)
	}
	return source
}

func withKeepAlive(source BroadcasterSource) BroadcasterSource {
	keepAliveSource := &KeepAliveBroadcasterSource{isSourcePaused: true, source: &source}
	err := keepAliveSource.Init()
	if err != nil {
		fmt.Printf("error while keepAlive-init for %q: %v\n", source.GetName(), err)
	}
	return keepAliveSource
}

func main() {
	stream := mjpeg.NewStreamWithInterval(50 * time.Millisecond)
	broadcaster := Broadcaster{}

	mjpegUrlSource := withKeepAlive(createMJpegUrlSource("http://admin:admin@192.168.0.178:8081/"))
	noSignalSource := createTextSource("NO SIGNAL")

	broadcaster.SwitchSource(noSignalSource)

	var wg sync.WaitGroup

	wg.Add(1)
	go broadcaster.Broadcast(&wg, stream)

	http.HandleFunc("/mjpeg", stream.ServeHTTP)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<img src="/mjpeg" />`))
	})

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		broadcaster.SwitchSource(mjpegUrlSource)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(`OK`))
	})
	http.HandleFunc("/pause", func(w http.ResponseWriter, r *http.Request) {
		broadcaster.SwitchSource(noSignalSource)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(`OK`))
	})

	server := &http.Server{Addr: ":8080"}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	go func() {
		<-sc
		fmt.Println("received signal, shutting down ...")
		server.Shutdown(context.Background())
		wg.Done()
	}()
	server.ListenAndServe()
	stream.Close()

	wg.Wait()
}
