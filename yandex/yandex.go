package yandex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

var (
	ErrNilContext    = errors.New("nil context cannot be used")
	ErrUnimplemented = errors.New("unimplemented")
)

type allLangsResponse struct {
	Languages []Language `json:"dirs,omitempty"`
}

var languageExistanceMarker = struct{}{}

const LanguageDirectionSeparator = "-"

type LanguageMapping map[Language]struct{}

type LanguageDirectory map[Language]LanguageMapping

type Context struct {
	Text []string `json:"text,omitempty"`
	From Language `json:"from,omitempty"`
	To   Language `json:"to,omitempty"`
}

type Response struct {
}

type Yandex struct {
	mu                 sync.Mutex
	apiKey, apiVersion string
	languages          LanguageDirectory
	primaryLanguage    Language     `json:"language,omitempty"`
	credentials        *Credentials `json:"credentials,omitempty"`
}

func New() (*Yandex, error) {
	creds, err := discoverCredentials()
	if err != nil {
		return nil, err
	}
	return NewWithCredentials(creds)
}

func NewWithCredentials(creds *Credentials) (*Yandex, error) {
	if creds == nil {
		return nil, ErrNilCredentials
	}
	yx := &Yandex{
		credentials: creds,
	}
	return yx, nil

}

func NewDiscoverContext() (*Yandex, error) {
	return nil, ErrUnimplemented
}

func (yt *Yandex) setAPIKey(apiKey string) {
	yt.apiKey = apiKey
}

func (yt *Yandex) setAPIVersion(apiVersion string) {
	yt.apiVersion = apiVersion
}

func (yt *Yandex) setPrimaryLanguage(primary Language) {
	yt.primaryLanguage = primary
}

func (yt *Yandex) SetPrimaryLanguage(primary Language) {
	yt.mu.Lock()
	defer yt.mu.Unlock()
	yt.primaryLanguage = primary
}

type DetectionResponse struct {
	Language Language   `json:"lang,omitempty"`
	Code     StatusCode `json:"code,omitempty"`
	Message  string     `json:"message,omitempty"`
}

func flattenTextForRequest(segments []string) string {
	return strings.Replace(strings.Join(segments, "+"), " ", "+", -1)
}

func (yx *Yandex) Detect(ctx *Context) (Language, error) {
	if ctx == nil {
		return UnknownLanguage, ErrNilContext
	}
	flatText := flattenTextForRequest(ctx.Text)
	if flatText == "" {
		return UnknownLanguage, nil
	}

	res, err := invokeAPI(&apiRequest{
		APIVersion: yx.apiVersion,
		Route:      RouteDetect,
		APIKey:     yx.credentials.ApiKey,
		Text:       flatText,
	})

	if err != nil {
		return UnknownLanguage, err
	}

	data, err := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	devLogPrintf("detect, data received=%s", data)
	if err != nil {
		return UnknownLanguage, err
	}

	detectResp := &DetectionResponse{}
	if err := json.Unmarshal(data, detectResp); err != nil {
		return UnknownLanguage, err
	}
	if detectResp.Message != "" {
		return UnknownLanguage, fmt.Errorf("%s", detectResp.Message)
	}
	return detectResp.Language, nil
}

func (yt *Yandex) Translate(ctx *Context) (*Response, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}
	return nil, ErrUnimplemented
}

func mapLanguages(directionSetLangs []Language) LanguageDirectory {
	langMapping := make(LanguageDirectory)
	for _, dirLang := range directionSetLangs {
		splits := strings.Split(string(dirLang), LanguageDirectionSeparator)
		if len(splits) < 2 {
			continue
		}

		primary, secondaries := Language(splits[0]), splits[1:]
		topLevel, topExists := langMapping[Language(primary)]
		if !topExists {
			topLevel = make(LanguageMapping)
			langMapping[primary] = topLevel
		}

		for _, secondary := range secondaries {
			secLang := Language(secondary)
			topLevel[secLang] = languageExistanceMarker
		}
	}

	return langMapping
}

// FetchLanguages retrieves the language directory that
// dictates valid language transitions
func (yx *Yandex) FetchLanguages() (LanguageDirectory, error) {
	if yx.languages != nil {
		return yx.languages, nil
	}

	// Now fetch them
	res, err := invokeAPI(&apiRequest{
		APIVersion: yx.apiVersion,
		Route:      "getLangs",
		APIKey:     yx.credentials.ApiKey,
	})

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	devLogPrintf("fetchLanguages, data received=%s", data)
	if err != nil {
		return nil, err
	}

	langsRes := &allLangsResponse{}
	if err := json.Unmarshal(data, langsRes); err != nil {
		return nil, err
	}

	yx.mu.Lock()
	defer yx.mu.Unlock()

	yx.languages = mapLanguages(langsRes.Languages)
	devLogPrintf("languages mapped %+v", yx.languages)
	return yx.languages, nil
}

// ValidPrimaryLanguage reports whether a language is an
// allowed primary language or a starting point language.
func (yx *Yandex) ValidPrimaryLanguage(lang Language) bool {
	langDirectory, err := yx.FetchLanguages()
	if err != nil {
		devLogPrintf("encountered err=%v when fetching lang directory", err)
		return false
	}
	_, presentAsPrimary := langDirectory[lang]
	return presentAsPrimary
}

// ValidTransition reports whether transitioning from one language
// to another is allowed e.g "en" -> "ro".
func (yx *Yandex) ValidTransition(from, to Language) bool {
	langDirectory, err := yx.FetchLanguages()
	if err != nil {
		return false
	}

	secondaries, primaryPresent := langDirectory[from]
	if !primaryPresent {
		return false
	}
	if secondaries == nil {
		return false
	}
	_, secondaryPresent := secondaries[to]
	return secondaryPresent
}
