var assert = require("assert")
var cookies = require('../index');

// these are simple tests
describe('Cookie', function() {
  describe('create', function() {
    it("from string", function(){
      var c = new cookies.Cookie('foo=bar');
      assert(c.getCookieHeaderString() === 'foo=bar');
    });

    it("from json", function(){
      var c = new cookies.Cookie({key:'foo', value:'bar'});
      assert(c.getCookieHeaderString() === 'foo=bar');
    });
  });

  describe('set', function() {
    it("should update k/v pair", function(){
      var c = new cookies.Cookie('foo=bar');
      c.set({key:'abc', value:'def'})
      assert(c.getCookieHeaderString() === 'abc=def');
    });
  });

  describe('toJSON', function() {
    it("should output a json representation", function(){
      var c = new cookies.Cookie('foo=bar');
      assert(c.toJSON().key === 'foo');
      assert(c.toJSON().value === 'bar');
    });
  });
});