/*
 * Copyright (c) 2018, salesforce.com, inc.
 * All rights reserved.
 * SPDX-License-Identifier: BSD-3-Clause
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/BSD-3-Clause
 */
package soql_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSoql(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Soql Suite")
}

type TestSoqlStruct struct {
	SelectClause NestedStruct      `soql:"selectClause,tableName=SM_Logical_Host__c"`
	WhereClause  TestQueryCriteria `soql:"whereClause"`
}

type TestQueryCriteria struct {
	IncludeNamePattern          []string `soql:"likeOperator,fieldName=Host_Name__c"`
	Roles                       []string `soql:"inOperator,fieldName=Role__r.Name"`
	ExcludeNamePattern          []string `soql:"notLikeOperator,fieldName=Host_Name__c"`
	AssetType                   string   `soql:"equalsOperator,fieldName=Tech_Asset__r.Asset_Type_Asset_Type__c"`
	Status                      string   `soql:"notEqualsOperator,fieldName=Status__c"`
	AllowNullLastDiscoveredDate *bool    `soql:"nullOperator,fieldName=Last_Discovered_Date__c"`
}

type NonSoqlStruct struct {
	Key   string
	Value string
}

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

type TestChildStruct struct {
	SelectClause ChildStruct        `soql:"selectClause,tableName=SM_Application_Versions__c"`
	WhereClause  ChildQueryCriteria `soql:"whereClause"`
}

type ChildStruct struct {
	Version string `soql:"selectColumn,fieldName=Version__c"`
}

type ChildQueryCriteria struct {
	Name string `soql:"equalsOperator,fieldName=Name__c"`
}

type ParentStruct struct {
	ID                string          `soql:"selectColumn,fieldName=Id"`
	Name              string          `soql:"selectColumn,fieldName=Name__c"`
	NonNestedStruct   NonNestedStruct `soql:"selectColumn,fieldName=NonNestedStruct__r"`
	ChildStruct       TestChildStruct `soql:"selectChild,fieldName=Application_Versions__r"`
	SomeNonSoqlMember string          `json:"some_nonsoql_member"`
}

type DefaultFieldNameParentStruct struct {
	ID                string          `soql:"selectColumn,fieldName=Id"`
	Name              string          `soql:"selectColumn,fieldName=Name__c"`
	NonNestedStruct   NonNestedStruct `soql:"selectColumn,fieldName=NonNestedStruct__r"`
	ChildStruct       TestChildStruct `soql:"selectChild"`
	SomeNonSoqlMember string          `json:"some_nonsoql_member"`
}

type InvalidTestChildStruct struct {
	WhereClause ChildQueryCriteria `soql:"whereClause"`
}
type InvalidParentStruct struct {
	ID          string                 `soql:"selectColumn,fieldName=Id"`
	Name        string                 `soql:"selectColumn,fieldName=Name__c"`
	ChildStruct InvalidTestChildStruct `soql:"selectChild,fieldName=Application_Versions__r"`
}

type ChildTagToNonStruct struct {
	ID          string `soql:"selectColumn,fieldName=Id"`
	Name        string `soql:"selectColumn,fieldName=Name__c"`
	ChildStruct string `soql:"selectChild,fieldName=Application_Versions__r"`
}

type MultipleSelectClause struct {
	SelectClause NestedStruct `soql:"selectClause,tableName=SM_Logical_Host__c"`
	ParentStruct ParentStruct `soql:"selectClause,tableName=SM_Table__c"`
}

type MultipleWhereClause struct {
	WhereClause1 ChildQueryCriteria `soql:"whereClause"`
	WhereClause2 ChildQueryCriteria `soql:"whereClause"`
}

type OnlyWhereClause struct {
	WhereClause TestQueryCriteria `soql:"whereClause"`
}

type EmptyStruct struct {
}

type QueryCriteriaWithInvalidTypes struct {
	IncludeNamePattern          []bool `soql:"likeOperator,fieldName=Host_Name__c"`
	Roles                       []int  `soql:"inOperator,fieldName=Role__r.Name"`
	ExcludeNamePattern          []bool `soql:"notLikeOperator,fieldName=Host_Name__c"`
	AssetType                   int    `soql:"equalsOperator,fieldName=Tech_Asset__r.Asset_Type_Asset_Type__c"`
	AllowNullLastDiscoveredDate string `soql:"nullOperator,fieldName=Last_Discovered_Date__c"`
}

type InvalidTagInStruct struct {
	SelectClause  NestedStruct       `soql:"selectClause,tableName=SM_Logical_Host__c"`
	WhereClause   ChildQueryCriteria `soql:"whereClause"`
	AnotherMember NestedStruct       `soql:"invalidClause,tableName=SM_Logical_Host__c"`
}

type DefaultFieldNameStruct struct {
	DefaultName string `soql:"selectColumn"`
	Description string `soql:"selectColumn,fieldName=Description__c"`
}

type DefaultTableNameStruct struct {
	SomeTableName NestedStruct      `soql:"selectClause"`
	WhereClause   TestQueryCriteria `soql:"whereClause"`
}

type DefaultFieldNameQueryCriteria struct {
	IncludeNamePattern []string `soql:"likeOperator,fieldName=Host_Name__c"`
	Role               []string `soql:"inOperator"`
}
