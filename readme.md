A small go program to purge old files, based on ionotify.  Heavily borrowed from https://github.com/fsnotify/fsnotify

## Motivation
In dealing with constantly churning data, we need to purge some of it to keep services happy and disk utilization under thresholds.  Often, we promise "days" of retention, but implementing that is messy and expensive.  It often looks like: 
```
find /path/to/directory/ -mindepth 1 -mtime +5 -delete
```

This puts a lot of load on the disk, as we're scanning through everything.  It's also inconsistent to the user, as "2 days ago" at noon can look like "2.5 days ago" if we run a purge cron 1x/day. 

The goal of this repo is to PoC using ionotify to remove files within a narrow time band of their expiration.  This should be made relatively efficient using a heap to track update times of files.  Though we may have multiple entries for subsequent writes, I'll need to collect data on if that's actually significant or if we can drop old write records.  Currently, we check the last modified time before deleting a file and place it back in the heap if a subsequent write has occured. 

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