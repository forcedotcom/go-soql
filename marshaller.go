/*
 * Copyright (c) 2018, salesforce.com, inc.
 * All rights reserved.
 * SPDX-License-Identifier: BSD-3-Clause
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/BSD-3-Clause
 */
package soql

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	openBrace                       = "("
	closeBrace                      = ")"
	orCondition                     = " OR "
	andCondition                    = " AND "
	singleQuote                     = "'"
	safeSingleQuote                 = "\\'"
	doubleQuote                     = "\""
	safeDoubleQuote                 = "\\\""
	backslash                       = "\\"
	safeBackslash                   = "\\\\"
	newLine                         = "\n"
	safeNewLine                     = "\\n"
	carriageReturn                  = "\r"
	safeCarriageReturn              = "\\r"
	tab                             = "\t"
	safeTab                         = "\\t"
	bell                            = "\b"
	safeBell                        = "\\b"
	formFeed                        = "\f"
	safeFormFeed                    = "\\f"
	underscore                      = "_"
	safeUnderscore                  = "\\_"
	percentSign                     = "%"
	safePercentSign                 = "\\%"
	comma                           = ","
	notOperator                     = "NOT "
	openLike                        = " LIKE '%"
	closeLike                       = "%'"
	inOperator                      = " IN "
	notInOperator                   = " NOT IN "
	equalsOperator                  = " = "
	period                          = "."
	null                            = "null"
	notEqualsOperator               = " != "
	greaterThanOperator             = " > "
	greaterThanOrEqualsToOperator   = " >= "
	lessThanOperator                = " < "
	lessThanOrEqualsToOperator      = " <= "
	greaterNextNDaysOperator        = " > NEXT_N_DAYS:"
	greaterOrEqualNextNDaysOperator = " >= NEXT_N_DAYS:"
	equalsNextNDaysOperator         = " = NEXT_N_DAYS:"
	lessNextNDaysOperator           = " < NEXT_N_DAYS:"
	lessOrEqualNextNDaysOperator    = " <= NEXT_N_DAYS:"
	greaterLastNDaysOperator        = " > LAST_N_DAYS:"
	greaterOrEqualLastNDaysOperator = " >= LAST_N_DAYS:"
	equalsLastNDaysOperator         = " = LAST_N_DAYS:"
	lessLastNDaysOperator           = " < LAST_N_DAYS:"
	lessOrEqualLastNDaysOperator    = " <= LAST_N_DAYS:"
	selectKeyword                   = "SELECT "
	whereKeyword                    = " WHERE "
	fromKeyword                     = " FROM "
	orderByKeyword                  = " ORDER BY "
	limitKeyword                    = " LIMIT "
	offsetKeyword                   = " OFFSET "
	ascKeyword                      = " ASC"
	descKeyword                     = " DESC"

	// DateFormat is the golang reference time in the soql dateTime fields format
	DateFormat = "2006-01-02T15:04:05.000-0700"

	// SoqlTag is the main tag name to be used to mark a struct field to be considered for soql marshaling
	SoqlTag = "soql"
	// SelectClause is the tag to be used when marking the struct to be considered for select clause
	SelectClause = "selectClause"
	// TableName is the parameter to be used to specify the name of the underlying SOQL table. It should be used
	// with SelectClause
	TableName = "tableName"
	// SelectColumn is the tag to be used for selecting a column in select clause
	SelectColumn = "selectColumn"
	// SelectChild is the tag to be used when selecting from child tables
	SelectChild = "selectChild"
	// FieldName is the parameter to be used to specify the name of the field in underlying SOQL object
	FieldName = "fieldName"
	// WhereClause is the tag to be used when marking the struct to be considered for where clause
	WhereClause = "whereClause"
	// Joiner is the parameter to be used to specify the joiner to use between properties within a where clause
	Joiner = "joiner"
	// OrderByClause is the tag to be used when marking the string slice to be considered for order by clause
	OrderByClause = "orderByClause"
	// LimitClause is the tag to be used when marking the int to be considered for limit clause
	LimitClause = "limitClause"
	// OffsetClause is the tag to be used when marking the int to be considered for offset clause
	OffsetClause = "offsetClause"
	// LikeOperator is the tag to be used for "like" operator in where clause
	LikeOperator = "likeOperator"
	// NotLikeOperator is the tag to be used for "not like" operator in where clause
	NotLikeOperator = "notLikeOperator"
	// InOperator is the tag to be used for "in" operator in where clause
	InOperator = "inOperator"
	// NotInOperator is the tag to be used for "not in" operator in where clause
	NotInOperator = "notInOperator"
	// EqualsOperator is the tag to be used for "=" operator in where clause
	EqualsOperator = "equalsOperator"
	// NotEqualsOperator is the tag to be used for "!=" operator in where clause
	NotEqualsOperator = "notEqualsOperator"
	// NullOperator is the tag to be used for " = null " or "!= null" operator in where clause
	NullOperator = "nullOperator"
	// GreaterThanOperator is the tag to be used for ">" operator in where clause
	GreaterThanOperator = "greaterThanOperator"
	// GreaterThanOrEqualsToOperator is the tag to be used for ">=" operator in where clause
	GreaterThanOrEqualsToOperator = "greaterThanOrEqualsToOperator"
	// LessThanOperator is the tag to be used for "<" operator in where clause
	LessThanOperator = "lessThanOperator"
	// LessThanOrEqualsToOperator is the tag to be used for "<=" operator in where clause
	LessThanOrEqualsToOperator = "lessThanOrEqualsToOperator"
	//GreaterNextNDaysOperator is the tag to be used for "> NEXT_N_DAYS:n" operator in where clause
	GreaterNextNDaysOperator = "greaterNextNDaysOperator"
	//GreaterNextNDaysOperator is the tag to be used for ">= NEXT_N_DAYS:n" operator in where clause
	GreaterOrEqualNextNDaysOperator = "greaterOrEqualNextNDaysOperator"
	//EqualsNextNDaysOperator is the tag to be used for "= NEXT_N_DAYS:n" operator in where clause
	EqualsNextNDaysOperator = "equalsNextNDaysOperator"
	//LessNextNDaysOperator is the tag to be used for "< NEXT_N_DAYS:n" operator in where clause
	LessNextNDaysOperator = "lessNextNDaysOperator"
	//LessOrEqualNextNDaysOperator is the tag to be used for "<= NEXT_N_DAYS:n" operator in where clause
	LessOrEqualNextNDaysOperator = "lessOrEqualNextNDaysOperator"
	//GreaterLastNDaysOperator is the tag to be used for "> LAST_N_DAYS:n" operator in where clause
	GreaterLastNDaysOperator = "greaterLastNDaysOperator"
	//GreaterOrEqualLastNDaysOperator is the tag to be used for ">= LAST_N_DAYS:n" operator in where clause
	GreaterOrEqualLastNDaysOperator = "greaterOrEqualLastNDaysOperator"
	//EqualsLastNDaysOperator is the tag to be used for "= LAST_N_DAYS:n" operator in where clause
	EqualsLastNDaysOperator = "equalsLastNDaysOperator"
	//LessLastNDaysOperator is the tag to be used for "< LAST_N_DAYS:n" operator in where clause
	LessLastNDaysOperator = "lessLastNDaysOperator"
	//LessOrEqualLastNDaysOperator is the tag to be used for "<= LAST_N_DAYS:n" operator in where clause
	LessOrEqualLastNDaysOperator = "lessOrEqualLastNDaysOperator"

	// Subquery is the tag to be used for a subquery in a where clause
	Subquery = "subquery"
)

var clauseBuilderMap = map[string]func(v interface{}, fieldName string) (string, error){
	LikeOperator:                    buildLikeClause,
	NotLikeOperator:                 buildNotLikeClause,
	InOperator:                      buildInClause,
	NotInOperator:                   buildNotInClause,
	EqualsOperator:                  buildEqualsClause,
	NullOperator:                    buildNullClause,
	NotEqualsOperator:               buildNotEqualsClause,
	GreaterThanOperator:             buildGreaterThanClause,
	GreaterThanOrEqualsToOperator:   buildGreaterThanOrEqualsToClause,
	LessThanOperator:                buildLessThanClause,
	LessThanOrEqualsToOperator:      buildLessThanOrEqualsToClause,
	GreaterNextNDaysOperator:        buildGreaterNextNDaysOperator,
	GreaterOrEqualNextNDaysOperator: buildGreaterOrEqualNextNDaysOperator,
	EqualsNextNDaysOperator:         buildEqualsNextNDaysOperator,
	LessNextNDaysOperator:           buildLessNextNDaysOperator,
	LessOrEqualNextNDaysOperator:    buildLessOrEqualNextNDaysOperator,
	GreaterLastNDaysOperator:        buildGreaterLastNDaysOperator,
	GreaterOrEqualLastNDaysOperator: buildGreaterOrEqualLastNDaysOperator,
	EqualsLastNDaysOperator:         buildEqualsLastNDaysOperator,
	LessLastNDaysOperator:           buildLessLastNDaysOperator,
	LessOrEqualLastNDaysOperator:    buildLessOrEqualLastNDaysOperator,
}

var (
	// ErrInvalidTag error is returned when invalid key is used in soql tag
	ErrInvalidTag = errors.New("ErrInvalidTag")

	// ErrNilValue error is returned when nil pointer is passed as argument
	ErrNilValue = errors.New("ErrNilValue")

	// ErrMultipleSelectClause error is returned when there are multiple selectClause in struct
	ErrMultipleSelectClause = errors.New("ErrMultipleSelectClause")

	// ErrNoSelectClause error is returned when there are No selectClause in struct
	ErrNoSelectClause = errors.New("ErrNoSelectClause")

	// ErrMultipleWhereClause error is returned when there are multiple whereClause in struct
	ErrMultipleWhereClause = errors.New("ErrMultipleWhereClause")

	// ErrInvalidOrderByClause error is returned when field with orderByClause tag is invalid
	ErrInvalidOrderByClause = errors.New("ErrInvalidOrderByClause")

	// ErrMultipleOrderByClause error is returned when there are multiple orderByClause in struct
	ErrMultipleOrderByClause = errors.New("ErrMultipleOrderByClause")

	// ErrInvalidSelectColumnOrderByClause error is returned when the selectColumn
	// associated with the order by clause is invalid
	ErrInvalidSelectColumnOrderByClause = errors.New("ErrInvalidSelectColumnOrderByClause")

	// ErrInvalidLimitClause error is returned when field with limitClause tag is invalid
	ErrInvalidLimitClause = errors.New("ErrInvalidLimitClause")

	// ErrMultipleLimitClause error is returned when there are multiple limitClause in struct
	ErrMultipleLimitClause = errors.New("ErrMultipleLimitClause")

	// ErrInvalidOffsetClause error is returned when field with offsetClause tag is invalid
	ErrInvalidOffsetClause = errors.New("ErrInvalidOffsetClause")

	// ErrMultipleOffsetClause error is returned when there are multiple offsetClause in struct
	ErrMultipleOffsetClause = errors.New("ErrMultipleOffsetClause")
)

// Order is the struct for defining the order by clause on a per column basis
// A slice of this struct tagged with the orderByClause tag in a soql struct
// specifies the columns from the selectClause struct to be included in the
// order by clause and their sort order
type Order struct {
	// Field contains the name of the field of the selectClause struct to be
	// included in the order by clause
	Field string
	// IsDesc indicates whether the ordering is DESC (true) or ASC (false)
	IsDesc bool
}

// https://developer.salesforce.com/docs/atlas.en-us.soql_sosl.meta/soql_sosl/sforce_api_calls_soql_select_quotedstringescapes.htm
var sanitizeCharacters = []string{
	singleQuote, safeSingleQuote,
	doubleQuote, safeDoubleQuote,
	backslash, safeBackslash,
	newLine, safeNewLine,
	carriageReturn, safeCarriageReturn,
	tab, safeTab,
	bell, safeBell,
	formFeed, safeFormFeed,
}

var sanitizeReplacer = strings.NewReplacer(sanitizeCharacters...)

var sanitizeLikeCharacters = append(
	sanitizeCharacters,
	underscore, safeUnderscore,
	percentSign, safePercentSign,
)

var sanitizeLikeReplacer = strings.NewReplacer(sanitizeLikeCharacters...)

func buildLikeClause(v interface{}, fieldName string) (string, error) {
	return constructLikeClause(v, fieldName, false)
}

func buildNotLikeClause(v interface{}, fieldName string) (string, error) {
	return constructLikeClause(v, fieldName, true)
}

func constructLikeClause(v interface{}, fieldName string, exclude bool) (string, error) {
	var buff strings.Builder
	patterns, ok := v.([]string)
	if !ok {
		return buff.String(), ErrInvalidTag
	}
	if len(patterns) > 1 {
		buff.WriteString(openBrace)
	}
	for indx, pattern := range patterns {
		if indx > 0 {
			if exclude {
				buff.WriteString(andCondition)
			} else {
				buff.WriteString(orCondition)
			}
		}
		if exclude {
			buff.WriteString(openBrace)
			buff.WriteString(notOperator)
		}
		buff.WriteString(fieldName)
		buff.WriteString(openLike)
		buff.WriteString(sanitizeLikeReplacer.Replace(pattern))
		buff.WriteString(closeLike)
		if exclude {
			buff.WriteString(closeBrace)
		}
	}
	if len(patterns) > 1 {
		buff.WriteString(closeBrace)
	}
	return buff.String(), nil
}

func buildInClause(v interface{}, fieldName string) (string, error) {
	return constructContainsClause(v, fieldName, inOperator)
}

func buildNotInClause(v interface{}, fieldName string) (string, error) {
	return constructContainsClause(v, fieldName, notInOperator)
}

func constructContainsClause(v interface{}, fieldName string, operator string) (string, error) {
	var buff strings.Builder
	var items []string
	useSingleQuotes := false

	switch u := v.(type) {
	case []string:
		useSingleQuotes = true
		items = u
	case []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []float32, []float64, []bool:
		items = strings.Fields(strings.Trim(fmt.Sprint(u), "[]"))
	case []time.Time:
		for _, item := range u {
			items = append(items, item.Format(DateFormat))
		}
	default:
		return buff.String(), ErrInvalidTag
	}

	if len(items) > 0 {
		buff.WriteString(fieldName)
		buff.WriteString(operator)
		buff.WriteString(openBrace)
	}
	for indx, item := range items {
		if indx > 0 {
			buff.WriteString(comma)
		}
		if useSingleQuotes {
			buff.WriteString(singleQuote)
			buff.WriteString(sanitizeReplacer.Replace(item))
			buff.WriteString(singleQuote)
		} else {
			buff.WriteString(item)
		}
	}
	if len(items) > 0 {
		buff.WriteString(closeBrace)
	}
	return buff.String(), nil
}

func buildNotEqualsClause(v interface{}, fieldName string) (string, error) {
	return constructComparisonClause(v, fieldName, notEqualsOperator)
}

func buildEqualsClause(v interface{}, fieldName string) (string, error) {
	return constructComparisonClause(v, fieldName, equalsOperator)
}

func buildGreaterThanClause(v interface{}, fieldName string) (string, error) {
	return constructComparisonClause(v, fieldName, greaterThanOperator)
}

func buildGreaterThanOrEqualsToClause(v interface{}, fieldName string) (string, error) {
	return constructComparisonClause(v, fieldName, greaterThanOrEqualsToOperator)
}

func buildLessThanClause(v interface{}, fieldName string) (string, error) {
	return constructComparisonClause(v, fieldName, lessThanOperator)
}

func buildLessThanOrEqualsToClause(v interface{}, fieldName string) (string, error) {
	return constructComparisonClause(v, fieldName, lessThanOrEqualsToOperator)
}

func constructComparisonClause(v interface{}, fieldName, operator string) (string, error) {
	var buff strings.Builder
	var value string
	useSingleQuotes := false

	switch u := v.(type) {
	case string:
		useSingleQuotes = true
		value = u
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		value = fmt.Sprint(u)
	case time.Time:
		value = u.Format(DateFormat)
	case *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64, *float32, *float64, *bool:
		if !reflect.ValueOf(u).IsNil() {
			value = fmt.Sprint(reflect.Indirect(reflect.ValueOf(u)))
		}
	case *time.Time:
		if !reflect.ValueOf(u).IsNil() {
			value = reflect.Indirect(reflect.ValueOf(u)).Interface().(time.Time).Format(DateFormat)
		}
	default:
		return buff.String(), ErrInvalidTag
	}

	if value != "" {
		buff.WriteString(fieldName)
		buff.WriteString(operator)
		if useSingleQuotes {
			buff.WriteString(singleQuote)
			buff.WriteString(sanitizeReplacer.Replace(value))
			buff.WriteString(singleQuote)
		} else {
			buff.WriteString(value)
		}
	}
	return buff.String(), nil
}

func buildGreaterNextNDaysOperator(v interface{}, fieldName string) (string, error) {
	return constructDateLiteralsClause(v, fieldName, greaterNextNDaysOperator)
}

func buildGreaterOrEqualNextNDaysOperator(v interface{}, fieldName string) (string, error) {
	return constructDateLiteralsClause(v, fieldName, greaterOrEqualNextNDaysOperator)
}

func buildEqualsNextNDaysOperator(v interface{}, fieldName string) (string, error) {
	return constructDateLiteralsClause(v, fieldName, equalsNextNDaysOperator)
}

func buildLessNextNDaysOperator(v interface{}, fieldName string) (string, error) {
	return constructDateLiteralsClause(v, fieldName, lessNextNDaysOperator)
}

func buildLessOrEqualNextNDaysOperator(v interface{}, fieldName string) (string, error) {
	return constructComparisonClause(v, fieldName, lessOrEqualNextNDaysOperator)

}

func buildGreaterLastNDaysOperator(v interface{}, fieldName string) (string, error) {
	return constructDateLiteralsClause(v, fieldName, greaterLastNDaysOperator)
}

func buildGreaterOrEqualLastNDaysOperator(v interface{}, fieldName string) (string, error) {
	return constructDateLiteralsClause(v, fieldName, greaterOrEqualLastNDaysOperator)
}

func buildEqualsLastNDaysOperator(v interface{}, fieldName string) (string, error) {
	return constructDateLiteralsClause(v, fieldName, equalsLastNDaysOperator)
}

func buildLessLastNDaysOperator(v interface{}, fieldName string) (string, error) {
	return constructDateLiteralsClause(v, fieldName, lessLastNDaysOperator)
}

func buildLessOrEqualLastNDaysOperator(v interface{}, fieldName string) (string, error) {
	return constructDateLiteralsClause(v, fieldName, lessOrEqualLastNDaysOperator)
}

func constructDateLiteralsClause(v interface{}, fieldName string, operator string) (string, error) {
	var buff strings.Builder
	var value string

	switch u := v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		value = fmt.Sprint(u)
	case *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64:
		if !reflect.ValueOf(u).IsNil() {
			value = fmt.Sprint(reflect.Indirect(reflect.ValueOf(u)))
		}
	default:
		return buff.String(), ErrInvalidTag
	}

	if value != "" {
		buff.WriteString(fieldName)
		buff.WriteString(operator)
		buff.WriteString(value)
	}
	return buff.String(), nil
}

func buildNullClause(v interface{}, fieldName string) (string, error) {
	reflectedValue, _, err := getReflectedValueAndType(v)
	if err == ErrNilValue {
		// Not an error case because nil value for *bool is valid
		return "", nil
	}
	val := reflectedValue.Interface()
	allowNull, ok := val.(bool)
	if !ok {
		return "", ErrInvalidTag
	}
	if allowNull {
		return fieldName + equalsOperator + null, nil
	}
	return fieldName + notEqualsOperator + null, nil
}

func getReflectedValueAndType(v interface{}) (reflect.Value, reflect.Type, error) {
	var reflectedValue reflect.Value
	var reflectedType reflect.Type
	if reflect.ValueOf(v).Kind() == reflect.Ptr {
		if reflect.ValueOf(v).IsNil() {
			return reflect.Value{}, nil, ErrNilValue
		}
		reflectedValue = reflect.Indirect(reflect.ValueOf(v))
	} else {
		reflectedValue = reflect.ValueOf(v)
	}
	reflectedType = reflectedValue.Type()
	return reflectedValue, reflectedType, nil
}

// mapSelectColumns maps the selectColumn field name in the soql tag to their
// corresponding field name in the struct needed by marshalOrderByClause
func mapSelectColumns(mappings map[string]string, parent string, gusParent string, v interface{}) error {
	reflectedValue, reflectedType, err := getReflectedValueAndType(v)
	if err != nil {
		return ErrInvalidSelectColumnOrderByClause
	}

	for i := 0; i < reflectedType.NumField(); i++ {
		field := reflectedType.Field(i)
		tag := field.Tag.Get(SoqlTag)
		if tag == "" {
			continue
		}
		// skip all fields that are not tagged as selectColumn
		if getClauseKey(tag) != SelectColumn {
			continue
		}

		fieldName := field.Name
		gusFieldName := getFieldName(tag, field.Name)

		// inside a nested struct, prepend parent to create full field names
		if parent != "" {
			fieldName = parent + period + fieldName
			gusFieldName = gusParent + period + gusFieldName
		}

		fieldValue := reflectedValue.Field(i)

		// the mapping for a struct field should be added regardless, to cover
		// the case of a struct field not being a nested field (e.g. time.Time)
		mappings[fieldName] = gusFieldName
		if fieldValue.Kind() == reflect.Struct {
			err := mapSelectColumns(mappings, fieldName, gusFieldName, fieldValue.Interface())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// v is the Order slice for specifying the columns and sort order
// s is the struct value containing fields with the selectColumn tag
func marshalOrderByClause(v interface{}, tableName string, s interface{}) (string, error) {
	reflectedValue, reflectedType, err := getReflectedValueAndType(v)
	if err != nil {
		return "", err
	}

	if reflectedType.Kind() != reflect.Slice {
		return "", ErrInvalidOrderByClause
	}

	if reflectedType.Elem() != reflect.TypeOf(Order{}) {
		return "", ErrInvalidOrderByClause
	}

	sReflectedValue, sReflectedType, err := getReflectedValueAndType(s)
	if err != nil {
		return "", err
	}

	if sReflectedType.Kind() != reflect.Struct {
		return "", ErrInvalidSelectColumnOrderByClause
	}

	columnMappings := make(map[string]string)

	err = mapSelectColumns(columnMappings, "", "", sReflectedValue.Interface())
	if err != nil {
		return "", err
	}

	if len(columnMappings) == 0 {
		return "", ErrInvalidSelectColumnOrderByClause
	}

	var buff strings.Builder

	previousConditionExists := false
	for i := 0; i < reflectedValue.Len(); i++ {
		order := reflectedValue.Index(i).Interface().(Order)
		fieldName := order.Field
		if strings.TrimSpace(fieldName) == "" {
			return "", ErrInvalidOrderByClause
		}

		columnName, ok := columnMappings[fieldName]
		if !ok {
			return "", ErrInvalidOrderByClause
		}

		if tableName != "" {
			columnName = tableName + period + columnName
		}
		orderString := ascKeyword
		if order.IsDesc {
			orderString = descKeyword
		}
		partialClause := columnName + orderString
		if previousConditionExists {
			buff.WriteString(comma)
		}
		buff.WriteString(partialClause)
		previousConditionExists = true
	}
	return buff.String(), nil

}

// v is the limit value provided
func marshalLimitClause(v interface{}) (string, error) {
	s, err := marshalIntValue(v)
	if err != nil {
		return "", ErrInvalidLimitClause
	}
	return s, nil
}

// v is the offset value provided
func marshalOffsetClause(v interface{}) (string, error) {
	s, err := marshalIntValue(v)
	if err != nil {
		return "", ErrInvalidOffsetClause
	}
	return s, nil
}

func marshalIntValue(v interface{}) (string, error) {
	vPtr, ok := v.(*int)
	if !ok {
		return "", errors.New("invalid type")
	}
	if vPtr == nil {
		return "", nil
	}

	vInt := *vPtr
	if vInt < 0 {
		return "", errors.New("invalid value")
	}

	vString := strconv.Itoa(vInt)

	return vString, nil
}

// MarshalOrderByClause returns a string representing the SOQL order by clause.
// Parameter v is a slice of the Order struct indicating the fields from
// parameter s, which is the value of the select column struct, that should be
// included in the clause and their respective ordering (i.e. ASC or DESC).
// Consider the following struct containing the fields with selectColumn tags:
// type SelectColumns struct {
// 	HostName          string  `soql:"selectColumn,fieldName=Host_Name__c"`
// 	RoleName          string  `soql:"selectColumn,fieldName=Role__r.Name"`
//  NumOfCPUCores     int     `soql:"selectColumn,fieldName=Num_of_CPU_Cores__c"`
// }
// s := SelectColumns{}
// And with an Order slice as follows:
// o := []Order{ Order{Field: "HostName", IsDesc: true},
// 	Order{Field: "NumOfCPUCores", IsDesc: false} }
// By calling MarshalOrderByClause() like the following:
// orderByClause, err := MarshalOrderByClause(o, s)
// if err != nil {
//		log.Warn("Error in marshaling order by clause")
// }
// fmt.Println(orderByClause)
// This will print the orderByClause as:
// Host_Name__c DESC,Num_of_CPU_Cores__c ASC
// For nested structs, specify the field name in the Order slice using a
// <parent>.<child> notation
// For example,
func MarshalOrderByClause(v interface{}, s interface{}) (string, error) {
	return marshalOrderByClause(v, "", s)
}

func marshalWhereClause(v interface{}, tableName, joiner string) (string, error) {
	var buff strings.Builder
	reflectedValue, reflectedType, err := getReflectedValueAndType(v)
	if err != nil {
		return "", err
	}
	previousConditionExists := false
	for i := 0; i < reflectedValue.NumField(); i++ {
		field := reflectedValue.Field(i)
		fieldType := reflectedType.Field(i)
		clauseTag := fieldType.Tag.Get(SoqlTag)
		clauseKey := getClauseKey(clauseTag)
		var partialClause string
		if clauseKey == Subquery {
			if field.Kind() != reflect.Struct && field.Kind() != reflect.Ptr {
				return "", ErrInvalidTag
			}
			if field.Kind() == reflect.Ptr {
				if reflect.ValueOf(field.Interface()).IsNil() {
					continue
				}
			}
			joiner, err := getJoiner(clauseTag)
			if err != nil {
				return "", err
			}
			partialClause, err = marshalWhereClause(field.Interface(), tableName, joiner)
			if err != nil {
				return "", err
			}

			partialClause = openBrace + partialClause + closeBrace
		} else {
			fieldName := getFieldName(clauseTag, fieldType.Name)
			if fieldName == "" {
				return "", ErrInvalidTag
			}
			fn, ok := clauseBuilderMap[clauseKey]
			if !ok {
				return "", ErrInvalidTag
			}
			columnName := fieldName
			if tableName != "" {
				columnName = tableName + period + fieldName
			}
			partialClause, err = fn(field.Interface(), columnName)
			if err != nil {
				return "", err
			}
		}
		if partialClause != "" {
			if previousConditionExists {
				buff.WriteString(joiner)
			}
			buff.WriteString(partialClause)
			previousConditionExists = true
		}
	}
	return buff.String(), nil
}

// MarshalWhereClause returns the string with all conditions that applies for SOQL where clause.
// As part of soql tag, you will need to specify the operator (one of the operators listed below) and
// then specify the name of the field using fieldName parameter.
// Following operators are currently supported:
// 1. LIKE: Like operator. E.g. Host_Name__c LIKE '%-db%'. Use likeOperator in as soql tag
// 2. NOT LIKE: Not like operator. E.g. (NOT Host_Name__c LIKE '%-db%'). Use notLikeOperator in soql tag
// 3. EQUALS (=): Equals operator. E.g. Asset_Type_Asset_Type__c = 'SERVER'. Use equalsOperator in soql tag
// 4. IN: In operator. E.g. Role__r.Name IN ('db','dbmgmt'). Use inOperator in soql tag
// 5. NULL ( = null ): Null operator. E.g. Last_Discovered_Date__c = null. Use nullOperator in soql tag
// 6. NOT NULL: Not null operator. E.g. Last_Discovered_Date__c != null. Use nullOperator in soql tag
// 7. GREATER THAN: Greater than operator. E.g. Last_Discovered_Date__c > 2006-01-02T15:04:05.000-0700. Use greaterThanOperator in soql tag
// 8. GREATER THAN OR EQUALS TO: Greater than or equals to operator. E.g. Num_of_CPU_Cores__c >= 16. Use greaterThanOrEqualsToOperator in soql tag
// 9. LESS THAN: Less than operator. E.g. Last_Discovered_Date__c < 2006-01-02T15:04:05.000-0700. Use lessThanOperator in soql tag
// 10. LESS THAN OR EQUALS TO: Less than or equals to operator. E.g. Num_of_CPU_Cores__c <= 16. Use lessThanOrEqualsToOperator in soql tag
// Consider following go struct
// type TestQueryCriteria struct {
// 	IncludeNamePattern          []string  `soql:"likeOperator,fieldName=Host_Name__c"`
// 	Roles                       []string  `soql:"inOperator,fieldName=Role__r.Name"`
// 	ExcludeNamePattern          []string  `soql:"notLikeOperator,fieldName=Host_Name__c"`
// 	AssetType                   string    `soql:"equalsOperator,fieldName=Tech_Asset__r.Asset_Type_Asset_Type__c"`
// 	AllowNullLastDiscoveredDate *bool     `soql:"nullOperator,fieldName=Last_Discovered_Date__c"`
//  NumOfCPUCores               int       `soql:"greaterThanOperator,fieldName=Num_of_CPU_Cores__c"`
// }
// allowNull := false
// t := TestQueryCriteria{
// 		AssetType:                   "SERVER",
// 		IncludeNamePattern:          []string{"-db", "-dbmgmt"},
// 		Roles:                       []string{"db", "dbmgmt"},
// 		ExcludeNamePattern:          []string{"-core", "-drp"},
// 		AllowNullLastDiscoveredDate: &allowNull,
// 		NumOfCPUCores:               16,
// }
// whereClause, err := MarshalWhereClause(t)
// if err  != nil {
//		log.Warn("Error in marshaling where clause")
// }
// fmt.Println(whereClause)
// This will print whereClause as:
// (Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND Role__r.Name IN ('db','dbmgmt') AND ((NOT Host_Name__c LIKE '%-core%') AND (NOT Host_Name__c LIKE '%-drp%')) AND Tech_Asset__r.Asset_Type_Asset_Type__c = 'SERVER' AND Last_Discovered_Date__c != null AND Num_of_CPU_Cores__c > 16
func MarshalWhereClause(v interface{}) (string, error) {
	return marshalWhereClause(v, "", andCondition)
}

func getClauseKey(clauseTag string) string {
	tagItems := strings.Split(clauseTag, ",")
	return tagItems[0]
}

func getJoiner(clauseTag string) (string, error) {
	tag := getTagValue(clauseTag, Joiner, "")
	switch strings.ToLower(tag) {
	case "or":
		return orCondition, nil
	case "and", "":
		return andCondition, nil
	default:
		return "", ErrInvalidTag
	}
}

func getTagValue(clauseTag, key, defaultValue string) string {
	tagItems := strings.Split(clauseTag, ",")
	value := defaultValue
	for _, tagItem := range tagItems {
		indx := strings.Index(tagItem, "=")
		if indx == -1 {
			continue
		}
		if tagItem[0:indx] == key {
			value = tagItem[indx+1:]
		}
	}
	return value
}

func getFieldName(clauseTag, defaultFieldName string) string {
	return getTagValue(clauseTag, FieldName, defaultFieldName)
}

func getTableName(clauseTag, defaultTableName string) string {
	return getTagValue(clauseTag, TableName, defaultTableName)
}

// MarshalSelectClause returns fields to be included in select clause. Child to parent and parent to child
// relationship is also supported.
// Using selectColumn and fieldName in soql tag lets you specify that the field should be included as part of
// select clause. It lets you specify the name of the field as it is named in SOQL object.
// type NonNestedStruct struct {
// 	Name          string `soql:"selectColumn,fieldName=Name"`
// 	SomeValue     string `soql:"selectColumn,fieldName=SomeValue__c"`
// }
// str, err := MarshalSelectClause(NonNestedStruct{}, "")
// if err  != nil {
//		log.Warn("Error in marshaling select clause")
// }
// fmt.Println(str)
// This will print selectClause as:
// Name,SomeValue__c
//
// Second argument to this function is the relationship name, typically used for parent relationships.
// So call to this function with relatonship name will result in marshaling as follows:
// str, err := MarshalSelectClause(NonNestedStruct{}, "NonNestedStruct__r")
// if err  != nil {
//		log.Warn("Error in marshaling select clause")
// }
// fmt.Println(str)
// This will print selectClause as:
// NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c
//
// You can specify parent relationships as nested structs and specify the relationship name as field name.
// type NestedStruct struct {
// 	ID              string          `soql:"selectColumn,fieldName=Id"`
// 	Name            string          `soql:"selectColumn,fieldName=Name__c"`
// 	NonNestedStruct NonNestedStruct `soql:"selectColumn,fieldName=NonNestedStruct__r"`
// }
// str, err := MarshalSelectClause(NestedStruct{}, "")
// if err  != nil {
//		log.Warn("Error in marshaling select clause")
// }
// fmt.Println(str)
// This will print selectClause as:
// Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c
//
// To specify child relationships in struct you will need to use selectChild tag as follows:
// type ParentStruct struct {
// 	ID              string          `soql:"selectColumn,fieldName=Id"`
// 	Name            string          `soql:"selectColumn,fieldName=Name__c"`
// 	ChildStruct     TestChildStruct `soql:"selectChild,fieldName=Application_Versions__r"`
// }
// type TestChildStruct struct {
// 	SelectClause ChildStruct `soql:"selectClause,tableName=SM_Application_Versions__c"`
// }
// type ChildStruct struct {
// 	Version string `soql:"selectColumn,fieldName=Version__c"`
// }
// str, err := MarshalSelectClause(ParentStruct{}, "")
// if err  != nil {
//		log.Warn("Error in marshaling select clause")
// }
// fmt.Println(str)
// This will print selectClause as:
// Id,Name__c,(SELCT SM_Application_Versions__c.Version__c FROM Application_Versions__r)
func MarshalSelectClause(v interface{}, relationShipName string) (string, error) {
	var buff strings.Builder
	prefix := relationShipName
	if prefix != "" {
		prefix += period
	}
	val, t, err := getReflectedValueAndType(v)
	if err != nil {
		return "", err
	}
	if t.Kind() == reflect.Struct {
		totalFields := t.NumField()
		for i := 0; i < totalFields; i++ {
			field := t.Field(i)
			clauseTag := field.Tag.Get(SoqlTag)
			if clauseTag == "" {
				continue
			}
			clauseKey := getClauseKey(clauseTag)
			isChildRelation := false
			switch clauseKey {
			case SelectColumn:
				isChildRelation = false
			case SelectChild:
				isChildRelation = true
			default:
				return "", ErrInvalidTag
			}
			fieldName := getFieldName(clauseTag, field.Name)
			if fieldName == "" {
				return "", ErrInvalidTag
			}
			if isChildRelation {
				subStr, err := marshal(val.Field(i), field.Type, prefix+fieldName)
				if err != nil {
					return "", err
				}
				buff.WriteString(subStr)
			} else {
				if field.Type.Kind() == reflect.Struct {
					v := reflect.New(field.Type)
					subStr, err := MarshalSelectClause(v.Elem().Interface(), prefix+fieldName)
					if err != nil {
						return "", err
					}
					buff.WriteString(subStr)
				} else {
					buff.WriteString(prefix)
					buff.WriteString(fieldName)
				}
			}
			buff.WriteString(comma)
		}
	} else {
		return "", ErrInvalidTag
	}
	return strings.TrimRight(buff.String(), comma), nil
}

func marshal(reflectedValue reflect.Value, reflectedType reflect.Type, childRelationName string) (string, error) {
	var buff strings.Builder
	if reflectedType.Kind() == reflect.Struct {
		totalFields := reflectedType.NumField()
		if totalFields == 0 {
			// Empty struct
			return "", nil
		}
		soqlTagPresent := false
		selectClausePresent := false
		whereClausePresent := false
		orderByClausePresent := false
		limitClausePresent := false
		offsetClausePresent := false
		var selectSubString strings.Builder
		var selectValue interface{}
		var whereValue interface{}
		var whereJoiner string
		var orderByValue interface{}
		var limitValue interface{}
		var offsetValue interface{}
		tableName := ""
		for i := 0; i < totalFields; i++ {
			field := reflectedType.Field(i)
			clauseTag := field.Tag.Get(SoqlTag)
			if clauseTag == "" {
				continue
			}
			soqlTagPresent = true
			clauseKey := getClauseKey(clauseTag)
			switch clauseKey {
			case SelectClause:
				if selectClausePresent {
					return "", ErrMultipleSelectClause
				}
				selectClausePresent = true
				selectValue = reflectedValue.Field(i).Interface()
				tableName = getTableName(clauseTag, field.Name)
				var relationName string
				if childRelationName == "" {
					relationName = ""
				} else {
					// This is child struct and we should use tableName as prefix for columns in select clause
					relationName = tableName
				}
				subStr, err := MarshalSelectClause(reflectedValue.Field(i).Interface(), relationName)
				if err != nil {
					return "", err
				}
				selectSubString.WriteString(selectKeyword)
				selectSubString.WriteString(subStr)
				selectSubString.WriteString(fromKeyword)
				if childRelationName == "" {
					// This is not a child struct and we should use table name as FROM
					selectSubString.WriteString(tableName)
				} else {
					// This is child struct and we should use relationship name as FROM
					selectSubString.WriteString(childRelationName)
				}
			case WhereClause:
				if whereClausePresent {
					return "", ErrMultipleWhereClause
				}
				whereClausePresent = true
				whereValue = reflectedValue.Field(i).Interface()
				var err error
				whereJoiner, err = getJoiner(clauseTag)
				if err != nil {
					return "", err
				}
			case OrderByClause:
				if orderByClausePresent {
					return "", ErrMultipleOrderByClause
				}
				orderByValue = reflectedValue.Field(i).Interface()
				orderByClausePresent = true
			case LimitClause:
				if limitClausePresent {
					return "", ErrMultipleLimitClause
				}
				limitValue = reflectedValue.Field(i).Interface()
				limitClausePresent = true
			case OffsetClause:
				if offsetClausePresent {
					return "", ErrMultipleOffsetClause
				}
				offsetValue = reflectedValue.Field(i).Interface()
				offsetClausePresent = true
			default:
				return "", ErrInvalidTag
			}
		}
		if !selectClausePresent && soqlTagPresent {
			return "", ErrNoSelectClause
		}
		if childRelationName != "" {
			buff.WriteString(openBrace)
		}
		buff.WriteString(selectSubString.String())
		if whereClausePresent {
			relationName := ""
			if childRelationName != "" {
				// This is child struct and we should use tableName as prefix for columns in where clause
				relationName = tableName
			}
			subStr, err := marshalWhereClause(whereValue, relationName, whereJoiner)
			if err != nil {
				return "", err
			}
			if subStr != "" {
				buff.WriteString(whereKeyword)
				buff.WriteString(subStr)
			}
		}
		if orderByClausePresent {
			relationName := ""
			if childRelationName != "" {
				// This is child struct and we should use tableName as prefix for columns in where clause
				relationName = tableName
			}
			subStr, err := marshalOrderByClause(orderByValue, relationName, selectValue)
			if err != nil {
				return "", err
			}
			if subStr != "" {
				buff.WriteString(orderByKeyword)
				buff.WriteString(subStr)
			}
		}
		if limitClausePresent {
			subStr, err := marshalLimitClause(limitValue)
			if err != nil {
				return "", err
			}
			if subStr != "" {
				buff.WriteString(limitKeyword)
				buff.WriteString(subStr)
			}
		}
		if offsetClausePresent {
			subStr, err := marshalOffsetClause(offsetValue)
			if err != nil {
				return "", err
			}
			if subStr != "" {
				buff.WriteString(offsetKeyword)
				buff.WriteString(subStr)
			}
		}
		if childRelationName != "" {
			buff.WriteString(closeBrace)
		}
	} else if childRelationName != "" {
		// Child relationship used for non struct member
		return "", ErrInvalidTag
	}
	return buff.String(), nil
}

// Marshal constructs the entire SOQL query based on the golang struct passed to it.
// Consider following example:
// type TestSoqlStruct struct {
// 	SelectClause NestedStruct      `soql:"selectClause,tableName=SM_Logical_Host__c"`
// 	WhereClause  TestQueryCriteria `soql:"whereClause"`
// }
// type TestQueryCriteria struct {
// 	IncludeNamePattern          []string `soql:"likeOperator,fieldName=Host_Name__c"`
// 	Roles                       []string `soql:"inOperator,fieldName=Role__r.Name"`
// }
// type NonSoqlStruct struct {
// 	Key   string
// 	Value string
// }
// type NonNestedStruct struct {
// 	Name          string `soql:"selectColumn,fieldName=Name"`
// 	SomeValue     string `soql:"selectColumn,fieldName=SomeValue__c"`
// 	NonSoqlStruct NonSoqlStruct
// }
// type NestedStruct struct {
// 	ID              string          `soql:"selectColumn,fieldName=Id"`
// 	Name            string          `soql:"selectColumn,fieldName=Name__c"`
// 	NonNestedStruct NonNestedStruct `soql:"selectColumn,fieldName=NonNestedStruct__r"`
// }
// soqlStruct := TestSoqlStruct {
//    SelectClause: NestedStruct{},
//    WhereClause: TestQueryCriteria{
//        IncludeNamePattern: []string{"-db", "-dbmgmt"},
//        Roles: []string{"db"},
//    }
// }
// str, err := Marshal(soqlStruct)
// if err  != nil {
//		log.Warn("Error in marshaling soql")
// }
// fmt.Println(str)
// This will print soql query as:
// SELECT Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c FROM SM_Logical_Host__c WHERE (Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND Role__r.Name IN ('db','dbmgmt')
func Marshal(v interface{}) (string, error) {
	rv, rt, err := getReflectedValueAndType(v)
	if err != nil {
		return "", err
	}
	return marshal(rv, rt, "")
}
