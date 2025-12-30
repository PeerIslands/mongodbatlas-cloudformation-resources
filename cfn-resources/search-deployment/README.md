# MongoDB::Atlas::SearchDeployment

## Description

Resource for managing [Search Nodes](https://www.mongodb.com/docs/atlas/cluster-config/multi-cloud-distribution/#search-nodes-for-workload-isolation).

## Requirements

Set up an AWS profile to securely give CloudFormation access to your Atlas credentials.
For instructions on setting up a profile, [see here](/README.md#mongodb-atlas-api-keys-credential-management).

## Attributes and Parameters

See the [resource docs](./docs/README.md).

## CloudFormation Examples

See the examples [CFN Template](/examples/search-deployment/search-deployment.json) for example resource.

## Submitting to Private Registry

To submit this resource to AWS CloudFormation Private Registry:

```bash
export AWS_DEFAULT_REGION=eu-west-1
export AWS_REGION=eu-west-1
source /Users/home/repos/PeerIslands/Mongo-TF-CFN-Converter/CONVERSION_PROMPTS/setup-credentials.sh /Users/home/repos/PeerIslands/Mongo-TF-CFN-Converter/CONVERSION_PROMPTS/credPersonalCfnDev.properties
export MONGODB_ATLAS_CLUSTER_NAME='cfn-test-search-deployment-20251229'
cd /Users/home/repos/PeerIslands/Mongo-TF-CFN-Converter/mongodbatlas-cloudformation-resources/cfn-resources
LOG_FILE="search-deployment/cfn-submit-$(date +%Y%m%d-%H%M%S).log"
script -q "$LOG_FILE" bash -c './cfn-submit-helper.sh search-deployment'
```

**Note**: Ensure all contract tests pass before submitting. Run `./cfn-testing-helper.sh search-deployment` first.
