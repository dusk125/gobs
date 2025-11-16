package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/dusk125/gobs"
	"github.com/spf13/cobra"
)

var (
	root = &cobra.Command{
		Use: "",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			gobs.SetLogHandler(nil)
			gobs.Startup("en-US")

			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			gobs.AddModulePath(wd+"/lib/obs-plugins", wd+"/data/obs/obs-plugins/%module%")

			toLoad := map[string]bool{
				"obs-ffmpeg":      false,
				"obs-browser":     false,
				"rtmp-services":   false,
				"obs-x264":        false,
				"obs-outputs":     false,
				"obs-transitions": false,
				"image-source":    false,
				"vlc-video":       false,
			}

			for info := range gobs.FindModules() {
				if _, has := toLoad[info.Name()]; has {
					mod, err := gobs.OpenModule(info.BinPath(), info.DataPath())
					if err != nil {
						panic(fmt.Errorf("failed to open module: %v from %v (data: %v): %v", info.Name(), info.BinPath(), info.DataPath(), err))
					}

					if err := mod.Init(); err != nil {
						panic(err)
					}
					toLoad[info.Name()] = true
				} else {
					slog.Warn("failed to find requested module", "module", info.Name())
				}
			}

			for k, v := range toLoad {
				if !v {
					slog.Warn("requested module was not loaded", "module", k)
				}
			}

			gobs.PostLoadModules()

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
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			gobs.Shutdown()
		},
	}
	version = &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), gobs.VersionString())
		},
	}
	inputList = &cobra.Command{
		Use: "input list",
		Run: func(cmd *cobra.Command, args []string) {
			for input := range gobs.InputTypes() {
				fmt.Println(input)
			}
		},
	}
	source = &cobra.Command{
		Use: "source",
	}
	sourceList = &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			for source := range gobs.SourceTypes() {
				fmt.Println(source)
			}
		},
	}
	sourceDefaults = &cobra.Command{
		Use:  "defaults",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			m := gobs.SourceDefaults(args[0])
			return printData(cmd, m)
		},
	}
	filterList = &cobra.Command{
		Use: "filter list",
		Run: func(cmd *cobra.Command, args []string) {
			for filter := range gobs.FilterTypes() {
				fmt.Println(filter)
			}
		},
	}
	transitionList = &cobra.Command{
		Use: "transition list",
		Run: func(cmd *cobra.Command, args []string) {
			for transition := range gobs.TransitionTypes() {
				fmt.Println(transition)
			}
		},
	}
	output = &cobra.Command{
		Use: "output",
	}
	outputList = &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			for output := range gobs.OutputTypes() {
				fmt.Println(output)
			}
		},
	}
	outputDefaults = &cobra.Command{
		Use:  "defaults",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			m := gobs.OutputDefaults(args[0])
			return printData(cmd, m)
		},
	}
	encoder = &cobra.Command{
		Use: "encoder",
	}
	encoderList = &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			for encoder := range gobs.EncoderTypes() {
				fmt.Println(encoder)
			}
		},
	}
	encoderDefaults = &cobra.Command{
		Use:  "defaults",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			m := gobs.EncoderDefaults(args[0])
			return printData(cmd, m)
		},
	}
	service = &cobra.Command{
		Use: "service",
	}
	serviceList = &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			for service := range gobs.ServiceTypes() {
				fmt.Println(service)
			}
		},
	}
	serviceDefaults = &cobra.Command{
		Use:  "defaults",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			m := gobs.ServiceDefaults(args[0])
			return printData(cmd, m)
		},
	}
	module = &cobra.Command{
		Use: "module",
	}
	moduleList = &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			for module := range gobs.FindModules() {
				fmt.Println(module.Name())
			}
		},
	}
)

func main() {
	source.AddCommand(sourceList, sourceDefaults)
	output.AddCommand(outputList, outputDefaults)
	encoder.AddCommand(encoderList, encoderDefaults)
	service.AddCommand(serviceList, serviceDefaults)
	module.AddCommand(moduleList)
	root.AddCommand(version, inputList, source, filterList, transitionList, output, encoder, service, module)

	_ = root.Execute()
}

func printData(cmd *cobra.Command, m gobs.Data) error {
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "\t")
	return enc.Encode(m)
}
