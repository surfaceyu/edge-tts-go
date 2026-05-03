package edgeTTS

import "strings"

const (
	TRUSTED_CLIENT_TOKEN  = "6A5AA1D4EAFF4E9FB37E23D68491D6F4"
	WSS_URL               = "wss://speech.platform.bing.com/consumer/speech/synthesize/readaloud/edge/v1?TrustedClientToken=" + TRUSTED_CLIENT_TOKEN
	VOICE_LIST            = "https://speech.platform.bing.com/consumer/speech/synthesize/readaloud/voices/list?trustedclienttoken=" + TRUSTED_CLIENT_TOKEN
	CHROMIUM_FULL_VERSION = "143.0.3650.75"
	SEC_MS_GEC_VERSION    = "1-" + CHROMIUM_FULL_VERSION
)

var CHROMIUM_MAJOR_VERSION = strings.SplitN(CHROMIUM_FULL_VERSION, ".", 2)[0]

// Locale
const (
	ZhCN = "zh-CN"
	EnUS = "en-US"
)

const (
	ChunkTypeAudio        = "Audio"
	ChunkTypeWordBoundary = "WordBoundary"
	ChunkTypeSessionEnd   = "SessionEnd"
	ChunkTypeEnd          = "ChunkEnd"
)
