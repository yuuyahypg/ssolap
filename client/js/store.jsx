import { applyMiddleware, createStore as reduxCreateStore } from 'redux';
import thunk from 'redux-thunk';
import rootReducer from './reducers/index.jsx';

const middlewares = [thunk];

export function createStore(state) {
  return reduxCreateStore(
    rootReducer,
    state,
    // applyMiddleware.apply(...middlewares),
    applyMiddleware(thunk),
  );
}

export let store = null;
export function getStore() { return store; }
export function setAsCurrentStore(s) {
  store = s;
  if (typeof window !== 'undefined') {
    window.store = store;
  }
}
