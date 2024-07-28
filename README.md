
<h1 align="center">
  <br>
  <a href="http://ufsit.clubhttps://ufsit.club/teams/blue.html"><img src="https://raw.githubusercontent.com/alexchristy/WazuhTest/main/assets/img/wazuhTest.png" alt="WazuhTest" width="200"></a>
  <br>
  WazuhTest
  <br>
</h1>

<h4 align="center">A testing framework built for <a href="https://wazuh.com" target="_blank">Wazuh</a>.</h4>

<p align="center">
  <a href="#key-features">Key Features</a> •
  <a href="#how-to-use">How To Use</a> •
  <a href="#download">Download</a> •
  <a href="#credits">Credits</a> •
  <a href="#related">Related</a> •
  <a href="#license">License</a>
</p>

![screenshot](https://raw.githubusercontent.com/alexchristy/WazuhTest/main/assets/img/wazuh_test_demo.gif)

## Key Features

* Test Decoders
  - Verify that your decoders are extracting info from logs.

* Test Rules
  - Ensure that the correct rules alert on logs.

* Programatically Verify Readiness

* Catch Errors Early

* Prevent Regression

* CI Pipeline Integration
  - Tests can be used with the [wazuh-pipeline](https://github.com/alexchristy/wazuh-pipeline) built to work with GitHub actions.

## How To Use

To clone and run this application, you'll need [Git](https://git-scm.com) and [Golang](https://go.dev/) installed on your computer. From your command line:

```bash
# Clone this repository
git clone https://github.com/alexchristy/WazuhTest

# Go into the repository
cd WazuhTest

# Compile tool
go build .
```

### Run Included Tests

```bash
./WazuhTest -d ./wazuh-tests/ -u {WAZUH_API_USERNAME} -p {WAZUH_API_PASSWORD} -t 3 -v {WAZUH_MANAGER_HOSTNAME}
```

> This step **REQUIRES** that you have a running Wazuh manager. The quickest way to do this is to use the official Wazuh manager [docker image](https://hub.docker.com/r/wazuh/wazuh-manager) and port forward port 55000.

## Writing Tests

Tests are organized by directories, each containing any number of JSON files defining the tests. Raw logs used for testing can be stored anywhere locally but are typically kept in the same directory as the test definition files.

For an example, see the `wazuh-tests/` directory in the repository. These tests work out of the box with any default Wazuh manager installation.

### Test Grouping

Tests are grouped by directories. The provided tests in `wazuh-tests/` are grouped into two categories: `ubuntu` and `centos`. The `.txt` files are the raw logs sent to the Wazuh manager API, while the `.json` files define the tests.

**Key Points:**
- JSON files define tests.
- Multiple JSON files can exist per directory.
- Logs can be in any format, but each file should contain **only** a single line.
- The location of the logs is defined in the JSON test files.

Example test structure:

```txt
wazuh-tests/
|
+-----ubuntu/
|     |
|     +-----203.txt 
|     +-----5710-local-net.txt  <--- (Raw logs can be any format)
|     +-----5170.txt
|     +-----test_agent.json     <--- (JSON files define tests)
|     +-----test_ssh.json
|
+-----centos/
|
(...)
```

### Test Syntax

Each JSON test file is a list of test objects under `tests`.

**Required Fields:**
* `RuleID` - An integer between 0 and 999999.
* `RuleLevel` - An integer between 0 and 16.
* `LogFilePath` - Path to the log file, which must exist, be readable, not empty, and contain only one line.
* `Format` - A valid format type such as "syslog", "json", "snort-full", etc.

**Optional Fields (warnings if not provided or empty):**
* `Version` - A string indicating the version of the test.
* `RuleDescription` - A string describing the rule.
* `Decoder` - A map of key-value pairs for the decoder.
* `Predecoder` - A map of key-value pairs for the predecoder.
* `TestDescription` - A string describing the test.

Example included test from `wazuh-tests/ubuntu/test_agent.json`:

```json
{
    "tests": [
        {
            "TestDescription": "Wazuh agent event queue is full",
            "RuleID": "203",
            "RuleLevel": "9",
            "Format": "syslog",
            "RuleDescription": "Agent event queue is full. Events may be lost.",
            "LogFilePath": "203.txt",
            "Predecoder": {},
            "Decoder": {}
        }
    ]
}
```

## Related

[wazuh-pipeline](https://github.com/alexchristy/wazuh-pipeline) - Wazuh CI pipeline that leverages this tool

## License

GNU General Public License v3.0
