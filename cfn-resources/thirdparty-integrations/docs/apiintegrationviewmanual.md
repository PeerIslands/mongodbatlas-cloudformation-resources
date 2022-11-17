# MongoDB::Atlas::Integrations ApiIntegrationViewManual

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "<a href="#apikeys" title="ApiKeys">ApiKeys</a>" : <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>,
    "<a href="#type" title="Type">Type</a>" : <i>String</i>
}
</pre>

### YAML

<pre>
<a href="#apikeys" title="ApiKeys">ApiKeys</a>: <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>
<a href="#type" title="Type">Type</a>: <i>String</i>
</pre>

## Properties

#### ApiKeys

_Required_: No

_Type_: <a href="apikeydefinition.md">apiKeyDefinition</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Type

Human-readable label that identifies the service to which you want to integrate with MongoDB Cloud. The value must match the third-party service integration type.

_Required_: No

_Type_: String

_Allowed Values_: <code>PAGER_DUTY</code> | <code>MICROSOFT_TEAMS</code> | <code>SLACK</code> | <code>DATADOG</code> | <code>NEW_RELIC</code> | <code>OPS_GENIE</code> | <code>VICTOR_OPS</code> | <code>FLOWDOCK</code> | <code>WEBHOOK</code> | <code>PROMETHEUS</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

