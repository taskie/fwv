# fwv

manipulate “Fixed Width Values”

![fwv](images/example.gif)

## Install

```sh
go get -u github.com/taskie/fwv/cmd/fwv
```

## Usage

### Pretty-print Fixed Width Values

```sh
fwv foo.txt
```

#### Input (foo.txt)

```txt
a        bb     あいう  ccc   dd     e
漢字    f     gg         h    ,      i
```

#### Output

```txt
a    bb あいう ccc dd e
漢字 f  gg     h   ,  i
```

### Convert Fixed Width Values to CSV

```sh
fwv foo.txt foo.csv
```

or

```sh
fwv -t csv <foo.txt >foo.csv
```

#### Output (foo.csv)

```csv
a,bb,あいう,ccc,dd,e
漢字,f,gg,h,",",i
```

### Convert CSV to Fixed Width Values

```sh
fwv foo.csv foo.txt
```

or

```sh
fwv -f csv <foo.csv >foo.txt
```

### Treat "Eastern Asian Ambiguous Width" as halfwidth

```sh
fwv -E foo.txt
```

or

```sh
env FWV_EAA_HALF_WIDTH=1 fwv foo.txt
```

### Ignore character width

Only the number of characters (runes) is used for calculating the width.

```sh
fwv -W foo.txt
```

### Specify a delimiter

```sh
fwv -d '│' foo.txt
```

```txt
a   │bb│あいう│ccc│dd│e
漢字│f │gg    │h  │, │i
```

### Do not trim whitespaces

```sh
fwv -d '│' -T foo.txt
```

```txt
a       │ bb   │  あいう  │ccc   │dd     │e
漢字    │f     │gg        │ h    │,      │i
```

## Dependencies

![dependency](images/dependency.png)

## License

Apache License 2.0
