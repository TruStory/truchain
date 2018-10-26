'use strict';

Object.defineProperty(exports, "__esModule", {
    value: true
});

var _base = require('./base');

var _base2 = _interopRequireDefault(_base);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

const chalk = require('chalk');

const testIndent = '  ';
const suiteIndent = ' ';

class SpecsReporter extends _base2.default {

    constructor(launcher) {
        let opts = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};

        super(launcher, opts);

        this.launcher = launcher;
        this.isParallelRun = opts.parallel;
        let self = this;

        launcher.on('start', function (suiteFiles) {
            console.log();
            console.log(`Launcher will run suites: ${suiteFiles}`);
        });

        launcher.on('end', function () {
            console.log();
        });

        launcher.on('new suite', function (suite) {
            self.addSuite(suite);
        });
    }

    addSuite(suite) {
        super.addSuite(suite);

        suite.on('test pass', function (test) {
            console.log();
            console.log(chalk.green(`${testIndent} ${_base.INDICATORS.ok} ${test.name} (${test.duration}ms)`));
        });

        suite.on('test fail', function (test, error) {
            console.log();
            console.log(chalk.red(`${testIndent} ${_base.INDICATORS.err} ${test.name} (${test.duration}ms)`));
        });

        suite.on('test skip', function (test) {
            console.log();
            console.log(chalk.cyan(`${testIndent} - ${test.name} (${test.duration}ms)`));
        });

        suite.on('end', function (suite, error) {
            console.log();

            if (suite.status === 'pass') {
                console.log(chalk.green(`${suiteIndent} Done: ${suite.file} (Passed)`));
            } else if (suite.status === 'skip') {
                console.log(chalk.cyan(`${suiteIndent} Done: ${suite.file} (Skipped)`));
            } else {
                console.log(chalk.red(`${suiteIndent} Done: ${suite.file} (Failed)`));
            }
        });
    }

}
exports.default = SpecsReporter;
module.exports = exports['default'];