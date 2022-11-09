# Mongodb::Atlas::ProjectInvitation ApiGroupInvitationView

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "<a href="#apikeys" title="ApiKeys">ApiKeys</a>" : <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>,
    "<a href="#createdat" title="CreatedAt">CreatedAt</a>" : <i>String</i>,
    "<a href="#expiresat" title="ExpiresAt">ExpiresAt</a>" : <i>String</i>,
    "<a href="#groupid" title="GroupId">GroupId</a>" : <i>String</i>,
    "<a href="#groupname" title="GroupName">GroupName</a>" : <i>String</i>,
    "<a href="#id" title="Id">Id</a>" : <i>String</i>,
    "<a href="#inviterusername" title="InviterUsername">InviterUsername</a>" : <i>String</i>,
    "<a href="#links" title="Links">Links</a>" : <i>[ <a href="link.md">Link</a>, ... ]</i>,
    "<a href="#roles" title="Roles">Roles</a>" : <i>[ String, ... ]</i>,
    "<a href="#username" title="Username">Username</a>" : <i>String</i>
}
</pre>

### YAML

<pre>
<a href="#apikeys" title="ApiKeys">ApiKeys</a>: <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>
<a href="#createdat" title="CreatedAt">CreatedAt</a>: <i>String</i>
<a href="#expiresat" title="ExpiresAt">ExpiresAt</a>: <i>String</i>
<a href="#groupid" title="GroupId">GroupId</a>: <i>String</i>
<a href="#groupname" title="GroupName">GroupName</a>: <i>String</i>
<a href="#id" title="Id">Id</a>: <i>String</i>
<a href="#inviterusername" title="InviterUsername">InviterUsername</a>: <i>String</i>
<a href="#links" title="Links">Links</a>: <i>
      - <a href="link.md">Link</a></i>
<a href="#roles" title="Roles">Roles</a>: <i>
      - String</i>
<a href="#username" title="Username">Username</a>: <i>String</i>
</pre>

## Properties

#### ApiKeys

_Required_: No

_Type_: <a href="apikeydefinition.md">apiKeyDefinition</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### CreatedAt

Date and time when MongoDB Cloud sent the invitation. This parameter expresses its value in ISO 8601 format in UTC.

_Required_: No

_Type_: String

_Pattern_: <code>^(?:[1-9]\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\d|2[0-3]):[0-5]\d:[0-5]\d(?:\.\d{1,9})?(?:Z)$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ExpiresAt

Date and time when MongoDB Cloud expires the invitation. This parameter expresses its value in ISO 8601 format in UTC.

_Required_: No

_Type_: String

_Pattern_: <code>^(?:[1-9]\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\d|2[0-3]):[0-5]\d:[0-5]\d(?:\.\d{1,9})?(?:Z)$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### GroupId

Unique 24-hexadecimal character string that identifies the project.

_Required_: No

_Type_: String

_Minimum_: <code>24</code>

_Maximum_: <code>24</code>

_Pattern_: <code>^([a-f0-9]{24})$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### GroupName

Human-readable label that identifies the project to which you invited the MongoDB Cloud user.

_Required_: No

_Type_: String

_Minimum_: <code>1</code>

_Maximum_: <code>64</code>

_Pattern_: <code>^[\p{L}\p{N}\-_.(),:&@+']{1,64}$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Id

Unique 24-hexadecimal character string that identifies the invitation.

_Required_: No

_Type_: String

_Minimum_: <code>24</code>

_Maximum_: <code>24</code>

_Pattern_: <code>^([a-f0-9]{24})$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### InviterUsername

Email address of the MongoDB Cloud user who sent the invitation.

_Required_: No

_Type_: String

_Pattern_: <code>^[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Links

_Required_: No

_Type_: List of <a href="link.md">Link</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Roles

One or more organization or project level roles to assign to the MongoDB Cloud user.

_Required_: No

_Type_: List of String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Username

Email address of the MongoDB Cloud user invited to join the project.

_Required_: No

_Type_: String

_Pattern_: <code>^[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

