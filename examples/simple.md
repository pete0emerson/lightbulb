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
Thu Mar  9 13:35:34 UTC 2023
```
