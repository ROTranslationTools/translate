package yandex

type StatusCode uint

const (
	// Reference https://tech.yandex.com/translate/doc/dg/reference/detect-docpage/
	StatusSuccess                            StatusCode = 200
	StatusInvalidAPIKey                      StatusCode = 401
	StatusBlockedAPIKey                      StatusCode = 402
	StatusDailyLimitExceededOnTranslatedText StatusCode = 404

	// Reference https://tech.yandex.com/translate/doc/dg/reference/translate-docpage/
	StatusMaximumTextSizeExceeded          StatusCode = 413
	StatusTextCannotBeTranslated           StatusCode = 422
	StatusTranslationDirectionNotSupported StatusCode = 501
)
