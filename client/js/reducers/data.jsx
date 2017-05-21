import * as ActionTypes from '../actions/data.jsx';

const initialState = {
  tuples: [],
};

const data = (state = initialState, action) => {
  switch (action.type) {
    case ActionTypes.SET_TUPLES:
      return {
        ...state,
      };
    default:
      return state;
  }
};

export default data;
