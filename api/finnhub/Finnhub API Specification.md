# Finnhub API Specification
This document details the API endpoints that we will be using from the Finhub API.

## Data Sources and Documentation
We based this information off of the [Finnhub API documentation](https://finnhub.io/docs/api) and their [swagger file](https://finnhub.io/api/v1/swagger.json).

## Metadata
### Base URL
`https://finnhub.io/api/v1`

### Authentication
All GET request require a token parameter `token=apiKey` in the URL or a header `X-Finnhub-Token: apiKey`.

## Endpoints
### Realtime Trade Stream (Websocket)
[Documentation](https://finnhub.io/docs/api/websocket-trades) \
Description: `Websocket endpoint for real-time trades.` \
URL: `wss://ws.finnhub.io?token=<YOUR_API_KEY>` \
Method: `Websocket` \
Response Schema (JSON):
```json
{
  "data": [
    {
      "p": 0, // Last price
      "s": "string", // Symbol
      "t": 0, // Unix timestamp
      "v": 0 // Volume
    }
  ],
  "type": "trade"
}
```

### Realtime Quote
[Documentation](https://finnhub.io/docs/api/quote) \
Description: `Realtime quote for a given symbol.` \
URL: `/quote?symbol=<SYMBOL>` \
Method: `GET` \
Response Schema (JSON):
```json
{
  "c": 0, // Current price
  "d": 0, // Change
  "dp": 0, // Percent change
  "h": 0, // Daily high
  "l": 0, // Daily low
  "o": 0, // Open price
  "pc": 0, // Previous close price
  "t": 0 // Unix timestamp
}
```

### Stock Candles
[Documentation](https://finnhub.io/docs/api/stock-candles) \
Description: `Candle data for a given symbol.` \
URL: `/stock/candle?symbol=<SYMBOL>&resolution=<RESOLUTION>&from=<FROM>&to=<TO>` \
Method: `GET` \
Response Schema (JSON):
```json
{
  "c": [
    0 // Close prices
  ],
  "h": [
    0 // High prices
  ],
  "l": [
    0 // Low prices
  ],
  "o": [
    0 // Open prices
  ],
  "s": "ok" or "no_data", // Status of the response. This field is optional but it is recommended to check it to make sure the data is good.
  "t": [
    0 // Unix timestamps
  ],
  "v": [
    0 // Volumes
  ]
}
```
