# embedded.template

embedded-template is a helper library for loading templates embedded in the binary. 
It uses the [go.rice](https://github.com/GeertJohan/go.rice) package for handling 
the embedded files.

## Installation

```
go get github.com/magna5/embedded-template
go get github.com/GeertJohan/go.rice/rice
```

## Usage

The library expects the templates to be in a `templates` directory.  Since `go.rice` 
requires loading the embedded files based on a static string and cannot use a 
constant or variable it has been hardcoded to templates.  All files in the templates
folder will be loaded and parsed.

```golang
data := {}
tmpl, err := embtemplate.LoadTemplates()
if err != nil {
    log.Error(err.Error())
    return
}
tmpl.Execute(os.Stdout, data)
```

The return from LoadTemplates is a `text/template` Template with all of the files
in `templates` parsed.

Since the returned template is a standard go template the embedded templates can 
be overridden by parsing files from the filesystem.  eg.

```golang
tmpl, err := embtemplate.LoadTemplates()
if err != nil {
    log.Error(err.Error())
    return
}
tmpl, err = tmpl.ParseGlob("overridedir/*")
if err != nil {
    log.Error(err.Error())
    return
}
tmpl.Execute(os.Stdout, data)
```

The `go.rice` package will load the files from the local `templates` directory during 
development.  To embed them in the binary you must use the `rice` cli at build time.

```bash
go build -o example
rice append --exec example
```
