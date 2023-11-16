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

## TODO
1. Add a timer and refactor code so policy eval isn't only triggered by the watch/event loop
1. Factor out the delay to be configurable
1. Factor out the configurable delay to get it's policy from an embeded OPA instance (with a Rego policy)