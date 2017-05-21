import _ from 'lodash';

const prefix = 'dimensions';

export const SET_DIMENSIONS = `${prefix}/SET_DIMENSIONS`;
export function setDimensions(json) {
  return {
    type: SET_DIMENSIONS,
    dimensions: json.dimensions,
    fact: json.fact,
    stateDimensions: getStateDimensions(json.dimensions, json.fact),
  };
}

function getStateDimensions(dimensions, fact) {
  const state = {};
  _.forEach(fact.dimensions, (value, key) => {
    state[value] = dimensions[key].name;
  });
  return state;
}

export function fetchDimensions() {
  return (dispatch) => {
    return fetch('./api/dimensions')
    .then(x => x.json())
    .then(json => dispatch(setDimensions(json)))
    .catch(error => console.log(error));
  };
}

export const SET_LEVEL = `${prefix}/SET_LEVEL`;
export function setLevel(dimension, level, state) {
  const copy = Object.assign({}, state);
  copy[dimension] = level;

  return {
    type: SET_LEVEL,
    dimensions: copy,
  };
}

export function selectLevel(dimension, level) {
  return (dispatch, getState) => {
    const state = getState();
    dispatch(setLevel(dimension, level, state.dimensions.stateDimensions));
  };
}
