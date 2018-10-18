copyto
======


**copyto** is a small commandline app written in Go that allows you to easily one way
sync files between source folder and target folder. i.e. all files from source that
matches (by relative path) corresponding files in target folder will be copied from source folder to
target one.

The app supports setting source and target paths directly from command line and configuration file
in TOML format also can be used. Using configuration file you can setup several sources and targets

Command line syntax:
--------------------
```
Usage: copyto [global options] <verb> [verb options]

Global options:
            --version Print version
        -v, --verbose Verbose output

Verbs:
    cmdline:
        -s, --source  Path to the source folder, to copy (sync) data from (*)
        -t, --target  Path to the target folder, to copy (sync) data to (*)
    config:
        -p, --path    Path to configuration file (*)
```
