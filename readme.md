# edge-tts-go

`edge-tts-go` is a golang module that allows you to use Microsoft Edge's online text-to-speech service from within your golang code or using the provided `edge-tts-go` command.

## Installation

To install it, run the following command:

    $ go install github.com/surfaceyu/edge-tts-go

## Usage

### Basic usage

If you want to use the `edge-tts-go` command, you can simply run it with the following command:

    $ edge-tts-go --text "Hello, world!" --write-media hello.mp3 --write-subtitles hello.vtt

### Changing the voice

If you want to change the language of the speech or more generally, the voice. 

You must first check the available voices with the `--list-voices` option:

    $ edge-tts-go --list-voices
    Name: Microsoft Server Speech Text to Speech Voice (af-ZA, AdriNeural)
    ShortName: af-ZA-AdriNeural
    Gender: Female
    Locale: af-ZA

    Name: Microsoft Server Speech Text to Speech Voice (am-ET, MekdesNeural)
    ShortName: am-ET-MekdesNeural
    Gender: Female
    Locale: am-ET

    Name: Microsoft Server Speech Text to Speech Voice (ar-EG, SalmaNeural)
    ShortName: ar-EG-SalmaNeural
    Gender: Female
    Locale: ar-EG

    Name: Microsoft Server Speech Text to Speech Voice (ar-SA, ZariyahNeural)
    ShortName: ar-SA-ZariyahNeural
    Gender: Female
    Locale: ar-SA

    ...

    $ edge-tts-go --voice zh-CN-XiaoxiaoNeural --text "秦时明月汉时关，万里长征人未还" --write-media hello_in_chinese.mp3

### Changing rate and volume

It is possible to make minor changes to the generated speech.

    $ edge-tts-go --rate=-50% --text "Hello, world!" --write-media hello_with_rate_halved.mp3
    $ edge-tts-go --volume=-50% --text "Hello, world!" --write-media hello_with_volume_halved.mp3

## go module

It is possible to use the `edge-tts-go` module directly from go. For a list of example applications:

* https://github.com/surfaceyu/edge-tts-go/edgeTTS

## thanks

* https://github.com/rany2/edge-tts
