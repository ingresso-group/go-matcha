go-matcha
=========

Features
--------

### Data formats

There are assertion methods for both JSON and XML.

Note that XML matching is currently fairly naïve in that it doesn't read XML schemas or check attributes. One particular limitation of this is that if you are expecting an array of elements back, and in the actual XML there is only one element in the array, the assertion will fail (since in the absence of a schema it is impossible to know if it is an array with one element or just a single element).

### Capturing values from the response

If you define a field with a `capture` tag then that field will be captured from the response. This is useful for more complex assertions.

For example, the following field would be captured as 'count':

```
Count    float64 `capture:"count"`
```

Captured values come back as a map of slices (type = `matcha.CapturedValues`). This is so that multiple values can be captured. For example, if we are expecting a list of objects, each with a "date" field, then `capturedValues["date"]` will be a list of dates, one for each object.

### Regex pattern matching

You can use a `pattern` tag on a string field if you are expecting to it to match a given regex.

For example if you have a date field such as `2016-09-12` then you could use the following field definition:

```
Date    string `pattern:"^[0-9]{4}-[0-9]{2}-[0-9]{2}"`
```

Note that it is not possible to use "complex" string literals in Go struct tags, therefore it is not possible to use some characters, such as `\`.

