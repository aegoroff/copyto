copyto
======

[![Build status](https://ci.appveyor.com/api/projects/status/21801axnoh4oxadf?svg=true)](https://ci.appveyor.com/project/aegoroff/copyto) [![codecov](https://codecov.io/gh/aegoroff/copyto/branch/master/graph/badge.svg)](https://codecov.io/gh/aegoroff/copyto) [![Go Report Card](https://goreportcard.com/badge/github.com/aegoroff/copyto)](https://goreportcard.com/report/github.com/aegoroff/copyto)

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

Examples:
---------

Let's do one way sync files between folder *D:\fSource* and *D:\fTarget*. Source folder content:
```
sub
  |- f3.pub
f1.docx
f2.pptx
f4.rar
```

Target folder content:
```
sub
  |- f3.pub
f1.docx
f2.pptx
f5.xlsx
```

So we use command line to run syncing:
```
copyto cmdline -s D:\fSource -t D:\fTarget
```

It will copy 3 files from *D:\fSource* to *D:\fTarget* and says that *f5.xlsx* not found in the source folder. File *f4.rar* will not be copied because it is not exists in the target folder:

```
   Found files that present in target but missing in source:
     \f5.xlsx

   Total copied:                              3
   Present in target but not found in source: 1
```

Now let's do the same task using config file. Create text file in UTF-8 encoding **without BOM** with the content like this:
```
# Example copyto config

title = "Exaample sync"

[sources]
 [sources.src1]
  source = 'D:\fSource'

[definitions]

  [definitions.def1]
  sourceLink = "src1"
  target = 'D:\fTarget'

  [definitions.def2]
  source = 'D:\fSource1'
  target = 'D:\fTarget2'
```
**IMPORTANT:** So as not to write double back slashes on Windows use string in '(apos) instead of "(quote).

You can use one source for several definitions using it's key (string after dot in square brackets) as value of *sourceLink* option. If both *source* and *sourceLink* defined in the
same definition *source* option wins.

And then use it using config verb:
```
copyto config -p D:\example.toml
```
The app will do sync and shows output like this:
```
 Section: def1
 Source: D:\fSource
 Target: D:\fTarget
   Found files that present in target but missing in source:
     \f5.xlsx

   Total copied:                              3
   Present in target but not found in source: 1

 Section: def2
 Source: D:\fSource1
 Target: D:\fTarget2

   Total copied:                              0
   Present in target but not found in source: 0
```