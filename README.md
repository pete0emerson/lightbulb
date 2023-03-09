# Lightbulb

Lightbulb is a program and a Golang module to sequentially execute and test Markdown code blocks. Lightbulb leverages formatted HTML comments to execute a code block, create a file, or prompt for environment variables. State (via environment variables) is maintained and carried over from block to block.

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


---
Read about the design and development approach for Lightbulb [here].(docs/design_and_development.md)