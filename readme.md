A small go program to purge old files, based on ionotify.  Heavily borrowed from https://github.com/fsnotify/fsnotify

## running
```
go run ./cmd watch /tmp/fsnotify/
```

in /tmp/fsnotify, do some commands:
```
$ echo "foo" > foo
$ echo "foo" > foo2
$ echo "foo" > foo3
$ echo "foo" > foo
$ echo "foo" > foo
$ ls
foo  foo3
```

## Policy
The time and path are being injected into the input as:
```go
	input := map[string]interface{}{
		"time": item.priority,
		"path": item.value,
	}
```

And there's an embedded example policy to expire/delete files after 10 seconds:
```Rego
package example.authz

import future.keywords.if
import future.keywords.in


default allow := false

delayTime := 10

allow if {
    input.time != null
    currentTimestamp := time.now_ns() / 1000000000  # Convert nanoseconds to seconds
    inputTimestamp := time.parse_rfc3339_ns(input.time) / 1000000000

    # Check if the input time is earlier than 10 seconds ago
    inputTimestamp < currentTimestamp - delayTime
}
```

From this, we could write more extensive rule sets (and conftest them) to give different expiration policies based on filename.  


## TODO
1. Allow for external policies to be ingested
1. Make this all work recursively