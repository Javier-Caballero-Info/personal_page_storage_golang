# Personal Page Storage - GoLang

> Golang Api service made with Gin Gonic framework and connecto to Amazon S3 service

## Table of Contents

-   [Overview](https://github.com/Javier-Caballero-Info/personal_page_storage_golangtree/master/README.md#overview)
-   [API Description](https://github.com/Javier-Caballero-Info/personal_page_storage_golangtree/master/README.md#api_description)
-   [Clone](https://github.com/Javier-Caballero-Info/personal_page_storage_golangtree/master/README.md#clone)
- [Requirements](https://github.com/Javier-Caballero-Info/personal_page_storage_golangtree/master#requirements)
- [Installation](https://github.com/Javier-Caballero-Info/personal_page_storage_golangtree/master#installation)
	- [GoLang](https://github.com/Javier-Caballero-Info/personal_page_storage_golangtree/master#golang)
- [Environment](https://github.com/Javier-Caballero-Info/personal_page_storage_golangtree/master#environment)
- [Developing](https://github.com/Javier-Caballero-Info/personal_page_storage_golangtree/master#developing)
- [Test](https://github.com/Javier-Caballero-Info/personal_page_storage_golangtree/master#test)
- [Build](https://github.com/Javier-Caballero-Info/personal_page_storage_golangtree/master#build)
- [Running with Docker](https://github.com/Javier-Caballero-Info/personal_page_storage_golangtree/master#running-with-docker)
	- [Building the image](https://github.com/Javier-Caballero-Info/personal_page_storage_golangtree/master#building-the-image)
	- [Starting up a container](https://github.com/Javier-Caballero-Info/personal_page_storage_golangtree/master#starting-up-a-container)
- [Contributing](https://github.com/Javier-Caballero-Info/personal_page_storage_golangtree/master#contributing)
- [Author](https://github.com/Javier-Caballero-Info/personal_page_storage_golangtree/master#author)
- [License](https://github.com/Javier-Caballero-Info/personal_page_storage_golangtree/master#license)

## Overview

This is a simple API to list, upload (create and update), and delete all types of files in AWS S3 service.
The service is very transparent, with littles setting could be connect to another object storage.

The response elements are maded thinking in the availability of get the download link easily.

This app validate the authorization token provide in Authorization header.

## API Description

For more information about the endpoints of the API please check the [apiary doc](https://personalpagestoragegolang.docs.apiary.io).

## Clone

```bash
git clone https://github.com/Javier-Caballero-Info/personal_page_storage_golang.git
git remote rm origin
git remote add origin <your-git-path>
```

## Requirements

* **GoLang:** 1.8 or above

## Installation

1. ### GoLang

    - Debian / Ubuntu

        - Ubuntu 16.04

            ```Bash
            sudo add-apt-repository ppa:longsleep/golang-backports
            ```

            ```bash
            sudo apt-get update
            sudo apt-get install golang-go
            ```

        - Ubuntu 16.10 or above

            ```bash
            sudo add-apt-repository ppa:ubuntu-lxc/lxd-stable
            ```

            ```bash
            sudo apt update
            sudo apt install golang
            ```
        - Snap

            ```Bash
            snap install --classic go
            ```
    - MacOS

        - Pkg installer

            ```bash
            curl -o go.pkg https://dl.google.com/go/go1.10.darwin-amd64.pkg
            shasum -a 256 go.pkg | grep 57510c51cb1f355f6723ac6bd7d9dd03facad474cbdb806db6ea2b616435dfdf
            sudo open go.pkg
            ```

            ```bash
            export PATH=$PATH:$GOPATH/bin
            ```

        - Brew
            ```bash
            brew install go
            ```

    - Windows

        - Installer

            Download the msi instalaller [https://golang.org/dl/](https://golang.org/dl/).

## Environment

Export the following environment variables:

```bash
PORT=3000
JWT_SIGN_ALGORITHM=HS256 # Signature to validate the JWT token
JWT_SECRET_KEY=secret # Secret key for jwt

# Firebase Credentials
AWS_REGION=us-east-1 # S3 region for the bucket
AWS_BUCKET=bucket_name # S3 bucket
AWS_BASE_PATH=secret # Path for subfolder inner the bucket, leave empty for root
AWS_ACCESS_KEY_ID=key_abc123 # Access key ID
AWS_SECRET_ACCESS_KEY=secret_abc123 # Secret access key
```

## Developing

>Setup the environment variables

After every change in the code you must stop the server and build the app again.

```
go run server.go
```

## Test

Only manual test, for now

## Build

```
go build server.go
```


## Running with Docker

To run the server on a Docker container, please execute the following from the root directory:

### Building the image
```bash
docker build -t personal_page_storage_golang .
```
### Starting up a container
```bash
docker run -p 3000:3000 -d \
-e JWT_SECRET_KEY="jwt-secret-string" \
-e JWT_SIGN_ALGORITHM="db.firebase.com" \
-e AWS_REGION="us-east-1" \
-e AWS_BUCKET="bucket_name" \
-e AWS_BASE_PATH="secret" \
-e AWS_ACCESS_KEY_ID="key_abc123" \
-e AWS_SECRET_ACCESS_KEY="secret_abc123" \
personal_page_storage_golang
```
## Contributing

Contributions welcome! See the  [Contributing Guide](https://github.com/Javier-Caballero-Info/personal_page_storage_golangblob/master/CONTRIBUTING.md).

## Author

Created and maintained by [Javier Hernán Caballero García](https://javiercaballero.info)).

## License

GNU General Public License v3.0

See  [LICENSE](https://github.com/Javier-Caballero-Info/personal_page_storage_golangblob/master/LICENSE)