package resource

import "github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/encoding"

/*
This file is autogenerated, do not edit;
changes will be undone by the next 'generate' command.

Updates to this type are made my editing the schema file
and executing the 'generate' command
*/

// Model is autogenerated from the json schema
type Model struct {
	ProjectId  *encoding.String `json:"ProjectId,omitempty"`
	Id         *encoding.String `json:"Id,omitempty"`
	Rules      []RuleDefinition `json:"Rules,omitempty"`
	TotalCount *encoding.Int    `json:"TotalCount,omitempty"`
}

// RuleDefinition is autogenerated from the json schema
type RuleDefinition struct {
	Comment   *encoding.String
	IpAddress *encoding.String
	CidrBlock *encoding.String
}
