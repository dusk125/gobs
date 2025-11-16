package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/dusk125/gobs"
	"golang.org/x/image/colornames"
)

func main() {
	defer func() {
		println("bye")
	}()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	gobs.Startup("en-US")
	defer func() {
		println("shutdown")
		gobs.Shutdown()
	}()

	if gobs.Initialized() {
		fmt.Println("Initialized!")
	} else {
		panic("Not initialized")
	}

	fmt.Println("Version", gobs.VersionString())

	vinfo := gobs.VideoInfo{
		GraphicsModule: gobs.GraphicsModuleOpenGL,
		FPSNum:         60,
		FPSDen:         1,
		BaseWidth:      1920,
		BaseHeight:     1080,
		OutputWidth:    1920,
		OutputHeight:   1080,
		ScaleType:      gobs.VIDEO_SCALE_BICUBIC,
		OutputFormat:   gobs.VIDEO_FORMAT_BGRA,
	}
	if err := gobs.ResetVideo(&vinfo); err != nil {
		panic(err)
	}

	ainfo := gobs.AudioInfo{
		SamplesPerSec: 44100,
		Speakers:      gobs.SPEAKERS_STEREO,
	}
	if err := gobs.ResetAudio(&ainfo); err != nil {
		panic(err)
	}

	toLoad := map[string]struct{}{
		"obs-ffmpeg":      {},
		"obs-browser":     {},
		"rtmp-services":   {},
		"obs-x264":        {},
		"obs-outputs":     {},
		"obs-transitions": {},
		"image-source":    {},
		"vlc-video":       {},
	}

	for info := range gobs.FindModules() {
		name := info.Name()
		if _, has := toLoad[name]; has {
			mod, err := gobs.OpenModule(info.BinPath(), info.DataPath())
			if err != nil {
				panic(fmt.Errorf("failed to open module: %v", name))
			}

			if err := mod.Init(); err != nil {
				panic(err)
			}
		}
	}

	gobs.PostLoadModules()
	gobs.LogLoadedModules()

	for id := range gobs.InputTypes() {
		fmt.Println(id)
		fmt.Println("\t", gobs.SourceDefaults(id))
		fmt.Println()
	}

	col := colornames.Cornflowerblue
	color := gobs.SourceCreate("color_source_v3", "hi", gobs.Data{
		"color": binary.NativeEndian.Uint32([]byte{col.R, col.G, col.B, col.A}),
	}, nil)
	defer color.Release()

	vlcData := gobs.SourceDefaults("vlc_source")
	vlcData["loop"] = false
	vlcData["playback_behavior"] = "always_play"
	vlcData["playlist"] = []map[string]any{
		{
			"value": "rtmp://localhost:1935/host/me",
		},
	}

	vlc := gobs.SourceCreate("vlc_source", "vlc", vlcData, nil)
	defer vlc.Release()

	browserData := gobs.SourceDefaults("browser_source")
	browserData["url"] = "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	browserData["width"] = vinfo.BaseWidth
	browserData["height"] = vinfo.BaseHeight
	browserData["fps"] = 60
	browserData["reroute_audio"] = true
	// browserData["restart_when_active"] = true

	browser := gobs.SourceCreate("browser_source", "browser", browserData, nil)
	defer browser.Release()

	scene := gobs.SceneCreate("scene")
	defer scene.Release()

	defer scene.Add(color).Remove()
	{
		vlc := scene.Add(vlc)
		defer vlc.Remove()

		vlc.SetBoundsType(gobs.BOUNDS_SCALE_INNER)
		vlc.SetBounds(gobs.Vec2{X: float32(vinfo.BaseWidth), Y: float32(vinfo.BaseHeight)})
	}
	{
		browser := scene.Add(browser)
		defer browser.Remove()

		browser.SetBoundsType(gobs.BOUNDS_SCALE_INNER)
		browser.SetBounds(gobs.Vec2{X: float32(vinfo.BaseWidth), Y: float32(vinfo.BaseHeight)})
	}

	gobs.SetOuputSource(0, scene.Source())

	videoData := gobs.Data{
		"width":   vinfo.BaseWidth,
		"height":  vinfo.BaseHeight,
		"fps_num": 60,
		"fps_den": 1,
		"bitrate": 6000,
		"keyint":  60,
		"preset":  "veryfast",
		"profile": "main",
		"tune":    "zerolatency",
	}
	video := gobs.VideoEncoderCreate("obs_x264", "video-enc", videoData, nil)
	defer video.Release()

	audioData := gobs.Data{
		"bitrate": 128,
	}
	audio := gobs.AudioEncoderCreate("ffmpeg_aac", "audio-enc", 0, audioData, nil)
	defer audio.Release()

	serviceData := gobs.Data{
		"server": "rtmp://localhost:1935/host/",
		"key":    "test",
	}
	serv := gobs.ServiceCreate("rtmp_custom", "service", serviceData, nil)
	defer serv.Release()

	out := gobs.OutputCreate("rtmp_output", "out", gobs.Data{}, nil)
	defer out.Release()

	video.SetVideo(gobs.GetVideo())
	audio.SetAudio(gobs.GetAudio())
	out.SetVideoEncoder(video)
	out.SetAudioEncoder(audio, 0)
	out.SetService(serv)
	out.SetPreferredSize(vinfo.BaseWidth, vinfo.BaseHeight)

	if !serv.CanTryToConnect() {
		panic("failed to connect")
	}

	running := make(chan struct{})
	go func() {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		for event := range out.Events(ctx) {
			switch event := event.(type) {
			case gobs.SignalOutputStart:
				fmt.Println("output started")
			case gobs.SignalOutputStop:
				fmt.Println("output stopped")
				code := event.Code

				if code != 0 {
					slog.Error("output failed", "error", out.Error(), "code", code)
				}
				close(running)
				println("closed")
			case gobs.SignalOutputPause:
				fmt.Println("output paused")
			case gobs.SignalOutputUnpause:
				fmt.Println("output unpaused")
			case gobs.SignalOutputStarting:
				fmt.Println("output starting")
			case gobs.SignalOutputStopping:
				fmt.Println("output stopping")
			case gobs.SignalOutputActivate:
				fmt.Println("output reconnect")
			case gobs.SignalOutputDeactivate:
				fmt.Println("output reconnect success")
			case gobs.SignalOutputReconnect:
				fmt.Println("output activate")
			case gobs.SignalOutputReconnectSuccess:
				fmt.Println("output deactivate")
			}
		}
	}()

	if err := out.Start(); err != nil {
		panic(fmt.Errorf("failed to start output: %v", err))
	}
	defer func() {
		println("stopping output")
		out.Stop()

		to, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		select {
		case <-to.Done():
			println("force stopping")
			out.ForceStop()
		case <-running:
			println("got stop signal")
			return
		}
	}()

	// ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ctx.Done():
			println("breaking")
			return
			// case <-ticker.C:
			// 	si := slices.Collect(scene.Items())
			// 	si[0], si[1] = si[1], si[0]
			// 	if !scene.ReorderItems(si) {
			// 		slog.Error("failed to reorder items")
			// 	}
		}
	}
}
