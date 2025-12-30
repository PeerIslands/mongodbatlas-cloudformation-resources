# MongoDB::Atlas::SearchDeployment

## Impact 
The following components use this resource and are potentially impacted by any changes. They should also be validated to ensure the changes do not cause a regression.
 - SearchDeployment L1 CDK constructor


## Prerequisites 
### Resources needed to run the manual QA
All resources are created as part of `cfn-testing-helper.sh`:

- Atlas Project
- Cluster

## Manual QA
Please follow the steps in [TESTING.md](../../../TESTING.md).


### Success criteria when testing the resource
1. After successful creation of the stack using a template from examples section, by seeing the details of the Cluster configuration in the Atlas UI (e.g. [https://cloud.mongodb.com/v2/<project-id>#/clusters/detail/<cluster-name>](https://cloud.mongodb.com/v2/<project-id>#/clusters/detail/<cluster-name>)) the isolated search nodes should appear together with the existing dedicated cluster nodes.

![image](https://github.com/mongodb/mongodbatlas-cloudformation-resources/assets/20469408/a08146ea-f6f2-4889-9576-b61ae97f01e7)

2. Ensure general [CFN resource success criteria](../../../TESTING.md#success-criteria-when-testing-the-resource) for this resource is met.


## Important Links
- [API Documentation](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/group/endpoint-atlas-search)
- [Resource Usage Documentation](https://www.mongodb.com/docs/atlas/cluster-config/multi-cloud-distribution/#search-nodes-for-workload-isolation)

## Running requests locally

To locally invoke requests, the AWS `sam local` and `cfn invoke` tools can be used:

```
sam local start-lambda --skip-pull-image
```
then in another shell:
```bash
repo_root=$(git rev-parse --show-toplevel)
cd ${repo_root}/cfn-resources/search-deployment
cfn invoke --function-name TestEntrypoint resource CREATE test/searchdeployment.sample-cfn-request.json 
cfn invoke --function-name TestEntrypoint resource UPDATE test/searchdeployment.sample-cfn-request.json
cfn invoke --function-name TestEntrypoint resource DELETE test/searchdeployment.sample-cfn-request.json
cd -
```

## Submitting to Private Registry

After all tests pass, submit the resource to AWS CloudFormation Private Registry:

```bash
export AWS_DEFAULT_REGION=eu-west-1
export AWS_REGION=eu-west-1
source /Users/home/repos/PeerIslands/Mongo-TF-CFN-Converter/CONVERSION_PROMPTS/setup-credentials.sh /Users/home/repos/PeerIslands/Mongo-TF-CFN-Converter/CONVERSION_PROMPTS/credPersonalCfnDev.properties
export MONGODB_ATLAS_CLUSTER_NAME='cfn-test-search-deployment-20251229'
cd /Users/home/repos/PeerIslands/Mongo-TF-CFN-Converter/mongodbatlas-cloudformation-resources/cfn-resources
LOG_FILE="search-deployment/cfn-submit-$(date +%Y%m%d-%H%M%S).log"
script -q "$LOG_FILE" bash -c './cfn-submit-helper.sh search-deployment'
```
