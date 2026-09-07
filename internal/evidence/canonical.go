package evidence

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"unicode/utf16"
	"unicode/utf8"
)

var (
	ErrInvalidJSON  = errors.New("invalid JSON")
	ErrDuplicateKey = errors.New("duplicate JSON object key")
)

// Canonicalize returns RFC 8785 JSON Canonicalization Scheme bytes.
func Canonicalize(input []byte) ([]byte, error) {
	if !utf8.Valid(input) {
		return nil, fmt.Errorf("%w: input is not UTF-8", ErrInvalidJSON)
	}
	if err := validateUnicodeEscapes(input); err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(input))
	decoder.UseNumber()
	value, err := decodeValue(decoder)
	if err != nil {
		return nil, err
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		if err == nil {
			return nil, fmt.Errorf("%w: trailing value", ErrInvalidJSON)
		}
		return nil, fmt.Errorf("%w: trailing data: %v", ErrInvalidJSON, err)
	}
	var output bytes.Buffer
	if err := writeCanonical(&output, value); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

// Digest computes a lower-case SHA-256 content identifier over domain and data.
func Digest(domain string, data []byte) string {
	hash := sha256.New()
	hash.Write([]byte(domain))
	hash.Write(data)
	return "sha256:" + hex.EncodeToString(hash.Sum(nil))
}

func decodeValue(decoder *json.Decoder) (any, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}
	switch typed := token.(type) {
	case json.Delim:
		switch typed {
		case '{':
			object := make(map[string]any)
			for decoder.More() {
				keyToken, err := decoder.Token()
				if err != nil {
					return nil, fmt.Errorf("%w: object key: %v", ErrInvalidJSON, err)
				}
				key, ok := keyToken.(string)
				if !ok {
					return nil, fmt.Errorf("%w: non-string object key", ErrInvalidJSON)
				}
				if _, exists := object[key]; exists {
					return nil, fmt.Errorf("%w: %q", ErrDuplicateKey, key)
				}
				value, err := decodeValue(decoder)
				if err != nil {
					return nil, err
				}
				object[key] = value
			}
			if _, err := decoder.Token(); err != nil {
				return nil, fmt.Errorf("%w: object close: %v", ErrInvalidJSON, err)
			}
			return object, nil
		case '[':
			var array []any
			for decoder.More() {
				value, err := decodeValue(decoder)
				if err != nil {
					return nil, err
				}
				array = append(array, value)
			}
			if _, err := decoder.Token(); err != nil {
				return nil, fmt.Errorf("%w: array close: %v", ErrInvalidJSON, err)
			}
			return array, nil
		default:
			return nil, fmt.Errorf("%w: unexpected delimiter %q", ErrInvalidJSON, typed)
		}
	default:
		return token, nil
	}
}

func writeCanonical(output *bytes.Buffer, value any) error {
	switch typed := value.(type) {
	case nil:
		output.WriteString("null")
	case bool:
		output.WriteString(strconv.FormatBool(typed))
	case string:
		writeString(output, typed)
	case json.Number:
		number, err := strconv.ParseFloat(string(typed), 64)
		if err != nil || math.IsInf(number, 0) || math.IsNaN(number) {
			return fmt.Errorf("%w: non-finite or unrepresentable number %q", ErrInvalidJSON, typed)
		}
		if number == 0 {
			output.WriteByte('0')
			break
		}
		encoded, _ := json.Marshal(number)
		output.Write(encoded)
	case []any:
		output.WriteByte('[')
		for index, item := range typed {
			if index > 0 {
				output.WriteByte(',')
			}
			if err := writeCanonical(output, item); err != nil {
				return err
			}
		}
		output.WriteByte(']')
	case map[string]any:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(left, right int) bool {
			return compareUTF16(keys[left], keys[right]) < 0
		})
		output.WriteByte('{')
		for index, key := range keys {
			if index > 0 {
				output.WriteByte(',')
			}
			writeString(output, key)
			output.WriteByte(':')
			if err := writeCanonical(output, typed[key]); err != nil {
				return err
			}
		}
		output.WriteByte('}')
	default:
		return fmt.Errorf("%w: unsupported value %T", ErrInvalidJSON, value)
	}
	return nil
}

func compareUTF16(left, right string) int {
	leftUnits := utf16.Encode([]rune(left))
	rightUnits := utf16.Encode([]rune(right))
	for index := 0; index < len(leftUnits) && index < len(rightUnits); index++ {
		if leftUnits[index] < rightUnits[index] {
			return -1
		}
		if leftUnits[index] > rightUnits[index] {
			return 1
		}
	}
	return len(leftUnits) - len(rightUnits)
}

func writeString(output *bytes.Buffer, value string) {
	output.WriteByte('"')
	for _, character := range value {
		switch character {
		case '\b':
			output.WriteString(`\b`)
		case '\t':
			output.WriteString(`\t`)
		case '\n':
			output.WriteString(`\n`)
		case '\f':
			output.WriteString(`\f`)
		case '\r':
			output.WriteString(`\r`)
		case '"', '\\':
			output.WriteByte('\\')
			output.WriteRune(character)
		default:
			if character < 0x20 {
				fmt.Fprintf(output, `\u%04x`, character)
			} else {
				output.WriteRune(character)
			}
		}
	}
	output.WriteByte('"')
}

func validateUnicodeEscapes(input []byte) error {
	inString := false
	escaped := false
	for index := 0; index < len(input); index++ {
		character := input[index]
		if !inString {
			if character == '"' {
				inString = true
			}
			continue
		}
		if escaped {
			escaped = false
			if character != 'u' {
				continue
			}
			if index+4 >= len(input) {
				return fmt.Errorf("%w: incomplete Unicode escape", ErrInvalidJSON)
			}
			value, err := strconv.ParseUint(string(input[index+1:index+5]), 16, 16)
			if err != nil {
				return fmt.Errorf("%w: invalid Unicode escape", ErrInvalidJSON)
			}
			index += 4
			if value >= 0xd800 && value <= 0xdbff {
				if index+6 >= len(input) || input[index+1] != '\\' || input[index+2] != 'u' {
					return fmt.Errorf("%w: unpaired high surrogate", ErrInvalidJSON)
				}
				low, lowErr := strconv.ParseUint(string(input[index+3:index+7]), 16, 16)
				if lowErr != nil || low < 0xdc00 || low > 0xdfff {
					return fmt.Errorf("%w: invalid surrogate pair", ErrInvalidJSON)
				}
				index += 6
			} else if value >= 0xdc00 && value <= 0xdfff {
				return fmt.Errorf("%w: unpaired low surrogate", ErrInvalidJSON)
			}
			continue
		}
		switch character {
		case '\\':
			escaped = true
		case '"':
			inString = false
		}
	}
	return nil
}
