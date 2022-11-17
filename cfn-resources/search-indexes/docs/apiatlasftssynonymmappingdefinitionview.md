# MongoDB::Atlas::SearchIndex ApiAtlasFTSSynonymMappingDefinitionView

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "<a href="#analyzer" title="Analyzer">Analyzer</a>" : <i>String</i>,
    "<a href="#apikeys" title="ApiKeys">ApiKeys</a>" : <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>,
    "<a href="#name" title="Name">Name</a>" : <i>String</i>,
    "<a href="#source" title="Source">Source</a>" : <i><a href="synonymsource.md">SynonymSource</a></i>
}
</pre>

### YAML

<pre>
<a href="#analyzer" title="Analyzer">Analyzer</a>: <i>String</i>
<a href="#apikeys" title="ApiKeys">ApiKeys</a>: <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>
<a href="#name" title="Name">Name</a>: <i>String</i>
<a href="#source" title="Source">Source</a>: <i><a href="synonymsource.md">SynonymSource</a></i>
</pre>

## Properties

#### Analyzer

Specific pre-defined method chosen to apply to the synonyms to be searched.

_Required_: No

_Type_: String

_Allowed Values_: <code>lucene.standard</code> | <code>lucene.simple</code> | <code>lucene.whitespace</code> | <code>lucene.keyword</code> | <code>lucene.arabic</code> | <code>lucene.armenian</code> | <code>lucene.basque</code> | <code>lucene.bengali</code> | <code>lucene.brazilian</code> | <code>lucene.bulgarian</code> | <code>lucene.catalan</code> | <code>lucene.chinese</code> | <code>lucene.cjk</code> | <code>lucene.czech</code> | <code>lucene.danish</code> | <code>lucene.dutch</code> | <code>lucene.english</code> | <code>lucene.finnish</code> | <code>lucene.french</code> | <code>lucene.galician</code> | <code>lucene.german</code> | <code>lucene.greek</code> | <code>lucene.hindi</code> | <code>lucene.hungarian</code> | <code>lucene.indonesian</code> | <code>lucene.irish</code> | <code>lucene.italian</code> | <code>lucene.japanese</code> | <code>lucene.korean</code> | <code>lucene.kuromoji</code> | <code>lucene.latvian</code> | <code>lucene.lithuanian</code> | <code>lucene.morfologik</code> | <code>lucene.nori</code> | <code>lucene.norwegian</code> | <code>lucene.persian</code> | <code>lucene.portuguese</code> | <code>lucene.romanian</code> | <code>lucene.russian</code> | <code>lucene.smartcn</code> | <code>lucene.sorani</code> | <code>lucene.spanish</code> | <code>lucene.swedish</code> | <code>lucene.thai</code> | <code>lucene.turkish</code> | <code>lucene.ukrainian</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ApiKeys

_Required_: No

_Type_: <a href="apikeydefinition.md">apiKeyDefinition</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Name

Human-readable label that identifies the synonym definition. Each **synonym.name** must be unique within the same index definition.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Source

_Required_: No

_Type_: <a href="synonymsource.md">SynonymSource</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

