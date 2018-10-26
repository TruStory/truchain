'use strict';

Object.defineProperty(exports, "__esModule", {
    value: true
});

var _fs = require('fs');

var _fs2 = _interopRequireDefault(_fs);

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

var _base = require('./base');

var _utils = require('./../utils');

var _he = require('he');

var _he2 = _interopRequireDefault(_he);

var _template = require('./html-src/template');

var _prettyMs = require('pretty-ms');

var _prettyMs2 = _interopRequireDefault(_prettyMs);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function loadCSS() {
    const stylesheet = _path2.default.join(__dirname, 'html-src/assets', 'html-report.css');
    return _fs2.default.readFileSync(stylesheet, 'utf-8');
}

function loadJS() {
    const js = _path2.default.join(__dirname, 'html-src/assets', 'html-report.js');
    return _fs2.default.readFileSync(js, 'utf-8');
}

class HTMLReporter {

    constructor(launcher) {
        let opts = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};

        this.launcher = launcher;
        this.isParallelRun = opts.parallel;
        this.reporterOptions = opts.reporterOptions;
        this.suites = [];
        let self = this;

        this.stats = {
            suites: 0,
            skippedSuites: 0,
            failedSuites: 0,
            passedSuites: 0,
            failedTests: 0,
            passedTests: 0,
            skippedTests: 0,
            tests: 0,
            start: process._api.startTime || new Date()
        };

        launcher.on('start', function (suiteFiles) {
            self.plannedSuites = suiteFiles.length;
        });

        launcher.on('end', function () {
            self.stats.end = new Date();
            self.stats.duration = (0, _prettyMs2.default)(self.stats.end.getTime() - self.stats.start.getTime());
            let consolidatedMarkupForAllSuites = '';
            self.suites.forEach(suite => {
                consolidatedMarkupForAllSuites += suite.html;
            });

            const html = (0, _template.buildHTML)({
                css: loadCSS(),
                js: loadJS(),
                report: consolidatedMarkupForAllSuites,
                stats: self.stats
            });

            try {
                let outputDirectory = this.reporterOptions.htmlReportDir || '';
                if (!(0, _utils.doesDirectoryExist)(outputDirectory)) {
                    console.error(`Directory ${outputDirectory} does not exist, cannot write html report`);
                    return;
                }

                let outputFile = this.reporterOptions.htmlReportName ? this.reporterOptions.htmlReportName + '.html' : 'report.html';
                const reportPath = _path2.default.resolve(process.cwd(), outputDirectory, outputFile);

                _fs2.default.writeFileSync(reportPath, html);
                console.log(`\nHTML report is written to "${reportPath}"`);
            } catch (err) {
                console.error(err);
            }
        });

        launcher.on('new suite', function (suite) {
            self.stats.suites++;
            self.addSuite(suite);
        });
    }

    addSuite(suite) {
        let self = this;

        let html = '';
        const suiteTitle = (0, _utils.escapeHTML)(suite.file);

        html += '<li class="suite">';
        html += `<h1 class="suite-name">${suiteTitle}</h1>`;
        html += '<ul>';

        this.suites.push({
            location: suite.file,
            status: null,
            html: html,
            error: null
        });

        suite.on('test pass', function (test) {
            let suite = self.getSuite(test.suite.file);
            const title = (0, _utils.escapeHTML)(test.name);

            suite.html += `<li class="test pass">`;
            suite.html += `<h2 class="test-name">${title} <span class="duration">( ${test.duration}ms )</span></h2>`;
            suite.html += '</li>';

            self.stats.passedTests++;
            self.stats.tests++;
        });

        suite.on('test fail', function (test, error) {
            let suite = self.getSuite(test.suite.file);
            const err = error;
            const title = (0, _utils.escapeHTML)(test.name);

            suite.html += '<li class="test fail">';
            suite.html += `<h2 class="test-name" onclick="toggleTestDetailedView(this)">${title} <span class="duration">( ${test.duration}ms )</span></h2>`;
            suite.html += `<pre class="error">${err.stack}</pre>`;

            if (self.reporterOptions.logRequests) {
                try {
                    if (test.requests && test.requests.length) {
                        let requestsLog = '';
                        for (let loggedRequestResponse of test.requests) {
                            let prettyLog = (0, _utils.prettifyRequestLog)(loggedRequestResponse);
                            prettyLog += '----------------------------------- \n';
                            requestsLog += prettyLog;
                        }
                        requestsLog = _he2.default.escape(requestsLog);
                        suite.html += `<pre class="test-requests">${requestsLog}</pre>`;
                    }
                } catch (e) {}
            }

            suite.html += '</li>';

            self.stats.failedTests++;
            self.stats.tests++;
        });

        suite.on('test skip', function (test) {
            let suite = self.getSuite(test.suite.file);

            const title = (0, _utils.escapeHTML)(test.name);
            const duration = test.duration || 1;

            suite.html += `<li class="test skip">`;
            suite.html += `<h2 class="test-name"">${title} <span class="duration">( ${test.duration}ms )</span></h2>`;
            suite.html += '</li>';

            self.stats.skippedTests++;
            self.stats.tests++;
        });

        suite.on('end', function (suite, error) {
            let suiteResultObject = self.getSuite(suite.file);
            suiteResultObject.status = suite.status;
            suiteResultObject.error = error || null;

            if (suite.status === 'fail') {
                if (error) {
                    suiteResultObject.html += `<pre class="suite-error">${error.stack}</pre>`;
                }

                if (self.reporterOptions.logRequests) {
                    try {
                        if (suite.requests && suite.requests.length) {
                            let requestsLog = '';

                            for (let loggedRequestResponse of suite.requests) {
                                let prettyLog = (0, _utils.prettifyRequestLog)(loggedRequestResponse);
                                prettyLog += '----------------------------------- \n';
                                requestsLog += prettyLog;
                            }

                            requestsLog = _he2.default.escape(requestsLog);
                            suiteResultObject.html += `<pre class="suite-requests">${requestsLog}</pre>`;
                        }
                    } catch (e) {}
                }

                self.stats.failedSuites++;
            } else if (suite.status === 'pass') {
                self.stats.passedSuites++;
            } else if (suite.status === 'skip') {
                self.stats.skippedSuites++;
            }

            if (suite.status === 'skip') {
                let html = '';
                html += '<li class="suite skip">';
                html += `<h1 class="suite-name">${(0, _utils.escapeHTML)(suite.file)} - Skipped</h1>`;
                html += '<ul>';

                if (suite.name && suite.name.length > 0) {
                    html += `<span class="suite-name-hidden">${suite.name}</span>`;
                }

                html += '</ul></li>';
                suiteResultObject.html = html;
            } else {
                if (suite.name && suite.name.length > 0) {
                    suiteResultObject.html += `<span class="suite-name-hidden">${suite.name}</span>`;
                }
                suiteResultObject.html += '</ul>';
                suiteResultObject.html += '</li>';
            }
        });
    }

    getSuite(location) {
        return this.suites.find(function (suiteObj) {
            return suiteObj.location === location;
        });
    }

}
exports.default = HTMLReporter;
module.exports = exports['default'];