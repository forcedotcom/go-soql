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
	"strings"
	"time"
)

const (
	openBrace                     = "("
	closeBrace                    = ")"
	orCondition                   = " OR "
	andCondition                  = " AND "
	singleQuote                   = "'"
	comma                         = ","
	notOperator                   = "NOT "
	openLike                      = " LIKE '%"
	closeLike                     = "%'"
	inOperator                    = " IN "
	equalsOperator                = " = "
	period                        = "."
	null                          = "null"
	notEqualsOperator             = " != "
	greaterThanOperator           = " > "
	greaterThanOrEqualsToOperator = " >= "
	lessThanOperator              = " < "
	lessThanOrEqualsToOperator    = " <= "
	selectKeyword                 = "SELECT "
	whereKeyword                  = " WHERE "
	fromKeyword                   = " FROM "

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
	// LikeOperator is the tag to be used for "like" operator in where clause
	LikeOperator = "likeOperator"
	// NotLikeOperator is the tag to be used for "not like" operator in where clause
	NotLikeOperator = "notLikeOperator"
	// InOperator is the tag to be used for "in" operator in where clause
	InOperator = "inOperator"
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
)

var clauseBuilderMap = map[string]func(v interface{}, fieldName string) string{
	LikeOperator:                  buildLikeClause,
	NotLikeOperator:               buildNotLikeClause,
	InOperator:                    buildInClause,
	EqualsOperator:                buildEqualsClause,
	NullOperator:                  buildNullClause,
	NotEqualsOperator:             buildNotEqualsClause,
	GreaterThanOperator:           buildGreaterThanClause,
	GreaterThanOrEqualsToOperator: buildGreaterThanOrEqualsToClause,
	LessThanOperator:              buildLessThanClause,
	LessThanOrEqualsToOperator:    buildLessThanOrEqualsToClause,
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
)

func buildLikeClause(v interface{}, fieldName string) string {
	return constructLikeClause(v, fieldName, false)
}

func buildNotLikeClause(v interface{}, fieldName string) string {
	return constructLikeClause(v, fieldName, true)
}

func constructLikeClause(v interface{}, fieldName string, exclude bool) string {
	var buff strings.Builder
	patterns, ok := v.([]string)
	if !ok {
		return buff.String()
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
		buff.WriteString(pattern)
		buff.WriteString(closeLike)
		if exclude {
			buff.WriteString(closeBrace)
		}
	}
	if len(patterns) > 1 {
		buff.WriteString(closeBrace)
	}
	return buff.String()
}

func buildInClause(v interface{}, fieldName string) string {
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
		return buff.String()
	}

	if len(items) > 0 {
		buff.WriteString(fieldName)
		buff.WriteString(inOperator)
		buff.WriteString(openBrace)
	}
	for indx, item := range items {
		if indx > 0 {
			buff.WriteString(comma)
		}
		if useSingleQuotes {
			buff.WriteString(singleQuote)
		}
		buff.WriteString(item)
		if useSingleQuotes {
			buff.WriteString(singleQuote)
		}
	}
	if len(items) > 0 {
		buff.WriteString(closeBrace)
	}
	return buff.String()
}

func buildNotEqualsClause(v interface{}, fieldName string) string {
	return constructComparisonClause(v, fieldName, notEqualsOperator)
}

func buildEqualsClause(v interface{}, fieldName string) string {
	return constructComparisonClause(v, fieldName, equalsOperator)
}

func buildGreaterThanClause(v interface{}, fieldName string) string {
	return constructComparisonClause(v, fieldName, greaterThanOperator)
}

func buildGreaterThanOrEqualsToClause(v interface{}, fieldName string) string {
	return constructComparisonClause(v, fieldName, greaterThanOrEqualsToOperator)
}

func buildLessThanClause(v interface{}, fieldName string) string {
	return constructComparisonClause(v, fieldName, lessThanOperator)
}

func buildLessThanOrEqualsToClause(v interface{}, fieldName string) string {
	return constructComparisonClause(v, fieldName, lessThanOrEqualsToOperator)
}

func constructComparisonClause(v interface{}, fieldName, operator string) string {
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
	default:
		return buff.String()
	}

	if value != "" {
		buff.WriteString(fieldName)
		buff.WriteString(operator)
		if useSingleQuotes {
			buff.WriteString(singleQuote)
		}
		buff.WriteString(value)
		if useSingleQuotes {
			buff.WriteString(singleQuote)
		}
	}
	return buff.String()
}

func buildNullClause(v interface{}, fieldName string) string {
	reflectedValue, _, err := getReflectedValueAndType(v)
	if err == ErrNilValue {
		return ""
	}
	val := reflectedValue.Interface()
	allowNull, ok := val.(bool)
	if !ok {
		return ""
	}
	if allowNull {
		return fieldName + equalsOperator + null
	}
	return fieldName + notEqualsOperator + null
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

func marshalWhereClause(v interface{}, tableName string) (string, error) {
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
		partialClause := fn(field.Interface(), columnName)
		if partialClause != "" {
			if previousConditionExists {
				buff.WriteString(andCondition)
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
	return marshalWhereClause(v, "")
}

func getClauseKey(clauseTag string) string {
	tagItems := strings.Split(clauseTag, ",")
	return tagItems[0]
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
		var selectSubString strings.Builder
		var whereValue interface{}
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
			subStr, err := marshalWhereClause(whereValue, relationName)
			if err != nil {
				return "", err
			}
			if subStr != "" {
				buff.WriteString(whereKeyword)
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
