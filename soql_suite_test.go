/*
 * Copyright (c) 2018, salesforce.com, inc.
 * All rights reserved.
 * SPDX-License-Identifier: BSD-3-Clause
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/BSD-3-Clause
 */
package soql_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/forcedotcom/go-soql"
)

func TestSoql(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Soql Suite")
}

type TestSoqlStruct struct {
	SelectClause NestedStruct      `soql:"selectClause,tableName=SM_Logical_Host__c"`
	WhereClause  TestQueryCriteria `soql:"whereClause"`
}

type TestSoqlMixedDataAndOperatorStruct struct {
	SelectClause NestedStruct                                `soql:"selectClause,tableName=SM_Logical_Host__c"`
	WhereClause  QueryCriteriaWithMixedDataTypesAndOperators `soql:"whereClause"`
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

type TestChildWithOrderByStruct struct {
	SelectClause  ChildStruct `soql:"selectClause,tableName=SM_Application_Versions__c"`
	OrderByClause []Order     `soql:"orderByClause"`
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

type OrderByParentStruct struct {
	ID          string                     `soql:"selectColumn,fieldName=Id"`
	Name        string                     `soql:"selectColumn,fieldName=Name__c"`
	ChildStruct TestChildWithOrderByStruct `soql:"selectChild,fieldName=Application_Versions__r"`
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

type InvalidSelectChildClause struct {
	ID          string `soql:"selectColumn,fieldName=Id"`
	Name        string `soql:"selectColumn,fieldName=Name__c"`
	ChildStruct int    `soql:"selectChild,fieldName=Application_Versions__r"`
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

type MultipleOrderByClause struct {
	OrderByClause1 []Order `soql:"orderByClause"`
	OrderByClause2 []Order `soql:"orderByClause"`
}

type OnlyWhereClause struct {
	WhereClause TestQueryCriteria `soql:"whereClause"`
}

type OnlyOrderByClause struct {
	OrderByClause []Order `soql:"orderByClause"`
}

type EmptyStruct struct {
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

type QueryCriteriaWithIntegerTypes struct {
	NumOfCPUCores                    int   `soql:"equalsOperator,fieldName=Num_of_CPU_Cores__c"`
	PhysicalCPUCount                 int8  `soql:"equalsOperator,fieldName=Physical_CPU_Count__c"`
	NumOfSuccessivePuppetRunFailures int16 `soql:"equalsOperator,fieldName=Number_Of_Successive_Puppet_Run_Failures__c"`
	NumOfCoolanLogFiles              int32 `soql:"equalsOperator,fieldName=Num_Of_Coolan_Log_Files__c"`
	PvtTestFailCount                 int64 `soql:"equalsOperator,fieldName=Pvt_Test_Fail_Count__c"`
}

type QueryCriteriaWithUnsignedIntegerTypes struct {
	NumOfCPUCores                    uint   `soql:"equalsOperator,fieldName=Num_of_CPU_Cores__c"`
	PhysicalCPUCount                 uint8  `soql:"equalsOperator,fieldName=Physical_CPU_Count__c"`
	NumOfSuccessivePuppetRunFailures uint16 `soql:"equalsOperator,fieldName=Number_Of_Successive_Puppet_Run_Failures__c"`
	NumOfCoolanLogFiles              uint32 `soql:"equalsOperator,fieldName=Num_Of_Coolan_Log_Files__c"`
	PvtTestFailCount                 uint64 `soql:"equalsOperator,fieldName=Pvt_Test_Fail_Count__c"`
}

type QueryCriteriaWithFloatTypes struct {
	NumOfCPUCores    float32 `soql:"equalsOperator,fieldName=Num_of_CPU_Cores__c"`
	PhysicalCPUCount float64 `soql:"equalsOperator,fieldName=Physical_CPU_Count__c"`
}

type QueryCriteriaWithBooleanType struct {
	NUMAEnabled   bool `soql:"equalsOperator,fieldName=NUMA_Enabled__c"`
	DisableAlerts bool `soql:"equalsOperator,fieldName=Disable_Alerts__c"`
}

type QueryCriteriaWithDateTimeType struct {
	CreatedDate time.Time `soql:"equalsOperator,fieldName=CreatedDate"`
}

type QueryCriteriaNumericComparisonOperators struct {
	NumOfCPUCores                    int `soql:"greaterThanOperator,fieldName=Num_of_CPU_Cores__c"`
	PhysicalCPUCount                 int `soql:"lessThanOperator,fieldName=Physical_CPU_Count__c"`
	NumOfSuccessivePuppetRunFailures int `soql:"greaterThanOrEqualsToOperator,fieldName=Number_Of_Successive_Puppet_Run_Failures__c"`
	NumOfCoolanLogFiles              int `soql:"lessThanOrEqualsToOperator,fieldName=Num_Of_Coolan_Log_Files__c"`
}

type QueryCriteriaWithMixedDataTypesAndOperators struct {
	BIOSType                         string    `soql:"equalsOperator,fieldName=BIOS_Type__c"`
	NumOfCPUCores                    int       `soql:"greaterThanOperator,fieldName=Num_of_CPU_Cores__c"`
	NUMAEnabled                      bool      `soql:"equalsOperator,fieldName=NUMA_Enabled__c"`
	PvtTestFailCount                 int64     `soql:"lessThanOrEqualsToOperator,fieldName=Pvt_Test_Fail_Count__c"`
	PhysicalCPUCount                 uint8     `soql:"greaterThanOrEqualsToOperator,fieldName=Physical_CPU_Count__c"`
	CreatedDate                      time.Time `soql:"equalsOperator,fieldName=CreatedDate"`
	DisableAlerts                    bool      `soql:"equalsOperator,fieldName=Disable_Alerts__c"`
	AllocationLatency                float64   `soql:"lessThanOperator,fieldName=Allocation_Latency__c"`
	MajorOSVersion                   string    `soql:"equalsOperator,fieldName=Major_OS_Version__c"`
	NumOfSuccessivePuppetRunFailures uint32    `soql:"equalsOperator,fieldName=Number_Of_Successive_Puppet_Run_Failures__c"`
	LastRestart                      time.Time `soql:"greaterThanOperator,fieldName=Last_Restart__c"`
}

type InvalidSelectClause struct {
	SelectClause string `soql:"selectClause,tableName=SM_Logical_Host__c"`
}

type TestSoqlOrderByStruct struct {
	SelectClause  NestedStruct `soql:"selectClause,tableName=SM_Logical_Host__c"`
	OrderByClause []Order      `soql:"orderByClause"`
}

type TestSoqlChildRelationOrderByStruct struct {
	SelectClause  OrderByParentStruct `soql:"selectClause,tableName=SM_Logical_Host__c"`
	OrderByClause []Order             `soql:"orderByClause"`
}

type TestSoqlLimitStruct struct {
	SelectClause NestedStruct      `soql:"selectClause,tableName=SM_Logical_Host__c"`
	WhereClause  TestQueryCriteria `soql:"whereClause"`
	Limit        int               `soql:"limitClause"`
}

type TestSoqlInvalidLimitStruct struct {
	SelectClause NestedStruct      `soql:"selectClause,tableName=SM_Logical_Host__c"`
	WhereClause  TestQueryCriteria `soql:"whereClause"`
	Limit        string            `soql:"limitClause"`
}

type TestSoqlMultipleLimitStruct struct {
	SelectClause NestedStruct      `soql:"selectClause,tableName=SM_Logical_Host__c"`
	WhereClause  TestQueryCriteria `soql:"whereClause"`
	Limit        int               `soql:"limitClause"`
	AlsoLimit    int               `soql:"limitClause"`
}

type ParentLimitStruct struct {
	ID          string           `soql:"selectColumn,fieldName=Id"`
	Name        string           `soql:"selectColumn,fieldName=Name__c"`
	ChildStruct ChildLimitStruct `soql:"selectChild,fieldName=Application_Versions__r"`
}

type ChildLimitStruct struct {
	SelectClause TestChildLimitSelect `soql:"selectClause,tableName=Application_Versions__c"`
	Limit        int                  `soql:"limitClause"`
}

type TestChildLimitSelect struct {
	ID      string `soql:"selectColumn,fieldName=Id"`
	Version string `soql:"selectColumn,fieldName=Version__c"`
}

type TestSoqlChildRelationLimitStruct struct {
	SelectClause ParentLimitStruct `soql:"selectClause,tableName=SM_Logical_Host__c"`
	Limit        int               `soql:"limitClause"`
}

type TestSoqlOffsetStruct struct {
	SelectClause NestedStruct      `soql:"selectClause,tableName=SM_Logical_Host__c"`
	WhereClause  TestQueryCriteria `soql:"whereClause"`
	Offset       int               `soql:"offsetClause"`
}

type TestSoqlInvalidOffsetStruct struct {
	SelectClause NestedStruct      `soql:"selectClause,tableName=SM_Logical_Host__c"`
	WhereClause  TestQueryCriteria `soql:"whereClause"`
	Offset       string            `soql:"offsetClause"`
}

type TestSoqlMultipleOffsetStruct struct {
	SelectClause NestedStruct      `soql:"selectClause,tableName=SM_Logical_Host__c"`
	WhereClause  TestQueryCriteria `soql:"whereClause"`
	Offset       int               `soql:"offsetClause"`
	AlsoOffset   int               `soql:"offsetClause"`
}
