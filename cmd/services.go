package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Nerzal/gocloak/v6"
	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"io/ioutil"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)


func LoadCnoConfig() (*CnoConfig, error){
	user, err := user.Current()
	if err != nil {
		return nil, err
	}
	workspace := filepath.Join(user.HomeDir, "/.cno")
	if _, err := os.Stat(filepath.Join(workspace, "/config")); os.IsNotExist(err) {
		fmt.Println("No cno CnoConfig file found! You have to run: cno config")
	}
	file, err := os.Open(filepath.Join(workspace, "/config"))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	config := &CnoConfig{}
	if err := json.Unmarshal(data, config); err!=nil{
		return nil, err
	}
	return config, nil
}

func RefreshToken(config *CnoConfig) error {
	prompt := promptui.Prompt{
		Label:    "username",
	}
	var err error = nil
	username, err := prompt.Run()
	if err != nil {
		return err
	}
	prompt.Label = "password"
	prompt.Mask = 1
	password, err := prompt.Run()
	if err != nil {
		return err
	}
	var token *gocloak.JWT
	client := gocloak.NewClient(config.OidcUrl)
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	token, err = client.Login(context.TODO(), config.OidcClientId, config.OidcClientSecret, "onboarding", username, password)
	s.Stop()
	if err != nil {
		return err
	}
	config.AccesToken = token.AccessToken

	return SaveConfigOnFileSystem(*config)
}


func GenerateKubeConfig(agentId, organizationId, defaultNamespace string) error{
	config , err := LoadCnoConfig()
	if err != nil {
		return err
	}
	client := &http.Client{}
	isTokenRefresed := false
loop:
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	req, _ := http.NewRequest("GET", config.ServerUrl+"/api/v1/agent/k8s-credentials/"+agentId+"/organization/"+organizationId, nil)
	req.Header.Add("Authorization", "Bearer "+config.AccesToken)
	res, err := client.Do(req)
	s.Stop()
	if err!=nil {
		return err
	}
	if  res.StatusCode == http.StatusUnauthorized {
		if isTokenRefresed {
			return errors.New("Token invalid. Maybe your cli does not use the same sso as the cno api")
		}
		fmt.Println("The token is expired or not correct. You Have to login again!")
		err := RefreshToken(config)
		if err != nil {
			return err
		}
		isTokenRefresed = true
		goto loop
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("GET k8s-credentials: "+string(response))
	}
	k8sConfig := k8sConfig{}
	if err = json.Unmarshal(response, &k8sConfig); err != nil {
		return err
	}
	decodedCert, err := base64.StdEncoding.DecodeString(k8sConfig.Cert)
	if err!=nil{
		return err
	}
	decodedKey, err := base64.StdEncoding.DecodeString(k8sConfig.Key)
	if  err!=nil{
		return err
	}
	k8sConfig.Cert = string(decodedCert)
	k8sConfig.Key = string(decodedKey)
	cnoCluster := api.NewCluster()
	cnoCluster.InsecureSkipTLSVerify = true
	cnoCluster.Server = k8sConfig.ServerIp

	cnoAuthInfo := api.NewAuthInfo()
	cnoAuthInfo.ClientCertificateData = decodedCert
	cnoAuthInfo.ClientKeyData = decodedKey

	cnoContext := api.NewContext()
	cnoContext.Namespace = defaultNamespace
	cnoContext.Cluster 	 = "cno"
	cnoContext.AuthInfo  = "cno"


	kubeConfig := *clientcmd.GetConfigFromFileOrDie("/Users/user/.kube/config")
	kubeConfig.CurrentContext = "cno"
	kubeConfig.Contexts["cno"] =  cnoContext
	kubeConfig.Clusters["cno"] =  cnoCluster
	kubeConfig.AuthInfos["cno"] =  cnoAuthInfo

	configAccess := clientcmd.NewDefaultClientConfig(kubeConfig, nil).ConfigAccess()

	err1 := clientcmd.ModifyConfig(configAccess, kubeConfig, true)
	if err1 != nil {
		return err1
	}
	return nil
}


func chooseProject() (*Project, error){
	cnoConfig, err := LoadCnoConfig()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	isTokenRefresed := false
loop:
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	req, _ := http.NewRequest("GET", cnoConfig.ServerUrl+"/api/v1/projects/organization/name/"+cnoConfig.Organization+"/user/me/", nil)
	req.Header.Add("Authorization", "Bearer "+cnoConfig.AccesToken)
	res, err := client.Do(req)
	s.Stop()
	if err!=nil {
		return nil, err
	}
	if  res.StatusCode == http.StatusUnauthorized {
		if isTokenRefresed {
			return nil, errors.New("Token invalid. Maybe your cli does not use the same sso as the cno api")
		}
		fmt.Println("The token is expired or not correct. You Have to login again!")
		err := RefreshToken(cnoConfig)
		if err != nil {
			return nil, err
		}
		isTokenRefresed = true
		goto loop
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("GET all project: "+string(response))
	}

	var projects []Project
	if err := json.Unmarshal(response, &projects); err!=nil {
		return nil, err
	}
	if len(projects)==0 {
		return nil, errors.New("You are not a member of any project  of this organization!")
	}
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   promptui.Styler(promptui.FGYellow)("▸")+" {{ .Name }}",
		Inactive: "  {{ .Name }}",
		Selected: promptui.IconGood+"  {{ .Name }}",

	}

	searcher := func(input string, index int) bool {
		project := projects[index]
		name := strings.Replace(strings.ToLower(project.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Select a project",
		Items:     projects,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	selectedGroup := projects[i]
	return &selectedGroup, nil
}

func getProject(projectName, orgName string) (*Project, error){
	cnoConfig, err := LoadCnoConfig()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	isTokenRefresed := false
loop:
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	req, _ := http.NewRequest("GET", cnoConfig.ServerUrl+"/api/v1/projects/name/"+projectName+"/organization/name/"+orgName, nil)
	req.Header.Add("Authorization", "Bearer "+cnoConfig.AccesToken)
	res, err := client.Do(req)
	s.Stop()
	if err!=nil {
		return nil, err
	}
	if  res.StatusCode == http.StatusUnauthorized {
		if isTokenRefresed {
			return nil, errors.New("Token invalid. Maybe your cli does not use the same sso as the cno api")
		}
		if len(cnoConfig.AccesToken)>1 {
			fmt.Println("The token is expired or not correct. You Have to login again!")
		}
		err := RefreshToken(cnoConfig)
		if err != nil {
			return nil, err
		}
		isTokenRefresed = true
		goto loop
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("GET project: "+string(response))
	}

	var project Project
	if err := json.Unmarshal(response, &project); err!=nil {
		return nil, err
	}

	return &project, nil
}

func getAllProjects() ([]Project, error){
	cnoConfig, err := LoadCnoConfig()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	isTokenRefresed := false
loop:
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	req, _ := http.NewRequest("GET", cnoConfig.ServerUrl+"/api/v1/projects/organization/name/"+cnoConfig.Organization+"/user/me", nil)
	req.Header.Add("Authorization", "Bearer "+cnoConfig.AccesToken)
	res, err := client.Do(req)
	s.Stop()
	if err!=nil {
		return nil, err
	}
	if  res.StatusCode == http.StatusUnauthorized {
		if isTokenRefresed {
			return nil, errors.New("Token invalid. Maybe your cli does not use the same sso as the cno api")
		}
		if len(cnoConfig.AccesToken)>1 {
			fmt.Println("The token is expired or not correct. You Have to login again!")
		}
		err := RefreshToken(cnoConfig)
		if err != nil {
			return nil, err
		}
		isTokenRefresed = true
		goto loop
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("GET all projects: "+string(response))
	}

	var projects []Project
	if err := json.Unmarshal(response, &projects); err!=nil {
		return nil, err
	}

	return projects, nil
}

func getEnvByProject(projectID string) ([]Environment, error){
	cnoConfig, err := LoadCnoConfig()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	isTokenRefresed := false
loop:
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	req, _ := http.NewRequest("GET", cnoConfig.ServerUrl+"/api/v1/environments/project/name/"+projectID+"/organization/name/"+cnoConfig.Organization, nil)
	req.Header.Add("Authorization", "Bearer "+cnoConfig.AccesToken)
	res, err := client.Do(req)
	s.Stop()
	if err!=nil {
		return nil, err
	}
	if  res.StatusCode == http.StatusUnauthorized {
		if isTokenRefresed {
			return nil, errors.New("Token invalid. Maybe your cli does not use the same sso as the cno api")
		}
		if len(cnoConfig.AccesToken)>1 {
			fmt.Println("The token is expired or not correct. You Have to login again!")
		}
		err := RefreshToken(cnoConfig)
		if err != nil {
			return nil, err
		}
		isTokenRefresed = true
		goto loop
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("GET environment by project: "+string(response))
	}

	var envs []Environment
	if err := json.Unmarshal(response, &envs); err!=nil {
		return nil, err
	}

	return envs, nil
}


func chooseEnv(projectName, organizationName string) (*Environment, error){
	cnoConfig, err := LoadCnoConfig()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	isTokenRefresed := false
loop:
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	req, _ := http.NewRequest("GET", cnoConfig.ServerUrl+"/api/v1/environments/project/name/"+projectName+"/organization/name/"+organizationName, nil)
	req.Header.Add("Authorization", "Bearer "+cnoConfig.AccesToken)
	res, err := client.Do(req)
	s.Stop()
	if err!=nil {
		return nil, err
	}
	if  res.StatusCode == http.StatusUnauthorized {
		if isTokenRefresed {
			return nil, errors.New("Token invalid. Maybe your cli does not use the same sso as the cno api")
		}
		fmt.Println("The token is expired or not correct. You Have to login again!")
		err := RefreshToken(cnoConfig)
		if err != nil {
			return nil, err
		}
		isTokenRefresed = true
		goto loop
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("GET environment for project: "+string(response))
	}

	var envs []Environment
	if err := json.Unmarshal(response, &envs); err!=nil {
		return nil, err
	}
	if len(envs) == 0 {
		return nil, errors.New("You have no access to any environments of this project!")
	}
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   promptui.Styler(promptui.FGYellow)("▸")+" {{ .Name }}",
		Inactive: "  {{ .Name }}",
		Selected: promptui.IconGood+"  {{ .Name }}",

	}

	searcher := func(input string, index int) bool {
		env := envs[index]
		name := strings.Replace(strings.ToLower(env.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Select an environment",
		Items:     envs,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	selectedEnv := envs[i]
	return &selectedEnv, nil
}

func getEnvByName(envName, projectName, organizationName string) (*Environment, error){
	cnoConfig, err := LoadCnoConfig()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	isTokenRefresed := false
loop:
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	req, _ := http.NewRequest("GET", cnoConfig.ServerUrl+"/api/v1/environments/name/"+envName+"/project/name/"+projectName+"/organization/name/"+organizationName, nil)
	req.Header.Add("Authorization", "Bearer "+cnoConfig.AccesToken)
	res, err := client.Do(req)
	s.Stop()
	if err!=nil {
		return nil, err
	}
	if  res.StatusCode == http.StatusUnauthorized {
		if isTokenRefresed {
			return nil, errors.New("Token invalid. Maybe your cli does not use the same sso as the cno api")
		}
		fmt.Println("The token is expired or not correct. You Have to login again!")
		err := RefreshToken(cnoConfig)
		if err != nil {
			return nil, err
		}
		isTokenRefresed = true
		goto loop
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, errors.New("environment not found: "+string(response))
	}else if res.StatusCode != http.StatusOK {
		return nil, errors.New("GET environment: "+string(response))
	}

	var env Environment
	if err := json.Unmarshal(response, &env); err!=nil {
		return nil, err
	}
	return &env, nil
}