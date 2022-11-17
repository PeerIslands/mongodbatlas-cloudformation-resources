# MongoDB::Atlas::AlertConfiguration ApiIntegerThresholdView

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "<a href="#apikeys" title="ApiKeys">ApiKeys</a>" : <i><a href="apikeydefinition.md">ApiKeyDefinition</a></i>,
    "<a href="#operator" title="Operator">Operator</a>" : <i>String</i>,
    "<a href="#threshold" title="Threshold">Threshold</a>" : <i>Integer</i>,
    "<a href="#units" title="Units">Units</a>" : <i>String</i>
}
</pre>

### YAML

<pre>
<a href="#apikeys" title="ApiKeys">ApiKeys</a>: <i><a href="apikeydefinition.md">ApiKeyDefinition</a></i>
<a href="#operator" title="Operator">Operator</a>: <i>String</i>
<a href="#threshold" title="Threshold">Threshold</a>: <i>Integer</i>
<a href="#units" title="Units">Units</a>: <i>String</i>
</pre>

## Properties

#### ApiKeys

_Required_: No

_Type_: <a href="apikeydefinition.md">ApiKeyDefinition</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Operator

Comparison operator to apply when checking the current metric value.

_Required_: No

_Type_: String

_Allowed Values_: <code>GREATER_THAN</code> | <code>LESS_THAN</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Threshold

Value of metric that, when exceeded, triggers an alert.

_Required_: No

_Type_: Integer

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Units

Element used to express the quantity. This can be an element of time, storage capacity, and the like.

_Required_: No

_Type_: String

_Allowed Values_: <code>BITS</code> | <code>BYTES</code> | <code>DAYS</code> | <code>GIGABITS</code> | <code>GIGABYTES</code> | <code>HOURS</code> | <code>KILOBITS</code> | <code>KILOBYTES</code> | <code>MEGABITS</code> | <code>MEGABYTES</code> | <code>MILLISECONDS</code> | <code>MINUTES</code> | <code>PETABYTES</code> | <code>RAW</code> | <code>SECONDS</code> | <code>TERABYTES</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

