package gojay

import (
	"time"
)

// DecodeTime decodes time with the given format
func (dec *Decoder) DecodeTime(v *time.Time, format string) error {
	if dec.isPooled == 1 {
		panic(InvalidUsagePooledDecoderError("Invalid usage of pooled decoder"))
	}
	return dec.decodeTime(v, format)
}

func (dec *Decoder) decodeTime(v *time.Time, format string) error {
	if format == time.RFC3339 {
		var ej = make(EmbeddedJSON, 0, 20)
		if err := dec.decodeEmbeddedJSON(&ej); err != nil {
			return err
		}
		if err := v.UnmarshalJSON(ej); err != nil {
			return err
		}
		return nil
	}
	var str string
	if err := dec.decodeString(&str); err != nil {
		return err
	}
	tt, err := time.Parse(format, str)
	if err != nil {
		return err
	}
	*v = tt
	return nil
}
