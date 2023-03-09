# Lightbulb

Lightbulb is a program and a Golang module to sequentially execute and test Markdown code blocks. Lightbulb leverages formatted HTML comments to execute a code block, create a file, or prompt for environment variables. State (via environment variables) is maintained and carried over from block to block.

## Features

* Create files with the contents of a Markdown code block
* Execute shell commands
* Create environment variables that persist from block to block

## Example Markdown file

In [this Markdown file](examples/simple.md), there are two code blocks that will be executed sequentially.

First, a file will be created, and then a shell command will be executed.

~~~

# Create `date.sh`

Make sure the file is executable.

<!-- lightbulb:createFile name:date.sh mode:0700 -->
```shell
#!/bin/bash

echo "The current date in UTC is $(date -u)."
```

# Run `date.sh`

<!-- lightbulb:runShell shell:bash -->
```console
./date.sh
```

This will output something like the following:
```console
The current date in UTC is Thu Mar  9 13:35:34 UTC 2023.
```
~~~

In this example, the file `date.sh` is created with specified permissions in the first step, and then the file is executed.

People who want to follow this Markdown step-by-step can do so (without even being aware that there are Lightbulb commands in place), but with the [contents of the above Markdown file](examples/simple.md) saved locally or sourced directly from the internet, Lightbulb can process the file and run each step for you.

## API

Lightbulb is driven by _actions_ via formatted HTML comments. Each action has _parameters_ with _keys_ and _values_ to configure those actions. The general format is:

~~~
<!-- lightbulb:ACTION name:VALUE [KEY:VALUE] [KEY:VALUE] -->
```
CODE BLOCK HERE
```
~~~

All actions must include the name of the action. This is stored in the state machine, and is useful when debugging the execution of a Markdown file when something has gone wrong.

All actions may include one or more tags. This is useful for testing purposes, as lightbulb can run only certain tags or skip certain tags if desired.

|name|description|default|required|
|----|-----------|-------|--------|
|`name`|The name of the block. This must be a single word.| | :white_check_mark:|
|`tags`|Tags allow lightbulb to control execution of actions.|all|:x:|

### createFile

`createFile` will create a file relative to the current directory.

#### Parameters

|name|description|default|required|
|----|-----------|-------|--------|
|`path`|The relative path where the file will be created. Absolute paths are not permitted.| |:white_check_mark:|
|`mode`|The unix [mode](https://linuxhandbook.com/chmod-command/) of the file.| 0700 | :x: |

#### Example

~~~
<!-- lightbulb:createDateFile name:sample_file path:sample/file.sh mode:0700 -->
```shell
#!/bin/bash

date
```
~~~

### runShell

`runShell` will create a temporary file, make it executable, and then run it from the working directory.

#### Parameters

|name|description|default|required|
|----|-----------|-------|--------|
|`command`|The command to run. By default, the command in the following code block is run. This is useful if you don't want to extract the command to be run from a code block.||:x:|
|`shell`|Set the name of the shell to use when executing the file. If the file contents specify a shell (`#!/bin/bash`), it will be used instead.|bash|:x:|
|`set`|Set `-x` and/or `-e` flags for the shell.||:x:|
|`exitOnError`|When the file runs, stop executing future Lightbulb actions if there is a non-0 exit code.|true|:x:|

#### Example

~~~
<!-- lightbulb:runShell name=showDate shell:bash set:x,e exitOnError:false -->
```bash

sample/file.sh
```
~~~

### setEnvironmentVars

`setEnvironmentVars` allows for the setting of environment variables that can persist throughout the running of all of the code blocks.

Mixing different configurations of environment variables isn't allowed, however, you can specify multiple lightbulb actions back-to-back to get the desired effect and / or order.

#### Parameters

|name|description|default|required|
|----|-----------|-------|--------|
|`keys`|A comma separated list of environment variable names. When prompted, `keys` will be presented in order.||:white_check_mark:|
|`prompt`|Specify when the user will be prompted for env var values. Options are: `never` (if the env var is not set, this will exit with a non-0 exit code), `missing` (prompt only if the env var isn't set), and `always` (prompt always).|missing|:x:|
|`secret`|If set to `true` and the user is prompted, text will not be echoed to the screen.|false|:x:|
|`persist`|If set to `true`, the environment variable value will persist through all subsequent code blocks.|true|:x:|
|`sensitive`|If set to `true`, the environment variable value will never be saved in the state machine.|false|:x:|

#### Example

~~~
<!-- lightbulb:setEnvironmentVars keys:FOO,BAR,BAZ prompt:always persist:false -->
<!-- lightbulb:setEnvironmentVars keys:GITHUB_TOKEN prompt:missing persist:false secret:true sensitive:true -->
~~~

---
Read about the design and development approach for Lightbulb [here](docs/design_and_development.md).