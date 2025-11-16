# gobs

Go bindings for [libobs](https://github.com/obsproject/obs-studio) (OBS Studio). Build headless OBS pipelines — video compositing, encoding, and streaming — entirely from Go.

## Features

- Full OBS lifecycle: initialize, configure video/audio, create sources, scenes, encoders, outputs, and services
- Module loading and discovery (`FindModules`, `OpenModule`, selective or bulk loading)
- Scene composition with transform and bounds control
- Signal-based event system with Go channels and context cancellation
- Encoder, output, source, filter, and transition type enumeration via `iter.Seq`
- Settings via `Data` (map[string]any) — matches OBS JSON settings

## Requirements

- OBS Studio 32.x built with `pkg-config` support (`libobs.pc`)
- Go 1.23+ (uses `iter.Seq`)
- CGO enabled
- Linux (Wayland or X11)

Set `PKG_CONFIG_PATH` to your OBS install if it's not in the default search path:

```sh
export PKG_CONFIG_PATH=/opt/obs/lib/pkgconfig
```

## Quick Start

```go
gobs.Startup("en-US")
defer gobs.Shutdown()

gobs.ResetVideo(&gobs.VideoInfo{
    GraphicsModule: gobs.GraphicsModuleOpenGL,
    FPSNum: 60, FPSDen: 1,
    BaseWidth: 1920, BaseHeight: 1080,
    OutputWidth: 1920, OutputHeight: 1080,
    OutputFormat: gobs.VIDEO_FORMAT_BGRA,
})

gobs.ResetAudio(&gobs.AudioInfo{
    SamplesPerSec: 44100,
    Speakers:      gobs.SPEAKERS_STEREO,
})

// Load modules
for info := range gobs.FindModules() {
    mod, _ := gobs.OpenModule(info.BinPath(), info.DataPath())
    mod.Init()
}
gobs.PostLoadModules()

// Create a scene with a color source
scene := gobs.SceneCreate("scene")
defer scene.Release()

color := gobs.SourceCreate("color_source_v3", "bg", gobs.Data{"color": uint32(0xFF6495ED)}, nil)
defer color.Release()

scene.Add(color)
```

## Docker

A Dockerfile is provided that builds OBS + [gamescope](https://github.com/ValveSoftware/gamescope) from source and runs the example app in a headless container:

```sh
docker build -f example/Dockerfile.gamescope -t gobs-gamescope .
docker run --device /dev/dri --cap-add SYS_NICE gobs-gamescope
```

See [`example/`](example/) for the full working example.
