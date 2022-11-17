# MongoDB::Atlas::CloudBackupSnapshotExportJob ApiAtlasDiskBackupBaseRestoreMemberView

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "<a href="#apikeys" title="ApiKeys">ApiKeys</a>" : <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>,
    "<a href="#replicasetname" title="ReplicaSetName">ReplicaSetName</a>" : <i>String</i>
}
</pre>

### YAML

<pre>
<a href="#apikeys" title="ApiKeys">ApiKeys</a>: <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>
<a href="#replicasetname" title="ReplicaSetName">ReplicaSetName</a>: <i>String</i>
</pre>

## Properties

#### ApiKeys

_Required_: No

_Type_: <a href="apikeydefinition.md">apiKeyDefinition</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ReplicaSetName

Human-readable label that identifies the replica set on the sharded cluster.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

