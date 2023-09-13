# Translator

Translator library with default `Json` and `Memory` driver.

## Requirements

### Translatable Structures

By default all translations read from config driver (json files or memory). Structures translations can resolved by structure itself and override global translations. For making structure translatable you must implement `Translatable` interface.

**Note:** Translatable functionality used by validator library.

**Caution:** use non-pointer implemantation for struct!

```go
type Person struct {
    Name string
    Age  string
}

func (p Person) GetTranslation(locale string, key string, field string) string {
    if  key == "required" {
        switch locale {
        case "en":
            if field == "Name" {
                return "Name is required"
            }else{
                return "Age is required"
            }
        }
    }
    return ""
}
```

## Create New Translator Driver

Translator library contains two different driver by default.

### Json Driver

JSON driver use json file for managing translations.

```go
// Signature:
NewJSONTranslator(fallbackLocale string, dir string) (Translator, error)

// Example:
import "github.com/gomig/translator"
jTrans, err := translator.NewJSONTranslator("en", "trans")
```

### Memory Driver

Use in-memory array for keep translations.

```go
// Signature:
NewMemoryTranslator(fallbackLocale string) Translator

// Example:
import "github.com/gomig/translator"
mTrans := translator.NewMemoryTranslator("en")
```

## Usage

Translator interface contains following methods:

### Register

Register new translation message for locale. Use placeholder in message for field name.

```go
// Signature:
Register(locale string, key string, message string)

// Example:
t.Register("en", "welcome", "Hello {name}, welcome!")
```

### Resolve

Resolve find translation for locale. If no translation found for locale return fallback translation or `""`.

```go
// Signature:
Resolve(locale string, key string) string

// Example:
trans := t.Resolve("en", "welcome") // Hello {name}, welcome!
```

### ResolveStruct

Find translation from translatable. If empty string returned from translatable or struct not translatable, default translation will resolved.

```go
// Signature:
ResolveStruct(s any, locale string, key string, field string) string
```

### Translate

Translate get translation for locale.

```go
// Signature:
Translate(locale string, key string, placeholders map[string]string) string

// Example:
t.Translate("en", "welcome", map[string]string{ "name": "John" }) // Hello John, welcome!
```

### TranslateStruct

Translate using translatable interface. if empty string returned from translatable or struct not translatable, default translation will resolved.

```go
// Signature:
TranslateStruct(s any, locale string, key string, field string, placeholders map[string]string) string
```
