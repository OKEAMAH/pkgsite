'use strict';
/*!
 * @license
 * Copyright 2019-2020 The Go Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
var _a;
const canonicalURLPath =
  (_a = document.querySelector('.js-canonicalURLPath')) === null || _a === void 0
    ? void 0
    : _a.dataset['canonicalUrlPath'];
if (canonicalURLPath && canonicalURLPath !== '') {
  document.addEventListener('keydown', e => {
    var _a, _b;
    const t = (_a = e.target) === null || _a === void 0 ? void 0 : _a.tagName;
    if (t === 'INPUT' || t === 'SELECT' || t === 'TEXTAREA') {
      return;
    }
    if ((_b = e.target) === null || _b === void 0 ? void 0 : _b.isContentEditable) {
      return;
    }
    if (e.metaKey || e.ctrlKey) {
      return;
    }
    switch (e.key) {
      case 'y':
        window.history.replaceState(null, '', canonicalURLPath);
        break;
    }
  });
}
//# sourceMappingURL=keyboard.js.map
