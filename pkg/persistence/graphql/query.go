package graphql

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/relvacode/iso8601"
)

func getIndent(
	indent int,
) string {
	output := ""
	for i := 0; i < indent; i++ {
		output += "  "
	}
	return output
}

// TODO: DRY this up- getFields, getObject and getWhere are very similar
func getFields(
	item interface{},
	indent int,
	ignoreEmpty bool,
) (string, error) {
	item = reflect.Indirect(reflect.ValueOf(item)).Interface()

	t := reflect.TypeOf(item)
	v := reflect.ValueOf(item)

	var err error

	fields := ""
	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)

		rawTag := strings.Trim(
			strings.Join(strings.Split(string(fieldType.Tag), ":")[1:], ":"),
			`"`,
		)

		omitEmpty := strings.Contains(rawTag, ",omitempty")

		tag := strings.Split(rawTag, ",")[0]

		if ignoreEmpty && omitEmpty && reflect.DeepEqual(fieldValue.Interface(), reflect.New(fieldType.Type).Elem().Interface()) {
			continue
		}

		extra := ""

		fieldTypeName := fieldType.Type.Name()

		if fieldTypeName != "UUID" && fieldTypeName != "Time" {
			if fieldType.Type.Kind() == reflect.Struct {
				nestedItem := fieldValue.Interface()

				extra, err = getFields(nestedItem, indent+1, ignoreEmpty)
				if err != nil {
					return "", err
				}

				extra = fmt.Sprintf(" %v", extra)
			}
		}

		fields += fmt.Sprintf("%v%v%v\n", getIndent(indent+1), tag, extra)
	}

	fields = strings.TrimRight(fields, "\n")

	return fmt.Sprintf(
		`{
%v
%v}`,
		strings.TrimRight(fields, "\n"),
		getIndent(indent),
	), nil
}

func getValueByField(
	item interface{},
) (map[string]string, error) {
	item = reflect.Indirect(reflect.ValueOf(item)).Interface()

	t := reflect.TypeOf(item)
	v := reflect.ValueOf(item)

	valueByField := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)

		rawTag := strings.Trim(
			strings.Join(strings.Split(string(fieldType.Tag), ":")[1:], ":"),
			`"`,
		)

		omitEmpty := strings.Contains(rawTag, ",omitempty")

		tag := strings.Split(rawTag, ",")[0]

		if omitEmpty && reflect.DeepEqual(fieldValue.Interface(), reflect.New(fieldType.Type).Elem().Interface()) {
			continue
		}

		fieldTypeName := fieldType.Type.Name()
		if fieldTypeName != "UUID" && fieldTypeName != "Time" {
			if fieldType.Type.Kind() == reflect.Struct {
				continue // TODO
			}
		}

		fieldValueInterface := fieldValue.Interface()

		value := ""

		switch fieldTypeName {
		case "UUID":
			value = fmt.Sprintf("%#v", fieldValueInterface.(uuid.UUID).String())
		case "Time":
			value = fmt.Sprintf("%#v", fieldValueInterface.(iso8601.Time).Format("2006-01-02T15:04:05-0700"))
		default:
			value = fmt.Sprintf("%#v", fieldValueInterface)
		}

		valueByField[tag] = value
	}

	return valueByField, nil
}

func getSortedKeys(
	valueByField map[string]string,
) []string {
	sortedKeys := make([]string, 0)
	for k := range valueByField {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)

	return sortedKeys
}

// TODO: DRY this up- getFields, getObject and getWhere are very similar
func getObject(
	item interface{},
	indent int,
) (string, error) {
	item = reflect.Indirect(reflect.ValueOf(item)).Interface()

	t := reflect.TypeOf(item)
	v := reflect.ValueOf(item)

	type FieldTypeAndValue struct {
		Type  reflect.StructField
		Value reflect.Value
	}

	fieldTypeAndValueByName := make(
		map[string]FieldTypeAndValue,
		0,
	)

	for i := 0; i < t.NumField(); i++ {
		fieldTypeAndValue := FieldTypeAndValue{
			Type:  t.Field(i),
			Value: v.Field(i),
		}

		tag := strings.Split(strings.Trim(
			strings.Join(strings.Split(string(fieldTypeAndValue.Type.Tag), ":")[1:], ":"),
			`"`,
		), ",")[0]

		fieldTypeAndValueByName[tag] = fieldTypeAndValue
	}

	sortedKeys := make([]string, 0)
	for key := range fieldTypeAndValueByName {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	var err error

	object := ""
	for _, key := range sortedKeys {
		fieldTypeAndValue := fieldTypeAndValueByName[key]
		fieldType := fieldTypeAndValue.Type
		fieldValue := fieldTypeAndValue.Value

		rawTag := strings.Trim(
			strings.Join(strings.Split(string(fieldType.Tag), ":")[1:], ":"),
			`"`,
		)

		omitEmpty := strings.Contains(rawTag, ",omitempty")

		tag := strings.Split(rawTag, ",")[0]

		if omitEmpty && reflect.DeepEqual(fieldValue.Interface(), reflect.New(fieldType.Type).Elem().Interface()) {
			continue
		}

		extra := ""

		fieldTypeName := fieldType.Type.Name()

		if fieldTypeName != "UUID" && fieldTypeName != "Time" {
			if fieldType.Type.Kind() == reflect.Struct {
				nestedItem := fieldValue.Interface()

				extra, err = getObject(nestedItem, indent+2)
				if err != nil {
					return "", err
				}

				// TODO: hack- relies on this DB
				parts := strings.Split(key, "_")
				constraint := fmt.Sprintf("%v_pkey", parts[len(parts)-1])

				// TODO: hack- use of JSON to get field names
				nestedItemJSON, err := json.Marshal(nestedItem)
				if err != nil {
					return "", err
				}

				// TODO: more hack
				nestedItemMap := make(map[string]interface{})
				err = json.Unmarshal(nestedItemJSON, &nestedItemMap)
				if err != nil {
					return "", err
				}

				// TODO: more hack
				columns := make([]string, 0)
				for k, v := range nestedItemMap {
					if reflect.TypeOf(v).Kind() == reflect.Map {
						continue
					}

					columns = append(columns, k)
				}

				sort.Strings(columns)

				extra = fmt.Sprintf(
					"{\n%vdata: %v,\n%von_conflict: {\n%vconstraint: %v\n%vupdate_columns: %v\n%v}\n%v}",
					getIndent(indent+2),
					extra,
					getIndent(indent+2),
					getIndent(indent+3),
					constraint,
					getIndent(indent+3),
					fmt.Sprintf("[%v]", strings.Join(columns, ", ")),
					getIndent(indent+2),
					getIndent(indent+1),
				)
			}
		}

		fieldValueInterface := fieldValue.Interface()

		value := ""

		if extra == "" {
			switch fieldTypeName {
			case "UUID":
				value = fmt.Sprintf("%#v", fieldValueInterface.(uuid.UUID).String())
			case "Time":
				value = fmt.Sprintf("%#v", fieldValueInterface.(iso8601.Time).Format("2006-01-02T15:04:05-0700"))
			default:
				value = fmt.Sprintf("%#v", fieldValueInterface)
			}
		} else {
			value = extra
		}

		object += fmt.Sprintf("%v%v: %v,\n", getIndent(indent+1), tag, value)
	}

	object = strings.TrimRight(object, ",\n")

	return fmt.Sprintf(
		`{
%v
%v}`,
		strings.TrimRight(object, "\n"),
		getIndent(indent),
	), nil
}

// TODO: DRY this up- getFields, getObject and getWhere are very similar
func getWhere(
	item interface{},
) (string, error) {
	valueByField, err := getValueByField(item)
	if err != nil {
		return "", err
	}

	sortedKeys := getSortedKeys(valueByField)

	whereParts := make([]string, 0)
	for _, field := range sortedKeys {
		whereParts = append(
			whereParts,
			fmt.Sprintf("%v: {_eq: %v}", field, valueByField[field]),
		)
	}

	return fmt.Sprintf(
		"{%v}",
		strings.Join(whereParts, ", "),
	), nil
}

func GetManyQuery(
	key string,
	item interface{},
	conditionKey string,
	conditionValue interface{},
	orderKey string,
	orderDirection string,
) (string, error) {
	fields, err := getFields(item, 1, false)
	if err != nil {
		return "", err
	}

	parts := make([]string, 0)

	if conditionKey != "" {
		parts = append(parts, fmt.Sprintf("where: {%v: {_eq: %#v}}", conditionKey, conditionValue))
	}

	if orderKey != "" && orderDirection != "" {
		parts = append(parts, fmt.Sprintf("order_by: {%v: %v}", orderKey, orderDirection))
	}

	joinedParts := ""
	if len(parts) > 0 {
		joinedParts = fmt.Sprintf(" (%v)", strings.Join(parts, ", "))
	}

	query := fmt.Sprintf(`
{
  %v%v %v
}
`, key, joinedParts, fields)

	return query, nil
}

func GetOneQuery(
	key string,
	item interface{},
	conditionKey string,
	conditionValue interface{},
) (string, error) {
	fields, err := getFields(item, 1, false)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf(`
{
  %v(where: {%v: {_eq: %#v}}, limit: 1, distinct_on: %v) %v
}
`, key, conditionKey, conditionValue, conditionKey, fields)

	return query, nil
}

func InsertQuery(
	key string,
	item interface{},
) (string, error) {
	fields, err := getFields(item, 1, false)
	if err != nil {
		return "", err
	}

	object, err := getObject(item, 1)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf(`
mutation {
  insert_%v_one(object: %v) %v
}
`, key, object, fields)

	return query, nil
}

func DeleteQuery(
	key string,
	item interface{},
) (string, error) {
	fields, err := getFields(item, 2, false)
	if err != nil {
		return "", err
	}

	where, err := getWhere(item)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf(`
mutation {
  delete_%v(where: %v) {
    returning %v
  }
}
`, key, where, fields)

	return query, nil
}
