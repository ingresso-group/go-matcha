go-matcha
=========

Features
--------

### Capturing values from the response

If you define a field with a `capture` tag then that field will be captured from the response. This is useful for more complex assertions.

For example, the following field would be captured as 'count':

```
Count    float64 `capture:"count"`
```

Note that it is currently only possible to capture one value for a given field. So, for example, if you're expecting to receive an array of JSON objects back and you wish to capture the value of an element in that object, then only the value of that field in the last object in the array will be captured.

TODO
----

- XML matching
- Capturing of outputs
- Pattern matching
