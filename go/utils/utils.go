package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"scraping_offers/go/constants"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type HTMLNode interface {
	Find(selector string) *goquery.Selection
}

func SaveJSON(filename string, data interface{}) error {
	folder := filepath.Dir(filename)
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(data)
}

func ReadJSON(filename string) (map[string]interface{}, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(file, &data)
	return data, err
}

func SaveHTML(htmlContent string, filename string) error {
	err := os.WriteFile(filename, []byte(htmlContent), 0644)
	if err == nil {
		fmt.Printf("HTML saved to %s\n", filename)
	}
	return err
}

func LoadHTML(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	return string(content), err
}

func NormalizeString(s string) string {
	return strings.Title(strings.TrimSpace(s))
}

func DetermineModality(title, description string) string {
	textToSearch := strings.ToLower(title + " " + description)
	keywords := map[string][]string{
		"Hybrid":  {"hybrid", "híbrido", "hibrido", "mixto"},
		"Remote":  {"remote", "remoto", "teletrabajo", "home office"},
		"On-site": {"on-site", "onsite", "presencial", "en oficina"},
	}

	for modality, terms := range keywords {
		for _, term := range terms {
			if strings.Contains(textToSearch, term) {
				return modality
			}
		}
	}
	return "N/A"
}

func IsForbidden(target string, forbiddens []string) bool {
	targetLower := strings.ToLower(target)

	for _, fc := range forbiddens {
		fcLower := strings.ToLower(fc)
		if strings.Contains(targetLower, fcLower) {
			return true
		}
	}
	return false
}

func FindKeywordInDescription(text string) string {
	pattern := regexp.MustCompile(`(?i)\b(crawl|scrap|acqui)\w+`)
	match := pattern.FindString(text)
	return match
}

func GetJSONFromHTML(node HTMLNode, target interface{}) error {
	for _, t := range constants.JsonTypes {
		selector := fmt.Sprintf("script[type='%s']", t)
		scriptTag := node.Find(selector).First()
		if scriptTag.Length() > 0 {
			return DecodeJSONP([]byte(scriptTag.Text()), target)
		}
	}

	var targetText string
	node.Find("script").Each(func(i int, sl *goquery.Selection) {
		text := sl.Text()
		if strings.Contains(text, "window.__INITIAL_PROPS__") {
			targetText = text
		}
	})

	if targetText != "" {
		return DecodeJSONP([]byte(targetText), target)
	}

	lastScript := node.Find("script").Last()
	if lastScript.Length() > 0 {
		return DecodeJSONP([]byte(lastScript.Text()), target)
	}

	return fmt.Errorf("no se encontró JSON válido en el HTML")
}
func DecodeJSONP(content []byte, target interface{}) error {
	text := string(content)
	if strings.Contains(text, "JSON.parse(") {
		re := regexp.MustCompile(`JSON\.parse\("(?P<content>.*)"\);?`)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			var unquoted string
			err := json.Unmarshal([]byte(`"`+matches[1]+`"`), &unquoted)
			if err != nil {
				return fmt.Errorf("error quitando escapes de JSON.parse: %v", err)
			}
			return json.Unmarshal([]byte(unquoted), target)
		}
	}
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start == -1 || end == -1 || end < start {
		return fmt.Errorf("no se encontró estructura JSON {} válida")
	}

	return json.Unmarshal([]byte(text[start:end+1]), target)
}

func FromTimestampToISOFormat(tsMs int64) string {
	t := time.Unix(0, tsMs*int64(time.Millisecond)).UTC()
	return t.Format(time.RFC3339)
}
func RandomImpersonation() string {
	return constants.ImpersonateList[rand.Intn(len(constants.ImpersonateList))]
}

func IsLocal() bool {
	return os.Getenv("ENV") == "local"
}

func SaveDebugHTML(prefix, id, html string) {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s/%s_%s_%s.html", "debug", prefix, timestamp, id)
	os.WriteFile(filename, []byte(html), 0644)
}
