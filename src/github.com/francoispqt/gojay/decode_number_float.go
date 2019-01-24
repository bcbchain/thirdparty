package gojay

// DecodeFloat64 reads the next JSON-encoded value from its input and stores it in the float64 pointed to by v.
//
// See the documentation for Unmarshal for details about the conversion of JSON into a Go value.
func (dec *Decoder) DecodeFloat64(v *float64) error {
	if dec.isPooled == 1 {
		panic(InvalidUsagePooledDecoderError("Invalid usage of pooled decoder"))
	}
	return dec.decodeFloat64(v)
}
func (dec *Decoder) decodeFloat64(v *float64) error {
	for ; dec.cursor < dec.length || dec.read(); dec.cursor++ {
		switch c := dec.data[dec.cursor]; c {
		case ' ', '\n', '\t', '\r', ',':
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			val, err := dec.getFloat()
			if err != nil {
				return err
			}
			*v = val
			return nil
		case '-':
			dec.cursor = dec.cursor + 1
			val, err := dec.getFloatNegative()
			if err != nil {
				return err
			}
			*v = -val
			return nil
		case 'n':
			dec.cursor++
			err := dec.assertNull()
			if err != nil {
				return err
			}
			dec.cursor++
			return nil
		default:
			dec.err = dec.makeInvalidUnmarshalErr(v)
			err := dec.skipData()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return dec.raiseInvalidJSONErr(dec.cursor)
}

func (dec *Decoder) getFloatNegative() (float64, error) {
	// look for following numbers
	for ; dec.cursor < dec.length || dec.read(); dec.cursor++ {
		switch dec.data[dec.cursor] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return dec.getFloat()
		default:
			return 0, dec.raiseInvalidJSONErr(dec.cursor)
		}
	}
	return 0, dec.raiseInvalidJSONErr(dec.cursor)
}

func (dec *Decoder) getFloat() (float64, error) {
	var end = dec.cursor
	var start = dec.cursor
	// look for following numbers
	for j := dec.cursor + 1; j < dec.length || dec.read(); j++ {
		switch dec.data[j] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			end = j
			continue
		case '.':
			// we get part before decimal as integer
			beforeDecimal := dec.atoi64(start, end)
			// then we get part after decimal as integer
			start = j + 1
			// get number after the decimal point
			// multiple the before decimal point portion by 10 using bitwise
			for i := j + 1; i < dec.length || dec.read(); i++ {
				c := dec.data[i]
				if isDigit(c) {
					end = i
					beforeDecimal = (beforeDecimal << 3) + (beforeDecimal << 1)
					continue
				} else if (c == 'e' || c == 'E') && j < i-1 {
					afterDecimal := dec.atoi64(start, end)
					dec.cursor = i + 1
					expI := end - start + 2
					if expI >= len(pow10uint64) || expI < 0 {
						return 0, dec.raiseInvalidJSONErr(dec.cursor)
					}
					pow := pow10uint64[expI]
					floatVal := float64(beforeDecimal+afterDecimal) / float64(pow)
					exp, err := dec.getExponent()
					if err != nil {
						return 0, err
					}
					pExp := (exp + (exp >> 31)) ^ (exp >> 31) + 1 // abs
					if pExp >= int64(len(pow10uint64)) || pExp < 0 {
						return 0, dec.raiseInvalidJSONErr(dec.cursor)
					}
					// if exponent is negative
					if exp < 0 {
						return float64(floatVal) * (1 / float64(pow10uint64[pExp])), nil
					}
					return float64(floatVal) * float64(pow10uint64[pExp]), nil
				}
				dec.cursor = i
				break
			}
			if end >= dec.length || end < start {
				return 0, dec.raiseInvalidJSONErr(dec.cursor)
			}
			// then we add both integers
			// then we divide the number by the power found
			afterDecimal := dec.atoi64(start, end)
			expI := end - start + 2
			if expI >= len(pow10uint64) || expI < 0 {
				return 0, dec.raiseInvalidJSONErr(dec.cursor)
			}
			pow := pow10uint64[expI]
			return float64(beforeDecimal+afterDecimal) / float64(pow), nil
		case 'e', 'E':
			dec.cursor = j + 1
			// we get part before decimal as integer
			beforeDecimal := uint64(dec.atoi64(start, end))
			// get exponent
			exp, err := dec.getExponent()
			if err != nil {
				return 0, err
			}
			pExp := (exp + (exp >> 31)) ^ (exp >> 31) + 1 // abs
			if pExp >= int64(len(pow10uint64)) || pExp < 0 {
				return 0, dec.raiseInvalidJSONErr(dec.cursor)
			}
			// if exponent is negative
			if exp < 0 {
				return float64(beforeDecimal) * (1 / float64(pow10uint64[pExp])), nil
			}
			return float64(beforeDecimal) * float64(pow10uint64[pExp]), nil
		case ' ', '\n', '\t', '\r', ',', '}', ']': // does not have decimal
			dec.cursor = j
			return float64(dec.atoi64(start, end)), nil
		}
		// invalid json we expect numbers, dot (single one), comma, or spaces
		return 0, dec.raiseInvalidJSONErr(dec.cursor)
	}
	return float64(dec.atoi64(start, end)), nil
}

// DecodeFloat32 reads the next JSON-encoded value from its input and stores it in the float32 pointed to by v.
//
// See the documentation for Unmarshal for details about the conversion of JSON into a Go value.
func (dec *Decoder) DecodeFloat32(v *float32) error {
	if dec.isPooled == 1 {
		panic(InvalidUsagePooledDecoderError("Invalid usage of pooled decoder"))
	}
	return dec.decodeFloat32(v)
}
func (dec *Decoder) decodeFloat32(v *float32) error {
	for ; dec.cursor < dec.length || dec.read(); dec.cursor++ {
		switch c := dec.data[dec.cursor]; c {
		case ' ', '\n', '\t', '\r', ',':
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			val, err := dec.getFloat32()
			if err != nil {
				return err
			}
			*v = val
			return nil
		case '-':
			dec.cursor = dec.cursor + 1
			val, err := dec.getFloat32Negative()
			if err != nil {
				return err
			}
			*v = -val
			return nil
		case 'n':
			dec.cursor++
			err := dec.assertNull()
			if err != nil {
				return err
			}
			dec.cursor++
			return nil
		default:
			dec.err = dec.makeInvalidUnmarshalErr(v)
			err := dec.skipData()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return dec.raiseInvalidJSONErr(dec.cursor)
}

func (dec *Decoder) getFloat32Negative() (float32, error) {
	// look for following numbers
	for ; dec.cursor < dec.length || dec.read(); dec.cursor++ {
		switch dec.data[dec.cursor] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return dec.getFloat32()
		default:
			return 0, dec.raiseInvalidJSONErr(dec.cursor)
		}
	}
	return 0, dec.raiseInvalidJSONErr(dec.cursor)
}

func (dec *Decoder) getFloat32() (float32, error) {
	var end = dec.cursor
	var start = dec.cursor
	// look for following numbers
	for j := dec.cursor + 1; j < dec.length || dec.read(); j++ {
		switch dec.data[j] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			end = j
			continue
		case '.':
			// we get part before decimal as integer
			beforeDecimal := dec.atoi32(start, end)
			// then we get part after decimal as integer
			start = j + 1
			// get number after the decimal point
			// multiple the before decimal point portion by 10 using bitwise
			for i := j + 1; i < dec.length || dec.read(); i++ {
				c := dec.data[i]
				if isDigit(c) {
					end = i
					beforeDecimal = (beforeDecimal << 3) + (beforeDecimal << 1)
					continue
				} else if (c == 'e' || c == 'E') && j < i-1 {
					afterDecimal := dec.atoi32(start, end)
					dec.cursor = i + 1
					expI := end - start + 2
					if expI >= len(pow10uint64) || expI < 0 {
						return 0, dec.raiseInvalidJSONErr(dec.cursor)
					}
					pow := pow10uint64[expI]
					floatVal := float32(beforeDecimal+afterDecimal) / float32(pow)
					exp, err := dec.getExponent()
					if err != nil {
						return 0, err
					}
					pExp := (exp + (exp >> 31)) ^ (exp >> 31) + 1 // abs
					if pExp >= int64(len(pow10uint64)) || pExp < 0 {
						return 0, dec.raiseInvalidJSONErr(dec.cursor)
					}
					// if exponent is negative
					if exp < 0 {
						return float32(floatVal) * (1 / float32(pow10uint64[pExp])), nil
					}
					return float32(floatVal) * float32(pow10uint64[pExp]), nil
				}
				dec.cursor = i
				break
			}
			if end >= dec.length || end < start {
				return 0, dec.raiseInvalidJSONErr(dec.cursor)
			}
			// then we add both integers
			// then we divide the number by the power found
			afterDecimal := dec.atoi32(start, end)
			expI := end - start + 2
			if expI >= len(pow10uint64) || expI < 0 {
				return 0, dec.raiseInvalidJSONErr(dec.cursor)
			}
			pow := pow10uint64[expI]
			return float32(beforeDecimal+afterDecimal) / float32(pow), nil
		case 'e', 'E':
			dec.cursor = j + 1
			// we get part before decimal as integer
			beforeDecimal := uint32(dec.atoi32(start, end))
			// get exponent
			exp, err := dec.getExponent()
			if err != nil {
				return 0, err
			}
			pExp := (exp + (exp >> 31)) ^ (exp >> 31) + 1
			// log.Print(exp, " after")
			if pExp >= int64(len(pow10uint64)) || pExp < 0 {
				return 0, dec.raiseInvalidJSONErr(dec.cursor)
			}
			// if exponent is negative
			if exp < 0 {
				return float32(beforeDecimal) * (1 / float32(pow10uint64[pExp])), nil
			}
			return float32(beforeDecimal) * float32(pow10uint64[pExp]), nil
		case ' ', '\n', '\t', '\r', ',', '}', ']': // does not have decimal
			dec.cursor = j
			return float32(dec.atoi32(start, end)), nil
		}
		// invalid json we expect numbers, dot (single one), comma, or spaces
		return 0, dec.raiseInvalidJSONErr(dec.cursor)
	}
	return float32(dec.atoi32(start, end)), nil
}
