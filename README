Anonymoe Web Email
=====================

![http://anony.moe/](https://github.com/philtyl/anonymoe/blob/master/public/img/logo.png?raw=true)

##### Current tip version: [`VERSION`](conf/VERSION) (see [Releases](https://github.com/philtyl/anonymoe/releases) for previous releases)

### Important Notes

1. The site [anony.moe](http://anony.moe/) is running under `master` branch.
2. If you think there are vulnerabilities in the project, please talk privately to **hostmaster@anony.moe**, and the name you want to be credited as. Thanks!

## Purpose

The goal of the Anonymoe mail server is to provide an no-identity and public box to receive mail to and retrieve through a webapp or REST API.

## Overview

- The Anonymoe project is hosted [online](http://anony.moe) and is free to use!
- Contribution, collaboration and suggestions are greatly welcomed

## Features

- Receives SMTP mail messages for any account under the hosting domain
- Fast and responsive WebUI for browsing mailboxes
- **WIP REST API for pragmatically pulling emails received to mailbox

## Browser Support

- Please see [Semantic UI](https://github.com/Semantic-Org/Semantic-UI#browser-support) for specific versions of supported browsers.

## Installation

#### Prerequisites

Building from source will require the following applications:
- go
- go-bindata
- go-sqlite3
- lessc

Running the application requires the following applications:
- go
- sqlite3

#### Building from Source

Build the anonymoe binary with: `make all`

#### Controlling Configuration Location

Control the installation location of application configurations and mail database with the `ANONY_CONFIG` environment variable.  This should point to a writable directory.  The default directory will be where the anonymoe binary is installed, usually `$GOPATH/bin/anonymoe-data/` Example:

```
ANONY_CONFIG=/mnt/storage/anonymoe/data/
```

#### Initializing Anonymoe Configurations

To create the initial configuration and database files, use the following command.  Note that if the configuration and database files already exist you will need to append `--force` to the command to overwrite those files.

```
anonymoe install --init
```

#### Configuring Anonymoe

Application configuration file will be located at `$ANONY_CONFIG/app.ini` and should be modified to the environment before starting the application.

#### Starting Anonymoe

To start Anonymoe server:

```
anonymoe server --start
```

## License

This project is under the MIT License. See the [LICENSE](https://github.com/philtyl/anonymoe/blob/master/LICENSE) file for the full license text.