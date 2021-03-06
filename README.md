softlayer-go [![Build Status](https://travis-ci.org/maximilien/softlayer-go.svg?branch=master)](https://travis-ci.org/maximilien/softlayer-go#)
============

An *incomplete* SoftLayer (SL) client API written in Go language.

## Getting Started
------------------

The best way to get started would be to look at the [integration](integration) tests for creating a [virtual guest](https://github.com/maximilien/softlayer-go/blob/master/integration/virtual_guest_lifecycle/virtual_guest_lifecycle_test.go) and the [test helpers](test_helpers). Here is a snippet of what is needed.

```go
//Add necessary imports, e.g., os, slclient, datatypes
// "os"
// slclient "github.com/maximilien/softlayer-go/client"
// datatypes "github.com/maximilien/softlayer-go/data_types"

//Access SoftLayer username and API key from environment variable or hardcode here
username := os.Getenv("SL_USERNAME")
apiKey := os.Getenv("SL_API_KEY")
	
//Create a softLayer-go client
client := slclient.NewSoftLayerClient(username, apiKey)

//Get the SoftLayer account service object
accountService, err := client.GetSoftLayer_Account_Service()
if err != nil {
  return err
}

//Create a template for the virtual guest (changing properties as needed)
virtualGuestTemplate := datatypes.SoftLayer_Virtual_Guest_Template{
  Hostname:  "some-hostname",
	Domain:    "softlayergo.com",
	StartCpus: 1,
	MaxMemory: 1024,
	Datacenter: datatypes.Datacenter{
		Name: "ams01",
	},
	SshKeys:                      [],  //or get the necessary keys and add here
	HourlyBillingFlag:            true,
	LocalDiskFlag:                true,
	OperatingSystemReferenceCode: "UBUNTU_LATEST",
}
	
//Get the SoftLayer virtual guest service
virtualGuestService, err := client.GetSoftLayer_Virtual_Guest_Service()
if err != nil {
  return err
}
	
//Create the virtual guest with the service
virtualGuest, err := virtualGuestService.CreateObject(virtualGuestTemplate)
if err != nil {
	return err
}
	
//Use the virtualGuest or other services...
```

### Overview Presentations (*)
--------------------------

TBD

### Cloning and Building
------------------------

Clone this repo and build it. Using the following commands on a Linux or Mac OS X system:

```
$ mkdir -p softlayer-go/src/github.com/maximilien
$ export GOPATH=$(pwd)/softlayer-go:$GOPATH
$ cd softlayer-go/src/github.com/maximilien
$ git clone https://github.com/maximilien/softlayer-go.git
$ cd softlayer-go
$ export SL_USERNAME=your-username@your-org.com
$ export SL_API_KEY=your-softlayer-api-key
$ ./bin/build
$ ./bin/test-unit
$ ./bin/test-integration
```

NOTE: if you get any dependency errors, then use `go get path/to/dependency` to get it, e.g., `go get github.com/onsi/ginkgo` and `go get github.com/onsi/gomega`

The executable output should now be located in: `out/slgo`. It does not do anything currently, expect printing a version number. In time this may change. For now, this project is intended to be a set of useful and reusable Golang libraries to access SoftLayer.

### Running Tests
-----------------

The [SoftLayer](http://www.softlayer.com) (SL) Golang client and associated tests and binary distribution depend on you having a real SL account. Get one for free for one month [here](http://www.softlayer.com/info/free-cloud). From your SL account you can get an API key. Using your account name and API key you will need to set two environment variables: `SL_USERNAME` and `SL_API_KEY`. You can do so as follows:

```
$ export SL_USERNAME=your-username@your-org.com
$ export SL_API_KEY=your-softlayer-api-key
```

You should run the tests to make sure all is well, do this with: `$ ./bin/test-unit` and `$ ./bin/test-integration` in your cloned repository. Please note that the `$ ./bin/test-integration` will spin up real SoftLayer virtual guests (VMs) and associated resources and will also delete them. This integration test may take up to 30 minutes (usually shorter)

The output should of `$ ./bin/test-unit` be similar to:

```
$ ./bin/build
$ bin/test-unit

 Cleaning build artifacts...

 Formatting packages...

 Unit Testing packages:
[1435782974] SoftLayer Client Suite - 11/11 specs - 7 nodes ••••••••••• SUCCESS! 7.280209ms
[1435782974] Services Suite - 89/89 specs - 7 nodes ••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••• SUCCESS! 50.117056ms

Ginkgo ran 2 suites in 2.006026409s
Test Suite Passed

 Vetting packages for potential issues...

SWEET SUITE SUCCESS
```

## Developing (*)
-----------------

1. Run `go get golang.org/x/tools/cmd/vet`
2. Run `go get github.com/xxx ...` to install test dependencies (as you see errors)
1. Write a [Ginkgo](https://github.com/onsi/ginkgo) test.
1. Run `bin/test` and watch the test fail.
1. Make the test pass.
1. Submit a pull request.

## Contributing
---------------

* We gratefully acknowledge and thank the [current contributors](https://github.com/maximilien/softlayer-go/graphs/contributors)
* We welcome any and all contributions as Pull Requests (PR)
* We also welcome issues and bug report and new feature request. We will address as time permits
* Follow the steps above in Developing to get your system setup correctly
* Please make sure your PR is passing Travis before submitting
* Feel free to email me or the current collaborators if you have additional questions about contributions
* Before submitting your first PR, please read and follow steps in [CONTRIBUTING.md](CONTRIBUTING.md)

### Managing dependencies
-------------------------

* All dependencies managed via [Godep](https://github.com/tools/godep). See [Godeps/_workspace](https://github.com/maximilien/softlayer-go/tree/master/Godeps/_workspace) directory on master

#### Short `godep` Guide
* If you ever import a new package `foo/bar` (after you `go get foo/bar`, so that foo/bar is in `$GOPATH`), you can type `godep save ./...` to add it to the `Godeps` directory.
* To restore dependencies from the `Godeps` directory, simply use `godep restore`. `restore` is the opposite of `save`.
* If you ever remove a dependency or a link becomes deprecated, the easiest way is probably to remove your entire `Godeps` directory and run `godep save ./...` again, after making sure all your dependencies are in your `$GOPATH`. Don't manually edit `Godeps.json`!
* To update an existing dependency, you can use `godep update foo/bar` or `godep update foo/...` (where `...` is a wildcard)
* The godep project [readme](https://github.com/tools/godep/README.md) is a pretty good resource: [https://github.com/tools/godep](https://github.com/tools/godep)

### Current conventions
-----------------------

* Basic Go conventions
* Strict TDD for any code added or changed
* Go fakes when needing to mock objects

(*) these items are in the works, we will remove the * once they are available

**NOTE**: this client is created to support the [bosh-softlayer-cpi](https://github.com/maximilien/bosh-softlayer-cpi) project and only implements the portion of the SL APIs needed to complete the implementation of the BOSH CPI. You are welcome to use it in your own projects and as you do if you find areas we have not yet implemented but that you need, please submit [Pull Requests](https://help.github.com/articles/using-pull-requests/) or engage with us in discussions.
