package edgeTTS

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"sort"

	"golang.org/x/crypto/ssh/terminal"
)

type EdgeTTS struct {
	communicator *Communicate
	texts        []CommunicateTextOption
	outCome      io.WriteCloser
}

type Args struct {
	Text           string
	Voice          string
	File           string
	Proxy          string
	Rate           string
	Volume         string
	WordsInCue     float64
	WriteMedia     string
	WriteSubtitles string
	ListVoices     bool
}

func isTerminal(file *os.File) bool {
	return terminal.IsTerminal(int(file.Fd()))
}

func PrintVoices(locale string) {
	// Print all available voices.
	voices, err := listVoices()
	if err != nil {
		return
	}
	sort.Slice(voices, func(i, j int) bool {
		return voices[i].ShortName < voices[j].ShortName
	})

	filterFieldName := map[string]bool{
		"SuggestedCodec": true,
		"FriendlyName":   true,
		"Status":         true,
		"VoiceTag":       true,
		"Language":       true,
	}

	for _, voice := range voices {
		if voice.Locale != locale {
			continue
		}
		fmt.Printf("\n")
		t := reflect.TypeOf(voice)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldName := field.Name
			if filterFieldName[fieldName] {
				continue
			}
			fieldValue := reflect.ValueOf(voice).Field(i).Interface()
			fmt.Printf("%s: %v\n", fieldName, fieldValue)
		}
	}
}

func NewTTS(args Args) *EdgeTTS {
	if isTerminal(os.Stdin) && isTerminal(os.Stdout) && args.WriteMedia == "" {
		fmt.Fprintln(os.Stderr, "Warning: TTS output will be written to the terminal. Use --write-media to write to a file.")
		fmt.Fprintln(os.Stderr, "Press Ctrl+C to cancel the operation. Press Enter to continue.")
		fmt.Scanln()
	}
	tts := NewCommunicate().WithVoice(args.Voice).WithRate(args.Rate).WithVolume(args.Volume)
	file, err := os.OpenFile(args.WriteMedia, os.O_APPEND|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v\n", err)
		return nil
	}
	return &EdgeTTS{
		communicator: tts,
		outCome:      file,
		texts:        []CommunicateTextOption{},
	}
}

func (eTTS *EdgeTTS) AddText(text string, voice string, rate string, volume string) *EdgeTTS {
	eTTS.texts = append(eTTS.texts, CommunicateTextOption{
		text:   text,
		voice:  voice,
		rate:   rate,
		volume: volume,
	})
	return eTTS
}

func (eTTS *EdgeTTS) Speak() {
	defer eTTS.communicator.Close()
	defer eTTS.outCome.Close()

	for _, text := range eTTS.texts {
		task := eTTS.communicator.WithText(text.text).WithVoice(text.voice).WithRate(text.rate).WithVolume(text.volume).stream()
		for {
			v, ok := <-task
			if ok {
				if v.Type == ChunkTypeAudio {
					eTTS.outCome.Write(v.Data)
					// } else if v.Type == ChunkTypeWordBoundary {
				} else if v.Type == ChunkTypeEnd {
					break
				}
			}
		}
	}
}
