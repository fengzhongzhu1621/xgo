package pflag

const (
	// OpenSourceDeployment is the open-source deployment method, do not rely on api gateway
	OpenSourceDeployment DeploymentMethod = "open_source"
	// BluekingDeployment is the deployment method for blueking, using api gateway
	BluekingDeployment DeploymentMethod = "blueking"
)

// DeploymentMethod is the deployment method
type DeploymentMethod string

// String get string value
func (d *DeploymentMethod) String() string {
	return string(*d)
}

// Set value
func (d *DeploymentMethod) Set(s string) error {
	*d = DeploymentMethod(s)
	return nil
}

// Type returns value type
func (d *DeploymentMethod) Type() string {
	return "DeploymentMethod"
}
