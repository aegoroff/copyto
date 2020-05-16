copyto
======

[![Build status](https://ci.appveyor.com/api/projects/status/21801axnoh4oxadf?svg=true)](https://ci.appveyor.com/project/aegoroff/copyto) [![codecov](https://codecov.io/gh/aegoroff/copyto/branch/master/graph/badge.svg)](https://codecov.io/gh/aegoroff/copyto) [![Go Report Card](https://goreportcard.com/badge/github.com/aegoroff/copyto)](https://goreportcard.com/report/github.com/aegoroff/copyto)

**copyto** is a small commandline app written in Go that allows you to easily one way
sync files between source folder and target folder. i.e. all files from source that
matches (by relative path) corresponding files in target folder will be copied from source folder to
target one.

The app supports setting source and target paths directly from command line and configuration file
in TOML format also can be used. Using configuration file you can setup several sources and targets.

Also you can setup file names filter using include or exclude options (or both) using either configuration
file or command line. Include filter allows only file matched to be copied. Exclude filter allows all files but 
those that matched to be copied.

Command line syntax:
--------------------
```
Usage:
  copyto [flags]
  copyto [command]

Available Commands:
  cmdline     Use command line to configure required application parameters
  config      Use TOML configuration file to configure required application parameters
  help        Help about any command
  version     Print the version number of copyto

Flags:
  -h, --help      help for copyto
  -v, --verbose   Verbose output

Use "copyto [command] --help" for more information about a command.
```

Command line mode syntax:
```
Use command line to configure required application parameters

Usage:
  copyto cmdline [flags]

Aliases:
  cmdline, cmd, l

Flags:
  -e, --exclude string   Exclude files whose names match pattern specified by the option
  -h, --help             help for cmdline
  -i, --include string   Include only files whose names match the pattern specified by the option
  -s, --source string    Path to the source folder, to copy (sync) data from (required)
  -t, --target string    Path to the target folder, to copy (sync) data to (required)

Global Flags:
  -v, --verbose   Verbose output
```

Config file mode syntax:
```
Use TOML configuration file to configure required application parameters

Usage:
  copyto config [flags]

Aliases:
  config, conf, c

Flags:
  -h, --help          help for config
  -p, --path string   Path to configuration file (required)

Global Flags:
  -v, --verbose   Verbose output
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

It will copy 3 files from *D:\fSource* to *D:\fTarget* and says that *f5.xlsx* not found in the source folder. 
File *f4.rar* will not be copied because it is not exists in the target folder:

```
   Found files that present in target but missing in source:
     \f5.xlsx

   Total copied:                              3
   Present in target but not found in source: 1
```

Now let's do the same task using config file. Create text file in UTF-8 encoding with the content like this:
```
# Example copyto config

title = "Exaample sync"

[sources]
 [sources.src1]
  source = 'D:\fSource'

[definitions]

  [definitions.def1]
  sourcelink = "src1"
  target = 'D:\fTarget'

  [definitions.def2]
  source = 'D:\fSource1'
  target = 'D:\fTarget2'
```
**IMPORTANT:** So as not to write double back slashes on Windows use string in '(apos) instead of "(quote).

**IMPORTANT:** All keys must be in lower case

You can use one source for several definitions using it's key (string after dot in square brackets) as value 
of *sourceLink* option. If both *source* and *sourceLink* defined in the
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

Filtering config example:
```
# Example copyto config

title = "Exaample sync"

[sources]
 [sources.src1]
  source = 'D:\fSource'

[definitions]

  [definitions.def1]
  sourcelink = "src1"
  target = 'D:\fTarget'
  exclude = '*.exe'

  [definitions.def2]
  source = 'D:\fSource1'
  target = 'D:\fTarget2'
  include = '*.txt'

  [definitions.def3]
  source = 'D:\fSource3'
  target = 'D:\fTarget4'
  include = '*.txt'
  exclude = 'bad*.txt'
```