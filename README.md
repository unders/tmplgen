# tmplgen
tmplgen generates static HTML pages from one directory to another directory

### Usage

###### Show help message
`tmplgen`

###### generates HTML files into a directory
`tmplgen -from=dir -to=dir files`

### Installing
See [releases](https://github.com/unders/tmplgen/releases)

**OSX binary download**
```
curl -L https://github.com/unders/tmplgen/releases/download/v1.0.0/tmplgen_1.0.1_darwin_amd64.tar.gz | tar -zxv
```

### Developer
Run `make` in the root directory

### How to create an release
```
export GITHUB_TOKEN=`YOUR_TOKEN`
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
goreleaser
```
