# MongoDB::Atlas::ServerlessInstance ApiAtlasServerlessClusterDescriptionViewManual

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "<a href="#apikeys" title="ApiKeys">ApiKeys</a>" : <i><a href="apikeydefinition.md">ApiKeyDefinition</a></i>,
    "<a href="#connectionstrings" title="ConnectionStrings">ConnectionStrings</a>" : <i><a href="serverlessinstanceconnectionstrings.md">ServerlessInstanceConnectionStrings</a></i>,
    "<a href="#createdate" title="CreateDate">CreateDate</a>" : <i>String</i>,
    "<a href="#groupid" title="GroupId">GroupId</a>" : <i>String</i>,
    "<a href="#id" title="Id">Id</a>" : <i>String</i>,
    "<a href="#links" title="Links">Links</a>" : <i>[ <a href="link.md">Link</a>, ... ]</i>,
    "<a href="#mongodbversion" title="MongoDBVersion">MongoDBVersion</a>" : <i>String</i>,
    "<a href="#name" title="Name">Name</a>" : <i>String</i>,
    "<a href="#providersettings" title="ProviderSettings">ProviderSettings</a>" : <i><a href="serverlessinstanceprovidersettings.md">ServerlessInstanceProviderSettings</a></i>,
    "<a href="#statename" title="StateName">StateName</a>" : <i>String</i>
}
</pre>

### YAML

<pre>
<a href="#apikeys" title="ApiKeys">ApiKeys</a>: <i><a href="apikeydefinition.md">ApiKeyDefinition</a></i>
<a href="#connectionstrings" title="ConnectionStrings">ConnectionStrings</a>: <i><a href="serverlessinstanceconnectionstrings.md">ServerlessInstanceConnectionStrings</a></i>
<a href="#createdate" title="CreateDate">CreateDate</a>: <i>String</i>
<a href="#groupid" title="GroupId">GroupId</a>: <i>String</i>
<a href="#id" title="Id">Id</a>: <i>String</i>
<a href="#links" title="Links">Links</a>: <i>
      - <a href="link.md">Link</a></i>
<a href="#mongodbversion" title="MongoDBVersion">MongoDBVersion</a>: <i>String</i>
<a href="#name" title="Name">Name</a>: <i>String</i>
<a href="#providersettings" title="ProviderSettings">ProviderSettings</a>: <i><a href="serverlessinstanceprovidersettings.md">ServerlessInstanceProviderSettings</a></i>
<a href="#statename" title="StateName">StateName</a>: <i>String</i>
</pre>

## Properties

#### ApiKeys

_Required_: No

_Type_: <a href="apikeydefinition.md">ApiKeyDefinition</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ConnectionStrings

_Required_: No

_Type_: <a href="serverlessinstanceconnectionstrings.md">ServerlessInstanceConnectionStrings</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### CreateDate

Date and time when MongoDB Cloud created this serverless instance. MongoDB Cloud represents this timestamp in ISO 8601 format in UTC.

_Required_: No

_Type_: String

_Pattern_: <code>^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d(?:\\.\\d{1,9})?(?:Z|[+-][01]\\d:[0-5]\\d)$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### GroupId

Unique 24-hexadecimal character string that identifies the project.

_Required_: No

_Type_: String

_Minimum_: <code>24</code>

_Maximum_: <code>24</code>

_Pattern_: <code>^([a-f0-9]{24})$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Id

Unique 24-hexadecimal digit string that identifies the serverless instance.

_Required_: No

_Type_: String

_Minimum_: <code>24</code>

_Maximum_: <code>24</code>

_Pattern_: <code>^([a-f0-9]{24})$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Links

_Required_: No

_Type_: List of <a href="link.md">Link</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### MongoDBVersion

Version of MongoDB that the serverless instance runs.

_Required_: No

_Type_: String

_Pattern_: <code>([\d]+\.[\d]+\.[\d]+)</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Name

Human-readable label that identifies the serverless instance.

_Required_: No

_Type_: String

_Minimum_: <code>1</code>

_Maximum_: <code>64</code>

_Pattern_: <code>^([a-zA-Z0-9]([a-zA-Z0-9-]){0,21}(?<!-)([\w]{0,42}))$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ProviderSettings

_Required_: No

_Type_: <a href="serverlessinstanceprovidersettings.md">ServerlessInstanceProviderSettings</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### StateName

Human-readable label that indicates the current operating condition of the serverless instance.

_Required_: No

_Type_: String

_Allowed Values_: <code>CREATING</code> | <code>DELETED</code> | <code>DELETING</code> | <code>IDLE</code> | <code>REPAIRING</code> | <code>UPDATING</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

