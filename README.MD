# [How to set up](https://go.dev/doc/install)
Download the GO version 1.18

Create .env variable file in the root containing below variables

    - DB_MaxLifetime=1
    - DB_MaxIdleConns=5
    - DB_MaxOpenConns=5
    - SERVER_PORT=5001
    - SERVER_TIMEOUT=50

Install dependecies using below GO command

    go mod tidy


# [How to run]
First build the module to check any build issues using below command
    
    go build

You can the run modules using vscode launch included in git or using below command

    go run main.go


# Stack used
Below is the technology stack used:

    - GO
    - SQLite


# How to test API
Install vscode extention with id:humao.rest-client which we shall use to test APIa

In the source there is a folder called: "test/api" which contains API collection testable with extension installed above

You shoul be able to send and receive request with the service running
