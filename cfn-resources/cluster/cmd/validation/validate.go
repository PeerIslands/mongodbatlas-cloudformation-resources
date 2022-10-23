package validation

type ModelValidator struct{}

var CreateRequiredFields = []string{"ApiKeys.PublicKey", "ApiKeys.PrivateKey", "Name", "ProjectId"}
var ReadRequiredFields = []string{"ApiKeys.PublicKey", "ApiKeys.PrivateKey", "Id"}
var UpdateRequiredFields = []string{"ApiKeys.PublicKey", "ApiKeys.PrivateKey", "Name", "ProjectId"}
var DeleteRequiredFields = []string{"ApiKeys.PublicKey", "ApiKeys.PrivateKey", "Name", "ProjectId"}
var ListRequiredFields = []string{"ApiKeys.PublicKey", "ApiKeys.PrivateKey", "ProjectId"}

func (m ModelValidator) GetCreateFields() []string {
	return CreateRequiredFields
}
func (m ModelValidator) GetReadFields() []string {
	return ReadRequiredFields
}
func (m ModelValidator) GetUpdateFields() []string {
	return UpdateRequiredFields
}
func (m ModelValidator) GetDeleteFields() []string {
	return DeleteRequiredFields
}
func (m ModelValidator) GetListFields() []string {
	return ListRequiredFields
}
