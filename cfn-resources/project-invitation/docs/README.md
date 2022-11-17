# MongoDB::Atlas::ProjectInvitation

Returns, adds, and edits collections of clusters and users in MongoDB Cloud.

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "Type" : "MongoDB::Atlas::ProjectInvitation",
    "Properties" : {
        "<a href="#apikeys" title="ApiKeys">ApiKeys</a>" : <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>,
        "<a href="#groupid" title="GroupId">GroupId</a>" : <i>String</i>,
        "<a href="#includecount" title="IncludeCount">IncludeCount</a>" : <i>Boolean</i>,
        "<a href="#invitationid" title="InvitationId">InvitationId</a>" : <i>String</i>,
        "<a href="#itemsperpage" title="ItemsPerPage">ItemsPerPage</a>" : <i>Integer</i>,
        "<a href="#pagenum" title="PageNum">PageNum</a>" : <i>Integer</i>,
        "<a href="#roles" title="Roles">Roles</a>" : <i>[ String, ... ]</i>,
        "<a href="#username" title="Username">Username</a>" : <i>String</i>
    }
}
</pre>

### YAML

<pre>
Type: MongoDB::Atlas::ProjectInvitation
Properties:
    <a href="#apikeys" title="ApiKeys">ApiKeys</a>: <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>
    <a href="#groupid" title="GroupId">GroupId</a>: <i>String</i>
    <a href="#includecount" title="IncludeCount">IncludeCount</a>: <i>Boolean</i>
    <a href="#invitationid" title="InvitationId">InvitationId</a>: <i>String</i>
    <a href="#itemsperpage" title="ItemsPerPage">ItemsPerPage</a>: <i>Integer</i>
    <a href="#pagenum" title="PageNum">PageNum</a>: <i>Integer</i>
    <a href="#roles" title="Roles">Roles</a>: <i>
      - String</i>
    <a href="#username" title="Username">Username</a>: <i>String</i>
</pre>

## Properties

#### ApiKeys

_Required_: No

_Type_: <a href="apikeydefinition.md">apiKeyDefinition</a>

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

#### InvitationId

Unique 24-hexadecimal digit string that identifies the invitation.

_Required_: No

_Type_: String

_Minimum_: <code>24</code>

_Maximum_: <code>24</code>

_Pattern_: <code>^([a-f0-9]{24})$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ItemsPerPage

Number of items that the response returns per page.

_Required_: No

_Type_: Integer

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### PageNum

Number of the page that displays the current set of the total objects that the response returns.

_Required_: No

_Type_: Integer

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Roles

One or more organization or project level roles to assign to the MongoDB Cloud user.

_Required_: No

_Type_: List of String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Username

Email address of the user account invited to this project.

_Required_: No

_Type_: String

_Pattern_: <code>^[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

## Return Values

### Fn::GetAtt

The `Fn::GetAtt` intrinsic function returns a value for a specified attribute of this type. The following are the available attributes and sample return values.

For more information about using the `Fn::GetAtt` intrinsic function, see [Fn::GetAtt](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html).

#### Results

List of returned documents that MongoDB Cloud provides when completing this request.

#### GroupName

Human-readable label that identifies the project to which you invited the MongoDB Cloud user.

#### Id

Unique 24-hexadecimal character string that identifies the invitation.

#### Links

List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.

#### ExpiresAt

Date and time when MongoDB Cloud expires the invitation. This parameter expresses its value in ISO 8601 format in UTC.

#### InviterUsername

Email address of the MongoDB Cloud user who sent the invitation.

#### TotalCount

Number of documents returned in this response.

#### CreatedAt

Date and time when MongoDB Cloud sent the invitation. This parameter expresses its value in ISO 8601 format in UTC.

#### ExpiresAt

Returns the <code>ExpiresAt</code> value.

#### Username

Returns the <code>Username</code> value.

#### CreatedAt

Returns the <code>CreatedAt</code> value.

#### GroupId

Returns the <code>GroupId</code> value.

#### InviterUsername

Returns the <code>InviterUsername</code> value.

#### Links

Returns the <code>Links</code> value.

#### GroupName

Returns the <code>GroupName</code> value.

#### Id

Returns the <code>Id</code> value.

