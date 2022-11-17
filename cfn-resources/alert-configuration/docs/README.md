# MongoDB::Atlas::AlertConfiguration

Returns and edits the conditions that trigger alerts and how MongoDB Cloud notifies users. This collection remains under revision and may change. Refer to the legacy documentation for this collection in the following link.

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "Type" : "MongoDB::Atlas::AlertConfiguration",
    "Properties" : {
        "<a href="#alertconfigid" title="AlertConfigId">AlertConfigId</a>" : <i>String</i>,
        "<a href="#apikeys" title="ApiKeys">ApiKeys</a>" : <i><a href="apikeydefinition.md">ApiKeyDefinition</a></i>,
        "<a href="#enabled" title="Enabled">Enabled</a>" : <i>Boolean</i>,
        "<a href="#eventtypename" title="EventTypeName">EventTypeName</a>" : <i>String</i>,
        "<a href="#groupid" title="GroupId">GroupId</a>" : <i>String</i>,
        "<a href="#includecount" title="IncludeCount">IncludeCount</a>" : <i>Boolean</i>,
        "<a href="#itemsperpage" title="ItemsPerPage">ItemsPerPage</a>" : <i>Integer</i>,
        "<a href="#matchers" title="Matchers">Matchers</a>" : <i>[ <a href="apimatcherview.md">ApiMatcherView</a>, ... ]</i>,
        "<a href="#metricthreshold" title="MetricThreshold">MetricThreshold</a>" : <i><a href="apimetricthresholdview.md">ApiMetricThresholdView</a></i>,
        "<a href="#notifications" title="Notifications">Notifications</a>" : <i>[ <a href="apinotificationview.md">ApiNotificationView</a>, ... ]</i>,
        "<a href="#pagenum" title="PageNum">PageNum</a>" : <i>Integer</i>,
        "<a href="#threshold" title="Threshold">Threshold</a>" : <i><a href="apiintegerthresholdview.md">ApiIntegerThresholdView</a></i>,
    }
}
</pre>

### YAML

<pre>
Type: MongoDB::Atlas::AlertConfiguration
Properties:
    <a href="#alertconfigid" title="AlertConfigId">AlertConfigId</a>: <i>String</i>
    <a href="#apikeys" title="ApiKeys">ApiKeys</a>: <i><a href="apikeydefinition.md">ApiKeyDefinition</a></i>
    <a href="#enabled" title="Enabled">Enabled</a>: <i>Boolean</i>
    <a href="#eventtypename" title="EventTypeName">EventTypeName</a>: <i>String</i>
    <a href="#groupid" title="GroupId">GroupId</a>: <i>String</i>
    <a href="#includecount" title="IncludeCount">IncludeCount</a>: <i>Boolean</i>
    <a href="#itemsperpage" title="ItemsPerPage">ItemsPerPage</a>: <i>Integer</i>
    <a href="#matchers" title="Matchers">Matchers</a>: <i>
      - <a href="apimatcherview.md">ApiMatcherView</a></i>
    <a href="#metricthreshold" title="MetricThreshold">MetricThreshold</a>: <i><a href="apimetricthresholdview.md">ApiMetricThresholdView</a></i>
    <a href="#notifications" title="Notifications">Notifications</a>: <i>
      - <a href="apinotificationview.md">ApiNotificationView</a></i>
    <a href="#pagenum" title="PageNum">PageNum</a>: <i>Integer</i>
    <a href="#threshold" title="Threshold">Threshold</a>: <i><a href="apiintegerthresholdview.md">ApiIntegerThresholdView</a></i>
</pre>

## Properties

#### AlertConfigId

Unique 24-hexadecimal digit string that identifies the alert configuration.

_Required_: No

_Type_: String

_Minimum_: <code>24</code>

_Maximum_: <code>24</code>

_Pattern_: <code>^([a-f0-9]{24})$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ApiKeys

_Required_: No

_Type_: <a href="apikeydefinition.md">ApiKeyDefinition</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Enabled

Flag that indicates whether someone enabled this alert configuration for the specified project.

_Required_: No

_Type_: Boolean

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### EventTypeName

Event type that triggers an alert.

_Required_: No

_Type_: String

_Allowed Values_: <code>AWS_ENCRYPTION_KEY_NEEDS_ROTATION</code> | <code>AZURE_ENCRYPTION_KEY_NEEDS_ROTATION</code> | <code>CLUSTER_MONGOS_IS_MISSING</code> | <code>CPS_RESTORE_FAILED</code> | <code>CPS_RESTORE_SUCCESSFUL</code> | <code>CPS_SNAPSHOT_BEHIND</code> | <code>CPS_SNAPSHOT_DOWNLOAD_REQUEST_FAILED</code> | <code>CPS_SNAPSHOT_FALLBACK_FAILED</code> | <code>CPS_SNAPSHOT_FALLBACK_SUCCESSFUL</code> | <code>CPS_SNAPSHOT_SUCCESSFUL</code> | <code>CREDIT_CARD_ABOUT_TO_EXPIRE</code> | <code>DAILY_BILL_OVER_THRESHOLD</code> | <code>GCP_ENCRYPTION_KEY_NEEDS_ROTATION</code> | <code>HOST_DOWN</code> | <code>JOINED_GROUP</code> | <code>NDS_X509_USER_AUTHENTICATION_CUSTOMER_CA_EXPIRATION_CHECK</code> | <code>NDS_X509_USER_AUTHENTICATION_MANAGED_USER_CERTS_EXPIRATION_CHECK</code> | <code>NO_PRIMARY</code> | <code>OUTSIDE_METRIC_THRESHOLD</code> | <code>OUTSIDE_SERVERLESS_METRIC_THRESHOLD</code> | <code>OUTSIDE_REALM_METRIC_THRESHOLD</code> | <code>PENDING_INVOICE_OVER_THRESHOLD</code> | <code>PRIMARY_ELECTED</code> | <code>REMOVED_FROM_GROUP</code> | <code>REPLICATION_OPLOG_WINDOW_RUNNING_OUT</code> | <code>TOO_MANY_ELECTIONS</code> | <code>USER_ROLES_CHANGED_AUDIT</code> | <code>USERS_WITHOUT_MULTIFACTOR_AUTH</code>

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

#### Matchers

List of rules that determine whether MongoDB Cloud checks an object for the alert configuration. You can filter using the matchers array if the **eventTypeName** specifies an event for a host, replica set, or sharded cluster.

_Required_: No

_Type_: List of <a href="apimatcherview.md">ApiMatcherView</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### MetricThreshold

_Required_: No

_Type_: <a href="apimetricthresholdview.md">ApiMetricThresholdView</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Notifications

List that contains the targets that MongoDB Cloud sends notifications.

_Required_: No

_Type_: List of <a href="apinotificationview.md">ApiNotificationView</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### PageNum

Number of the page that displays the current set of the total objects that the response returns.

_Required_: No

_Type_: Integer

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Threshold

_Required_: No

_Type_: <a href="apiintegerthresholdview.md">ApiIntegerThresholdView</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

## Return Values

### Fn::GetAtt

The `Fn::GetAtt` intrinsic function returns a value for a specified attribute of this type. The following are the available attributes and sample return values.

For more information about using the `Fn::GetAtt` intrinsic function, see [Fn::GetAtt](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html).

#### Id

Unique 24-hexadecimal digit string that identifies this alert configuration.

#### Links

List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.

#### TypeName

Human-readable label that displays the alert type.

#### Created

Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.

#### TotalCount

Number of documents returned in this response.

#### Updated

Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.

#### Results

List of returned documents that MongoDB Cloud provides when completing this request.

#### GroupId

Returns the <code>GroupId</code> value.

#### Id

Returns the <code>Id</code> value.

#### Updated

Returns the <code>Updated</code> value.

#### Created

Returns the <code>Created</code> value.

#### Links

Returns the <code>Links</code> value.

#### TypeName

Returns the <code>TypeName</code> value.

#### Number

Returns the <code>Number</code> value.

#### Units

Returns the <code>Units</code> value.

#### GroupId

Returns the <code>GroupId</code> value.

#### Id

Returns the <code>Id</code> value.

#### ClusterName

Returns the <code>ClusterName</code> value.

#### Created

Returns the <code>Created</code> value.

#### EventTypeName

Returns the <code>EventTypeName</code> value.

#### Status

Returns the <code>Status</code> value.

#### TypeName

Returns the <code>TypeName</code> value.

#### Updated

Returns the <code>Updated</code> value.

#### HostnameAndPort

Returns the <code>HostnameAndPort</code> value.

#### LastNotified

Returns the <code>LastNotified</code> value.

#### MetricName

Returns the <code>MetricName</code> value.

#### AlertConfigId

Returns the <code>AlertConfigId</code> value.

#### ReplicaSetName

Returns the <code>ReplicaSetName</code> value.

#### CurrentValue

Returns the <code>CurrentValue</code> value.

#### Resolved

Returns the <code>Resolved</code> value.

#### AcknowledgingUsername

Returns the <code>AcknowledgingUsername</code> value.

#### Links

Returns the <code>Links</code> value.

