'use strict';

Object.defineProperty(exports, "__esModule", {
    value: true
});

var _fs = require('fs');

var _fs2 = _interopRequireDefault(_fs);

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

var _glob = require('glob');

var _glob2 = _interopRequireDefault(_glob);

var _suite = require('./suite');

var _suite2 = _interopRequireDefault(_suite);

var _launcher = require('./launcher');

var _launcher2 = _interopRequireDefault(_launcher);

var _utils = require('./utils');

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

const MAX_PARALLEL_LIMIT = 24;

class JustAPI {

    constructor(suiteFiles, options) {
        this.suites = suiteFiles;
        this.options = options;
        this.max_parallel_limit = MAX_PARALLEL_LIMIT;
    }

    init() {

        if (!this.suites.length) {
            console.error('No test suites found');
            process.exit(1);
        } else {
            this.files = this.suites.map(suite => _path2.default.resolve(suite));
            let files = this.files.map(file => `  - ${file}`);
            console.info(`Found suites:\n${files.join('\n')}`);
        }

        this.setupReporters(this.options.reporter);
        this.setReporterOptions(this.options.reporterOptions);
        this.configureParallelism(this.options.parallel);
        this.configurePatternMatching(this.options.grep);

        const launcherOptions = {
            parallel: this.parallel,
            max_parallel: this.max_parallel,
            max_parallel_limit: this.max_parallel_limit,
            grep: this.grep,
            reporters: this.reporters,
            reporterOptions: this.reporterOptions
        };

        const launcher = new _launcher2.default(this.files, launcherOptions);
        launcher.go();
    }

    setupReporters(input) {
        this.reporters = [];

        const reporterIdentified = input || 'specs';
        let chosenReporters = reporterIdentified.split(',');

        if (!(chosenReporters.indexOf('specs') !== -1)) {
            chosenReporters.push('specs');
        }

        chosenReporters = Array.from(new Set(chosenReporters));

        for (let reporterInput of chosenReporters) {
            let Reporter;

            switch (reporterInput.trim()) {
                case 'specs':
                    Reporter = require('./reporters/specs.js');
                    break;
                case 'html':
                    Reporter = require('./reporters/html.js');
                    break;
                case 'json':
                    Reporter = require('./reporters/json.js');
                    break;
                default:
                    try {
                        Reporter = require(reporterInput.trim());
                        break;
                    } catch (err) {
                        console.error(`\nError occurred while loading the reporter: ${reporterInput.trim()}`);
                        console.error(`\n  ${err.stack}`);
                        process.exit(1);
                    }
            }

            this.reporters.push(Reporter);
        }
    }

    setReporterOptions(input) {
        this.reporterOptions = {};

        if (input !== undefined) {
            let options = input.split(',');

            for (let option of options) {
                const item = option.split('=');

                if (item.length > 2 || item.length === 0) {
                    throw new Error(`invalid reporter option '${option}'`);
                } else if (item.length === 2) {
                    this.reporterOptions[item[0]] = item[1];
                } else {
                    this.reporterOptions[item[0]] = true;
                }
            }
        }
    }

    configureParallelism(input) {
        const userInputForParallelism = input;

        if (userInputForParallelism) {
            if ((0, _utils.isNumber)(userInputForParallelism) && parseInt(userInputForParallelism) > 1 && parseInt(userInputForParallelism) <= MAX_PARALLEL_LIMIT) {
                this.parallel = true;
                this.max_parallel = parseInt(userInputForParallelism);
            } else {
                console.error(`Given argument for parallel option ${input} is invalid, Please provide a number (> 1 and <= ${MAX_PARALLEL_LIMIT})`);
                process.exit(1);
            }
        } else {
            this.parallel = false;
            this.max_parallel = 1;
        }
    }

    configurePatternMatching(input) {
        const grepRegexInput = input;

        if (grepRegexInput) {
            if (typeof grepRegexInput === 'string') {
                let arg = grepRegexInput.match(/^\/(.*)\/(g|i|)$|.*/);
                this.grep = new RegExp(arg[1] || arg[0], arg[2]);
            } else {
                this.grep = grepRegexInput;
            }
        }
    }
}
exports.default = JustAPI;
module.exports = exports['default'];