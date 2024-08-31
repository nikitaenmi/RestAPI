package speller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	_ "github.com/lib/pq"
)

type Mistakes struct {
	Word string   `json:"word"`
	S    []string `json:"s"`
}

func CheckText(text string) (string, error) {

	resp, err := http.Get("https://speller.yandex.net/services/spellservice.json/checkText" + "?text=S" + url.QueryEscape(text))
	if err != nil {
		return text, fmt.Errorf("error get: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return text, fmt.Errorf("error read: %w", err)
	}

	var mistakes []Mistakes
	err = json.Unmarshal(body, &mistakes)
	if err != nil {
		return text, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	for _, Mistakes := range mistakes {
		if len(Mistakes.S) > 0 {
			text = strings.Replace(text, Mistakes.Word, Mistakes.S[0], -1)
		}
	}
	return text, nil

}
