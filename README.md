# User Guide
## Install

    go get -u github.com/dakario/cno
    
###### NB:
Make sure you have $GOPATH/bin in your path

## Commandes
### config
This command allows you to configure the cno cli so that it can communicate correctly with the sso and the cno api

Flags:

|        name          |        description                             |           value example          |
|----------------------|------------------------------------------------|----------------------------------|
| server-url           |   application api URL                          | https://cno.beopenit.com         |
| oidc-url             |   sso server url                               | https://sso.beopenit.com         |   
| oidc-client-id       |   oidc client-id                               | public                           |
| oidc-client-secret   |   secret of the oidc client-id (maybe empty)   |                                  |


Use:
    
    cno config --server-url=https://cno.beopenit.com --oidc-url=https://sso.beopenit.com --oidcclient-id=public --oidc-client-secret=""


NB: If an flags is not setted, the cli will invite you to enter his value 

### select
This command have two other subcommands: project and env

#### select project
This command allows you to have a valid kubeconfig allowing you to interact with the cluster on which the project is deployed.
The generated kubeconfig contains a certificate with your username as CN signed by the cluster k8s.
Which will allow the k8s cluster to identify you. 

Flags:

|        name        |        description                      |
|--------------------|-----------------------------------------|
| project-id         |   id of the project you want to select  |

Use:

    cno select project --project-id <your-project-id>
    
NB: If project-id flag not set, the cli will invite you to select first the company and the organization where your project is located and then the project as such.

#### select env
This command allows you to configure an environment of the selected project as the default namespace of you kubeconfig.

Use:

    cno select env --env-id <your-env-id>
    
NB: If env-id flag not set, the cli will invite you to select an environment from the list of project environments to which you have access

#Developer's Guide

##Create and publish a new release
1. Create a new tag and publish them

        git tag -a v0.1.0 -m "First release"
        git push origin v0.1.0

2. Create a github access token and export GITHUB_TOKEN variable

        export GITHUB_TOKEN=<acces-token>
        
3. Create and publish the release

        goreleaser --rm-dist