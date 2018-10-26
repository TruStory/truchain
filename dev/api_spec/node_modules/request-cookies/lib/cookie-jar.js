var URL = require('url');
var tc = require('tough-cookie');
var Cookie = require('./cookie');

var noop = function(){};


/**
 * Holds cookies
 * @param {Object} store [optional] the toJSON value of another CookieJar
 * @constructor
 */
var CookieJar = function(store) {
  this._jar = new tc.CookieJar();
  if (store) this._jar.store.idx = store;
};

/**
 * Add a cookie to the cookie jar
 * @param {String|Object} data can be a Cookie|json object or a set-cookie
 *                        header string
 * @param {String} url
 * @param {Object} options [optional]
 */
CookieJar.prototype.add = function(data, url, options) {
  if (!options) options = {};

  if (data instanceof Cookie) {
    this._jar.setCookieSync(data._cookie, url, options);
  } else if (typeof data === 'object' || typeof data === 'string') {
    this._jar.setCookieSync(new Cookie(data)._cookie, url, options);
  }
};

/**
 * Remove a cookie by name for a given domain. If the key then all cookies for
 * the given domain and path (if given) will be removed. Note: when cookies
 * are added without a path specified the default path is "/" in accordance
 * with http://tools.ietf.org/search/rfc6265#section-5.1.4 - so you should use
 * that for the path param when appropriate.
 * @param  {String} url
 * @param  {String} key
 */
CookieJar.prototype.remove = function(url, key) {
  var domain = null;
  var path ='/';

  var urlInfo = URL.parse(url);

  if (!urlInfo.protocol) {
    domain = urlInfo;
  } else {
    domain = urlInfo.host;
    path = urlInfo.path;
  }

  if (key) {
    this._jar.store.removeCookie(domain, path, key, noop);
  } else {
    this._jar.store.removeCookies(domain, path, noop);
  }
};

/**
 * Get cookies that match the given properties
 * @param  {String} url
 * @param  {Object} options [optional]
 * @return {Array} array of Cookie objects
 */
CookieJar.prototype.getCookies = function(url, options) {
  if (!options) options = {};

  // get tough cookies, wrap w/Cookie, return array of Cookie objects
  var tcCookies = this._jar.getCookiesSync(url, options);
  var cookies = [];
  for (var i=0, len = tcCookies.length; i<len; i++) {
    cookies.push(new Cookie(tcCookies[i]));
  };
  return cookies;
};


/**
 * Get HTTP Cookie header string
 * @param  {String}   url
 * @param  {Object} options [optional]
 * @return {String} HTTP Cookie header string
 */
CookieJar.prototype.getCookieHeaderString = function(url, options) {
  if (!options) options = {};

  return this._jar.getCookieStringSync(url, options);
};

/**
 * Get cookie jar as a JSON object
 * @return {Object}
 */
CookieJar.prototype.toJSON = function() {
  return this._jar.store.idx;
};

module.exports = CookieJar;