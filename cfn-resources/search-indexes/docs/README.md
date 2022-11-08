# Mongodb::Atlas::indexes

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "Type" : "Mongodb::Atlas::indexes",
    "Properties" : {
        "<a href="#analyzer" title="Analyzer">Analyzer</a>" : <i>String</i>,
        "<a href="#analyzers" title="Analyzers">Analyzers</a>" : <i>[ <a href="apiatlasftsanalyzersviewmanual.md">ApiAtlasFTSAnalyzersViewManual</a>, ... ]</i>,
        "<a href="#apikeys" title="ApiKeys">ApiKeys</a>" : <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>,
        "<a href="#clustername" title="ClusterName">ClusterName</a>" : <i>String</i>,
        "<a href="#collectionname" title="CollectionName">CollectionName</a>" : <i>String</i>,
        "<a href="#database" title="Database">Database</a>" : <i>String</i>,
        "<a href="#databasename" title="DatabaseName">DatabaseName</a>" : <i>String</i>,
        "<a href="#groupid" title="GroupId">GroupId</a>" : <i>String</i>,
        "<a href="#indexid" title="IndexId">IndexId</a>" : <i>String</i>,
        "<a href="#mappings" title="Mappings">Mappings</a>" : <i><a href="apiatlasftsmappingsviewmanual.md">ApiAtlasFTSMappingsViewManual</a></i>,
        "<a href="#name" title="Name">Name</a>" : <i>String</i>,
        "<a href="#searchanalyzer" title="SearchAnalyzer">SearchAnalyzer</a>" : <i>String</i>,
        "<a href="#synonyms" title="Synonyms">Synonyms</a>" : <i>[ <a href="apiatlasftssynonymmappingdefinitionview.md">ApiAtlasFTSSynonymMappingDefinitionView</a>, ... ]</i>
    }
}
</pre>

### YAML

<pre>
Type: Mongodb::Atlas::indexes
Properties:
    <a href="#analyzer" title="Analyzer">Analyzer</a>: <i>String</i>
    <a href="#analyzers" title="Analyzers">Analyzers</a>: <i>
      - <a href="apiatlasftsanalyzersviewmanual.md">ApiAtlasFTSAnalyzersViewManual</a></i>
    <a href="#apikeys" title="ApiKeys">ApiKeys</a>: <i><a href="apikeydefinition.md">apiKeyDefinition</a></i>
    <a href="#clustername" title="ClusterName">ClusterName</a>: <i>String</i>
    <a href="#collectionname" title="CollectionName">CollectionName</a>: <i>String</i>
    <a href="#database" title="Database">Database</a>: <i>String</i>
    <a href="#databasename" title="DatabaseName">DatabaseName</a>: <i>String</i>
    <a href="#groupid" title="GroupId">GroupId</a>: <i>String</i>
    <a href="#indexid" title="IndexId">IndexId</a>: <i>String</i>
    <a href="#mappings" title="Mappings">Mappings</a>: <i><a href="apiatlasftsmappingsviewmanual.md">ApiAtlasFTSMappingsViewManual</a></i>
    <a href="#name" title="Name">Name</a>: <i>String</i>
    <a href="#searchanalyzer" title="SearchAnalyzer">SearchAnalyzer</a>: <i>String</i>
    <a href="#synonyms" title="Synonyms">Synonyms</a>: <i>
      - <a href="apiatlasftssynonymmappingdefinitionview.md">ApiAtlasFTSSynonymMappingDefinitionView</a></i>
</pre>

## Properties

#### Analyzer

Specific pre-defined method chosen to convert database field text into searchable words. This conversion reduces the text of fields into the smallest units of text. These units are called a **term** or **token**. This process, known as tokenization, involves a variety of changes made to the text in fields:

- extracting words
- removing punctuation
- removing accents
- changing to lowercase
- removing common words
- reducing words to their root form (stemming)
- changing words to their base form (lemmatization)
 MongoDB Cloud uses the selected process to build the Atlas Search index.

_Required_: No

_Type_: String

_Allowed Values_: <code>lucene.standard</code> | <code>lucene.simple</code> | <code>lucene.whitespace</code> | <code>lucene.keyword</code> | <code>lucene.arabic</code> | <code>lucene.armenian</code> | <code>lucene.basque</code> | <code>lucene.bengali</code> | <code>lucene.brazilian</code> | <code>lucene.bulgarian</code> | <code>lucene.catalan</code> | <code>lucene.chinese</code> | <code>lucene.cjk</code> | <code>lucene.czech</code> | <code>lucene.danish</code> | <code>lucene.dutch</code> | <code>lucene.english</code> | <code>lucene.finnish</code> | <code>lucene.french</code> | <code>lucene.galician</code> | <code>lucene.german</code> | <code>lucene.greek</code> | <code>lucene.hindi</code> | <code>lucene.hungarian</code> | <code>lucene.indonesian</code> | <code>lucene.irish</code> | <code>lucene.italian</code> | <code>lucene.japanese</code> | <code>lucene.korean</code> | <code>lucene.kuromoji</code> | <code>lucene.latvian</code> | <code>lucene.lithuanian</code> | <code>lucene.morfologik</code> | <code>lucene.nori</code> | <code>lucene.norwegian</code> | <code>lucene.persian</code> | <code>lucene.portuguese</code> | <code>lucene.romanian</code> | <code>lucene.russian</code> | <code>lucene.smartcn</code> | <code>lucene.sorani</code> | <code>lucene.spanish</code> | <code>lucene.swedish</code> | <code>lucene.thai</code> | <code>lucene.turkish</code> | <code>lucene.ukrainian</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Analyzers

List of user-defined methods to convert database field text into searchable words.

_Required_: No

_Type_: List of <a href="apiatlasftsanalyzersviewmanual.md">ApiAtlasFTSAnalyzersViewManual</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ApiKeys

_Required_: No

_Type_: <a href="apikeydefinition.md">apiKeyDefinition</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ClusterName

Name of the cluster that contains the database and collection with one or more Application Search indexes.

_Required_: No

_Type_: String

_Minimum_: <code>1</code>

_Maximum_: <code>64</code>

_Pattern_: <code>^([a-zA-Z0-9]([a-zA-Z0-9-]){0,21}(?<!-)([\w]{0,42}))$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### CollectionName

Name of the collection that contains one or more Atlas Search indexes.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Database

Human-readable label that identifies the database that contains the collection with one or more Atlas Search indexes.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### DatabaseName

Human-readable label that identifies the database that contains the collection with one or more Atlas Search indexes.

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

#### IndexId

Unique 24-hexadecimal digit string that identifies the Atlas Search index. Use the [Get All Atlas Search Indexes for a Collection API](https://docs.atlas.mongodb.com/reference/api/fts-indexes-get-all/) endpoint to find the IDs of all Atlas Search indexes.

_Required_: No

_Type_: String

_Minimum_: <code>24</code>

_Maximum_: <code>24</code>

_Pattern_: <code>^([a-f0-9]{24})$</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Mappings

_Required_: No

_Type_: <a href="apiatlasftsmappingsviewmanual.md">ApiAtlasFTSMappingsViewManual</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Name

Human-readable label that identifies this index. Within each namespace, names of all indexes in the namespace must be unique.

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### SearchAnalyzer

Method applied to identify words when searching this index.

_Required_: No

_Type_: String

_Allowed Values_: <code>lucene.standard</code> | <code>lucene.simple</code> | <code>lucene.whitespace</code> | <code>lucene.keyword</code> | <code>lucene.arabic</code> | <code>lucene.armenian</code> | <code>lucene.basque</code> | <code>lucene.bengali</code> | <code>lucene.brazilian</code> | <code>lucene.bulgarian</code> | <code>lucene.catalan</code> | <code>lucene.chinese</code> | <code>lucene.cjk</code> | <code>lucene.czech</code> | <code>lucene.danish</code> | <code>lucene.dutch</code> | <code>lucene.english</code> | <code>lucene.finnish</code> | <code>lucene.french</code> | <code>lucene.galician</code> | <code>lucene.german</code> | <code>lucene.greek</code> | <code>lucene.hindi</code> | <code>lucene.hungarian</code> | <code>lucene.indonesian</code> | <code>lucene.irish</code> | <code>lucene.italian</code> | <code>lucene.japanese</code> | <code>lucene.korean</code> | <code>lucene.kuromoji</code> | <code>lucene.latvian</code> | <code>lucene.lithuanian</code> | <code>lucene.morfologik</code> | <code>lucene.nori</code> | <code>lucene.norwegian</code> | <code>lucene.persian</code> | <code>lucene.portuguese</code> | <code>lucene.romanian</code> | <code>lucene.russian</code> | <code>lucene.smartcn</code> | <code>lucene.sorani</code> | <code>lucene.spanish</code> | <code>lucene.swedish</code> | <code>lucene.thai</code> | <code>lucene.turkish</code> | <code>lucene.ukrainian</code>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Synonyms

Rule sets that map words to their synonyms in this index.

_Required_: No

_Type_: List of <a href="apiatlasftssynonymmappingdefinitionview.md">ApiAtlasFTSSynonymMappingDefinitionView</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

## Return Values

### Fn::GetAtt

The `Fn::GetAtt` intrinsic function returns a value for a specified attribute of this type. The following are the available attributes and sample return values.

For more information about using the `Fn::GetAtt` intrinsic function, see [Fn::GetAtt](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html).

#### IndexID

Unique 24-hexadecimal digit string that identifies this Atlas Search index.

#### Status

Condition of the search index when you made this request.

| Status | Index Condition |
 |---|---|
 | IN_PROGRESS | Atlas is building or re-building the index after an edit. |
 | STEADY | You can use this search index. |
 | FAILED | Atlas could not build the index. |
 | MIGRATING | Atlas is upgrading the underlying cluster tier and migrating indexes. |


