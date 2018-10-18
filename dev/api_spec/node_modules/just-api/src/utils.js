'use strict';

Object.defineProperty(exports, "__esModule", {
    value: true
});
exports.doesDirectoryExist = doesDirectoryExist;
exports.findSuiteFiles = findSuiteFiles;
exports.assertFileValidity = assertFileValidity;
exports.loadModule = loadModule;
exports.runModuleFunction = runModuleFunction;
exports.runInlineFunction = runInlineFunction;
exports.convertMillisToHumanReadableFormat = convertMillisToHumanReadableFormat;
exports.wait = wait;
exports.isNumber = isNumber;
exports.escapeHTML = escapeHTML;
exports.prettifyRequestLog = prettifyRequestLog;
exports.equals = equals;

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

var _fs = require('fs');

var _fs2 = _interopRequireDefault(_fs);

var _glob = require('glob');

var _glob2 = _interopRequireDefault(_glob);

var _he = require('he');

var _he2 = _interopRequireDefault(_he);

var _errors = require('./errors');

var _isEqual = require('lodash/isEqual');

var _isEqual2 = _interopRequireDefault(_isEqual);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function doesDirectoryExist(dirPath) {
    const fullPath = _path2.default.resolve(process.cwd(), dirPath);

    try {
        const stat = _fs2.default.statSync(fullPath);
        return stat.isDirectory();
    } catch (err) {
        return false;
    }
}

function findSuiteFiles(basePath) {
    let digDeep = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : false;
    let extensions = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : ['yml', 'yaml'];

    let files = [];

    if (!_fs2.default.existsSync(basePath)) {

        if (_fs2.default.existsSync(basePath + '.yml')) {
            basePath += '.yml';
        } else if (_fs2.default.existsSync(basePath + '.yaml')) {
            basePath += '.yaml';
        } else {
            files = _glob2.default.sync(basePath);

            if (!files.length) {
                throw new Error(`No suites found using path/pattern ${basePath}`);
            }
            return files;
        }
    }

    try {
        let stat = _fs2.default.statSync(basePath);

        if (stat.isFile()) {
            return basePath;
        }
    } catch (err) {
        return;
    }

    _fs2.default.readdirSync(basePath).forEach(function (fileOrDir) {
        let file = _path2.default.join(basePath, fileOrDir);

        try {
            var stat = _fs2.default.statSync(file);

            if (stat.isDirectory()) {
                if (digDeep) {
                    files = files.concat(findSuiteFiles(file, digDeep, extensions));
                }
                return;
            }
        } catch (err) {
            return;
        }

        let re = new RegExp('\\.(?:' + extensions.join('|') + ')$');

        if (!stat.isFile() || !re.test(file) || _path2.default.basename(file)[0] === '.') {
            return;
        }

        files.push(file);
    });

    return files;
}

function assertFileValidity(filePath, fileContext) {
    const absPath = filePath;

    if (!_fs2.default.existsSync(absPath)) {
        const FileDoesNotExistError = (0, _errors.customError)('FileDoesNotExistError');
        throw new FileDoesNotExistError(`${fileContext} file doesn't exist at '${absPath}'`);
    }

    if (!_fs2.default.statSync(absPath).isFile()) {
        throw new Error(`${fileContext} at: ${absPath} is not a file`);
    }

    return absPath;
}

function loadModule(modulePath) {
    try {
        return require(modulePath);
    } catch (e) {
        throw e;
    }
}

async function runModuleFunction(module, fnName, context, args) {
    let CustomFunctionNotFoundInModuleError = (0, _errors.customError)('CustomFunctionNotFoundInModuleError');
    let NotAFunctionError = (0, _errors.customError)('NotAFunctionError');

    try {
        let func = module[fnName];

        if (!func) {
            throw new CustomFunctionNotFoundInModuleError(`'${fnName}' function not found in module`);
        }

        if (typeof func !== 'function') {
            throw new NotAFunctionError(`'${fnName}', Provide valid javascript function`);
        }

        let result = await module[fnName].call(context);

        return result;
    } catch (error) {
        throw error;
    }
}

async function runInlineFunction(fn, context, args) {
    let NotAFunctionError = (0, _errors.customError)('NotAFunctionError');

    if (typeof fn !== 'function') {
        throw new NotAFunctionError(`'${fn}' is not a function, Provide valid inline javascript function`);
    }

    try {
        let result = await fn.call(context);
        return result;
    } catch (error) {
        throw error;
    }
}

function convertMillisToHumanReadableFormat(duration) {
    let milliseconds = parseInt(duration % 1000);
    let seconds = parseInt(duration / 1000 % 60);
    let minutes = parseInt(duration / (1000 * 60) % 60);
    let hours = parseInt(duration / (1000 * 60 * 60) % 24);

    if (hours === 0) {
        return minutes + "m" + seconds + "s." + milliseconds + "ms";
    }

    return hours + "h" + minutes + "m" + seconds + "s." + milliseconds + "ms";
}

async function wait(durationInMillis) {
    return new Promise(resolve => setTimeout(resolve, durationInMillis));
}

function isNumber(number) {
    return !isNaN(parseFloat(number)) && isFinite(number);
}

function escapeHTML(html) {
    return _he2.default.escape(String(html));
}

function prettifyRequestLog(reqResInfo) {
    let result = '';

    const request = reqResInfo.request;
    const response = reqResInfo.response;
    const error = reqResInfo.error;

    result += 'Request: \n\n';
    result += `${request.method.toUpperCase()} ${request.uri} \n`;

    for (let headerKey in request.headers) {
        result += `${headerKey}: ${request.headers[headerKey]}\n`;
    }

    result += '\n';

    if (request.body && request.formRequest) {
        result += "It's a form/multipart-form request, this may or may not be the actual raw body \n";
        result += JSON.stringify(request.body) + '\n';
    }

    if (request.body && !request.formRequest) {
        result += request.body + '\n';
    }

    if (error) {
        result += '\n--Encountered following error \n\n';
        result += `${error}`;
        result += '\n';
    } else {
        result += '\nResponse: \n\n';
        result += `Status code: ${response.statusCode} \n`;

        for (let headerKey in response.headers) {
            result += `${headerKey}: ${response.headers[headerKey]}\n`;
        }

        result += '\n';

        if (response.headers['content-type'].indexOf('application/json') !== -1) {
            try {
                //TODO send json as pretty multiline string so it's easy to read
                result += response.body.toString();
            } catch (e) {
                result += response.body.toString();
            }
        } else {
            result += response.body.toString();
        }

        result += '\n';
        result += `\nRequest duration: ${response.timings.total}ms \n`;
    }

    return result;
}

function equals(value, other) {
    return (0, _isEqual2.default)(value, other);
}