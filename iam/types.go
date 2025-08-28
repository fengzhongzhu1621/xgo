package iam

type IamAuthOperation string

const (
	// IamGrantOperation TODO
	IamGrantOperation = "grant"
	// IamRevokeOperation TODO
	IamRevokeOperation = "revoke"
)

type IamPermission struct {
	SystemID   string      `json:"system_id"`
	SystemName string      `json:"system_name"`
	Actions    []IamAction `json:"actions"`
}

type IamAction struct {
	ID                   string            `json:"id"`
	Name                 string            `json:"name"`
	RelatedResourceTypes []IamResourceType `json:"related_resource_types"`
}

type IamResourceType struct {
	SystemID   string                  `json:"system_id"`
	SystemName string                  `json:"system_name"`
	Type       string                  `json:"type"`
	TypeName   string                  `json:"type_name"`
	Instances  [][]IamResourceInstance `json:"instances,omitempty"`
	Attributes []IamResourceAttribute  `json:"attributes,omitempty"`
}

type IamResourceInstance struct {
	Type     string `json:"type"`
	TypeName string `json:"type_name"`
	ID       string `json:"id"`
	Name     string `json:"name"`
}

type IamResourceAttribute struct {
	ID     string                      `json:"id"`
	Values []IamResourceAttributeValue `json:"values"`
}

type IamResourceAttributeValue struct {
	ID string `json:"id"`
}

type IamInstanceWithCreator struct {
	System    string                `json:"system"`
	Type      string                `json:"type"`
	ID        string                `json:"id"`
	Name      string                `json:"name"`
	Creator   string                `json:"creator"`
	Ancestors []IamInstanceAncestor `json:"ancestors,omitempty"`
}

type IamInstances struct {
	System    string        `json:"system"`
	Type      string        `json:"type"`
	Instances []IamInstance `json:"instances"`
}

type IamInstancesWithCreator struct {
	IamInstances `json:",inline"`
	Creator      string `json:"creator"`
}

type IamInstance struct {
	ID        string                `json:"id"`
	Name      string                `json:"name"`
	Ancestors []IamInstanceAncestor `json:"ancestors,omitempty"`
}

type IamInstanceAncestor struct {
	System string `json:"system"`
	Type   string `json:"type"`
	ID     string `json:"id"`
}

type IamCreatorActionPolicy struct {
	Action   ActionWithID `json:"action"`
	PolicyID int64        `json:"policy_id"`
}

type ActionWithID struct {
	ID string `json:"id"`
}

type IamBatchOperateInstanceAuthReq struct {
	Asynchronous bool             `json:"asynchronous"`
	Operate      IamAuthOperation `json:"operate"`
	System       string           `json:"system"`
	Actions      []ActionWithID   `json:"actions"`
	Subject      IamSubject       `json:"subject"`
	Resources    []IamInstances   `json:"resources"`
	ExpiredAt    int64            `json:"expired_at"`
}

type IamSubject struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type IamBatchOperateInstanceAuthRes struct {
	Action   ActionWithID `json:"action"`
	PolicyID int64        `json:"policy_id"`
}

type Permission struct {
	SystemID      string `json:"system_id"`
	SystemName    string `json:"system_name"`
	ScopeType     string `json:"scope_type"`
	ScopeTypeName string `json:"scope_type_name"`
	ScopeID       string `json:"scope_id"`
	ScopeName     string `json:"scope_name"`
	ActionID      string `json:"action_id"`
	ActionName    string `json:"action_name"`
	// newly added two field.
	ResourceTypeName string `json:"resource_type_name"`
	ResourceType     string `json:"resource_type"`

	Resources [][]Resource `json:"resources"`
}

type Resource struct {
	ResourceTypeName string `json:"resource_type_name"`
	ResourceType     string `json:"resource_type"`
	ResourceName     string `json:"resource_name"`
	ResourceID       string `json:"resource_id"`
}
