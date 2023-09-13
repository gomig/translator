package translator_test

import (
	"testing"

	"github.com/gomig/translator"
)

type person struct {
	name string
}

func (this person) GetTranslation(locale string, key string, field string) string {
	if key == "required" && field == "name" {
		switch locale {
		case "en":
			return "Name is required"
		case "fa":
			return "نام الزامی است"
		}
	}
	if key == "placeholder" && field == "name" {
		return "Text with placeholder of {name}"
	}
	return ""
}

func TestRegisterResolve(t *testing.T) {
	tr := translator.NewMemoryTranslator("en")
	tr.Register("en", "welcome", "Hello {name}, welcome!")
	if tr.Resolve("en", "welcome") != "Hello {name}, welcome!" {
		t.Error("failed to register/resolve translation!")
	}
}

func TestResolveStruct(t *testing.T) {
	p := person{name: "john"}
	tr := translator.NewMemoryTranslator("en")
	if tr.ResolveStruct(p, "fa", "required", "name") != "نام الزامی است" {
		t.Error("failed to resolve struct!")
	}
}

func TestTranslate(t *testing.T) {
	tr := translator.NewMemoryTranslator("en")
	tr.Register("en", "welcome", "Hello {name}, welcome!")
	placeholders := map[string]string{"name": "John"}
	if tr.Translate("en", "welcome", placeholders) != "Hello John, welcome!" {
		t.Error("failed to translate!")
	}
}

func TestTranslateStruct(t *testing.T) {
	p := person{name: "john"}
	tr := translator.NewMemoryTranslator("en")
	placeholders := map[string]string{"name": "John"}
	if tr.TranslateStruct(p, "en", "placeholder", "name", placeholders) != "Text with placeholder of John" {
		t.Error("failed to translate struct!")
	}
}
