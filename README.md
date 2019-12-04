# SOQL

This package supports marshalling a golang struct into SOQL. Like `json` tags, this package provides `soql` tags that you can use to annotate your golang structs. Once tagged `Marshal` method will return SOQL query that will let you query the required Salesforce object using Salesforce API.
