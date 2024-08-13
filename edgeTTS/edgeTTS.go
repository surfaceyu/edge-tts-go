package edgeTTS

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"

	"golang.org/x/crypto/ssh/terminal"
)

type EdgeTTS struct {
	communicator *Communicate
	texts        []*CommunicateTextTask
	outCome      io.WriteCloser
}

type Args struct {
	Text           string
	Voice          string
	Proxy          string
	Rate           string
	Volume         string
	WordsInCue     float64
	WriteMedia     string
	WriteSubtitles string
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
	if _, err := os.Stat(args.WriteMedia); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(args.WriteMedia), 0755)
		if err != nil {
			log.Fatalf("Failed to create dir: %v\n", err)
			return nil
		}
	}
	tts := NewCommunicate().WithVoice(args.Voice).WithRate(args.Rate).WithVolume(args.Volume)
	file, err := os.OpenFile(args.WriteMedia, os.O_WRONLY|os.O_APPEND|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v\n", err)
		return nil
	}
	tts.openWs()
	return &EdgeTTS{
		communicator: tts,
		outCome:      file,
		texts:        []*CommunicateTextTask{},
	}
}

func (eTTS *EdgeTTS) task(text string, voice string, rate string, volume string) *CommunicateTextTask {
	return &CommunicateTextTask{
		text: text,
		option: CommunicateTextOption{
			voice:  voice,
			rate:   rate,
			volume: volume,
		},
	}
}

func (eTTS *EdgeTTS) AddTextDefault(text string) *EdgeTTS {
	eTTS.texts = append(eTTS.texts, eTTS.task(text, "", "", ""))
	return eTTS
}

func (eTTS *EdgeTTS) AddTextWithVoice(text string, voice string) *EdgeTTS {
	eTTS.texts = append(eTTS.texts, eTTS.task(text, voice, "", ""))
	return eTTS
}

func (eTTS *EdgeTTS) AddText(text string, voice string, rate string, volume string) *EdgeTTS {
	eTTS.texts = append(eTTS.texts, eTTS.task(text, voice, rate, volume))
	return eTTS
}

func (eTTS *EdgeTTS) Speak() {
	defer eTTS.communicator.close()
	defer eTTS.outCome.Close()

	go eTTS.communicator.allocateTask(eTTS.texts)
	eTTS.communicator.createPool()
		_, err := eTTS.outCome.Write(text.speechData)
		if err != nil {
			log.Fatalln("Failed to write to file:", err)
		}
	}
}
