request-cookies [![Build Status](https://travis-ci.org/lalitkapoor/request-cookies.png?branch=master)](https://travis-ci.org/lalitkapoor/request-cookies)
==================

Cookie management for node's request library

see tests for examples


<!-- Start lib/cookie-jar.js -->

# CookieJar

## CookieJar()

Holds cookies

## add(data, url, options)

Add a cookie to the cookie jar

### Params:

* **String|Object** *data* can be a Cookie, json object, or a set-cookie header string

* **String** *url*

* **Object** *options* [optional]

## remove(url, key)

Remove a cookie by name for a given domain. If the key then all cookies for
the given domain and path (if given) will be removed. Note: when cookies
are added without a path specified the default path is &quot;/&quot; in accordance
with http://tools.ietf.org/search/rfc6265#section-5.1.4 - so you should use
that for the path param when appropriate.

### Params:

* **String** *url*

* **String** *key*

## getCookies(url, options)

Get cookies that match the given properties

### Params:

* **String** *url*

* **Object** *options* [optional]

## getCookieHeaderString(url, options)

Get HTTP Cookie header string

### Params:

* **String** *url*

* **Object** *options* [optional]

## toJSON()

Get cookie jar as a JSON object

### Return:

* **Object**

<!-- End lib/cookie-jar.js -->

<!-- Start lib/cookie.js -->

# Cookie

## Cookie(data)

HTTP Cookie

### Params:

* **String|Object** *data* can be a cookie header string or a json object

## getCookieHeaderString()

Get cookie as a cookie header string

### Return:

* **String**

## set(data)

Set multiple cookie properties at once

### Params:

* **String|Object** *data*

## toJSON()

Get cookie as a JSON object

### Return:

* **Object**

<!-- End lib/cookie.js -->

<!-- generated using markdox -->
