package yandex

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

var (
	debug = os.Getenv("DEBUG") != ""
)

type Route string

// Format is representative of the output format
// Possible values:
// + plain - Text without markup (default value).
// + html  - Text in HTML format.
type Format string

type Language string

var apiRequestURLTemplate = template.Must(template.New("apiRequestURLTemplate").Parse("https://translate.yandex.net/api/{{ if .APIVersion }}{{.APIVersion}}{{ else }}v1.5{{ end }}/tr.json/{{.Route}}?key={{.APIKey}}"))

func devLogPrintf(fmt_ string, args ...interface{}) {
	if debug {
		log.Printf(fmt_, args...)
	}
}

type apiRequest struct {
	Format     Format `json:"format"`
	Route      Route  `json:"route"`
	APIVersion string `json:"api_version"`
	APIKey     string `json:"api_key"`
}

const (
	RouteDetect    Route = "detect"
	RouteTranslate Route = "translate"
RouteGetLanguages Route = "getLangs"
)

const (
	FormatPlain Format = "plain"
	FormatHTML  Format = "html"
)

func (f Format) String() string {
	switch f {
	default:
		return string(FormatPlain)
	case FormatHTML:
		return string(FormatHTML)
	}
}

func (r Route) String() string {
	switch r {
	default:
		return ""
	case RouteDetect:
		return string(RouteDetect)
	case RouteTranslate:
		return string(RouteTranslate)
	case RouteGetLanguages:
		return string(RouteGetLanguages)
	}
}

func (f *Format) UnmarshalJSON(bRepr []byte) error {
	str := string(bRepr)
	fRepr := Format(str)

	switch fRepr {
	default:
		return fmt.Errorf("unknown format %q", str)
	case FormatPlain, FormatHTML:
		*f = fRepr
		return nil
	}
}

func (f *Format) MarshalJSON() ([]byte, error) {
	return []byte(f.String()), nil
}

func (r *Route) UnmarshalJSON(bRepr []byte) error {
	str := string(bRepr)
	rRepr := Route(str)

	switch rRepr {
	default:
		return fmt.Errorf("unsupported route %q", str)
	case RouteTranslate, RouteDetect:
		*r = rRepr
		return nil
	}
}

func (r *Route) MarshalJSON() ([]byte, error) {
	return []byte(r.String()), nil
}

func invokeAPI(areq *apiRequest) (*http.Response, error) {
	prc, pwc := io.Pipe()
	errsChan := make(chan error)
	go func() {
		err := apiRequestURLTemplate.Execute(pwc, areq)
		_ = pwc.Close()
		errsChan <- err
	}()

	urlData, err := ioutil.ReadAll(prc)
	if err != nil {
		return nil, err
	}

	if err := <-errsChan; err != nil {
		return nil, err
	}
	reqURL := string(urlData)
	devLogPrintf("reqURL=%q apiRequest=%+v", reqURL, areq)
	return http.Get(reqURL)
}
