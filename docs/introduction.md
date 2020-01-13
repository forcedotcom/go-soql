# Introduction

If you are a Golang developer intending to write client for interacting with Salesforce, this post is for you. Learn how you can annotate your Golang structs to generate SOQL queries to be used in Salesforce APIs and how this will make it easy for you to write Golang client for Salesforce.

Salesforce has REST [API](https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/intro_rest_resources.htm) that allows any third party to integrate with it using their language of choice. One of the APIs (/services/data/<version>/query) allows developers to query Salesforce object using [SOQL](https://developer.salesforce.com/docs/atlas.en-us.soql_sosl.meta/soql_sosl/sforce_api_calls_soql_sosl_intro.htm) query. SOQL is very powerful, SQL like query language that allows developers to query Salesforce objects. It has extensive support for conditional expression to filter the objects being retrieved as well as very powerful way of performing join operations using relationship queries.

Golang is gaining popularity in recent years and is now very widely adopted languages in enterprises. As more and more developers embrace Golang it is important to provide easier means and ways to interact with Salesforce APIs. One such step is to provide soql annotation library that will allow Golang developers to tag their structs very similar to how they tag their structs for JSON marshaling/unmarshaling. And that is the aim of new [go-soql](https://github.com/forcedotcom/go-soql) library. It allows Golang developers to tag their structs and then marshal it into SOQL query that can be used in the query API of Salesforce. Lets take a look at how it works.

You will start by defining the Golang struct representing the Salesforce object(s) you want to query. Lets say the Salesforce object is Account and it has following fields (subset of fields shown here from the actual [list](https://developer.salesforce.com/docs/atlas.en-us.api.meta/api/sforce_api_objects_account.htm)):

| FieldName          | Type     |
| ------------------ | -------- |
| Name               | string   |
| AccountNumber      | picklist |
| AccountSource      | string   |
| BillingAddress     | string   |
| HasOptedOutOfEmail | boolean  |
| LastActivityDate   | Date     |
| NumberOfEmployees  | int      |

Golang struct representation for it looks as follows:

```
type Account struct {
    Name               string
    AccountNumber      string
    AccountSource      string
    BillingAddress     string
    HasOptedOutOfEmail bool
    LastActivityDate   time.Time
    NumberOfEmployees  int
}
```

Now if you want to query all Account where AccountSource is one of ‘Advertisement’ or Data.com’ then the SOQL query will look as follows:

```

SELECT Name,AccountNumber,AccountSource,BillingAddress,HasOptedOutOfEmail,LastActivityDate,NumberOfEmployees FROM Account WHERE AccountSource IN ('Advertisement', 'Data.com')

```

Now wouldn’t it be great if instead of hardcoding this query we can generate it from our Account struct directly. Not only it will help us automatically change SOQL query based on addition/removal of attributes from our model but also it will help us directly unmarshal the response from Salesforce into the struct! Lets see how we can achieve it. We start by annotating our Account struct as follows:

```
type Account struct {
    Name               string    `soql:"selectColumn,fieldName=Name" json:"Name"`
    AccountNumber      string    `soql:"selectColumn,fieldName=AccountNumber" json:"AccountNumber"`
    AccountSource      string    `soql:"selectColumn,fieldName=AccountSource" json:"AccountSource"`
    BillingAddress     string    `soql:"selectColumn,fieldName=BillingAddress" json:"BillingAddress"`
    HasOptedOutOfEmail bool      `soql:"selectColumn,fieldName=HasOptedOutOfEmail" json:"HasOptedOutOfEmail"`
    LastActivityDate   time.Time `soql:"selectColumn,fieldName=LastActivityDate" json:"LastActivityDate"`
    NumberOfEmployees  int       `soql:"selectColumn,fieldName=NumberOfEmployees" json:"NumberOfEmployees"`
}
```

If we just want to generate the select clause without conditional expression then it can be done as follows:

```
soqlQuery, err := soql.MarshalSelectClause(Account{}, "")
```

And this will result in following

```

Name,AccountNumber,AccountSource,BillingAddress,HasOptedOutOfEmail,LastActivityDate,NumberOfEmployees

```

This is, of course, of limited use. We want full SOQL query to be automatically generated. This needs us to define few more structs to model this:

```
type AccountQueryCriteria struct {
    AccountSource []string `soql:"inOperator,fieldName=AccountSource"`
}

type AccountSoqlQuery struct {
    SelectClause  Account              `soql:"selectClause,tableName=Account"`
    WhereClause   AccountQueryCriteria `soql:"whereClause"`
}
```

Now we can generate complete SOQL query using the above two structs:

```
soqlStruct := AccountSoqlQuery{
    SelectClause: Account{},
    WhereClause: AccountQueryCriteria{
        AccountSource: []string{"Advertisement", "Data.com"},
    },
}
soqlQuery, err := soql.Marshal(soqlStruct)
```

And viola! This will generate SOQL query that we expect. Now we can use that in our call to /services/data/<version>/query API. Note that you can add/remove/change the AccountQueryCriteria struct to add/remove/change the conditional expression. If you add more than one field to AccountQueryCriteria then they are combined using AND logical operator. What is cool is that the returned JSON data from Salesforce can be directly unmarshalled into Account struct. Here’s the struct that represents response from query API

```
type QueryResponse struct {
    Done            bool      `json:"done"`
    NextRecordsURL  string    `json:"nextRecordsUrl"`
    Accounts        []Account `json:"records"`
    TotalSize       int       `json:"totalSize"`
}
```

Your client code will look something like this

```
values := url.Values{}
values.Set("q", soqlQuery) // soqlQuery variable defined in code snippet above
path := fmt.Sprintf("/services/data/v44.0/query?%s",values.Encode())
serverURL := "https://<your domain>.salesforce.com"
req, err := http.NewRequest(http.MethodGet, serverURL+path, nil)
if err != nil {
   // Handle error case
}
req.Header.Add("Authorization", "Bearer <session ID>")
req.Header.Add("Content-Type", "application/json")
httpClient := &http.Client{}
resp, err := httpClient.Do(req)
if err != nil {
   // Handle error case
}
payload, err := ioutil.ReadAll(resp.Body)
if err != nil {
   // Handle error case
}
var queryResponse QueryResponse
err = json.Unmarshal(payload, &queryResponse)
if err != nil {
   // Handle error case
}
```

As you can note above using this writing Golang client for Salesforce can now be much more easy affair with new soql tag library. All you need to do is define the struct and annotate it.

Hopefully now you have a feel of what go-soql library can do for you. However, this blog post just scratches the surface of what go-sosql library can do. It has very extensive support for logical operators as well as child to parent and parent to child relationships. More details on how to use that can be found out in README of the repo.
