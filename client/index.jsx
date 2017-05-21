import { render } from 'react-dom';
import ReactDOMServer from 'react-dom/server';
import * as React from 'react';
import { Provider } from 'react-redux';
import { Router, browserHistory, match, RouterContext } from 'react-router';
import Helmet from 'react-helmet';
import injectTapEventPlugin from "react-tap-event-plugin";
import { rootReducer, initialState } from './js/reducers/index.jsx';
import routes from './js/routes.jsx';
import { createStore, setAsCurrentStore } from './js/store.jsx';

require("babel-polyfill");
require('whatwg-fetch');

if (typeof window !== 'undefined') {
  // init promise polyfill
  window.Promise = window.Promise || Promise;
  // init fetch polyfill
  window.self = window;
  const store = createStore(window['--app-initial']);
  setAsCurrentStore(store);

  match({ history: browserHistory, routes }, (error, redirectLocation, renderProps) => {
    render(
      <Provider store={ store }>
        <Router { ...renderProps }/>
      </Provider>,
      document.getElementById('app'),
    );
  });
} else {
  Object.assign = null;
  Object.assign = require("object-assign");

  global.main = (options, callback) => {
    // console.log('render server side', JSON.stringify(options));
    const result = {
      uuid: options.uuid,
      app: null,
      title: null,
      meta: null,
      initial: null,
      error: null,
      redirect: null,
    };

    const store = createStore();
    setAsCurrentStore(store);

    try {
      match({ routes, location: options.url }, (error, redirectLocation, renderProps) => {
        try {
          if (error) {
            result.error = error;
          } else if (redirectLocation) {
            result.redirect = redirectLocation.pathname + redirectLocation.search;
          } else {
            result.app = ReactDOMServer.renderToString(
              <Provider store={ store }>
                <RouterContext { ...renderProps }/>
              </Provider>,
            );
            const { title, meta } = Helmet.rewind();
            result.title = title.toString();
            result.meta = meta.toString();
            result.initial = JSON.stringify(store.getState());
          }
        } catch (e) {
          result.error = e;
        }
        return callback(JSON.stringify(result));
      });
    } catch (e) {
      result.error = e;
      return callback(JSON.stringify(result));
    }
    return null;
  };
}
