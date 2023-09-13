package translator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/gomig/utils"
	"github.com/tidwall/gjson"
)

// JSONDriver json driver
type JSONDriver struct {
	fallback string
	dir      string
	jsonData string
	data     map[string]string
}

func (this *JSONDriver) init(fallbackLocale string, dir string) error {
	if this.data == nil {
		this.data = make(map[string]string)
	}
	this.fallback = fallbackLocale
	this.dir = dir
	return this.Load()
}

func (JSONDriver) err(pattern string, params ...any) error {
	return utils.TaggedError([]string{"JsonTranslator"}, pattern, params...)
}

// Load load translations file to memory
func (this *JSONDriver) Load() error {
	var resolveFiles = func(dir string) (map[string]string, error) {
		dir = filepath.Dir(path.Join(dir, "some.txt"))
		res := make(map[string]string)
		files := utils.FindFile(dir, ".json")
		for _, f := range files {
			if filepath.Dir(f) != dir {
				continue
			}
			// get file info
			filePath := f
			fileName := filepath.Base(f)
			fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

			// read file
			bytes, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, err
			}

			// validate json
			content := string(bytes)
			if !gjson.Valid(content) {
				return nil, errors.New("Invalid json for " + filePath)
			}

			// append to files list
			res[fileName] = content
		}
		return res, nil
	}

	var unwrapJson = func(jsonText string) (string, error) {
		var res map[string]any
		if err := json.Unmarshal([]byte(jsonText), &res); err != nil {
			return "", err
		}

		if byts, err := json.Marshal(res); err != nil {
			return "", err
		} else {
			jsonStr := string(byts)
			return strings.TrimSuffix(strings.TrimPrefix(jsonStr, "{"), "}"), nil
		}
	}

	locales, err := utils.GetSubDirectory(this.dir)
	if err != nil {
		return this.err(err.Error())
	}

	locales = append(locales, "")
	contents := make([]string, 0)

	for _, locale := range locales {
		files, err := resolveFiles(path.Join(this.dir, locale))
		if err != nil {
			return this.err(err.Error())
		}

		if len(files) == 1 {
			for _, cnt := range files {
				if locale == "" {
					unwrappedJson, err := unwrapJson(cnt)
					if err != nil {
						return this.err(err.Error())
					}
					contents = append(contents, unwrappedJson)
				} else {
					contents = append(contents, `"`+locale+`":`+cnt)
				}
			}
		} else {
			if locale == "" {
				for file, cnt := range files {
					contents = append(contents, `"`+file+`":`+cnt)
				}
			} else {
				subContent := make([]string, 0)
				for file, cnt := range files {
					subContent = append(subContent, `"`+file+`":`+cnt)
				}
				contents = append(contents, `"`+locale+`":{`+strings.Join(subContent, ",")+"}")
			}
		}
	}

	this.jsonData = "{" + strings.Join(contents, ",") + "}"
	return nil
}

// Register new translation message for locale
// Use placeholder in message for field name
// @example:
// t.Register("en", "welcome", "Hello {name}, welcome!")
func (this *JSONDriver) Register(locale string, key string, message string) {
	this.data[fmt.Sprintf("[%s].%s", locale, key)] = message
}

// Resolve find translation for locale
// if no translation found for locale return fallback translation or nil
func (this JSONDriver) Resolve(locale string, key string) string {
	if m, ok := this.data[fmt.Sprintf("[%s].%s", locale, key)]; ok {
		return m
	}

	value := gjson.Get(this.jsonData, locale+"."+key)
	if !value.Exists() {
		value = gjson.Get(this.jsonData, this.fallback+"."+key)
	}
	return value.String()
}

// ResolveStruct find translation from translatable
// if empty string returned from translatable or struct not translatable, default translation will resolved
func (this JSONDriver) ResolveStruct(s any, locale string, key string, field string) string {
	if tr := resolveTranslatable(s); tr != nil {
		tr := tr.GetTranslation(locale, key, field)
		if tr != "" {
			return tr
		}
	}
	return this.Resolve(locale, key)
}

// Translate get translation for locale
// @example:
// t.Translate("en", "welcome", map[string]string{ "name": "John" })
func (this JSONDriver) Translate(locale string, key string, placeholders map[string]string) string {
	message := this.Resolve(locale, key)
	for p, v := range placeholders {
		message = strings.ReplaceAll(message, "{"+p+"}", v)
	}
	return message
}

// TranslateStruct translate using translatable interface
// if empty string returned from translatable or struct not translatable, default translation will resolved
// Caution: use non-pointer implemantation for struct
func (this JSONDriver) TranslateStruct(s any, locale string, key string, field string, placeholders map[string]string) string {
	message := this.ResolveStruct(s, locale, key, field)
	for p, v := range placeholders {
		message = strings.ReplaceAll(message, "{"+p+"}", v)
	}
	return message
}
