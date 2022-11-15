package aws_vpc_endpoint

import (
	"fmt"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"go.mongodb.org/atlas/mongodbatlas"
)

func CreateVpcEndpoint(peCon mongodbatlas.PrivateEndpointConnection, region string, subnetId string, VpcId string) (*string, *handler.ProgressEvent) {
	mySession := session.Must(session.NewSession())

	// Create a EC2 client from just a session.
	svc := ec2.New(mySession, aws.NewConfig().WithRegion(region))

	subnetIds := []*string{&subnetId}

	vcpType := "Interface"

	connection := ec2.CreateVpcEndpointInput{
		VpcId:           &VpcId,
		ServiceName:     &peCon.EndpointServiceName,
		VpcEndpointType: &vcpType,
		SubnetIds:       subnetIds,
	}

	vpcE, err := svc.CreateVpcEndpoint(&connection)
	if err != nil {
		fpe := handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          fmt.Sprintf("Error creating vcp Endpoint: %s", err.Error()),
			HandlerErrorCode: cloudformation.HandlerErrorCodeGeneralServiceException}
		return nil, &fpe
	}

	return vpcE.VpcEndpoint.VpcEndpointId, nil
}
