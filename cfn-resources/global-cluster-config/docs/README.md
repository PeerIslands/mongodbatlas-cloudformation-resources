# MongoDB::Atlas::GlobalClusterConfig

Returns, adds, and removes Global Cluster managed namespaces and custom zone mappings.

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "Type" : "MongoDB::Atlas::GlobalClusterConfig",
    "Properties" : {
        "<a href="#apikeys" title="ApiKeys">ApiKeys</a>" : <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>,
        "<a href="#clustername" title="ClusterName">ClusterName</a>" : <i>String</i>,
        "<a href="#collection" title="Collection">Collection</a>" : <i>String</i>,
        "<a href="#id" title="Id">Id</a>" : <i>String</i>,
        "<a href="#db" title="Db">Db</a>" : <i>String</i>,
        "<a href="#groupid" title="GroupId">GroupId</a>" : <i>String</i>
    }
}
</pre>

### YAML

<pre>
Type: MongoDB::Atlas::GlobalClusterConfig
Properties:
    <a href="#apikeys" title="ApiKeys">ApiKeys</a>: <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>
    <a href="#clustername" title="ClusterName">ClusterName</a>: <i>String</i>
    <a href="#collection" title="Collection">Collection</a>: <i>String</i>
    <a href="#id" title="Id">Id</a>: <i>String</i>
    <a href="#db" title="Db">Db</a>: <i>String</i>
    <a href="#groupid" title="GroupId">GroupId</a>: <i>String</i>
</pre>

## Properties

#### ApiKeys

_Required_: No

_Type_: <a href="apikeydefinition.md">apiKeyDefinition</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ClusterName

Human-readable label that identifies this advanced cluster.

_Required_: No

_Type_: String

_Minimum_: <code>1</code>

_Maximum_: <code>64</code>

_Pattern_: <code>^([a-zA-Z0-9]([a-zA-Z0-9-]){0,21}(?<!-)([\w]{0,42}))$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Collection

Human-readable label that identifies the collection associated with the managed namespace.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Id

.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Db

Human-readable label that identifies the database that contains the collection.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### GroupId

Unique 24-hexadecimal digit string that identifies your project.

_Required_: No

_Type_: String

_Minimum_: <code>24</code>

_Maximum_: <code>24</code>

_Pattern_: <code>^([a-f0-9]{24})$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

