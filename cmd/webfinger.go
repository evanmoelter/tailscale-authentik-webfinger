package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type webfinger struct {
	Subject *string `json:"subject"`
	Links   []link  `json:"links"`
}

type link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

func handler(config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resource := r.URL.Query().Get("resource")
		var subject *string
		if resource != "" && strings.HasPrefix(resource, "acct:") && strings.Contains(resource, "@") {
			trimmedSubject := strings.TrimPrefix(resource, "acct:")
			subject = &trimmedSubject
		}

		issuer := (&url.URL{
			Scheme: "https",
			Host:   config.AuthentikHost,
			Path:   path.Join("application/o", config.AuthentikApp) + "/",
		}).String()

		response := webfinger{
			Subject: subject,
			Links: []link{
				{
					Rel:  "http://openid.net/specs/connect/1.0/issuer",
					Href: issuer,
				},
				{
					Rel:  "authorization_endpoint",
					Href: issuer + "oauth2/authorize/",
				},
				{
					Rel:  "token_endpoint",
					Href: issuer + "oauth2/token/",
				},
				{
					Rel:  "userinfo_endpoint",
					Href: issuer + "oauth2/userinfo/",
				},
				{
					Rel:  "jwks_uri",
					Href: issuer + "oauth2/jwks/",
				},
			},
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(response); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(buf.Bytes())
	}
}
