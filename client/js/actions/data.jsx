import _ from 'lodash';

const prefix = 'data';
const request = require('superagent');

export const SET_TUPLES = `${prefix}/SET_TUPLES`;
export function setTuples(json, stateDimensions) {
  const valueList = {};
  const selectedValue = {};
  _.map(stateDimensions, (value, key) => {
    if (key !== "region" && value !== "none") {
      const list = json.tuples.map((tuple) => {
        return tuple[value];
      });
      valueList[value] = _.uniq(list);
      selectedValue[value] = "";
    }
  });
  return {
    type: SET_TUPLES,
    tuples: json.tuples,
    valueList,
    selectedValue,
  };
}

export function fetchRequest() {
  return (dispatch, getState) => {
    const state = getState();

    return request
      .get('./api/request')
      .query(state.dimensions.stateDimensions)
      .end((err, res) => {
        dispatch(setTuples(res.body, state.dimensions.stateDimensions));
      });
  };
}

export const SWITCH_TAB = `${prefix}/SWITCH_TAB`;
export function switchTab(state) {
  return {
    type: SWITCH_TAB,
    state,
  };
}

export function fetchSwitchTab(state) {
  return (dispatch) => {
    return dispatch(switchTab(state));
  };
}

export const SELECT_VALUE = `${prefix}/SELECT_VALUE`;
export function selectValue(level, value) {
  return {
    type: SELECT_VALUE,
    level,
    value,
  };
}

export function fetchSelectValue(level, value) {
  return (dispatch) => {
    return dispatch(selectValue(level, value));
  };
}
