# MongoDB::Atlas::LDAPVerify ApiAtlasNDSLDAPVerifyConnectivityJobRequestValidationView

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "<a href="#apikeys" title="ApiKeys">ApiKeys</a>" : <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>,
    "<a href="#status" title="Status">Status</a>" : <i>String</i>,
    "<a href="#validationtype" title="ValidationType">ValidationType</a>" : <i>String</i>
}
</pre>

### YAML

<pre>
<a href="#apikeys" title="ApiKeys">ApiKeys</a>: <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>
<a href="#status" title="Status">Status</a>: <i>String</i>
<a href="#validationtype" title="ValidationType">ValidationType</a>: <i>String</i>
</pre>

## Properties

#### ApiKeys

_Required_: No

_Type_: <a href="apikeydefinition.md">apiKeyDefinition</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Status

Human-readable string that indicates the result of this verification test.

_Required_: No

_Type_: String

_Allowed Values_: <code>FAIL</code> | <code>OK</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ValidationType

Human-readable label that identifies this verification test that MongoDB Cloud runs.

_Required_: No

_Type_: String

_Allowed Values_: <code>AUTHENTICATE</code> | <code>AUTHORIZATION_ENABLED</code> | <code>CONNECT</code> | <code>PARSE_AUTHZ_QUERY</code> | <code>QUERY_SERVER</code> | <code>SERVER_SPECIFIED</code> | <code>TEMPLATE</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

