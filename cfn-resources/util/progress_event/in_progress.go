package progress_events

import "github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"

func GetInProgressProgressEvent(message string, callBackContext map[string]interface{}) handler.ProgressEvent {
	return handler.ProgressEvent{
		OperationStatus:      handler.InProgress,
		Message:              message,
		CallbackDelaySeconds: 10,
		CallbackContext:      callBackContext}
}
