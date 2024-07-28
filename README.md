
<h1 align="center">
  <br>
  <a href="https://ufsit.club/teams/blue.html"><img src="https://raw.githubusercontent.com/alexchristy/WazuhTest/main/assets/img/wazuhTest.png" alt="WazuhTest" width="200"></a>
  <br>
  WazuhTest
  <br>
</h1>

<h4 align="center">A testing framework built for <a href="https://wazuh.com" target="_blank">Wazuh</a>.</h4>

<p align="center">
  <a href="#key-features">Key Features</a> •
  <a href="#how-to-use">How To Use</a> •
  <a href="#what-are-the-tests">What are the tests?</a> •
  <a href="#organizing-tests">Organizing Tests</a> •
  <a href="#test-syntax">Test Syntax</a> •
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

## What are the tests?

Tests are organized by directories, each containing any number of JSON files defining the tests. Raw logs used for testing can be stored anywhere locally but are typically kept in the same directory as the test definition files.

For an example, see the `wazuh-tests/` directory in the repository. These tests work out of the box with any default Wazuh manager installation.

## Organizing Tests

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

## Test Syntax

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

Example included tests from `wazuh-tests/ubuntu/test_ssh.json`:

```json
{
    "tests": [
        {
            "TestDescription": "SSH login to a non-existent user",
            "RuleID": "5710",
            "RuleLevel": "5",
            "Format": "syslog",
            "RuleDescription": "sshd: Attempt to login using a non-existent user",
            "LogFilePath": "5710.txt",
            "Predecoder": {},
            "Decoder": {}
        },
        {
            "TestDescription": "SSH login to a non-existent user from a local network",
            "RuleID": "5710",
            "RuleLevel": "5",
            "Format": "syslog",
            "RuleDescription": "sshd: Attempt to login using a non-existent user",
            "LogFilePath": "5710-local-net.txt",
            "Predecoder": {},
            "Decoder": {
                "srcip": "10.0.0.4",
                "srcport": "59528",
                "srcuser": "non-existent"
            }
        }
    ]
}
```

### Decoder Fields

The `Predecoder` and `Decoder` fields in a test accept arbitrary key-value pairs that are checked against the Wazuh output.

**Example:**

Below is an example test using both `Predecoder` and `Decoder` fields where the appropriate values for `Predecoder` and `Decoder` are determined based on the Wazuh Ruleset Test output. The key-value pairs in the test must match the Wazuh Ruleset Test output key-value pairs exactly, and the order does not matter.

**Finding Wazuh Ruleset Test**:

![Visual steps to access Wazuh Ruleset Test tool in dashboard](https://github.com/user-attachments/assets/0960d86a-86b0-4819-a749-6acb510fc13c)

**Raw Log:**

```txt
Mar  5 13:49:34 ip-10-0-0-10 sshd[1602]: Invalid user non-existent from 10.0.0.4 port 59528
```

**Run Test:**

![image](https://github.com/user-attachments/assets/e7bba0ae-8de5-4723-a096-b4a6c0b42b54)

**Predecoder Output:**

![Predecoder values from Wazuh Ruleset Test with red arrows pointing out program_name and hostname values](https://github.com/user-attachments/assets/3f34adf0-4374-42ae-8aa0-4ba747fb836e)

**Decoder Output:**

![Decoder values from Wazuh Ruleset Test with red arrows pointing out the src* values](https://github.com/user-attachments/assets/cca6aef5-150a-44fd-8363-dc79b2d38e10)

**Test:**

```json
{
    "tests": [
        {
            "TestDescription": "SSH login to a non-existent user from a local network",
            "RuleID": "5710",
            "RuleLevel": "5",
            "Format": "syslog",
            "RuleDescription": "sshd: Attempt to login using a non-existent user",
            "LogFilePath": "5710-local-net.txt",
            "Predecoder": {
              "program_name": "sshd",
              "hostname": "ip-10-0-0-10"
            },
            "Decoder": {
                "srcip": "10.0.0.4",
                "srcport": "59528",
                "srcuser": "non-existent"
            }
        }
    ]
}
```

> **Note:** When using the Wazuh log test API, if a value is not found in decoder or predecoder, the tool automatically checks for it in the data field before reporting it as missing.

## Related

[wazuh-pipeline](https://github.com/alexchristy/wazuh-pipeline) - Wazuh CI pipeline that leverages this tool

## License

GNU General Public License v3.0
