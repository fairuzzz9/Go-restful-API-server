# Track and Trace backend

**Track and Trace** is the backend portion of the Tracking Page UI. It will reside at AWS and invoke through AWS API Gateway.

# Contact

NG BOO JIUN - boojiun@pos.com.my

Project Link: [https://gitlab.com/pos_malaysia/go-skeletons/go-skeleton-rest-app](https://gitlab.com/pos_malaysia/go-skeletons/go-skeleton-rest-app)

# Getting Started

Please install Golang first before running this skeleton program. 

To build the binary, run the following script from the root directory.
```
./scripts/build_server.sh
```

If everything goes well, you should get the binary in the ./bin directory. 

Run 
```
./bin/server
```


# Testing

To test all the packages, run the test script from the root directory.


```
./scripts/run_all_tests.sh
```



# Folders and files 

## Description of main.go

The main.go file is divided into 3 parts - setup, initialize and wait for shutdown signals.

* Setup Echo server. 

* Initialize the Echo’s routes.

* Echo server will be started by an anonymous goroutine. The **Graceful()** function will listen to any interrupt or kill signal and shutdown the server gracefully. 

## Description of setup.go

* SQL statements that will be used by "gitlab.com/pos_malaysia/golib/database" get initialized here.
  
* Setup Echo to use "gitlab.com/pos_malaysia/golib/logs".  The Zerolog logger is adapted into the Echo web framework through the **lecho** adapter. 

## Description of /pkg/handlers directory

One handler per file instead of putting every handlers into a single file.  
For example, Home handler will be in home.go and HealthCheck handler codes go into healthcheck.go


## Description of /pkg/http/response directory


Where custom response messages are defined.

## Description of /pkg/utilities directory

Place for common utilities functions.


## Description of /cmd directory

Location of main.go and setup.go files. Try to keep main.go as small as possible and place initialization codes into setup.go

## Description of /internal/routes directory

Where routes to handlers relationships are organized.


## Description of /scripts directory

Scripts to perform various test, build, install, analysis, etc operations.


# Sample application 

Please refer to https://gitlab.com/pos_malaysia/go-recipes/sample-database-application for sample application that uses this skeleton.
