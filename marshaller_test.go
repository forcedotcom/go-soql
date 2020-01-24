/*
 * Copyright (c) 2018, salesforce.com, inc.
 * All rights reserved.
 * SPDX-License-Identifier: BSD-3-Clause
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/BSD-3-Clause
 */
package soql_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/forcedotcom/go-soql"
)

var _ = Describe("Marshaller", func() {
	Describe("MarshalWhereClause", func() {
		var (
			clause         string
			expectedClause string
			err            error
		)
		Context("when non pointer value is passed as argument", func() {
			var (
				critetria TestQueryCriteria
			)

			JustBeforeEach(func() {
				clause, err = MarshalWhereClause(critetria)
				Expect(err).ToNot(HaveOccurred())
			})

			Context("when there are no fields populated", func() {
				It("returns empty where clause", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(clause).To(BeEmpty())
				})
			})

			Context("when only like clause pattern is populated", func() {
				Context("when there is only one item in the like clause array", func() {
					BeforeEach(func() {
						critetria = TestQueryCriteria{
							IncludeNamePattern: []string{"-db"},
						}
						expectedClause = "Host_Name__c LIKE '%-db%'"
					})

					It("returns where clause with only one condition", func() {
						Expect(clause).To(Equal(expectedClause))
					})
				})

				Context("when there is more than one item in the like clause array", func() {
					BeforeEach(func() {
						critetria = TestQueryCriteria{
							IncludeNamePattern: []string{"-db", "-dbmgmt", "-dgdb"},
						}
						expectedClause = "(Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%' OR Host_Name__c LIKE '%-dgdb%')"
					})

					It("returns where clause with OR condition", func() {
						Expect(clause).To(Equal(expectedClause))
					})
				})
			})

			Context("when only not like clause is populated", func() {
				Context("when there is only one item in the not like clause array", func() {
					BeforeEach(func() {
						critetria = TestQueryCriteria{
							ExcludeNamePattern: []string{"-db"},
						}
						expectedClause = "(NOT Host_Name__c LIKE '%-db%')"
					})

					It("returns where clause with only one condition", func() {
						Expect(clause).To(Equal(expectedClause))
					})
				})

				Context("when there is more than one item in the not like clause array", func() {
					BeforeEach(func() {
						critetria = TestQueryCriteria{
							ExcludeNamePattern: []string{"-db", "-dbmgmt", "-dgdb"},
						}
						expectedClause = "((NOT Host_Name__c LIKE '%-db%') AND (NOT Host_Name__c LIKE '%-dbmgmt%') AND (NOT Host_Name__c LIKE '%-dgdb%'))"
					})

					It("returns where clause with OR condition", func() {
						Expect(clause).To(Equal(expectedClause))
					})
				})
			})

			Context("when only equalsClause is populated", func() {
				BeforeEach(func() {
					critetria = TestQueryCriteria{
						AssetType: "SERVER",
					}
					expectedClause = "Tech_Asset__r.Asset_Type_Asset_Type__c = 'SERVER'"
				})

				It("returns appropriate where clause", func() {
					Expect(clause).To(Equal(expectedClause))
				})
			})

			Context("when only notEqualsClause is populated", func() {
				BeforeEach(func() {
					critetria = TestQueryCriteria{
						Status: "InActive",
					}
					expectedClause = "Status__c != 'InActive'"
				})

				It("returns appropriate where clause", func() {
					Expect(clause).To(Equal(expectedClause))
				})
			})

			Context("when only inClause is populated", func() {
				Context("when there is only one item in the inClause array", func() {
					BeforeEach(func() {
						critetria = TestQueryCriteria{
							Roles: []string{"db"},
						}
						expectedClause = "Role__r.Name IN ('db')"
					})
					It("returns where clause with only one item in IN clause", func() {
						Expect(clause).To(Equal(expectedClause))
					})
				})

				Context("when there is more than one item in the inClause array", func() {
					BeforeEach(func() {
						critetria = TestQueryCriteria{
							Roles: []string{"db", "dbmgmt"},
						}
						expectedClause = "Role__r.Name IN ('db','dbmgmt')"
					})
					It("returns where clause with all the items in IN clause", func() {
						Expect(clause).To(Equal(expectedClause))
					})
				})
			})

			Context("when only null clause is populated", func() {
				Context("when null is allowed", func() {
					BeforeEach(func() {
						allowNull := true
						critetria = TestQueryCriteria{
							AllowNullLastDiscoveredDate: &allowNull,
						}

						expectedClause = "Last_Discovered_Date__c = null"
					})

					It("returns appropriate where clause", func() {
						Expect(clause).To(Equal(expectedClause))
					})
				})

				Context("when null is not allowed", func() {
					BeforeEach(func() {
						allowNull := false
						critetria = TestQueryCriteria{
							AllowNullLastDiscoveredDate: &allowNull,
						}

						expectedClause = "Last_Discovered_Date__c != null"
					})

					It("returns appropriate where clause", func() {
						Expect(clause).To(Equal(expectedClause))
					})
				})
			})

			Context("when likeOperator and inClause are populated", func() {
				BeforeEach(func() {
					critetria = TestQueryCriteria{
						IncludeNamePattern: []string{"-db", "-dbmgmt"},
						Roles:              []string{"db"},
					}

					expectedClause = "(Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND Role__r.Name IN ('db')"
				})

				It("returns properly formed clause for name and role joined by AND clause", func() {
					Expect(clause).To(Equal(expectedClause))
				})
			})

			Context("when likeOperator and equalsClause are populated", func() {
				BeforeEach(func() {
					critetria = TestQueryCriteria{
						IncludeNamePattern: []string{"-db", "-dbmgmt"},
						AssetType:          "SERVER",
					}

					expectedClause = "(Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND Tech_Asset__r.Asset_Type_Asset_Type__c = 'SERVER'"
				})

				It("returns properly formed clause for likeOperator and inClause joined by AND clause", func() {
					Expect(clause).To(Equal(expectedClause))
				})
			})

			Context("when both likeOperator and notlikeOperator are populated", func() {
				BeforeEach(func() {
					critetria = TestQueryCriteria{
						IncludeNamePattern: []string{"-db", "-dbmgmt"},
						ExcludeNamePattern: []string{"-core", "-drp"},
					}

					expectedClause = "(Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND ((NOT Host_Name__c LIKE '%-core%') AND (NOT Host_Name__c LIKE '%-drp%'))"
				})

				It("returns properly formed clause for likeOperator and notlikeOperator joined by AND clause", func() {
					Expect(clause).To(Equal(expectedClause))
				})
			})

			Context("when all clauses are populated", func() {
				BeforeEach(func() {
					allowNull := false
					critetria = TestQueryCriteria{
						AssetType:                   "SERVER",
						IncludeNamePattern:          []string{"-db", "-dbmgmt"},
						Roles:                       []string{"db", "dbmgmt"},
						ExcludeNamePattern:          []string{"-core", "-drp"},
						AllowNullLastDiscoveredDate: &allowNull,
					}

					expectedClause = "(Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND Role__r.Name IN ('db','dbmgmt') AND ((NOT Host_Name__c LIKE '%-core%') AND (NOT Host_Name__c LIKE '%-drp%')) AND Tech_Asset__r.Asset_Type_Asset_Type__c = 'SERVER' AND Last_Discovered_Date__c != null"
				})

				It("returns properly formed clause joined by AND clause", func() {
					Expect(clause).To(Equal(expectedClause))
				})
			})
		})

		Context("when pointer is passed as argument", func() {
			var (
				critetria *TestQueryCriteria
			)

			JustBeforeEach(func() {
				clause, err = MarshalWhereClause(critetria)
			})

			Context("when nil is passed as argument", func() {
				It("returns empty where clause", func() {
					Expect(err).To(Equal(ErrNilValue))
					Expect(clause).To(BeEmpty())
				})
			})

			Context("when empty value is passed as argument", func() {
				BeforeEach(func() {
					critetria = &TestQueryCriteria{}
				})

				It("returns empty where clause", func() {
					Expect(clause).To(BeEmpty())
				})
			})

			Context("when all values are populated", func() {
				BeforeEach(func() {
					critetria = &TestQueryCriteria{
						AssetType:          "SERVER",
						IncludeNamePattern: []string{"-db", "-dbmgmt"},
						Roles:              []string{"db", "dbmgmt"},
						ExcludeNamePattern: []string{"-core", "-drp"},
					}

					expectedClause = "(Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND Role__r.Name IN ('db','dbmgmt') AND ((NOT Host_Name__c LIKE '%-core%') AND (NOT Host_Name__c LIKE '%-drp%')) AND Tech_Asset__r.Asset_Type_Asset_Type__c = 'SERVER'"
				})

				It("returns properly formed clause joined by AND clause", func() {
					Expect(clause).To(Equal(expectedClause))
				})
			})
		})

		Context("when all clauses are signed integer data types", func() {
			var criteria QueryCriteriaWithIntegerTypes
			BeforeEach(func() {
				criteria = QueryCriteriaWithIntegerTypes{
					NumOfCPUCores:                    16,
					PhysicalCPUCount:                 4,
					NumOfSuccessivePuppetRunFailures: -1,
					NumOfCoolanLogFiles:              1024,
					PvtTestFailCount:                 9223372036854775807,
				}

				expectedClause = "Num_of_CPU_Cores__c = 16 AND Physical_CPU_Count__c = 4 AND Number_Of_Successive_Puppet_Run_Failures__c = -1 AND Num_Of_Coolan_Log_Files__c = 1024 AND Pvt_Test_Fail_Count__c = 9223372036854775807"
			})

			It("returns properly formed clause joined by AND clause", func() {
				clause, err = MarshalWhereClause(criteria)
				Expect(err).ToNot(HaveOccurred())
				Expect(clause).To(Equal(expectedClause))
			})
		})

		Context("when all clauses are unsigned integer data types", func() {
			var criteria QueryCriteriaWithUnsignedIntegerTypes
			BeforeEach(func() {
				criteria = QueryCriteriaWithUnsignedIntegerTypes{
					NumOfCPUCores:                    16,
					PhysicalCPUCount:                 4,
					NumOfSuccessivePuppetRunFailures: 0,
					NumOfCoolanLogFiles:              1024,
					PvtTestFailCount:                 9223372036854775807,
				}

				expectedClause = "Num_of_CPU_Cores__c = 16 AND Physical_CPU_Count__c = 4 AND Number_Of_Successive_Puppet_Run_Failures__c = 0 AND Num_Of_Coolan_Log_Files__c = 1024 AND Pvt_Test_Fail_Count__c = 9223372036854775807"
			})

			It("returns properly formed clause joined by AND clause", func() {
				clause, err = MarshalWhereClause(criteria)
				Expect(err).ToNot(HaveOccurred())
				Expect(clause).To(Equal(expectedClause))
			})
		})

		Context("when all clauses are float data types", func() {
			var criteria QueryCriteriaWithFloatTypes
			BeforeEach(func() {
				criteria = QueryCriteriaWithFloatTypes{
					NumOfCPUCores:    16.00000000,
					PhysicalCPUCount: -4.12345678,
				}

				expectedClause = "Num_of_CPU_Cores__c = 16 AND Physical_CPU_Count__c = -4.12345678"
			})

			It("returns properly formed clause joined by AND clause", func() {
				clause, err = MarshalWhereClause(criteria)
				Expect(err).ToNot(HaveOccurred())
				Expect(clause).To(Equal(expectedClause))
			})
		})

		Context("when all clauses are boolean data types", func() {
			var criteria QueryCriteriaWithBooleanType
			BeforeEach(func() {
				criteria = QueryCriteriaWithBooleanType{
					NUMAEnabled:   true,
					DisableAlerts: false,
				}

				expectedClause = "NUMA_Enabled__c = true AND Disable_Alerts__c = false"
			})

			It("returns properly formed clause joined by AND clause", func() {
				clause, err = MarshalWhereClause(criteria)
				Expect(err).ToNot(HaveOccurred())
				Expect(clause).To(Equal(expectedClause))
			})
		})

		Context("when all clauses are date time data types", func() {
			var criteria QueryCriteriaWithDateTimeType
			var currentTime time.Time
			BeforeEach(func() {
				currentTime = time.Now()
				criteria = QueryCriteriaWithDateTimeType{
					CreatedDate: currentTime,
				}

				expectedClause = "CreatedDate = " + currentTime.Format(DateFormat)
			})

			It("returns properly formed clause", func() {
				clause, err = MarshalWhereClause(criteria)
				Expect(err).ToNot(HaveOccurred())
				Expect(clause).To(Equal(expectedClause))
			})
		})

		Context("when all clauses are mixed data types and operators", func() {
			var criteria QueryCriteriaWithMixedDataTypesAndOperators
			var currentTime time.Time
			BeforeEach(func() {
				currentTime = time.Now()
				criteria = QueryCriteriaWithMixedDataTypesAndOperators{
					BIOSType:                         "98.7.654a",
					NumOfCPUCores:                    32,
					NUMAEnabled:                      true,
					PvtTestFailCount:                 256,
					PhysicalCPUCount:                 4,
					CreatedDate:                      currentTime,
					DisableAlerts:                    false,
					AllocationLatency:                10.5,
					MajorOSVersion:                   "20",
					NumOfSuccessivePuppetRunFailures: 0,
					LastRestart:                      currentTime,
				}

				expectedClause = "BIOS_Type__c = '98.7.654a' AND Num_of_CPU_Cores__c > 32 AND NUMA_Enabled__c = true AND Pvt_Test_Fail_Count__c <= 256 AND Physical_CPU_Count__c >= 4 AND CreatedDate = " + currentTime.Format(DateFormat) + " AND Disable_Alerts__c = false AND Allocation_Latency__c < 10.5 AND Major_OS_Version__c = '20' AND Number_Of_Successive_Puppet_Run_Failures__c = 0 AND Last_Restart__c > " + currentTime.Format(DateFormat)
			})

			It("returns properly formed clause", func() {
				clause, err = MarshalWhereClause(criteria)
				Expect(err).ToNot(HaveOccurred())
				Expect(clause).To(Equal(expectedClause))
			})
		})

		Context("when clauses contains gt, gte, lt and lte operators", func() {
			var criteria QueryCriteriaNumericComparisonOperators
			BeforeEach(func() {
				criteria = QueryCriteriaNumericComparisonOperators{
					NumOfCPUCores:                    16,
					PhysicalCPUCount:                 4,
					NumOfSuccessivePuppetRunFailures: 0,
					NumOfCoolanLogFiles:              1024,
				}

				expectedClause = "Num_of_CPU_Cores__c > 16 AND Physical_CPU_Count__c < 4 AND Number_Of_Successive_Puppet_Run_Failures__c >= 0 AND Num_Of_Coolan_Log_Files__c <= 1024"
			})

			It("returns properly formed clause", func() {
				clause, err = MarshalWhereClause(criteria)
				Expect(err).ToNot(HaveOccurred())
				Expect(clause).To(Equal(expectedClause))
			})
		})

		Context("when no fieldName parameter is specified in tag", func() {
			var defaultFieldNameCriteria DefaultFieldNameQueryCriteria
			BeforeEach(func() {
				defaultFieldNameCriteria = DefaultFieldNameQueryCriteria{
					IncludeNamePattern: []string{"-db", "-dbmgmt"},
					Role:               []string{"foo", "bar"},
				}
				expectedClause = "(Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND Role IN ('foo','bar')"
			})

			It("returns properly formed clause joined by AND clause", func() {
				clause, err = MarshalWhereClause(defaultFieldNameCriteria)
				Expect(err).ToNot(HaveOccurred())
				Expect(clause).To(Equal(expectedClause))
			})
		})

		Context("when tag is invalid", func() {
			Context("when struct has invalid tag key", func() {
				type InvalidCriteriaStruct struct {
					SomePattern      []string `soql:"likeOperator,fieldName=Some_Pattern__c"`
					SomeOtherPattern string   `soql:"invalidClause,fieldName=Some_Other_Field"`
				}

				It("returns ErrInvalidTag error", func() {
					str, err := MarshalWhereClause(InvalidCriteriaStruct{})
					Expect(err).To(Equal(ErrInvalidTag))
					Expect(str).To(BeEmpty())
				})
			})

			Context("when struct has missing fieldName", func() {
				type MissingFieldName struct {
					SomePattern      []string `soql:"likeOperator,fieldName=Some_Pattern__c"`
					SomeOtherPattern string   `soql:"equalsOperator,fieldName="`
				}

				It("returns ErrInvalidTag error", func() {
					str, err := MarshalWhereClause(MissingFieldName{
						SomePattern:      []string{"test"},
						SomeOtherPattern: "foo",
					})
					Expect(err).To(Equal(ErrInvalidTag))
					Expect(str).To(BeEmpty())
				})
			})

			Context("when struct has invalid types", func() {
				It("returns empty string", func() {
					str, err := MarshalWhereClause(QueryCriteriaWithInvalidTypes{})
					Expect(err).ToNot(HaveOccurred())
					Expect(str).To(BeEmpty())
				})
			})
		})
	})

	Describe("MarshalSelectClause", func() {
		Context("when non pointer value is passed as argument", func() {
			Context("when no relationship name is passed", func() {
				Context("when no nested struct is passed", func() {
					It("returns just the json tag names of fields concatenanted by comma", func() {
						str, err := MarshalSelectClause(NonNestedStruct{}, "")
						Expect(err).ToNot(HaveOccurred())
						Expect(str).To(Equal("Name,SomeValue__c"))
					})
				})

				Context("when no fieldName parameter is specified in tag", func() {
					It("returns propery resolved list of field names by using defaults", func() {
						str, err := MarshalSelectClause(DefaultFieldNameStruct{}, "")
						Expect(err).ToNot(HaveOccurred())
						Expect(str).To(Equal("DefaultName,Description__c"))
					})
				})

				Context("when nested struct is passed", func() {
					It("returns properly resolved list of field names", func() {
						str, err := MarshalSelectClause(NestedStruct{}, "")
						Expect(err).ToNot(HaveOccurred())
						Expect(str).To(Equal("Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c"))
					})
				})
			})

			Context("when relationship name is passed", func() {
				Context("when no nested struct is passed", func() {
					It("returns just the json tag names of fields concatenanted by comma and prefixed by relationship name", func() {
						str, err := MarshalSelectClause(NonNestedStruct{}, "Role__r")
						Expect(err).ToNot(HaveOccurred())
						Expect(str).To(Equal("Role__r.Name,Role__r.SomeValue__c"))
					})
				})
			})

			Context("when struct has invalid tag key", func() {
				type InvalidStruct struct {
					Id  string `soql:"selectColumn,fieldName=Id"`
					Foo string `soql:"invalidClause,fieldName=Foo"`
				}

				It("returns ErrInvalidTag error", func() {
					str, err := MarshalSelectClause(InvalidStruct{}, "")
					Expect(err).To(Equal(ErrInvalidTag))
					Expect(str).To(BeEmpty())
				})
			})

			Context("when struct has missing fieldName", func() {
				type MissingFieldName struct {
					SomePattern      []string `soql:"selectColumn,fieldName=Some_Pattern__c"`
					SomeOtherPattern string   `soql:"selectColumn,fieldName="`
				}

				It("returns ErrInvalidTag error", func() {
					str, err := MarshalSelectClause(MissingFieldName{}, "")
					Expect(err).To(Equal(ErrInvalidTag))
					Expect(str).To(BeEmpty())
				})
			})

			Context("when struct has child relationship", func() {
				Context("when child struct has select clause only", func() {
					It("returns properly constructed select clause", func() {
						str, err := MarshalSelectClause(ParentStruct{}, "")
						Expect(err).ToNot(HaveOccurred())
						Expect(str).To(Equal("Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c,(SELECT SM_Application_Versions__c.Version__c FROM Application_Versions__r)"))
					})
				})

				Context("when child struct has select clause and where clause", func() {
					It("returns properly constructed select clause", func() {
						str, err := MarshalSelectClause(ParentStruct{
							ChildStruct: TestChildStruct{
								WhereClause: ChildQueryCriteria{
									Name: "sfdc-release",
								},
							},
						}, "")
						Expect(err).ToNot(HaveOccurred())
						Expect(str).To(Equal("Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c,(SELECT SM_Application_Versions__c.Version__c FROM Application_Versions__r WHERE SM_Application_Versions__c.Name__c = 'sfdc-release')"))
					})
				})

				Context("when selectChild tag does not have fieldName parameter", func() {
					It("returns properly constructed select clause", func() {
						str, err := MarshalSelectClause(DefaultFieldNameParentStruct{
							ChildStruct: TestChildStruct{
								WhereClause: ChildQueryCriteria{
									Name: "sfdc-release",
								},
							},
						}, "")
						Expect(err).ToNot(HaveOccurred())
						Expect(str).To(Equal("Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c,(SELECT SM_Application_Versions__c.Version__c FROM ChildStruct WHERE SM_Application_Versions__c.Name__c = 'sfdc-release')"))
					})
				})

				Context("when child struct does not have select clause", func() {
					It("returns error", func() {
						_, err := MarshalSelectClause(InvalidParentStruct{}, "")
						Expect(err).To(Equal(ErrNoSelectClause))
					})
				})

				Context("when selectChild tag is applied to non struct member", func() {
					It("returns error", func() {
						_, err := MarshalSelectClause(ChildTagToNonStruct{}, "")
						Expect(err).To(Equal(ErrInvalidTag))
					})
				})
			})
		})

		Context("when pointer value is passed as argument", func() {
			Context("when nil is passed", func() {
				It("returns ErrNilValue error", func() {
					var r *NestedStruct
					str, err := MarshalSelectClause(r, "")
					Expect(err).To(Equal(ErrNilValue))
					Expect(str).To(BeEmpty())
				})
			})

			Context("when nested struct is passed", func() {
				It("returns properly resolved list of field names", func() {
					str, err := MarshalSelectClause(&NestedStruct{}, "")
					Expect(err).ToNot(HaveOccurred())
					Expect(str).To(Equal("Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c"))
				})
			})
		})
	})

	Describe("Marshal", func() {
		var (
			soqlStruct    interface{}
			expectedQuery string
			actualQuery   string
			err           error
		)

		JustBeforeEach(func() {
			actualQuery, err = Marshal(soqlStruct)
		})

		Context("when empty struct is passed as argument", func() {
			BeforeEach(func() {
				soqlStruct = EmptyStruct{}
			})

			It("returns empty string", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(actualQuery).To(BeEmpty())
			})
		})

		Context("when valid value is passed as argument", func() {
			BeforeEach(func() {
				soqlStruct = TestSoqlStruct{
					SelectClause: NestedStruct{},
					WhereClause: TestQueryCriteria{
						IncludeNamePattern: []string{"-db", "-dbmgmt"},
						Roles:              []string{"db", "dbmgmt"},
					},
				}
				expectedQuery = "SELECT Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c FROM SM_Logical_Host__c WHERE (Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND Role__r.Name IN ('db','dbmgmt')"
			})

			It("returns properly constructed soql query", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(actualQuery).To(Equal(expectedQuery))
			})
		})

		Context("when valid value with mixed data type and operator is passed as argument", func() {
			BeforeEach(func() {
				currentTime := time.Now()
				soqlStruct = TestSoqlMixedDataAndOperatorStruct{
					SelectClause: NestedStruct{},
					WhereClause: QueryCriteriaWithMixedDataTypesAndOperators{
						BIOSType:                         "98.7.654a",
						NumOfCPUCores:                    32,
						NUMAEnabled:                      true,
						PvtTestFailCount:                 256,
						PhysicalCPUCount:                 4,
						CreatedDate:                      currentTime,
						DisableAlerts:                    false,
						AllocationLatency:                10.5,
						MajorOSVersion:                   "20",
						NumOfSuccessivePuppetRunFailures: 0,
						LastRestart:                      currentTime,
					},
				}
				expectedQuery = "SELECT Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c FROM SM_Logical_Host__c WHERE BIOS_Type__c = '98.7.654a' AND Num_of_CPU_Cores__c > 32 AND NUMA_Enabled__c = true AND Pvt_Test_Fail_Count__c <= 256 AND Physical_CPU_Count__c >= 4 AND CreatedDate = " + currentTime.Format(DateFormat) + " AND Disable_Alerts__c = false AND Allocation_Latency__c < 10.5 AND Major_OS_Version__c = '20' AND Number_Of_Successive_Puppet_Run_Failures__c = 0 AND Last_Restart__c > " + currentTime.Format(DateFormat)
			})

			It("returns properly constructed soql query", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(actualQuery).To(Equal(expectedQuery))
			})
		})

		Context("when valid pointer is passed as argument", func() {
			BeforeEach(func() {
				soqlStruct = &TestSoqlStruct{
					SelectClause: NestedStruct{},
					WhereClause: TestQueryCriteria{
						IncludeNamePattern: []string{"-db", "-dbmgmt"},
						Roles:              []string{"db", "dbmgmt"},
					},
				}
				expectedQuery = "SELECT Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c FROM SM_Logical_Host__c WHERE (Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND Role__r.Name IN ('db','dbmgmt')"
			})

			It("returns properly constructed soql query", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(actualQuery).To(Equal(expectedQuery))
			})
		})

		Context("when struct with no soql tags is passed", func() {
			BeforeEach(func() {
				soqlStruct = NonSoqlStruct{}
			})

			It("returns emptyString", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(actualQuery).To(BeEmpty())
			})
		})

		Context("when struct with multiple selectClause is passed", func() {
			BeforeEach(func() {
				soqlStruct = MultipleSelectClause{}
			})

			It("returns error", func() {
				Expect(err).To(Equal(ErrMultipleSelectClause))
			})
		})

		Context("when struct with multiple whereClause is passed", func() {
			BeforeEach(func() {
				soqlStruct = MultipleWhereClause{}
			})

			It("returns error", func() {
				Expect(err).To(Equal(ErrMultipleWhereClause))
			})
		})

		Context("when struct with only whereClause is passed", func() {
			BeforeEach(func() {
				soqlStruct = OnlyWhereClause{
					WhereClause: TestQueryCriteria{
						AssetType:          "SERVER",
						IncludeNamePattern: []string{"-db", "-dbmgmt"},
						Roles:              []string{"db", "dbmgmt"},
					},
				}
			})

			It("returns error", func() {
				Expect(err).To(Equal(ErrNoSelectClause))
			})
		})

		Context("when struct with multiple whereClause is passed", func() {
			BeforeEach(func() {
				soqlStruct = MultipleWhereClause{
					WhereClause1: ChildQueryCriteria{
						Name: "foo",
					},
					WhereClause2: ChildQueryCriteria{
						Name: "bar",
					},
				}
			})

			It("returns error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when nil pointer is passed", func() {
			BeforeEach(func() {
				var ptr *TestSoqlStruct
				soqlStruct = ptr
			})

			It("returns ErrNilValue error", func() {
				Expect(err).To(Equal(ErrNilValue))
			})
		})

		Context("when struct with invalid tag is passed", func() {
			BeforeEach(func() {
				soqlStruct = InvalidTagInStruct{}
			})

			It("returns error", func() {
				Expect(err).To(Equal(ErrInvalidTag))
			})
		})

		Context("when no table name is specified for selectClause", func() {
			BeforeEach(func() {
				soqlStruct = DefaultTableNameStruct{
					SomeTableName: NestedStruct{},
					WhereClause: TestQueryCriteria{
						IncludeNamePattern: []string{"-db", "-dbmgmt"},
						Roles:              []string{"db", "dbmgmt"},
					},
				}
				expectedQuery = "SELECT Id,Name__c,NonNestedStruct__r.Name,NonNestedStruct__r.SomeValue__c FROM SomeTableName WHERE (Host_Name__c LIKE '%-db%' OR Host_Name__c LIKE '%-dbmgmt%') AND Role__r.Name IN ('db','dbmgmt')"
			})

			It("uses name of the field as table name and returns properly constructed soql query", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(actualQuery).To(Equal(expectedQuery))
			})
		})
	})
})
