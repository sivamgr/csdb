# Candlestick-DB
Database for 1-minute Candle-Stick Data.

# Data Upload Format

- Upload Endpoint
```
http://localhost:8087/upload
```

    
- Zip of txt/csv files 
- Zip of Zip files

- CSV Format A

```
symbol,yyyymmdd,hh:mm,o,h,l,c,volume,open-intrest
```

- CSV Format B with file name as `<symbol>.txt` or `<symbol>.csv`

```
yyyymmdd,hh:mm,o,h,l,c,volume,open-intrest
```

# Read

## List Symbols
```
http://localhost:8087/symbols
```

-   returns JSON Array of Symbols, 
```
["NIFTY","NIFTY_1","ACC"]
```

## Fetch Candles

```
http://localhost:8087/data?s=<symbol>?from=utc&to=utc
```

-   returns JSON Array of Candle Stick Data,
```
[
    {"s":"NIFTY","dt":"2020-01-03 09:29:00 +0530" "o":15000.00,"h":15100.00,"l":15000.00,"c":15000.00,v:"0",oi:"0"},
    {"s":"NIFTY_1","dt":"2020-01-03 09:29:00 +0530" "o":15100.00,"h":15190.00,"l":15100.00,"c":15110.00,v:"3200",oi:"1000000"},
    ...
]
 ```



# Run 
- Run once,

```docker run --name csdb -v $PWD/csdb:/opt sivamgr/csdb:latest```


- Run Forever,

```docker run -d --restart unless-stopped --name csdb -v $PWD/csdb:/opt sivamgr/csdb:latest```
