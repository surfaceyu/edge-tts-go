package edgeTTS

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Voice struct {
	Name           string `json:"Name"`
	ShortName      string `json:"ShortName"`
	Gender         string `json:"Gender"`
	Locale         string `json:"Locale"`
	SuggestedCodec string `json:"SuggestedCodec"`
	FriendlyName   string `json:"FriendlyName"`
	Status         string `json:"Status"`
	Language       string
	VoiceTag       VoiceTag `json:"VoiceTag"`
}
type VoiceTag struct {
	ContentCategories  []string `json:"ContentCategories"`
	VoicePersonalities []string `json:"VoicePersonalities"`
}

func listVoices() ([]Voice, error) {
	// Send GET request to retrieve the list of voices.
	client := http.Client{}
	req, err := http.NewRequest("GET", VOICE_LIST, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authority", "speech.platform.bing.com")
	req.Header.Set("Sec-CH-UA", `" Not;A Brand";v="99", "Microsoft Edge";v="91", "Chromium";v="91"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36 Edg/91.0.864.41")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response.
	var voices []Voice
	err = json.Unmarshal(body, &voices)
	if err != nil {
		return nil, err
	}

	return voices, nil
}

type VoicesManager struct {
	voices       []Voice
	calledCreate bool
}

func (vm *VoicesManager) create(customVoices []Voice) error {
	vm.voices = customVoices
	if customVoices == nil {
		voices, err := listVoices()
		if err != nil {
			return err
		}
		vm.voices = voices
	}
	for i, voice := range vm.voices {
		locale := voice.Locale
		if locale == "" {
			return errors.New("Invalid voice locale")
		}
		language := locale[:2]
		vm.voices[i].Language = language
	}
	vm.calledCreate = true
	return nil
}

func (vm *VoicesManager) find(attributes Voice) []Voice {
	if !vm.calledCreate {
		panic("VoicesManager.find() called before VoicesManager.create()")
	}

	var matchingVoices []Voice
	for _, voice := range vm.voices {
		matched := true
		if attributes.Language != "" && attributes.Language != voice.Language {
			matched = false
		}
		if attributes.Name != "" && attributes.Name != voice.Name {
			matched = false
		}
		if attributes.Gender != "" && attributes.Gender != voice.Gender {
			matched = false
		}
		if attributes.Locale != "" && attributes.Locale != voice.Locale {
			matched = false
		}
		if matched {
			matchingVoices = append(matchingVoices, voice)
		}
	}
	return matchingVoices
}
