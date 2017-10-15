# master
Simple command line utility for running several command copies in parallel

## Usage

### Flags

```-n int``` - number of parallel process

```-r``` - set this flag to repeat command after it exits

### Example

``` 
master -n 10 echo hello
```

Will run ```echo hello``` ten times 
