import _ from 'lodash';

const prefix = 'dimensions';

export const SET_DIMENSIONS = `${prefix}/SET_DIMENSIONS`;
export function setDimensions(json) {
  const stateDimensions = getStateDimensions(json.dimensions);
  return {
    type: SET_DIMENSIONS,
    dimensions: json.dimensions,
    fact: json.fact,
    isDBConnected: json.isDBConnected,
    stateDimensions,
    columns: getColumns(json.dimensions, stateDimensions),
  };
}

function getStateDimensions(dimensions) {
  const state = {};
  _.forEach(dimensions, (dimension) => {
    state[dimension.name] = dimension.rollUp[0][0];
  });
  return state;
}

function getColumns(dimensions, stateDimensions) {
  const columns = [];
  _.forEach(dimensions, (dimension) => {
    if (stateDimensions[dimension.name] !== "none") {
      const stateLevel = stateDimensions[dimension.name];

      _.forEach(dimension.rollUp, (branch) => {
        let f = false;
        _.forEach(branch, (level) => {
          if (level === stateLevel) {
            f = true;
          }

          if (f) {
            columns.push(level);
          }
        });
      });
    }
  });
  return columns;
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
export function setLevel(stateDimensions, columns) {
  return {
    type: SET_LEVEL,
    stateDimensions,
    columns,
  };
}

export function selectLevel(dimension, level) {
  return (dispatch, getState) => {
    const state = getState();
    const copy = Object.assign({}, state.dimensions.stateDimensions);
    copy[dimension] = level;
    const columns = getColumns(state.dimensions.rootDimensions, copy);

    dispatch(setLevel(copy, columns));
  };
}
