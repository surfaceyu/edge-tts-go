package edgeTTS

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func uuidWithOutDashes() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func stringToBytes(text interface{}) []byte {
	var textBytes []byte
	switch v := text.(type) {
	case string:
		encoder := unicode.UTF8.NewEncoder()
		encodedText, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(v), encoder))
		if err != nil {
			panic(fmt.Sprintf("Error encoding text: %s", err.Error()))
		}
		textBytes = encodedText
	case []byte:
		textBytes = v
	default:
		panic("str must be string or []byte")
	}
	return textBytes
}

func bytesToString(text interface{}) string {
	var testBytes string
	switch v := text.(type) {
	case string:
		testBytes = v
	case []byte:
		testBytes = string(v)
	default:
		panic("str must be string or []byte")
	}
	return testBytes
}

func splitTextByByteLength(text interface{}, byteLength int) []string {
	// 将字符串转换为字节数组
	textBytes := stringToBytes(text)
	// 按照字节长度拆分字符串
	var result []string
	currentByteLength := 0
	currentString := ""
	for _, b := range textBytes {
		if currentByteLength+len(string(b)) <= byteLength {
			currentString += string(b)
			currentByteLength += len(string(b))
		} else {
			result = append(result, currentString)
			currentString = string(b)
			currentByteLength = len(string(b))
		}
	}
	if currentString != "" {
		result = append(result, currentString)
	}
	return result
}

func mkssml(text interface{}, voice string, rate string, volume string) string {
	textStr := bytesToString(text)
	ssml := fmt.Sprintf("<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xml:lang='en-US'><voice name='%s'><prosody pitch='+0Hz' rate='%s' volume='%s'>%s</prosody></voice></speak>", voice, rate, volume, textStr)
	return ssml
}

func ssmlHeadersPlusData(requestID string, timestamp string, ssml string) string {
	return fmt.Sprintf(
		"X-RequestId:%s\r\n"+
			"Content-Type:application/ssml+xml\r\n"+
			"X-Timestamp:%sZ\r\n"+
			"Path:ssml\r\n\r\n"+
			"%s",
		requestID, timestamp, ssml)
}

func removeIncompatibleCharacters(text interface{}) string {
	cleanedStr := bytesToString(text)
	runes := []rune(cleanedStr)
	for i, r := range runes {
		code := int(r)
		if (0 <= code && code <= 8) || (11 <= code && code <= 12) || (14 <= code && code <= 31) {
			runes[i] = ' '
		}
	}

	return string(runes)
}

func dateToString() string {
	return time.Now().UTC().Format("Mon Jan 02 2006 15:04:05 GMT-0700 (Coordinated Universal Time)")
}

func calcMaxMsgSize(voice string, rate string, volume string) int {
	websocketMaxSize := int(math.Pow(2, 16))
	overheadPerMessage := len(ssmlHeadersPlusData(uuidWithOutDashes(), dateToString(), mkssml("", voice, rate, volume))) + 50
	return websocketMaxSize - overheadPerMessage
}

func getHeadersAndData(data interface{}) (map[string]string, []byte, error) {
	var dataBytes []byte
	switch v := data.(type) {
	case string:
		dataBytes = []byte(v)
	case []byte:
		dataBytes = v
	default:
		return nil, nil, fmt.Errorf("data must be string or []byte")
	}

	headers := make(map[string]string)
	lines := bytes.Split(dataBytes[:bytes.Index(dataBytes, []byte("\r\n\r\n"))], []byte("\r\n"))
	for _, line := range lines {
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) < 2 {
			continue
		}
		key := string(parts[0])
		value := strings.TrimSpace(string(parts[1]))
		headers[key] = value
	}

	return headers, dataBytes[bytes.Index(dataBytes, []byte("\r\n\r\n"))+4:], nil
}
