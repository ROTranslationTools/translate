package yandex

import (
	"testing"
)

func TestFetchLanguages(t *testing.T) {
	t.Parallel()

	yx, err := New()
	if err != nil {
		t.Fatalf("initializing new err=%v", err)
	}

	languages, err := yx.FetchLanguages()
	if err != nil {
		t.Fatalf("fetching languages, err=%v", err)
	}

	if len(languages) < 1 {
		t.Errorf("expected atleast one language fetched")
	}
}

func TestValidPrimaryLanguages(t *testing.T) {
	t.Parallel()

	yx, err := New()
	if err != nil {
		t.Fatalf("initializing new err=%v", err)
	}

	testCases := [...]struct {
		data string
		want bool
	}{
		0:  {"bogus", false},
		1:  {"en", true},
		2:  {"", false},
		3:  {" ", false},
		4:  {"  ", false},
		5:  {"sv", true},
		6:  {"ru", true},
		7:  {"en", true},
		8:  {"uk", true},
		9:  {"sl", true},
		10: {"be", true},
		11: {"sr", true},
		12: {"pl", true},
		13: {"it", true},
		14: {"es", true},
		15: {"lt", true},
	}

	for i, tt := range testCases {
		lang := Language(tt.data)
		got, want := yx.ValidPrimaryLanguage(lang), tt.want
		if got != want {
			t.Errorf("#%d: got=%v want=%v; lang=%s", i, got, want, lang)
		}
	}
}

func TestValidTransitionsFromTo(t *testing.T) {
	t.Parallel()

	yx, err := New()
	if err != nil {
		t.Fatalf("initializing new err=%v", err)
	}

	testCases := [...]struct {
		from, to string
		want     bool
	}{
		0:  {"bogus", "bogus", false},
		1:  {"en", "en", false},
		2:  {"uk", "", false},
		3:  {" ", "", false},
		4:  {"  ", "", false},
		5:  {"sv", "ru", true},
		6:  {"ru", "sv", true},
		7:  {"en", "uk", true},
		8:  {"az", "ru", true},
		9:  {"sl", "lt", false},
		10: {"be", "fr", true},
		11: {"sr", "fr", false},
		12: {"ru", "da", true},
		13: {"it", "ru", true},
		14: {"es", "fr", false},
		15: {"sl", "en", true},
	}

	for i, tt := range testCases {
		fromLang, toLang := Language(tt.from), Language(tt.to)
		got, want := yx.ValidTransition(fromLang, toLang), tt.want
		if got != want {
			t.Errorf("#%d: got=%v want=%v; from=%q to=%q", i, got, want, fromLang, toLang)
		}
	}
}
