# left

## Description

_left_ generates pdf formatted letters from plain text files. 
Hence the name _left_ as in "LEtter From Txt".

## License

_left_ is licensed under the [GPL, version 3](LICENSE).  
_left_ uses [go-fpdf](https://github.com/go-pdf/fpdf) to render PDF files, which is licensed under the [MIT license](https://github.com/go-pdf/fpdf/blob/main/LICENSE)

## Usage

To get information on how to use _left_, run 
```
left -help
```

## Configuration

_left_ reads configuration files in the following order, where settings in later files override settings read from earlier files:
- /etc/left/defaults.json (when running on linux)
- ${UserConfigDir}/left/defaults.json (also see [UserConfigDir documentation](https://pkg.go.dev/os#UserConfigDir))
- optionally the config file specified via command line argument
- the configuration in the letter input file

_left_ can dump a sample configuration to stdout that can be used as a starting point:
```
left -dump-config
```

## Creating letters

As stated above, _left_ creates letters from simple text input files.
To get started, run the following command and _left_ will output a working letter input file to stdout:
```
left -create
```

## Tips & Tooling

### vim

I edit my letters with vim. While _left_ does not care about the file extension, 
I find it useful to give my letters the extension .left.
Then I add the following section to my vim config and get reasonable syntax highlighting:
```
augroup left_ft
au!
  autocmd BufNewFile,BufRead *.left set syntax=javascript
augroup END
```
 