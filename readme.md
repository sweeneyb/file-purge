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