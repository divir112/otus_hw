package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

type ValidationError struct {
	Field string
	Err   error
}

func (ve ValidationError) Error() string {
	return ve.Err.Error()
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	res := strings.Builder{}
	for _, e := range v {
		res.WriteString(e.Err.Error() + "\n")
	}
	return res.String()
}

func checkMin(fieldType reflect.StructField, fieldI interface{}, rule string) error {
	if fieldType.Type.Kind() == reflect.Int {
		value := fieldI.(int)
		expectedValue, _ := strconv.Atoi(rule)
		if value < expectedValue {
			return ValidationError{
				Field: fieldType.Name,
				Err:   fmt.Errorf("field %s is not valid %d, expected min %d", fieldType.Name, value, expectedValue),
			}
		}
		return nil
	}
	return ValidationError{
		Field: fieldType.Name,
		Err:   fmt.Errorf("invalid validation tag %s for field %s", fieldType.Tag.Get("validate"), fieldType.Name),
	}
}

func checkMax(fieldType reflect.StructField, fieldI interface{}, rule string) error {
	if fieldType.Type.Kind() == reflect.Int {
		value := fieldI.(int)
		expectedValue, _ := strconv.Atoi(rule)
		if value > expectedValue {
			return ValidationError{
				Field: fieldType.Name,
				Err:   fmt.Errorf("field %s is not valid %d, expected max %d", fieldType.Name, value, expectedValue),
			}
		}

		return nil
	}

	return ValidationError{
		Field: fieldType.Name,
		Err:   fmt.Errorf("invalid validation tag %s for field %s", fieldType.Tag.Get("validate"), fieldType.Name),
	}
}

func checkLen(fieldType reflect.StructField, fieldValue reflect.Value, rule string) error {
	if fieldType.Type.Kind() == reflect.Slice || fieldType.Type.Kind() == reflect.String {
		expectedValue, _ := strconv.Atoi(rule)
		if fieldValue.Len() != expectedValue {
			return ValidationError{
				Field: fieldType.Name,
				Err:   fmt.Errorf("field %s is not valid %d, expected len %d", fieldType.Name, fieldValue.Len(), expectedValue),
			}
		}

		return nil
	}
	return ValidationError{
		Field: fieldType.Name,
		Err:   fmt.Errorf("invalid validation tag %s for field %s", fieldType.Tag.Get("validate"), fieldType.Name),
	}
}

func checkRegexp(fieldType reflect.StructField, fieldValue reflect.Value, rule string) error {
	if fieldType.Type.Kind() == reflect.String {
		value := fieldValue.String()
		r, err := regexp.Compile(rule)
		if err != nil {
			return ValidationError{
				Field: fieldType.Name,
				Err:   fmt.Errorf("invalid regexp"),
			}
		}

		res := r.MatchString(value)
		if !res {
			return ValidationError{
				Field: fieldType.Name,
				Err:   fmt.Errorf("field %s is not valid %s, expected regexp %s", fieldType.Name, value, rule),
			}
		}

		return nil
	}

	return ValidationError{
		Field: fieldType.Name,
		Err:   fmt.Errorf("invalid validation tag %s for field %s", fieldType.Tag.Get("validate"), fieldType.Name),
	}
}

func checkIn(fieldType reflect.StructField, fieldValue reflect.Value, rule string) error {
	if fieldType.Type.Kind() == reflect.String {
		value := fieldValue.String()
		setValues := strings.Split(rule, ",")
		res := slices.Contains(setValues, value)
		if !res {
			return ValidationError{
				Field: fieldType.Name,
				Err:   fmt.Errorf("field %s must not contain %s, expected values %s", fieldType.Name, value, setValues),
			}
		}
		return nil
	}

	if fieldType.Type.Kind() == reflect.Int {
		fieldI := fieldValue.Interface()
		value := fieldI.(int)
		setValues := strings.Split(rule, ",")
		strV := strconv.Itoa(value)
		res := slices.Contains(setValues, strV)
		if !res {
			return ValidationError{
				Field: fieldType.Name,
				Err:   fmt.Errorf("field %s must not contain %d, expected values %s", fieldType.Name, value, setValues),
			}
		}
		return nil
	}
	return ValidationError{
		Field: fieldType.Name,
		Err:   fmt.Errorf("invalid validation tag %s for field %s", fieldType.Tag.Get("validate"), fieldType.Name),
	}
}

func Validate(v interface{}) error {
	r := reflect.ValueOf(v)
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("argument is not struct")
	}
	errs := ValidationErrors{}
	for i := range r.NumField() {
		fieldValue := r.Field(i)
		fieldType := t.Field(i)
		fieldName := t.Field(i).Name
		if !unicode.IsUpper(rune(fieldName[0])) {
			continue
		}
		fieldI := fieldValue.Interface()

		tag := fieldType.Tag.Get("validate")
		rules := strings.Split(tag, "|")
		for _, r := range rules {
			if r == "" {
				continue
			}
			rule := strings.Split(r, ":")
			ruleName := rule[0]
			ruleValue := rule[1]
			var err error
			switch ruleName {
			case "min":
				err = checkMin(fieldType, fieldI, ruleValue)
			case "max":
				err = checkMax(fieldType, fieldI, ruleValue)
			case "len":
				err = checkLen(fieldType, fieldValue, ruleValue)
			case "regexp":
				err = checkRegexp(fieldType, fieldValue, ruleValue)
			case "in":
				err = checkIn(fieldType, fieldValue, ruleValue)
			}
			var validationError ValidationError
			if errors.As(err, &validationError) {
				errs = append(errs, validationError)
			}
		}
	}
	if len(errs) != 0 {
		return errs
	}
	return nil
}
