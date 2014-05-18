# WatchFS

WatchFS is a tool that monitors file system for changes and runs specified command after the change.

## Install
The tool can be installed using "go get" command:
```
  $ go get -u github.com/igm/watchfs
```

## Usage

Simple notification after any change on the file system (ignoring hidden files starting with "."):
```
  $ watchfs echo "change!!!"
```
 As the original purpose of this tool was to execute go tests automatically after the source code change
 the format for that purpose is:
```
  $ watchfs -f ".\*\.go$" go test 
```
To specify a different time wait period after a change (or set of changes) occurs use timeout parameter:
```
  $ watchfs -t 250ms go test
```
For the full ist of avaialble parameters run command without any parameters:
```
  $ watchfs
```

### Limitations

Current version supports current directory monitoring only. Recursive file system monitoring is not supported.

## Timeout

Usually file system emits several various events when editing a file in a editor. For that purpose timeout parameter defines how long to wait before executing the command. The wait period starts again after any filesystem event.

