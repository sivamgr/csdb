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

- Example 
```
curl http://localhost:8087/inbox -F file=@NSE_EQT_1MIN_20220909.zip
```


# Read

## List Symbols
```
http://localhost:8087/symbols
```
-   returns JSON Array of Symbols, 
```
["20MICRONS","21STCENMGM","3IINFOLTD","3MINDIA","3PLAND","4THDIM","5PAISA","63MOONS","A2ZINFRA","AAKASH","AAREYDRUGS","AARON","AARTIDRUGS","AARTIIND","AARTISURF","AARVEEDEN","AARVI","AAVAS","ABAN","ABB","ABBOTINDIA","ABCAPITAL","ABFRL","ABMINTLLTD","ABSLAMC","ABSLBANETF","ABSLNN50ET","ACC","ACCELYA","ACCURACY","ACE","ACRYSIL","ADANIENT","ADANIGREEN","ADANIPORTS","ADANIPOWER","ADANITRANS","ADFFOODS","ADL","ADORWELD","ADROITINFO"]
```

## Fetch Candles

```
http://localhost:8087/data?s=<symbol>?f=<utc>&t=<utc>
```

s : Symbol list
f : from datetime in UTC Seconds (Default is Now()- 7 days)
t : to datetime in UTC Seconds  (Default is f + 7 days)

-   returns JSON Array of Candle Stick Data,

- Example, 
```
curl http://localhost:8087/data?s=IOC

```
curl http://localhost:8087/data?s=IOC
```
[{"s":"IOC","c":[{"t":1662455700,"o":71.15,"h":71.15,"l":71,"c":71.05,"v":93847,"oi":0},{"t":1662455760,"o":71.1,"h":71.15,"l":71,"c":71.15,"v":66433,"oi":0},{"t":1662455820,"o":71.1,"h":71.15,"l":71.05,"c":71.1,"v":24120,"oi":0},{"t":1662455880,"o":71.1,"h":71.1,"l":71.05,"c":71.05,"v":47174,"oi":0},{"t":1662455940,"o":71.1,"h":71.1,"l":71.05,"c":71.05,"v":32065,"oi":0},{"t":1662456000,"o":71.1,"h":71.15,"l":71.05,"c":71.15,"v":32908,"oi":0},{"t":1662456060,"o":71.1,"h":71.15,"l":71.1,"c":71.1,"v":31662,"oi":0},{"t":1662456120,"o":71.15,"h":71.15,"l":71,"c":71.05,"v":47202,"oi":0},{"t":1662456180,"o":71.1,"h":71.15,"l":71.05,"c":71.1,"v":57350,"oi":0},{"t":1662456240,"o":71.1,"h":71.15,"l":71.1,"c":71.1,"v":19531,"oi":0},{"t":1662456300,"o":71.15,"h":71.15,"l":71.05,"c":71.1,"v":26871,"oi":0},{"t":1662456360,"o":71.15,"h":71.15,"l":71.05,"c":71.05,"v":26335,"oi":0},{"t":1662456420,"o":71.1,"h":71.1,"l":71.05,"c":71.1,"v":8847,"oi":0},{"t":1662456480,"o":71.1,"h":71.15,"l":71.05,"c":71.15,"v":10747,"oi":0},{"t":1662456540,"o":71.15,"h":71.15,"l":71.1,"c":71.15,"v":6841,"oi":0},{"t":1662456600,"o":71.15,"h":71.15,"l":71.1,"c":71.1,"v":15954,"oi":0},{"t":1662456660,"o":71.15,"h":71.15,"l":71.1,"c":71.15,"v":11516,"oi":0},{"t":1662456720,"o":71.1,"h":71.15,"l":71.05,"c":71.1,"v":23230,"oi":0},{"t":1662456780,"o":71.1,"h":71.1,"l":71.05,"c":71.05,"v":9692,"oi":0}]
 ```


# Run 
- Run once,

```docker run --name csdb -p 8087:8087 -v $PWD/csdb:/opt sivamgr/csdb:latest```


- Run Forever,

```docker run -d --restart unless-stopped --name csdb -p 8087:8087 -v $PWD/csdb:/opt sivamgr/csdb:latest```
