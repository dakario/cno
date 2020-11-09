package cmd

type CnoConfig struct {
	ServerUrl        string `json:"serverUrl"`
	OidcUrl          string `json:"oidcUrl"`
	OidcClientId     string `json:"oidcClientId"`
	OidcClientSecret string `json:"oidcClientSecret"`
	AccesToken       string `json:"accesToken"`
	Organization     string `json:"organizationId"`
	GroupId          string `json:"groupId"`
	Project          string `json:"projectId"`
	Environment      string `json:"environmentId"`
}

type k8sConfig struct {
	ServerIp		string `json:"serverIp"`
	Cert       		string `json:"cert"`
	Key 			string `json:"key"`
}

type Organization struct{
	ID             	string     		`json:"id"`
	Name            string     		`json:"name"`
}

type Group struct {
	ID           		string    `json:"id"`
	Name         		string    `json:"name" validate:"required"`
	OrganizationID    	string    `json:"organizationId"`
}

type Environment struct {
	ID                    	string            `json:"id"`
	Name                  	string            `json:"name" validate:"required"`
	Status                	string            `json:"status"`
	CpuUsedPercentage     	uint              `json:"cpuUsedPercentage"`
	MemoryUsedPercentage  	uint              `json:"memoryUsedPercentage"`
	StorageUsedPercentage 	uint              `json:"storageUsedPercentage"`
	ProjectID             	string            `json:"projectId"`
	GroupID        		  	string            `json:"groupId"`
	OrganizationID        	string            `json:"organizationId"`
	AgentID				  	string            `json:"AgentId"`
	ProjectName        	  	string            `json:"projectName"`
	OrganizationName      	string            `json:"organizationName"`
}

type Project struct {
	ID                   	string                `json:"id" gorm:"primary_key"`
	Name                 	string                `json:"name" validate:"required"`
	TypeCluster          	string                `json:"typeCluster" validate:"required"`
	GroupID       		 	string                `json:"groupId"`
	OrganizationID       	string                `json:"organizationId"`
	Environments         	[]Environment         `json:"environments" gorm:"foreignkey:ProjectID"`
	UidAgent             	string                `json:"uidAgent"`
}

type ProjectPaginate struct {
	TotalRecord int         `json:"totalRecord"`
	Records     []Project 	`json:"records"`
	Limit       int         `json:"limit"`
	Page        int         `json:"page"`
}

type EnvironmentPaginate struct {
	TotalRecord int         	`json:"totalRecord"`
	Records     []Environment 	`json:"records"`
	Limit       int         	`json:"limit"`
	Page        int         	`json:"page"`
}
