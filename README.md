# SOQL

This package supports marshalling a golang struct into SOQL. Like `json` tags, this package provides `soql` tags that you can use to annotate your golang structs. Once tagged `Marshal` method will return SOQL query that will let you query the required Salesforce object using Salesforce API.

## Introduction

Please refer to [introduction](./docs/introduction.md) to understand the basics. Once you read through it you can refer to documentation below that covers features of this repo in more depth.

## How to use

Start with using `soql` tags on members of your golang structs. `soql` is the main tag. There are following subtags supported:

```
    selectClause // is the tag to be used when marking the struct to be considered for select clause in soql.
    whereClause // is the tag to be used when marking the struct to be considered for where clause in soql.
    selectColumn // is the tag to be used for selecting a column in select clause. It should be used on members of struct that have been tagged with selectClause.
    selectChild // is the tag to be used when selecting from child tables. It should be used on members of struct that have been tagged with selectClause.
    likeOperator // is the tag to be used for "like" operator in where clause. It should be used on members of struct that have been tagged with whereClause.
    notLikeOperator // is the tag to be used for "not like" operator in where clause. It should be used on members of struct that have been tagged with whereClause.
    inOperator // is the tag to be used for "in" operator in where clause. It should be used on members of struct that have been tagged with whereClause.
    equalsOperator // is the tag to be used for "=" operator in where clause. It should be used on members of struct that have been tagged with whereClause.
    notEqualsOperator // is the tag to be used for "!=" operator in where clause. It should be used on members of struct that have been tagged with whereClause.
    nullOperator // is the tag to be used for " = null " or "!= null" operator in where clause. It should be used on members of struct that have been tagged with whereClause.
```

Following are supported parameters:

```
    fieldName // is the parameter to be used to specify the name of the field in underlying Salesforce object. It can be used with all tags listed above other than selectClause and whereClause.
    tableName // is the parameter to be used to specify the name of the table of underlying Salesforce Object. It can be be used only with selectClause.

```

If `fieldName` and `tableName` parameters are not provided then the name of the field will be used as default.

Lets take a look at one example of a simple non-nested struct and how it can be used to construct a soql query:

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

To use above structs to create SOQL query

```
soqlStruct := TestSoqlStruct{
    WhereClause: TestQueryCriteria {
        IncludeNamePattern: []string{"foo", "bar"},
        Roles: []string{"admin", "user"},
    }
}
soqlQuery, err := Marshal(soqlStruct)
if err != nil {
    fmt.Printf("Error in marshaling: %s\n", err.Error())
}
fmt.Println(soqlQuery)
```

Above struct will result in following SOQL query:

```
SELECT Name,SomeValue__c FROM SM_SomeObject__C WHERE (Name__c LIKE '%foo%' OR Name__c LIKE '%bar%') AND Role__c IN ('admin','user')
```

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

Intended users of this package are developers writing clients to interact with Salesforce. They can now define golang structs, annotate them and generate SOQL queries to be passed to Salesforce API. Great thing about this is that the json structure of returned response matches with selectClause, so you can just unmarshal response into the golang struct that was annotated with `selectClause` and now you have your query response directly available in golang struct.

## License

go-soql is BSD3 licensed. Here is the link to license [file](./LICENSE.txt)

## Contributing

You are welcome to contribute to this repo. Please create PR and send it for review. Please follow code of conduct as documented [here](./CODE_OF_CONDUCT.md)

If you have a question, comment, bug report, feature request, etc. please open a GitHub issue.
