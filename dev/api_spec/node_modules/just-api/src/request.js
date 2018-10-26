'use strict';

Object.defineProperty(exports, "__esModule", {
    value: true
});

var _extends = Object.assign || function (target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i]; for (var key in source) { if (Object.prototype.hasOwnProperty.call(source, key)) { target[key] = source[key]; } } } return target; };

var _slicedToArray = function () { function sliceIterator(arr, i) { var _arr = []; var _n = true; var _d = false; var _e = undefined; try { for (var _i = arr[Symbol.iterator](), _s; !(_n = (_s = _i.next()).done); _n = true) { _arr.push(_s.value); if (i && _arr.length === i) break; } } catch (err) { _d = true; _e = err; } finally { try { if (!_n && _i["return"]) _i["return"](); } finally { if (_d) throw _e; } } return _arr; } return function (arr, i) { if (Array.isArray(arr)) { return arr; } else if (Symbol.iterator in Object(arr)) { return sliceIterator(arr, i); } else { throw new TypeError("Invalid attempt to destructure non-iterable instance"); } }; }();

var _response = require('./response');

var _response2 = _interopRequireDefault(_response);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

const request = require('request');
class JustAPIRequest {
    constructor(options, test) {
        this.options = options;
        this.test = test;
    }

    buildCookieJar() {
        if (this.options.cookiesData && Object.keys(this.options.cookiesData).length > 0) {

            let cookieJar = request.jar();

            for (let _ref of Object.entries(this.options.cookiesData)) {
                var _ref2 = _slicedToArray(_ref, 2);

                let key = _ref2[0];
                let value = _ref2[1];

                cookieJar.setCookie(request.cookie(`${key}=${value}`), this.options.baseUrl || this.options.url);
            }

            this.options.jar = cookieJar;
        }
    }

    async send() {
        const self = this;
        self.options.time = true;
        self.buildCookieJar();

        try {
            let response = await new Promise(function (resolve, reject) {
                request(self.options, function (err, response, body) {
                    let log = {};

                    log.request = {
                        uri: this.uri.href,
                        method: this.method,
                        headers: _extends({}, this.headers)
                    };

                    if (this.body) {
                        log.request.body = this.body.toString('utf8');
                    } else if (this.formData) {
                        log.request.formRequest = true;
                        log.request.body = _extends({}, this.formData);
                    }

                    if (err !== null) {
                        log.error = err;
                        self.test.suite.requests.push(log);
                        return reject(err);
                    }

                    log.response = {
                        headers: _extends({}, response.headers),
                        statusCode: response.statusCode,
                        body: body,
                        timings: _extends({}, response.timingPhases)
                    };

                    self.test.suite.requests.push(log);

                    let justAPIResponse = new _response2.default(this, err, response, body);
                    return resolve(justAPIResponse);
                });
            });

            return response;
        } catch (error) {
            throw error;
        }
    }

}
exports.default = JustAPIRequest;
module.exports = exports['default'];