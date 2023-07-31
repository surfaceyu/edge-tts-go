package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/pflag"
	"github.com/surfaceyu/edge-tts-go/edgeTTS"
)

func usage() {
	fmt.Println("Microsoft Edge TTS")
	pflag.PrintDefaults()
}

func main() {
	listVoices := pflag.BoolP("list-voices", "l", false, "lists available voices and exits")
	locale := pflag.StringP("locale", "", "", "locale for TTS ex: zh-CN, en-US")
	text := pflag.StringP("text", "t", "", "what TTS will say")
	file := pflag.StringP("file", "f", "", "same as --text but read from file")
	voice := pflag.StringP("voice", "v", "zh-CN-XiaoxiaoNeural", "voice for TTS")
	volume := pflag.String("volume", "+0%", "set TTS volume")
	rate := pflag.String("rate", "+0%", "set TTS rate")
	writeMedia := pflag.String("write-media", "", "send media output to file instead of stdout")
	// proxy := pflag.String("proxy", "", "use a proxy for TTS and voice list")
	pflag.Usage = usage
	pflag.Parse()

	if *listVoices {
		edgeTTS.PrintVoices(*locale)
		os.Exit(0)
	}
	if *file != "" {
		if *file == "/dev/stdin" {
			reader := bufio.NewReader(os.Stdin)
			*text, _ = reader.ReadString('\n')
		} else {
			data, _ := ioutil.ReadFile(*file)
			*text = string(data)
		}
	}
	if *text != "" {
		args := edgeTTS.Args{
			Text:       *text,
			Voice:      *voice,
			File:       *file,
			ListVoices: *listVoices,
			Rate:       *rate,
			Volume:     *volume,
			WriteMedia: *writeMedia,
		}
		edgeTTS.NewTTS(args).AddText(args.Text, args.Voice, args.Rate, args.Volume).Speak()
	}
}
