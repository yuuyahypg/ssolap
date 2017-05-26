import { combineReducers } from 'redux';
import dimensions from './dimensions.jsx';
import data from './data.jsx';
import map from './map.jsx';

const rootReducer = combineReducers({
  dimensions,
  data,
  map,
});

export default rootReducer;
