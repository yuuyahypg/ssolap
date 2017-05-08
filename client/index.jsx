import { render } from 'react-dom';
import ReactDOMServer from 'react-dom/server';
import * as React from 'react';
import App from 'js/application.jsx';

require("babel-polyfill");
require('whatwg-fetch');

if (typeof window !== 'undefined') {
  render(<App />, document.getElementById('app'));
} else {
  global.main = (options, callback) => {
    // console.log('render server side', JSON.stringify(options))
    const s = ReactDOMServer.renderToString(React.createElement(App, {}));

    callback(JSON.stringify({
      uuid: options.uuid,
      app: s,
      title: null,
      meta: null,
      initial: null,
      error: null,
      redirect: null,
    }));
  };
}
