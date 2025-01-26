package godm

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	CR   = "\r"
	LF   = "\n"
	CRLF = CR + LF
	LFCR = LF + CR
)

func detectLineEnding(raw string) string {
	if strings.Contains(raw, CRLF) {
		return CRLF
	}
	if strings.Contains(raw, LFCR) {
		return LFCR
	}
	if strings.Contains(raw, CR) {
		return CR
	}
	return LF
}

type KV struct {
	Key   string
	Value string
}

func parseIntoKVs(raw string) ([]KV, error) {
	var (
		kvs        = []KV{}
		trimmed    = strings.TrimSpace(raw)
		lineEnding = detectLineEnding(raw)
	)

	lines := strings.Split(trimmed, lineEnding)
	for _, line := range lines {
		// skip empty lines
		if line == "" {
			continue
		}
		k, v, err := parseLine(line)
		if err != nil {
			return nil, err
		}

		kvs = append(kvs, KV{k, v})
	}

	return kvs, nil
}

func parseLine(line string) (k string, v string, err error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return "", "", fmt.Errorf("can't parse empty line")
	}

	// parse comment
	if strings.HasPrefix(line, "COMMENT") {
		return "COMMENT", strings.TrimSpace(line[7:]), nil
	}

	// parse KVN
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid input: %s", line)
	}

	value := strings.TrimSpace(parts[1])

	// remove units if present
	unitIndex := strings.LastIndex(value, "[")
	if unitIndex != -1 && strings.HasSuffix(value, "]") {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(value[:unitIndex]), nil
	}

	return strings.TrimSpace(parts[0]), value, nil
}

// TODO: write tests
func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// TODO: write tests
func parseTime(s string) (time.Time, error) {
	//  handle this: 2020-065T16:00:00
	if len(strings.Split(s, "-")) == 2 {
		dayOfYear, err := strconv.Atoi(strings.Split(s, "-")[1][:3])
		if err != nil {
			return time.Time{}, err
		}
		year, err := strconv.Atoi(strings.Split(s, "-")[0])
		if err != nil {
			return time.Time{}, err
		}
		timePart := strings.Split(s, "T")[1]
		layout := "2006-002T15:04:05"
		if strings.Contains(timePart, ".") {
			layout = "2006-002T15:04:05.999999999"
		}
		t, err := time.Parse(layout, fmt.Sprintf("%04d-%03dT%s", year, dayOfYear, timePart))
		if err != nil {
			return time.Time{}, err
		}
		return t, nil
	}

	if !strings.Contains(s, "Z") {
		s += "Z"
	}
	return time.Parse(time.RFC3339, s)
}

type Field struct {
	Name       string
	ReflectVal reflect.Value
	Required   bool
}

// pure reflection insanity
func getODMFields(v interface{}) ([]Field, error) {
	var (
		fields = []Field{}
		val    = reflect.ValueOf(v).Elem()
	)

	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		switch field.Type.Kind() {
		case reflect.Struct:
			if field.Type.Name() == "Time" {
				tag, required := parseOdmTag(field)
				if tag == "" {
					continue
				}
				fields = append(fields, Field{
					Name:       tag,
					ReflectVal: val.Field(i),
					Required:   required,
				})
				continue
			}
			f := val.Field(i).Addr().Interface()
			subFields, err := getODMFields(f)
			if err != nil {
				return nil, err
			}
			fields = append(fields, subFields...)
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.Struct {
				fields = append(fields, Field{
					Name:       field.Name,
					ReflectVal: val.Field(i),
					Required:   false,
				})
				break
			}
			fallthrough
		default:
			tag, required := parseOdmTag(field)
			if tag == "" {
				continue
			}
			fields = append(fields, Field{
				Name:       tag,
				ReflectVal: val.Field(i),
				Required:   required,
			})
		}
	}

	return fields, nil
}

// pure reflection insanity
func populateODMFields(fields []Field, kvs []KV) error {
	kvIndex := 0
	for _, f := range fields {
		if kvIndex >= len(kvs) {
			if f.Required {
				return fmt.Errorf("expected key %s, got none, more fields than KVs", f.Name)
			}
			continue
		}

		kv := kvs[kvIndex]
		if f.Name != kv.Key && f.Required {
			return fmt.Errorf("expected key %s, got %s", f.Name, kv.Key)
		}

		if f.ReflectVal.Kind() == reflect.Slice && f.ReflectVal.Type().Elem().Kind() == reflect.Struct {
		outer:
			for {
				instance := reflect.New(f.ReflectVal.Type().Elem()).Interface()
				subFields, err := getODMFields(instance)
				if err != nil {
					return err
				}

			inner:
				for _, subField := range subFields {
					if kvIndex >= len(kvs) {
						break outer
					}
					localKv := kvs[kvIndex]

					if subField.Name != localKv.Key && subField.Name != "COMMENT" {
						break outer
					}

					if subField.Name != localKv.Key {
						continue inner
					}

					switch subField.ReflectVal.Kind() {
					case reflect.Int:
						n, err := strconv.Atoi(localKv.Value)
						if err != nil {
							return err
						}
						subField.ReflectVal.SetInt(int64(n))
						kvIndex++
					case reflect.String:
						subField.ReflectVal.SetString(localKv.Value)
						kvIndex++
					case reflect.Float64:
						n, err := parseFloat(localKv.Value)
						if err != nil {
							return err
						}
						subField.ReflectVal.SetFloat(n)
						kvIndex++
					case reflect.Struct:
						if subField.ReflectVal.Type().Name() == "Time" {
							t, err := parseTime(localKv.Value)
							if err != nil {
								return errors.Wrapf(err, "failed to parse time for key %s", localKv.Key)
							}
							subField.ReflectVal.Set(reflect.ValueOf(t))
							kvIndex++
						}
					case reflect.Slice:
						for kvs[kvIndex].Key == subField.Name {
							subField.ReflectVal.Set(reflect.Append(subField.ReflectVal, reflect.ValueOf(kvs[kvIndex].Value)))
							kvIndex++
						}
					default:
						return fmt.Errorf("unsupported type %s", subField.ReflectVal.Kind())
					}
				}

				f.ReflectVal.Set(reflect.Append(f.ReflectVal, reflect.ValueOf(instance).Elem()))
			}
		}

		if f.Name == kv.Key {
			switch f.ReflectVal.Kind() {
			case reflect.Int:
				println("parsing int", kv.Value)
				n, err := strconv.Atoi(kv.Value)
				if err != nil {
					return err
				}
				f.ReflectVal.SetInt(int64(n))
				kvIndex++
			case reflect.String:
				f.ReflectVal.SetString(kv.Value)
				kvIndex++
			case reflect.Float64:
				n, err := parseFloat(kv.Value)
				if err != nil {
					return err
				}
				f.ReflectVal.SetFloat(n)
				kvIndex++
			case reflect.Slice:
				for kvs[kvIndex].Key == f.Name {
					f.ReflectVal.Set(reflect.Append(f.ReflectVal, reflect.ValueOf(kvs[kvIndex].Value)))
					kvIndex++
				}
			case reflect.Struct:
				if f.ReflectVal.Type().Name() == "Time" {
					t, err := parseTime(kv.Value)
					if err != nil {
						return errors.Wrapf(err, "failed to parse time for key %s", kv.Key)
					}
					f.ReflectVal.Set(reflect.ValueOf(t))
					kvIndex++
				}
			default:
				return fmt.Errorf("unsupported type %s", f.ReflectVal.Kind())
			}
		}
	}
	return nil
}

func parseOdmTag(f reflect.StructField) (string, bool) {
	tag := f.Tag.Get("odm")
	if tag == "" {
		return "", false
	}

	var (
		parts    = strings.Split(tag, ",")
		name     = parts[0]
		required = false
	)

	if len(parts) > 1 && parts[1] == "required" {
		required = true
	}

	return name, required
}
