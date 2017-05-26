import * as ActionTypes from '../actions/data.jsx';

const initialState = {
  tuples: [],
  tabState: "table",
  valueList: {},
  selectedValue: {},
};

const data = (state = initialState, action) => {
  switch (action.type) {
    case ActionTypes.SET_TUPLES:
      return {
        ...state,
        tuples: action.tuples,
        valueList: action.valueList,
        selectedValue: action.selectedValue,
      };
    case ActionTypes.SWITCH_TAB:
      return {
        ...state,
        tabState: action.state,
      };
    case ActionTypes.SELECT_VALUE:
      return {
        ...state,
        selectedValue: {
          ...state.selectedValue,
          [action.level]: action.value,
        },
      };
    default:
      return state;
  }
};

export default data;
