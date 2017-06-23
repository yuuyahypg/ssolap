import * as ActionTypes from '../actions/dimensions.jsx';

const initialState = {
  rootDimensions: [],
  fact: {},
  isDBConnected: false,
  stateDimensions: {},
  columns: [],
};

const dimensions = (state = initialState, action) => {
  switch (action.type) {
    case ActionTypes.SET_DIMENSIONS:
      return {
        ...state,
        rootDimensions: action.dimensions,
        fact: action.fact,
        isDBConnected: action.isDBConnected,
        stateDimensions: action.stateDimensions,
        columns: action.columns,
      };
    case ActionTypes.SET_LEVEL:
      return {
        ...state,
        stateDimensions: action.stateDimensions,
        columns: action.columns,
      };
    default:
      return state;
  }
};

export default dimensions;
