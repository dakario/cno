package cmd

type CnoConfig struct {
	ServerUrl        	string `json:"serverUrl"`
	OidcUrl          	string `json:"oidcUrl"`
	OidcClientId     	string `json:"oidcClientId"`
	OidcClientSecret 	string `json:"oidcClientSecret"`
	AccesToken       	string `json:"accesToken"`
	CompanyId        	string `json:"companyId"`
	OrganizationId     	string `json:"organizationId"`
	ProjectId          	string `json:"projectId"`
	EnvironmentId      	string `json:"environmentId"`
}

type k8sConfig struct {
	ServerIp	string `json:"serverIp"`
	Cert       	string `json:"cert"`
	Key 		string `json:"key"`
}

type Company struct{
	ID             	string     		`json:"id"`
	Name            string     		`json:"name"`
}

type Organization struct {
	ID           string    `json:"id"`
	Name         string    `json:"name" validate:"required"`
	CompanyID    string    `json:"companyId"`
}

type Environment struct {
	ID                    string            `json:"id"`
	Name                  string            `json:"name" validate:"required"`
	Status                string            `json:"status"`
	CpuUsedPercentage     uint              `json:"cpuUsedPercentage"`
	MemoryUsedPercentage  uint              `json:"memoryUsedPercentage"`
	StorageUsedPercentage uint              `json:"storageUsedPercentage"`
	ProjectID             string            `json:"projectId"`
	OrganizationID        string            `json:"organizationId"`
	CompanyID             string            `json:"companyId"`

}

type Project struct {
	ID                   string                `json:"id" gorm:"primary_key"`
	Name                 string                `json:"name" validate:"required"`
	TypeCluster          string                `json:"typeCluster" validate:"required"`
	OrganizationID       string                `json:"organizationId"`
	CompanyID            string                `json:"companyId"`
	Environments         []Environment         `json:"environments" gorm:"foreignkey:ProjectID"`
	UidAgent             string                `json:"uidAgent"`
}

type ProjectPaginate struct {
	TotalRecord int         `json:"totalRecord"`
	Records     []Project `json:"records"`
	Limit       int         `json:"limit"`
	Page        int         `json:"page"`
}

type EnvironmentPaginate struct {
	TotalRecord int         `json:"totalRecord"`
	Records     []Environment `json:"records"`
	Limit       int         `json:"limit"`
	Page        int         `json:"page"`
}
