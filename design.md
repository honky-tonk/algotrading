# Design Of Architecture 
# Design Of Interfaces
## /api/stock/query
**Interface path**: `/api/stock/query`

**Request Method**: `POST`

**Description**: front-end send stock metadata to back-end to query stock data 

**header**: `Content-Type: application/json`

**Request Body**: 
```json
{
  "data_type":  //Int type for stock data type like daily data,
  "stock_name":  //String type for stock name,
  "period": //int type for period you query(start now)
  "indic":    //Indicator type 4 is non indicator
  "indic_period":   //Indicator period
}
```

"data_type":

- 1:Daily
- 2:Weekly
- 3: Monthly

"indic":

- 4: NoneIndicator
- 5: SMA
- 6: EMA
- 7: MACD
- 8: KDJ


**Respond header**: `Content-Type: application/json`

**Respond Body**:

if request body `indic` field is 4 
```json
{
"prices": [
        {
            "stock_price": {
                "open": //float,
                "high": , // float
                "low": , //float
                "close": ,  //float
                "volume": //int
            },
            "time": //string
        },
        {

        }，
        ...
        ...
    ]
}
```

if request body  `indic` field NOT equal  4 
```json
{
    "prices":[
        {
            "stock_price": {
                "open": //float,
                "high": , // float
                "low": , //float
                "close": ,  //float
                "volume": //int
            },
            "time": //string
        },
        {

        }，
        ...
        ...
    ],
    "period":,
    "indic": {

    }

}
```
if `indic` is ```SMA``` Then indic field is 

```json
{
    "indic_period": ,
    "indic_values": [
      {
        "time":  ,
        "value":
      }
      ...
      ...
      ...
    ]
}
```

if `indic` is ```EMA``` Then indic field is 

```json
{
    "indic_period": ,
    "indic_values": [
      {
        "time":  ,
        "value":
      }
      ...
      ...
      ...
    ]
}
```