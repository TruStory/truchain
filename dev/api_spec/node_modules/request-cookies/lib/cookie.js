var tc = require('tough-cookie');

var noop = function(){};

var modifiable = [
  'key'
, 'value'
, 'expires'
, 'maxAge'
, 'domain'
, 'path'
, 'secure'
, 'httpOnly'
, 'extensions'
];

var readOnly = [
  'hostOnly'
, 'creation'
, 'lastAccessed'
, 'pathIsDefault'
];


/**
 * HTTP Cookie
 * @constructor
 * @param {String|Object} data can be a cookie header string or a json object
 */
var Cookie = function(data) {
  var self = this;
  this._cookie = new tc.Cookie();

  modifiable.map(function(attribute) {
    Object.defineProperty(self, attribute, {
    enumerable: true
    , get: function() {
        return self._cookie[attribute];
      }
    , set: function(val) {
        self._cookie[attribute] = val;
      }
    });
  });

  readOnly.map(function(attribute) {
    Object.defineProperty(self, attribute, {
      get: function(){
        return this._cookie[attribute];
      }
    , set: noop
    , enumerable: true
    });
  });

  this.set(data);
};

/**
 * Get cookie as a cookie header string
 * @return {String}
 */
Cookie.prototype.getCookieHeaderString = function() {
  return this._cookie.cookieString();
};

/**
 * Set multiple cookie properties at once
 * @param {String|Object} data
 */
Cookie.prototype.set = function(data) {
  if (!data) return;

  if (typeof data === 'string') {
    this._cookie = tc.Cookie.parse(data);
    return;
  }

  if (typeof data === 'object') {
    for (key in data) {
      if (!data.hasOwnProperty(key)) continue;
      if (modifiable.indexOf(key) < 0) continue;
      this[key] = data[key];
    }
    return;
  }

  throw Error("data parameter must be string or object.");
};

/**
 * Get cookie as a JSON object
 * @return {Object}
 */
Cookie.prototype.toJSON = function() {
  return {
    key: this.key
  , value: this.value
  , expires: this.expires
  , maxAge: this.maxAge
  , domain: this.domain
  , path: this.path
  , secure: this.secure
  , httpOnly: this.httpOnly
  , extensions: this.extensions
  };
};

module.exports = Cookie;