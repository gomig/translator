package translator

import "strings"

type record struct {
	locale  string
	key     string
	message string
}

type memTranslator struct {
	fallback string
	records  []record
}

func (this *memTranslator) init(fallbackLocale string) {
	this.fallback = fallbackLocale
	this.records = make([]record, 0)
}

// Register new translation message for locale
// Use placeholder in message for field name
// @example:
// t.Register("en", "welcome", "Hello {name}, welcome!")
func (this *memTranslator) Register(locale string, key string, message string) {
	this.records = append(this.records, record{
		locale:  locale,
		key:     key,
		message: message,
	})
}

// Resolve find translation for locale
// if no translation found for locale return fallback translation or nil
func (this memTranslator) Resolve(locale string, key string) string {
	for _, r := range this.records {
		if r.locale == locale && r.key == key {
			return r.message
		}
	}

	if locale != this.fallback {
		return this.Resolve(this.fallback, key)
	}

	return ""
}

// ResolveStruct find translation from translatable
// if empty string returned from translatable or struct not translatable, default translation will resolved
func (this memTranslator) ResolveStruct(s any, locale string, key string, field string) string {
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
func (this memTranslator) Translate(locale string, key string, placeholders map[string]string) string {
	message := this.Resolve(locale, key)
	for p, v := range placeholders {
		message = strings.ReplaceAll(message, "{"+p+"}", v)
	}
	return message
}

// TranslateStruct translate using translatable interface
// if empty string returned from translatable or struct not translatable, default translation will resolved
// Caution: use non-pointer implemantation for struct
func (this memTranslator) TranslateStruct(s any, locale string, key string, field string, placeholders map[string]string) string {
	message := this.ResolveStruct(s, locale, key, field)
	for p, v := range placeholders {
		message = strings.ReplaceAll(message, "{"+p+"}", v)
	}
	return message
}
