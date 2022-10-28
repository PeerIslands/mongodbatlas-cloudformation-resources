# Mongodb::Atlas::maintenancewindow

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "Type" : "Mongodb::Atlas::maintenancewindow",
    "Properties" : {
        "<a href="#apikeys" title="ApiKeys">ApiKeys</a>" : <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>,
        "<a href="#autodeferonceenabled" title="AutoDeferOnceEnabled">AutoDeferOnceEnabled</a>" : <i>Boolean</i>,
        "<a href="#dayofweek" title="DayOfWeek">DayOfWeek</a>" : <i>Double</i>,
        "<a href="#groupid" title="GroupId">GroupId</a>" : <i>String</i>,
        "<a href="#hourofday" title="HourOfDay">HourOfDay</a>" : <i>Double</i>,
        "<a href="#startasap" title="StartASAP">StartASAP</a>" : <i>Boolean</i>
    }
}
</pre>

### YAML

<pre>
Type: Mongodb::Atlas::maintenancewindow
Properties:
    <a href="#apikeys" title="ApiKeys">ApiKeys</a>: <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>
    <a href="#autodeferonceenabled" title="AutoDeferOnceEnabled">AutoDeferOnceEnabled</a>: <i>Boolean</i>
    <a href="#dayofweek" title="DayOfWeek">DayOfWeek</a>: <i>Double</i>
    <a href="#groupid" title="GroupId">GroupId</a>: <i>String</i>
    <a href="#hourofday" title="HourOfDay">HourOfDay</a>: <i>Double</i>
    <a href="#startasap" title="StartASAP">StartASAP</a>: <i>Boolean</i>
</pre>

## Properties

#### ApiKeys

_Required_: No

_Type_: <a href="apikeydefinition.md">apiKeyDefinition</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### AutoDeferOnceEnabled

Flag that indicates whether MongoDB Cloud should defer all maintenance windows for one week after you enable them.

_Required_: No

_Type_: Boolean

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### DayOfWeek

One-based integer that represents the day of the week that the maintenance window starts.

| Value | Day of Week |
|---|---|
| `1` | Sunday |
| `2` | Monday |
| `3` | Tuesday |
| `4` | Wednesday |
| `5` | Thursday |
| `6` | Friday |
| `7` | Saturday |


_Required_: No

_Type_: Double

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### GroupId

Unique 24-hexadecimal digit string that identifies your project.

_Required_: No

_Type_: String

_Minimum_: <code>24</code>

_Maximum_: <code>24</code>

_Pattern_: <code>^([a-f0-9]{24})$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### HourOfDay

Zero-based integer that represents the hour of the of the day that the maintenance window starts according to a 24-hour clock. Use `0` for midnight and `12` for noon.

_Required_: No

_Type_: Double

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### StartASAP

Flag that indicates whether MongoDB Cloud starts the maintenance window immediately upon receiving this request. To start the maintenance window immediately for your project, MongoDB Cloud must have maintenance scheduled and you must set a maintenance window. This flag resets to `false` after MongoDB Cloud completes maintenance.

_Required_: No

_Type_: Boolean

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

## Return Values

### Ref

When you pass the logical ID of this resource to the intrinsic `Ref` function, Ref returns the GroupId.
