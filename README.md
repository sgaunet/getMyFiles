# getMyFiles

Recursive copy of a folder to another with every files. 

If a file copy fails, there will be 100 tries, with pause of 30s between every try. This little program is motivated by the fact that I would like to get back a lot of files from pcloud (where lot of files need to be downloaded again before making a copy, and with my little bandwith, there is a lot of errors... ).

And I'm trying goreleaser also.

# Build

```
go build getMyFiles.go
```

# Run

```
getMyFiles -o origin -d dest
```

Get help : 

```
getMyFiles --help
```

# Release

```
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
goreleaser --snapshot  # Check
goreleaser 
```