/******/ (() => { // webpackBootstrap
/******/  "use strict";
/******/  // The require scope
/******/  var __webpack_require__ = {};
/******/  
/************************************************************************/
/******/  /* webpack/runtime/compat get default export */
/******/  (() => {
/******/    // getDefaultExport function for compatibility with non-harmony modules
/******/    __webpack_require__.n = (module) => {
/******/      var getter = module && module.__esModule ?
/******/        () => (module['default']) :
/******/        () => (module);
/******/      __webpack_require__.d(getter, { a: getter });
/******/      return getter;
/******/    };
/******/  })();
/******/  
/******/  /* webpack/runtime/define property getters */
/******/  (() => {
/******/    // define getter functions for harmony exports
/******/    __webpack_require__.d = (exports, definition) => {
/******/      for(var key in definition) {
/******/        if(__webpack_require__.o(definition, key) && !__webpack_require__.o(exports, key)) {
/******/          Object.defineProperty(exports, key, { enumerable: true, get: definition[key] });
/******/        }
/******/      }
/******/    };
/******/  })();
/******/  
/******/  /* webpack/runtime/hasOwnProperty shorthand */
/******/  (() => {
/******/    __webpack_require__.o = (obj, prop) => (Object.prototype.hasOwnProperty.call(obj, prop))
/******/  })();
/******/  
/******/  /* webpack/runtime/make namespace object */
/******/  (() => {
/******/    // define __esModule on exports
/******/    __webpack_require__.r = (exports) => {
/******/      if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/        Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/      }
/******/      Object.defineProperty(exports, '__esModule', { value: true });
/******/    };
/******/  })();
/******/  
/************************************************************************/
var __webpack_exports__ = {};
// ESM COMPAT FLAG
__webpack_require__.r(__webpack_exports__);

;// CONCATENATED MODULE: external "express"
const external_express_namespaceObject = require("express");
var external_express_default = /*#__PURE__*/__webpack_require__.n(external_express_namespaceObject);
;// CONCATENATED MODULE: external "body-parser"
const external_body_parser_namespaceObject = require("body-parser");
var external_body_parser_default = /*#__PURE__*/__webpack_require__.n(external_body_parser_namespaceObject);
;// CONCATENATED MODULE: external "fs"
const external_fs_namespaceObject = require("fs");
var external_fs_default = /*#__PURE__*/__webpack_require__.n(external_fs_namespaceObject);
;// CONCATENATED MODULE: external "cors"
const external_cors_namespaceObject = require("cors");
var external_cors_default = /*#__PURE__*/__webpack_require__.n(external_cors_namespaceObject);
;// CONCATENATED MODULE: external "path"
const external_path_namespaceObject = require("path");
var external_path_default = /*#__PURE__*/__webpack_require__.n(external_path_namespaceObject);
;// CONCATENATED MODULE: ./src/backend.ts





const app = external_express_default()();
const port = 3000;
app.use(external_body_parser_default().json());
app.use(external_cors_default()());
app.post('/generate-css', (req, res) => {
    const { logo, main_color, background_color } = req.body;
    const generatedCSS = `
  .sidebar__logo__character {
    visibility: hidden; 
  }
  .sidebar__logo::before {
    content:"";
    position: absolute;
    left: 20px;
    top: 30px;
    width: 40px; 
    height: 40px; 
    background: url("${logo}");
    background-size: 40px 40px;
  }
    .sidebar {  
      background-color: ${background_color};
    }
    .sidebar__nav-item {
      color: ${main_color};
    }
  `;
    const outputDir = process.env.PATHS || '/shared/app/css';
    const fileName = 'generated.css';
    const filePath = external_path_default().join(outputDir, fileName);
    external_fs_default().writeFile(filePath, generatedCSS, (err) => {
        if (err) {
            console.error(err);
            res.status(500).json({ error: 'Error generating and saving CSS' });
        }
        else {
            console.log('CSS generated and saved successfully!');
            res.json({ success: true });
        }
    });
});
app.get('/generated-css', (req, res) => {
    const filePath = './packages/kndp/charts/kndp/files/generated.css';
    external_fs_default().readFile(filePath, 'utf8', (err, data) => {
        if (err) {
            console.error(err);
            res.status(500).json({ error: 'Error reading CSS file' });
        }
        else {
            console.log('CSS file read successfully!');
            res.send(data);
        }
    });
});
app.listen(port, () => {
    console.log(`Backend API listening at http://localhost:${port}`);
});

var __webpack_export_target__ = exports;
for(var i in __webpack_exports__) __webpack_export_target__[i] = __webpack_exports__[i];
if(__webpack_exports__.__esModule) Object.defineProperty(__webpack_export_target__, "__esModule", { value: true });
/******/ })()
;