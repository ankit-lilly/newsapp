package utils

import (
	"net/url"
)

func UrlToString(original string) string {
	escaped := url.QueryEscape(original)
	return escaped
}

func StringToUrl(original string) (string, error) {
	return url.QueryUnescape(original)
}
