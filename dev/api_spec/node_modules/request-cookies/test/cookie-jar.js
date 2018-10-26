var assert = require("assert")
var cookies = require('../index');

var url = 'http://www.example.com';

// these are simple tests
describe('CookieJar', function() {
  describe('create', function() {
    it("from the exported data of another jar", function(){
      var jar1 = new cookies.CookieJar();
      var c = new cookies.Cookie('foo=bar');
      jar1.add(c, url);
      var jarData = jar1.toJSON();

      var jar2 = new cookies.CookieJar(jarData);
      var header = jar2.getCookieHeaderString(url)
      assert(header === 'foo=bar');
    });
  });

  describe('add', function() {
    it("using a Cookie object", function(){
      var jar = new cookies.CookieJar();
      var c = new cookies.Cookie('foo=bar');
      jar.add(c, url);

      var header = jar.getCookieHeaderString(url)
      assert(header === 'foo=bar');
    });

    it("using a cookie string", function(){
      var jar = new cookies.CookieJar();
      jar.add('foo=bar', url);

      var header = jar.getCookieHeaderString(url)
      assert(header === 'foo=bar');
    });

    it("using a json object", function(){
      var jar = new cookies.CookieJar();
      jar.add({key:'foo', value:'bar'}, url);

      var header = jar.getCookieHeaderString(url);
      assert(header === 'foo=bar');
    });
  });

  describe('remove', function() {
    it('should not exist in jar', function(){
      var jar = new cookies.CookieJar();
      jar.add('foo=bar', url);

      var header = jar.getCookieHeaderString(url);
      assert(header === 'foo=bar');
      jar.remove(url);

      var header2 = jar.getCookieHeaderString(url);
      assert(header2 === '');
    });
  });

  describe('getCookies', function() {
    it('should find cookies that match the domain', function(){
      var jar = new cookies.CookieJar();
      jar.add('foo=bar', url);

      var matchingCookies = jar.getCookies(url);
      assert(matchingCookies.length === 1);
      assert(matchingCookies[0].getCookieHeaderString() === 'foo=bar');
    });
  });
});