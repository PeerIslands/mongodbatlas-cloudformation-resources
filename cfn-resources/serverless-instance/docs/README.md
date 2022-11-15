# Mongodb::Atlas::ServerlessInstance

Returns, adds, edits, and removes serverless instances.

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "Type" : "Mongodb::Atlas::ServerlessInstance",
    "Properties" : {
        "<a href="#apikeys" title="ApiKeys">ApiKeys</a>" : <i><a href="apikeydefinition.md">ApiKeyDefinition</a></i>,
        "<a href="#connectionstrings" title="ConnectionStrings">ConnectionStrings</a>" : <i><a href="serverlessinstanceconnectionstrings.md">ServerlessInstanceConnectionStrings</a></i>,
        "<a href="#groupid" title="GroupId">GroupId</a>" : <i>String</i>,
        "<a href="#includecount" title="IncludeCount">IncludeCount</a>" : <i>Boolean</i>,
        "<a href="#itemsperpage" title="ItemsPerPage">ItemsPerPage</a>" : <i>Integer</i>,
        "<a href="#name" title="Name">Name</a>" : <i>String</i>,
        "<a href="#pagenum" title="PageNum">PageNum</a>" : <i>Integer</i>,
        "<a href="#providersettings" title="ProviderSettings">ProviderSettings</a>" : <i><a href="serverlessinstanceprovidersettings.md">ServerlessInstanceProviderSettings</a></i>,
    }
}
</pre>

### YAML

<pre>
Type: Mongodb::Atlas::ServerlessInstance
Properties:
    <a href="#apikeys" title="ApiKeys">ApiKeys</a>: <i><a href="apikeydefinition.md">ApiKeyDefinition</a></i>
    <a href="#connectionstrings" title="ConnectionStrings">ConnectionStrings</a>: <i><a href="serverlessinstanceconnectionstrings.md">ServerlessInstanceConnectionStrings</a></i>
    <a href="#groupid" title="GroupId">GroupId</a>: <i>String</i>
    <a href="#includecount" title="IncludeCount">IncludeCount</a>: <i>Boolean</i>
    <a href="#itemsperpage" title="ItemsPerPage">ItemsPerPage</a>: <i>Integer</i>
    <a href="#name" title="Name">Name</a>: <i>String</i>
    <a href="#pagenum" title="PageNum">PageNum</a>: <i>Integer</i>
    <a href="#providersettings" title="ProviderSettings">ProviderSettings</a>: <i><a href="serverlessinstanceprovidersettings.md">ServerlessInstanceProviderSettings</a></i>
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

#### GroupId

Unique 24-hexadecimal digit string that identifies your project.

_Required_: No

_Type_: String

_Minimum_: <code>24</code>

_Maximum_: <code>24</code>

_Pattern_: <code>^([a-f0-9]{24})$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### IncludeCount

Flag that indicates whether the response returns the total number of items (**totalCount**) in the response.

_Required_: No

_Type_: Boolean

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ItemsPerPage

Number of items that the response returns per page.

_Required_: No

_Type_: Integer

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Name

Human-readable label that identifies the serverless instance.

_Required_: No

_Type_: String

_Minimum_: <code>1</code>

_Maximum_: <code>64</code>

_Pattern_: <code>^([a-zA-Z0-9]([a-zA-Z0-9-]){0,21}(?<!-)([\w]{0,42}))$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### PageNum

Number of the page that displays the current set of the total objects that the response returns.

_Required_: No

_Type_: Integer

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ProviderSettings

_Required_: No

_Type_: <a href="serverlessinstanceprovidersettings.md">ServerlessInstanceProviderSettings</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

## Return Values

### Ref

When you pass the logical ID of this resource to the intrinsic `Ref` function, Ref returns the Id.

### Fn::GetAtt

The `Fn::GetAtt` intrinsic function returns a value for a specified attribute of this type. The following are the available attributes and sample return values.

For more information about using the `Fn::GetAtt` intrinsic function, see [Fn::GetAtt](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html).

#### Results

List of returned documents that MongoDB Cloud provides when completing this request.


#### CreateDate

Date and time when MongoDB Cloud created this serverless instance. MongoDB Cloud represents this timestamp in ISO 8601 format in UTC.

#### Id

Unique 24-hexadecimal digit string that identifies the serverless instance.

#### TotalCount

Number of documents returned in this response.

#### ConnectionStrings

Returns the <code>ConnectionStrings</code> value.

#### Links

List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.

#### StateName

Human-readable label that indicates the current operating condition of the serverless instance.

#### MongoDBVersion

Version of MongoDB that the serverless instance runs.

#### Endpoints

Returns the <code>Endpoints</code> value.

#### SrvConnectionString

Returns the <code>SrvConnectionString</code> value.

#### Type

Returns the <code>Type</code> value.

#### PrivateEndpoint

Returns the <code>PrivateEndpoint</code> value.

#### StandardSrv

Returns the <code>StandardSrv</code> value.

#### CreateDate

Returns the <code>CreateDate</code> value.

#### MongoDBVersion

Returns the <code>MongoDBVersion</code> value.

#### StateName

Returns the <code>StateName</code> value.

#### ConnectionStrings

Returns the <code>ConnectionStrings</code> value.

#### Links

Returns the <code>Links</code> value.

#### GroupId

Returns the <code>GroupId</code> value.

#### Id

Returns the <code>Id</code> value.

#### EndpointId

Returns the <code>EndpointId</code> value.

#### ProviderName

Returns the <code>ProviderName</code> value.

#### Region

Returns the <code>Region</code> value.

