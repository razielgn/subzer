# Subzer [![Build Status](https://travis-ci.org/razielgn/subzer.png?branch=master)](https://travis-ci.org/razielgn/subzer)

Subzer is a small CLI utility written in Go.  
Its purpose is to convert Subrip subtitle files to a custom subtitle format.

## Building

To build it, you're going to need a golang 1.1 installation.  
Then a `make` will compile it for you.  
Tests are run by `make test`.  
Cross compiling awesomeness provided by `make cross` (win32, win64, linux32, linux64, osx64).

## Usage

* **Just an input file**:

``` bash
subzer -i input.srt # Output will be input.txt
```

* **Recusively scan a directory for .srt files**

``` bash
subzer -r folder/ # Output will be filename.txt
```
