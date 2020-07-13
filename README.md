# SOQL

This package supports marshalling a golang struct into SOQL. Like `json` tags, this package provides `soql` tags that you can use to annotate your golang structs. Once tagged `Marshal` method will return SOQL query that will let you query the required Salesforce object using Salesforce API.

## Introduction

Please refer to [introduction](./docs/introduction.md) to understand the basics. This blog [post](https://developer.salesforce.com/blogs/2020/02/soql-tags-for-golang.html) also captures basics in little more detail.

Once you read through it you can refer to documentation below that covers features of this repo in more depth.

## How to use

Start with using `soql` tags on members of your golang structs. `soql` is the main tag. There are following subtags supported:

```
    selectClause // is the tag to be used when marking the struct to be considered for select clause in soql.
    whereClause // is the tag to be used when marking the struct to be considered for where clause in soql.
    orderByClause // is the tag to be used when marking the Order slice to be considered for order by clause in soql.
    limitClause // is the tag to be used when marking the *int to be considered for limit clause in soql.
    offsetClause // is the tag to be used when marking the *int to be considered for offset clause in soql.
    selectColumn // is the tag to be used for selecting a column in select clause. It should be used on members of struct that have been tagged with selectClause.
    selectChild // is the tag to be used when selecting from child tables. It should be used on members of struct that have been tagged with selectClause.
    likeOperator // is the tag to be used for "like" operator in where clause. It should be used on members of struct that have been tagged with whereClause.
    notLikeOperator // is the tag to be used for "not like" operator in where clause. It should be used on members of struct that have been tagged with whereClause.
    inOperator // is the tag to be used for "in" operator in where clause. It should be used on members of struct that have been tagged with whereClause.
    equalsOperator // is the tag to be used for "=" operator in where clause. It should be used on members of struct that have been tagged with whereClause.
    notEqualsOperator // is the tag to be used for "!=" operator in where clause. It should be used on members of struct that have been tagged with whereClause.
    nullOperator // is the tag to be used for " = null " or "!= null" operator in where clause. It should be used on members of struct that have been tagged with whereClause.
    greaterThanOperator // is the tag to be used for ">" operator in where clause. It should be used on members of struct that have been tagged with whereClause.
    lessThanOperator // is the tag to be used for "<" operator in where clause. It should be used on members of struct that have been tagged with whereClause.
    greaterThanOrEqualsToOperator // is the tag to be used for ">=" operator in where clause. It should be used on members of struct that have been tagged with whereClause.
    lessThanOrEqualsToOperator // is the tag to be used for "<=" operator in where clause. It should be used on members of struct that have been tagged with whereClause.
```

Following are supported parameters:

```
    fieldName // is the parameter to be used to specify the name of the field in underlying Salesforce object. It can be used with all tags listed above other than selectClause and whereClause.
    tableName // is the parameter to be used to specify the name of the table of underlying Salesforce Object. It can be be used only with selectClause.

```

If `fieldName` and `tableName` parameters are not provided then the name of the field will be used as default.

### Basic Usage

Lets take a look at one example of a simple non-nested struct and how it can be used to construct a soql query:

```
type TestSoqlStruct struct {
	SelectClause  NonNestedStruct   `soql:"selectClause,tableName=SM_SomeObject__c"`
	WhereClause   TestQueryCriteria `soql:"whereClause"`
	OrderByClause []Order           `soql:"orderByClause"`
	LimitClause   *int              `soql:"limitClause"`
	OffsetClause  *int              `soql:"offsetClause"`
}
type TestQueryCriteria struct {
	IncludeNamePattern          []string `soql:"likeOperator,fieldName=Name__c"`
	Roles                       []string `soql:"inOperator,fieldName=Role__c"`
}
type NonNestedStruct struct {
	Name          string `soql:"selectColumn,fieldName=Name__c"`
	SomeValue     string `soql:"selectColumn,fieldName=SomeValue__c"`
}
```

To use above structs to create SOQL query

```
limit := 5
offset := 10
soqlStruct := TestSoqlStruct{
    WhereClause: TestQueryCriteria {
        IncludeNamePattern: []string{"foo", "bar"},
        Roles: []string{"admin", "user"},
    },
    OrderByClause: []Order{Order{Field:"Name", IsDesc:true}},
    LimitClause: &limit,
    OffsetClause: &offset,
}
soqlQuery, err := Marshal(soqlStruct)
if err != nil {
    fmt.Printf("Error in marshaling: %s\n", err.Error())
}
fmt.Println(soqlQuery)
```

Above struct will result in following SOQL query:

```
SELECT Name__c,SomeValue__c FROM SM_SomeObject__C WHERE (Name__c LIKE '%foo%' OR Name__c LIKE '%bar%') AND Role__c IN ('admin','user') ORDER BY Name__c DESC LIMIT 5 OFFSET 10
```

### Advanced usage

#### Relationships

This package supports child to parent as well as parent to child relationships. Here's a more complex example that includes both the relationships and how the soql query is marshalled:

```
type ComplexSoqlStruct struct {
	SelectClause ParentStruct  `soql:"selectClause,tableName=SM_Parent__c"`
	WhereClause  QueryCriteria `soql:"whereClause"`
}

type QueryCriteria struct {
	IncludeNamePattern []string `soql:"likeOperator,fieldName=Name__c"`
	Roles              []string `soql:"inOperator,fieldName=Role__r.Name"`
	ExcludeNamePattern []string `soql:"notLikeOperator,fieldName=Name__c"`
	SomeType           string   `soql:"equalsOperator,fieldName=Some_Parent__r.Some_Type__c"`
	Status             string   `soql:"notEqualsOperator,fieldName=Status__c"`
	AllowNullValue     *bool    `soql:"nullOperator,fieldName=Value__c"`
}

type ParentStruct struct {
	ID                string          `soql:"selectColumn,fieldName=Id"`
	Name              string          `soql:"selectColumn,fieldName=Name__c"`
	NonNestedStruct   NonNestedStruct `soql:"selectColumn,fieldName=NonNestedStruct__r"` // child to parent relationship
	ChildStruct       TestChildStruct `soql:"selectChild,fieldName=Child__r"`            // parent to child relationship
	SomeNonSoqlMember string          `json:"some_nonsoql_member"`
}

type NonNestedStruct struct {
	Name          string `soql:"selectColumn,fieldName=Name"`
	SomeValue     string `soql:"selectColumn,fieldName=SomeValue__c"`
	NonSoqlStruct NonSoqlStruct
}

type NonSoqlStruct struct {
	Key   string
	Value string
}

type TestChildStruct struct {
	SelectClause ChildStruct        `soql:"selectClause,tableName=SM_Child__c"`
	WhereClause  ChildQueryCriteria `soql:"whereClause"`
}

type ChildStruct struct {
	Version string `soql:"selectColumn,fieldName=Version__c"`
}

type ChildQueryCriteria struct {
	Name string `soql:"equalsOperator,fieldName=Name__c"`
}

allowNull := false
soqlStruct := ComplexSoqlStruct{
    SelectClause: ParentStruct{
        ChildStruct: TestChildStruct{
            WhereClause: ChildQueryCriteria{
                Name: "some-name",
            },
        },
    },
    WhereClause: QueryCriteria{
        SomeType:           "typeA",
        IncludeNamePattern: []string{"-foo", "-bar"},
        Roles:              []string{"admin", "user"},
        ExcludeNamePattern: []string{"-far", "-baz"},
        Status:             "InActive",
        AllowNullValue:     &allowNull,
    },
}
soqlQuery, err := soql.Marshal(soqlStruct)
if err != nil {
    fmt.Printf("Error in marshalling: %s\n", err.Error())
}
fmt.Println(soqlQuery)
```

Above struct will result in following SOQL query:

```
SELECT Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c,(SELECT SM_Child__c.Version__c FROM Child__r WHERE SM_Child__c.Name__c = 'some-name') FROM SM_Parent__c WHERE (Name__c LIKE '%-foo%' OR Name__c LIKE '%-bar%') AND Role__r.Name IN ('admin','user') AND ((NOT Name__c LIKE '%-far%') AND (NOT Name__c LIKE '%-baz%')) AND Some_Parent__r.Some_Type__c = 'typeA' AND Status__c != 'InActive' AND Value__c != null
```

You can find detailed usage in `marshaller_test.go`.

#### Subqueries

This package supports nested conditions within `WHERE` clauses as well. For example:

``` go

type contact struct {
	Name  string `soql:"selectColumn,fieldName=Name" json:"Name"`
	Email string `soql:"selectColumn,fieldName=Email" json:"Email"`
	Phone string `soql:"selectColumn,fieldName=Phone" json:"Phone"`
}

type soqlQuery struct {
	SelectClause contact       `soql:"selectClause,tableName=Contact"`
	WhereClause  queryCriteria `soql:"whereClause"`
}

type queryCriteria struct {
	Position    positionCriteria    `soql:"subquery,joiner=OR"`
	Contactable contactableCriteria `soql:"subquery,joiner=OR"`
}

type positionCriteria struct {
	Title             string              `soql:"equalsOperator,fieldName=Title"`
	DepartmentManager deptManagerCriteria `soql:"subquery"`
}

type deptManagerCriteria struct {
	Department string   `soql:"equalsOperator,fieldName=Department"`
	Title      []string `soql:"likeOperator,fieldName=Title"`
}

type contactableCriteria struct {
	EmailOK emailCheck `soql:"subquery,joiner=and"`
	PhoneOK phoneCheck `soql:"subquery,joiner=and"`
}

type emailCheck struct {
	Email         bool `soql:"nullOperator,fieldName=Email"`
	EmailOptedOut bool `soql:"equalsOperator,fieldName=HasOptedOutOfEmail"`
}

type phoneCheck struct {
	Phone     bool `soql:"nullOperator,fieldName=Phone"`
	DoNotCall bool `soql:"equalsOperator,fieldName=DoNotCall"`
}
soqlStruct := soqlQuery{
    WhereClause: queryCriteria{
        Position: positionCriteria{
            Title: "Purchasing Manager",
            DepartmentManager: deptManagerCriteria{
                Department: "Accounting",
                Title:      []string{"Manager"},
            },
        },
        Contactable: contactableCriteria{
            EmailOK: emailCheck{
                Email:         false,
                EmailOptedOut: false,
            },
            PhoneOK: phoneCheck{
                Phone:     false,
                DoNotCall: false,
            },
        },
    },
}
query, err := soql.Marshal(soqlStruct)
if err != nil {
     fmt.Printf("Error in marshalling: %s\n", err.Error())
}
fmt.Println(soqlQuery)

```

The above code will generate this SOQL query:

``` sql
SELECT Name,Email,Phone
FROM Contact
WHERE (Title = 'Purchasing Manager' OR (Department = 'Accounting' AND Title LIKE '%Manager%')) AND ((Email != null AND HasOptedOutOfEmail = false) OR (Phone != null AND DoNotCall = false))
```

#### Advantages

Intended users of this package are developers writing clients to interact with Salesforce. They can now define golang structs, annotate them and generate SOQL queries to be passed to Salesforce API. Great thing about this is that the json structure of returned response matches with selectClause, so you can just unmarshal response into the golang struct that was annotated with `selectClause` and now you have your query response directly available in golang struct.

## Tags explained

This section explains each of the supported tags in detail

### Top level tags

This section explains top level tags used in constructing SOQL query. Following snippet will be used as example for explaining these tags:

```
type TestSoqlStruct struct {
    SelectClause NonNestedStruct   `soql:"selectClause,tableName=SM_SomeObject__c"`
    WhereClause  TestQueryCriteria `soql:"whereClause"`
}
type TestQueryCriteria struct {
    IncludeNamePattern          []string `soql:"likeOperator,fieldName=Name__c"`
    Roles                       []string `soql:"inOperator,fieldName=Role__c"`
}
type NonNestedStruct struct {
    Name          string `soql:"selectColumn,fieldName=Name__c"`
    SomeValue     string `soql:"selectColumn,fieldName=SomeValue__c"`
}
```

1. `selectClause`: This tag is used on the struct which should be considered for generating part of SOQL query that contains columns/fields that should be selected. It should be used only on `struct` type. If used on types other than `struct` then `ErrInvalidTag` error will be returned. This tag is associated with `tableName` parameter. It specifies the name of the table (Salesforce object) from which the columns should be selected. If not specified name of the field is used as table name (Salesforce object). In the snippet above `SelectClause` member of `TestSoqlStruct` is tagged with `selectClause` to indicate that members in `NonNestedStruct` should be considered as fields to be selected from Salesforce object `SM_SomeObject__c`.

1. `whereClause`: This tag is used on the struct which encapsulates the query criteria for SOQL query. There is an optional parameter `joiner` for this tag. In the snippet above `WhereClause` member of `TestSoqlStruct` is tagged with `whereClause` to indicate that members in `TestQueryCriteria` should be considered for generating `WHERE` clause in SOQL query. If there are more than one field in `TestQueryCriteria` struct then they will be combined using `AND` logical operator. If the `joiner` parameter is set to `or` (case insensitive), then the fields will be combined using `OR` logical operator.  If the `joiner` parameter is set to `and` (case insensitive) or is not set, then the fields will be combined using `AND` logical operator. If any other value is provided, then `ErrInvalidTag` error will be returned. The `joiner` parameter is only supported when using `Marshal`; when calling `MarshalWhereClause`, the fields will always be combined with the `AND` logical operator.

1. `orderByClause`: This tag is used on the slice of `Order` to capture the ordering of columns and sort order. There are no parameters for this tag. Clients using this library can expose `Order` struct from this library to their users if they wish to allow users of the client to control ordering of the result.

1. `limitClause`: This tag is used on the \*int that describes the limit value for SOQL query. There are no parameters for this tag. Passing `nil` here will omit the `LIMIT` clause from the generated query. Passing a pointer to an integer value less than zero will cause an error.

1. `offsetClause`: This tag is used on the \*int that describes the offset value for SOQL query. There are no parameters for this tag. Passing `nil` here will omit the `OFFSET` clause from the generated query. Passing a pointer to an integer value less than zero will cause an error.

### Second level tags

This section explains the tags that should be used on members of struct tagged with `selectClause` and `whereClause`. These tags indicate how the members of the struct should be used in generating `SELECT` and `WHERE` clause.

#### Tags to be used on selectClause structs

This section explains the list of tags that can be used on members tagged with `selectClause`. Following snippet will be used explaining these tags:

```
type ParentStruct struct {
	ID                string          `soql:"selectColumn,fieldName=Id"`
	Name              string          `soql:"selectColumn,fieldName=Name__c"`
	NonNestedStruct   NonNestedStruct `soql:"selectColumn,fieldName=NonNestedStruct__r"` // child to parent relationship
	ChildStruct       TestChildStruct `soql:"selectChild,fieldName=Child__r"`            // parent to child relationship
	SomeNonSoqlMember string          `json:"some_nonsoql_member"`
}

type NonNestedStruct struct {
	Name          string `soql:"selectColumn,fieldName=Name"`
	SomeValue     string `soql:"selectColumn,fieldName=SomeValue__c"`
	NonSoqlStruct NonSoqlStruct
}

type TestChildStruct struct {
	SelectClause ChildStruct        `soql:"selectClause,tableName=SM_Child__c"`
	WhereClause  ChildQueryCriteria `soql:"whereClause"`
}
```

1. `selectColumn`: Members that are tagged with this tag will be considered in generating select clause of SOQL query. This tag is associated with `fieldName` parameter. It specifies the name of the field of underlying Salesforce object. If not specified the name of the field is used as underlying Salesforce object field name. This tag can be used on primitive data types as well as user defined structs. If used on user defined structs like `NonNestedStruct` member in `ParentStruct` it will be treated as child to parent relationship and the value specified in `fieldName` parameter (or default value of name of the member itself) will be prefixed to the members of that struct (`NonNestedStruct` in case of our example above).
1. `selectChild`: This tag is used on members which should be modelled as parent to child relation. It should be used on `struct` type only. If used on any other type then `ErrInvalidTag` error will be returned. The member on which this tag is used should in turn consist of members tagged with `selectClause` and `whereClause`. Please refer to `ChildStruct` member of `ParentStruct`.

#### Tags to be used on whereClause structs

This section explains the list of tags that can be used on members tagged with `whereClause`. Following snippet will be used as example for explaining these tags:

```
type QueryCriteria struct {
	IncludeNamePattern []string `soql:"likeOperator,fieldName=Name__c"`
	ExcludeNamePattern []string `soql:"notLikeOperator,fieldName=Name__c"`
	Roles              []string `soql:"inOperator,fieldName=Role__r.Name"`
	SomeType           string   `soql:"equalsOperator,fieldName=Some_Type__c"`
	SomeBoolType       *bool    `soql:"equalsOperator,fieldName=Some_Bool_Type__c"`
	Status             string   `soql:"notEqualsOperator,fieldName=Status__c"`
	AllowNullValue     *bool    `soql:"nullOperator,fieldName=Value__c"`
	NumOfCPUCores      int      `soql:"greaterThanOperator,fieldName=Num_of_CPU_Cores__c"`
	PhysicalCPUCount   uint8    `soql:"greaterThanOrEqualsToOperator,fieldName=Physical_CPU_Count__c"`
	AllocationLatency  float64  `soql:"lessThanOperator,fieldName=Allocation_Latency__c"`
	PvtTestFailCount   int64    `soql:"lessThanOrEqualsToOperator,fieldName=Pvt_Test_Fail_Count__c"`
    Subquery sub `soql:"subquery"`
}

type sub struct {
    NumOfCPUCores      int      `soql:"lessThanOperator,fieldName=Num_of_CPU_Cores__c"`
	PhysicalCPUCount   uint8    `soql:"lessThanOrEqualsToOperator,fieldName=Physical_CPU_Count__c"`
}
```

1. `likeOperator`: This tag is used on members which should be considered to construct field expressions in where clause using `LIKE` comparison operator. This tag should be used on member of type `[]string`. Used on any other type, `ErrInvalidTag` error will be returned. If there are more than one item in the slice then they will be combined using `OR` logical operator. Example will clarify this more:

   ```
   whereClause, _ := MarshalWhereClause(QueryCriteria{
       IncludeNamePattern: []string{"-foo", "-bar"},
   })
   // whereClause will be: WHERE (Name__c LIKE '%-foo%' OR Name__c LIKE '%-bar%')
   ```

1. `notLikeOperator`: This tag is used on members which should be considered to construct field expressions in where clause using `NOT LIKE` comparison operator. This tag should be used on member of type `[]string`. Used on any other type, `ErrInvalidTag` error will be returned. If there are more than one item in the slice then they will be combined using `AND` logical operator. Example will clarify this more:

   ```
   whereClause, _ := MarshalWhereClause(QueryCriteria{
       ExcludeNamePattern: []string{"-far", "-baz"},
   })
   // whereClause will be: WHERE ((NOT Name__c LIKE '%-far%') AND (NOT Name__c LIKE '%-baz%'))
   ```

1. `inOperator`: This tag is used on members which should be considered to construct field expressions in where clause using `IN` comparison operator. This tag should be used on member of type `[]string`, `[]int`, `[]int8`, `[]int16`, `[]int32`, `[]int64`, `[]uint`, `[]uint8`, `[]uint16`, `[]uint32`, `[]uint64`, `[]float32`, `[]float64`, `[]bool` or `[]time.Time`. Used on any other type, `ErrInvalidTag` error will be returned. Example will clarify this more:

   ```
   whereClause, _ := MarshalWhereClause(QueryCriteria{
       Roles: []string{"admin", "user"},
   })
   // whereClause will be: WHERE Role__r.Name IN ('admin','user')
   ```

1. `equalsOperator`: This tag is used on members which should be considered to construct field expressions in where clause using `=` comparison operator. This tag should be used on member of type `string`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `bool`, `*int`, `*int8`, `*int16`, `*int32`, `*int64`, `*uint`, `*uint8`, `*uint16`, `*uint32`, `*uint64`, `*float32`, `*float64`, `*bool` or `time.Time`. Used on any other type, `ErrInvalidTag` error will be returned. Example will clarify this more:

   ```
   whereClause, _ := MarshalWhereClause(QueryCriteria{
       SomeType: "SomeValue",
   })
   // whereClause will be: WHERE Some_Type__c = 'SomeValue'
   ```

   If pointers are used as data type then the field will be included in WHERE clause only if the variable is initialized. So in case below it will be included because the variable is initialized unlike example above.

   ```
   b := true
   whereClause, _ := MarshalWhereClause(QueryCriteria{
       SomeType:     "SomeValue",
       SomeBoolType: &b,
   })
   // whereClause will be: WHERE Some_Type__c = 'SomeValue' AND Some_Bool_Type__c = true
   ```

1. `notEqualsOperator`: This tag is used on members which should be considered to construct field expressions in where clause using `!=` comparison operator. This tag should be used on member of type `string`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `bool`, `*int`, `*int8`, `*int16`, `*int32`, `*int64`, `*uint`, `*uint8`, `*uint16`, `*uint32`, `*uint64`, `*float32`, `*float64`, `*bool` or `time.Time`. Used on any other type, `ErrInvalidTag` error will be returned. Example will clarify this more:

   ```
   whereClause, _ := MarshalWhereClause(QueryCriteria{
       Status: "DOWN",
   })
   // whereClause will be: WHERE Status__c != 'DOWN'
   ```

   Fields that are pointers will only be included if they are initialized else they will be skipped from WHERE clause.

1. `nullOperator`: This tag is used on members which should be considered to construct field expressions in where clause using `= null` or `!= null` comparison operator. This tag should be used on member of type `*bool` or `bool`. Used on any other type, `ErrInvalidTag` error will be returned. Recommended to use `*bool` as `bool` will always be initialized by golang to `false` and will result in `!= null` check even if not intended. Example will clarify this more:

   ```
   allowNull := true
   whereClause, _ := MarshalWhereClause(QueryCriteria{
       AllowNullValue: &allowNull,
   })
   // whereClause will be: WHERE Value__c = null
   allowNull = false
   whereClause, _ := MarshalWhereClause(QueryCriteria{
       AllowNullValue: &allowNull,
   })
   // whereClause will be: WHERE Value__c != null
   ```

1. `greaterThanOperator`: This tag is used on members which should be considered to construct field expressions in where clause using `>` comparison operator. This tag should be used on member of type `string`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `bool`, `*int`, `*int8`, `*int16`, `*int32`, `*int64`, `*uint`, `*uint8`, `*uint16`, `*uint32`, `*uint64`, `*float32`, `*float64`, `*bool` or `time.Time`. Used on any other type, `ErrInvalidTag` error will be returned. Example will clarify this more:

   ```
   whereClause, _ := MarshalWhereClause(QueryCriteria{
       NumOfCPUCores: 8,
   })
   // whereClause will be: WHERE Num_of_CPU_Cores__c > 8
   ```

   Fields that are pointers will only be included if they are initialized else they will be skipped from WHERE clause.

1. `greaterThanOrEqualsToOperator`: This tag is used on members which should be considered to construct field expressions in where clause using `>=` comparison operator. This tag should be used on member of type `string`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `bool`, `*int`, `*int8`, `*int16`, `*int32`, `*int64`, `*uint`, `*uint8`, `*uint16`, `*uint32`, `*uint64`, `*float32`, `*float64`, `*bool` or `time.Time`. Used on any other type, `ErrInvalidTag` error will be returned. Example will clarify this more:

   ```
   whereClause, _ := MarshalWhereClause(QueryCriteria{
       PhysicalCPUCount: 4,
   })
   // whereClause will be: WHERE Physical_CPU_Count__c >= 4
   ```

   Fields that are pointers will only be included if they are initialized else they will be skipped from WHERE clause.

1. `lessThanOperator`: This tag is used on members which should be considered to construct field expressions in where clause using `<` comparison operator. This tag should be used on member of type `string`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `bool`, `*int`, `*int8`, `*int16`, `*int32`, `*int64`, `*uint`, `*uint8`, `*uint16`, `*uint32`, `*uint64`, `*float32`, `*float64`, `*bool` or `time.Time`. Used on any other type, `ErrInvalidTag` error will be returned. Example will clarify this more:

   ```
   whereClause, _ := MarshalWhereClause(QueryCriteria{
       AllocationLatency: 28.9,
   })
   // whereClause will be: WHERE Allocation_Latency__c < 28.9
   ```

   Fields that are pointers will only be included if they are initialized else they will be skipped from WHERE clause.

1. `lessThanOrEqualsToOperator`: This tag is used on members which should be considered to construct field expressions in where clause using `<=` comparison operator. This tag should be used on member of type `string`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `bool`, `*int`, `*int8`, `*int16`, `*int32`, `*int64`, `*uint`, `*uint8`, `*uint16`, `*uint32`, `*uint64`, `*float32`, `*float64`, `*bool` or `time.Time`. Used on any other type, `ErrInvalidTag` error will be returned. Example will clarify this more:

   ```
   whereClause, _ := MarshalWhereClause(QueryCriteria{
       PvtTestFailCount: 32,
   })
   // whereClause will be: WHERE Pvt_Test_Fail_Count__c <= 32
   ```

   Fields that are pointers will only be included if they are initialized else they will be skipped from WHERE clause.

If there are more than one fields in the struct tagged with `whereClause` then they will be combined using `AND` logical operator. This has been demonstrated in the code snippets in [Advanced usage](#advanced-usage).

1. `subquery`: This tag is used on members which should be used to construct related sets of conditions wrapped in `()` in the query. This tag should only be used on members of type `struct`. Used on any other type, `ErrInvalidTag` error will be returned. Any of the above property tags (including `subquery`) may be used in the designated `struct`.

#### The Order struct and orderByClause

This section explains the Order struct to be used for the `orderByClause`.

```
type Order struct {
	Field string
	IsDesc bool
}
```

The `Order` struct is part of this library and has two fields, `Field` which is of type string and `IsDesc` of type bool. Each struct represents a column from the select column list that should be included in the `ORDER BY` clause, as well as the sort order on that column. The value of `Field` should be the name of the struct field with the `selectColumn` tag, and set `IsDesc` to `true` to specify the sort order on that column as `DESC` (set it to `false` for `ASC`). Create a slice of `Order` structs and tag it with the `orderByClause` soql tag to define the `ORDER BY` clause. Using the following `NestedStruct` as an example of the `selectClause`:

```
type NonNestedStruct struct {
	Name          string `soql:"selectColumn,fieldName=Name"`
	SomeValue     string `soql:"selectColumn,fieldName=SomeValue__c"`
	NonSoqlStruct NonSoqlStruct
}

type NestedStruct struct {
	ID              string          `soql:"selectColumn,fieldName=Id"`
	Name            string          `soql:"selectColumn,fieldName=Name__c"`
	NonNestedStruct NonNestedStruct `soql:"selectColumn,fieldName=NonNestedStruct__r"`
}
```

To order the query results using the `Name__c` field of the NestedStruct in `ASC` order and `SomeValue__c` of the `NonNestedStruct` field in `DESC` order, the following `Order` slice should be used:

```
order := []Order{Order{Field:"Name",IsDesc:false},Order{Field:"NonNestedStruct.SomeValue",IsDesc:true}}
```

To specify fields in nested structs, use the `<parent>.<field>` dot notation.

The final soql query struct would look like:

```
type TestSoqlStruct struct {
    SelectClause  NestedStruct   `soql:"selectClause,tableName=SM_SomeObject__c"`
    OrderByClause []Order        `soql:"orderByClause"`
}
```

## License

go-soql is BSD3 licensed. Here is the link to license [file](./LICENSE.txt)

## Contributing

You are welcome to contribute to this repo. Please create PR and send it for review. Please follow code of conduct as documented [here](./CODE_OF_CONDUCT.md)

If you have a question, comment, bug report, feature request, etc. please open a GitHub issue.
