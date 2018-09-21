# fwv

manipulate “Fixed Width Values”

## install

```sh
go get github.com/taskie/fwv/cmd/fwv
```

## usage

### convert CSV to Fixed Width Values

```sh
fwv -o foo.txt foo.csv
```

or

```sh
fwv <foo.csv >foo.txt
```

#### foo.csv (input)

```
a,bb,あいう,ccc
漢字,d,eee,f
```

#### foo.txt (output)

```
a    bb あいう ccc
漢字 d  eee    f
```

### convert Fixed Width Values to CSV

```sh
fwv -r -o foo.csv foo.txt
```

or

```sh
fwv -r <foo.txt >foo.csv
```

### treat "Eastern Asian Ambiguous Width" as halfwidth

```sh
fwv -E
```

or

```sh
env FWV_EAA_HALF_WIDTH=1 fwv
```

### ignore character width

Only the number of characters (runes) are considered.

```sh
fwv -W
```

### specify a delimiter

```sh
fwv -d '│'
```

#### foo.txt (output)

```
a   │bb│あいう│ccc
漢字│d │eee   │f
```

## dependency

![dependency](images/dependency.png)

## license

Apache License 2.0
