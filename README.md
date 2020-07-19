deligt
======================

**deligt**   - Scraping stock data

# install

## Go Get
```bash
$ go get github.com/ybalexdp/deligt  
```

# Config


## conf/deligt.yaml

Set fullpath to output the result.   
ex )
```yaml
path: "/Users/ybalexdp/deligt/stock.csv"
```


# Command Line Options
```bash
deligt update [--number <number>] [--all]
```

*--number <number\>*  
　Specify the company code for which you want to get stock information.

*--all*  
　The data for all companies listed in the file will be updated.