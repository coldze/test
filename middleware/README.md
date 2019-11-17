###Test solution

[Go to main](../README.md)

* package`source` contains interfaces and implementations of data-source
    * `httpDataSource` - this data-source is able to get data from external API via http-calls.
    * `redisCacheSource` - gets data from redis, using provided key.
    * `cachedDataSource` - uses both data-sources from above to get data and cache it.
* package `handlers` contains handlers for incoming http calls.
* root of this package contains some common interfaces and implementations.