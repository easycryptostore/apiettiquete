# Apiettiquete

a polite way to ask for API data, from external, rate-limiting sources

in short: you can "hammer" this API as much as you want, while you have a X seconds guarantee of fresh data.

example:

calling 2 URLs:

URI_1=https://api.nomics.com/v1/currencies/ticker?key=2018-09-demo-dont-deploy-b69315e440beb145&ids=BTC,ETH&interval=1h&convert=USD

URI_2=https://api.nomics.com/v1/currencies/ticker?key=2018-09-demo-dont-deploy-b69315e440beb145&ids=BTC,ETH&interval=1h&convert=EUR

they do some rate limiting, and I don't need it 1 second fresh, so I get them every 30 second, and cache the result.

Calling:
~~~~
GET    /api/USD
GET    /api/EUR
~~~~

returns a response instantly, while making sure the data is at least X seconds fresh
also you can "hammer" this API as much as you want


### License: MIT
