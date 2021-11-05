package common

import (
	"net/http"
	"net/url"
)

func SetCookie(w http.ResponseWriter, r *http.Request, key, value string, maxAge int) error {
	originUrl, err := url.Parse(r.Header.Get("Origin"))

	if err != nil {
		return err
	}

	useSecure := originUrl.Scheme == "https"

	cookie := &http.Cookie{
		Name:       key,
		Value:     	url.QueryEscape(value),
		Path:       "/",
		Domain:     originUrl.Hostname(),
		MaxAge:     maxAge,
		Secure:     useSecure,
		HttpOnly:   false,
		SameSite:   http.SameSiteStrictMode,
		Unparsed:   nil,
	}

	http.SetCookie(w, cookie)

	return nil
}