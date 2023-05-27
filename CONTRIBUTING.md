# Contributing

For general contribution and community guidelines, please see the [community repo](https://github.com/cyberark/community).

## Development

Before getting started, you should install some developer tools.
These are not required to deploy the Conjur Service Broker but
they will let you develop using a standardized, expertly configured
environment.

1.  [git][get-git] to manage source code
2.  [Docker][get-docker] to manage dependencies and runtime environments
3.  [Tilt][get-tilt] to orchestrate Docker environments

To test the usage of the Conjur Service Broker within a CF deployment, you can
follow the demo scripts in the [Cloud Foundry demo repo](https://github.com/conjurinc/cloudfoundry-conjur-demo).

## Development Environment

The `Tiltfile` configuration file sets up a development environment that allows you
to selectively run unit and integration tests interactively against local,
containerized instances of the Conjur Service Broker and Conjur.

In this development environment, the Service Broker source code is
volume mounted in the Service Broker instances, so that any changes that
you make to Service Broker code is immediately reflected in the
Service Broker instances. In other words, there is no need to rebuild
and restart containers when code changes are made.

To start the Service Broker development environment, simply run:

```sh-session
tilt up
```

After starting up Service Broker and Conjur container instances, tilt builds project and runs unit and integration tests.
You can rerun any step such as build, test or unit tests from the Tilt dashboard.

## Non-Interactive Testing

### Running Unit Tests

To run the Conjur Service Broker unit tests, first deploy the app using tilt:

```sh-session
tilt up
```

Then, unit tests will execute automatically. If you make any changes to the code, it will automatically update.
Then re-run test\_in\_docker in the dashboard.

### Running Local Integration Tests

The [test/integration/main\_test.go](./test/integration/main_test.go) file provides a full
suite of integration tests for testing Service Broker functionality
against Conjur. To run these test application must be deployed using Tilt.

To run the Service Broker local integration tests, first run Tilt:

```sh-session
tilt up
```

Then, run the tests with the following command:

```sh-session
go test ./test/integration/main_test.go
```

Alternatively, you can re-run integration-test in the Tilt dashboard.

### End-to-End (E2E) Integration Testing

End-to-End testing is automatically triggered during pipeline. <br>
It is not supported and possibly not impossible to run these locally due to infrastructure dependencies.

E2E integration tests are being executed in Tanzu virtual machine based on binary that is built during Jenkins pipeline.

## Releases

### Verify and update dependencies

1.  Review the changes to `go.mod` since the last release and make any needed
    updates to [NOTICES.txt](./NOTICES.txt):
    *   Verify that dependencies fit into supported licenses types:
        ```shell
         go-licenses check ./... --allowed_licenses="MIT,ISC,Apache-2.0,BSD-3-Clause"
        ```
        If there is new dependency having unsupported license, such license should be included to [notices.tpl](./notices.tpl)
        file in order to get generated in NOTICES.txt.

    *   If no errors occur, proceed to generate updated NOTICES.txt:
        ```shell
         go-licenses report ./... --template notices.tpl > NOTICES.txt
        ```

### Update the version and changelog

1.  Create a new branch for the version bump.

2.  Based on the unreleased content, determine the new version number and update
    the [VERSION](VERSION) file. This project uses [semantic versioning](https://semver.org/).

3.  Ensure the [changelog](CHANGELOG.md) is up to date with the changes included in the release.

4.  Ensure the [open source acknowledgements](NOTICES.txt) are up to date with the dependencies,
    and update the file if there have been any new or changed dependencies since the last release.

5.  Commit these changes - `Bump version to x.y.z` is an acceptable commit message - and open a PR
    for review. Your PR should include updates to
    `CHANGELOG.md`, and if there are any license updates, to `NOTICES.txt`.

### Release and Promote

1.  Jenkins build parameters can be utilized to release and promote successful builds.

2.  Merging into main/master branches will automatically trigger a release.

3.  Reference the [internal automated release doc](https://github.com/conjurinc/docs/blob/master/reference/infrastructure/automated_releases.md#release-and-promotion-process)
    for releasing and promoting.

## Contributing steps

1.  Fork it
2.  Create your feature branch (`git checkout -b my-new-feature`)
3.  Commit your changes (`git commit -am 'Added some feature'`)
4.  Push to the branch (`git push origin my-new-feature`)
5.  Create new Pull Request

[get-docker]: https://docs.docker.com/engine/installation
[get-git]: https://git-scm.com/downloads
[get-tilt]: https://docs.tilt.dev/install.html
