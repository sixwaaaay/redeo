package resp

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Append implements ResponseWriter
func (b *bufioW) Append(v interface{}) error {
	switch v := v.(type) {
	case nil:
		b.AppendNil()
	case CustomResponse:
		v.AppendTo(b)
	case error:
		msg := v.Error()
		if !strings.HasPrefix(msg, "ERR ") {
			msg = "ERR " + msg
		}
		b.AppendError(msg)
	case bool:
		if v {
			b.AppendInt(1)
		} else {
			b.AppendInt(0)
		}
	case int:
		b.AppendInt(int64(v))
	case int8:
		b.AppendInt(int64(v))
	case int16:
		b.AppendInt(int64(v))
	case int32:
		b.AppendInt(int64(v))
	case int64:
		b.AppendInt(v)
	case uint:
		b.AppendInt(int64(v))
	case uint8:
		b.AppendInt(int64(v))
	case uint16:
		b.AppendInt(int64(v))
	case uint32:
		b.AppendInt(int64(v))
	case uint64:
		b.AppendInt(int64(v))
	case string:
		b.AppendBulkString(v)
	case []byte:
		b.AppendBulk(v)
	case CommandArgument:
		b.AppendBulk(v)
	case float32:
		b.AppendInlineString(strconv.FormatFloat(float64(v), 'f', -1, 32))
	case float64:
		b.AppendInlineString(strconv.FormatFloat(v, 'f', -1, 64))
	default:
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(v)

			b.AppendArrayLen(s.Len())
			for i := 0; i < s.Len(); i++ {
				if err := b.Append(s.Index(i).Interface()); err != nil {
					return err
				}
			}
		case reflect.Map:
			s := reflect.ValueOf(v)

			b.AppendArrayLen(s.Len() * 2)
			for _, key := range s.MapKeys() {
				if err := b.Append(key.Interface()); err != nil {
					return err
				}
				if err := b.Append(s.MapIndex(key).Interface()); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("resp: unsupported type %T", v)
		}
	}
	return nil
}
