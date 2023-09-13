package translator

// Translatable interface for struct
type Translatable interface {
	GetTranslation(locale string, key string, field string) string
}

func resolveTranslatable(s any) Translatable {
	if v, ok := s.(Translatable); ok {
		return v
	}
	return nil
}
