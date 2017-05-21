import * as ActionTypes from '../actions/dimensions.jsx';

const initialState = {
  rootDimensions: [],
  fact: {},
  stateDimensions: {},
};

const dimensions = (state = initialState, action) => {
  switch (action.type) {
    case ActionTypes.SET_DIMENSIONS:
      return {
        ...state,
        rootDimensions: action.dimensions,
        fact: action.fact,
        stateDimensions: action.stateDimensions,
      };
    case ActionTypes.SET_LEVEL:
      return {
        ...state,
        stateDimensions: action.dimensions,
      };
    default:
      return state;
  }
};

export default dimensions;
