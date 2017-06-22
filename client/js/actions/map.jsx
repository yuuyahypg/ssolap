import _ from 'lodash';

const prefix = 'map';
const request = require('superagent');

export const SET_GEO = `${prefix}/SET_GEO`;
export function setGeo(param, geojson) {
  return {
    type: SET_GEO,
    geometries: geojson,
    bounds: param,
  };
}

export function fetchSetGeo(param) {
  return (dispatch) => {
    return request
      .get('./api/geometry')
      .query(param)
      .end((err, res) => {
        dispatch(setGeo(param, res.body.geojson));
      });
  };
}

export const UPDATE_GEOMETRIES = `${prefix}/UPDATE_GEOMETRIES`;
export function updateGeometries(geojson, geometries) {
  return {
    type: UPDATE_GEOMETRIES,
    geometries: geojson,
    margedGeometries: geometries,
  };
}

export function fetchUpdateGeo(param) {
  return (dispatch, getState) => {
    return request
      .get('./api/geometry')
      .query(param)
      .end((err, res) => {
        const state = getState();
        const geometries = getMargedGeometries(res.body.geojson.features, state.map.tuples, state.map.max);
        dispatch(updateGeometries(res.body.geojson, geometries));
      });
  };
}

export const GET_TUPLES = `${prefix}/GET_TUPLES`;
export function getTuples(tuples, max, geometries) {
  return {
    type: GET_TUPLES,
    tuples,
    max,
    geometries,
  };
}

export function fetchGetTuples() {
  return (dispatch, getState) => {
    const state = getState();

    if (!isSelectValue(state.data.selectedValue)) {
      return null;
    }

    const info = getMatchTuplesInfo(state.data.tuples, state.data.selectedValue);
    const geometries = getMargedGeometries(state.map.geometries.features, info.tuples, info.max);

    return dispatch(getTuples(info.tuples, info.max, geometries));
  };
}

function isSelectValue(values) {
  let flag = true;
  _.forEach(values, (value) => {
    if (value === "") {
      flag = false;
    }
  });

  return flag;
}

function getMatchTuplesInfo(tuples, values) {
  const newTuples = [];
  let max = 0;

  _.forEach(tuples, (tuple) => {
    let f = true;
    _.forIn(values, (v, k) => {
      if (v !== tuple[k]) {
        f = false;
      }
    });

    if (f) {
      newTuples.push(tuple);
      if (tuple.count > max) {
        max = tuple.count;
      }
    }
  });

  return {
    tuples: newTuples,
    max,
  };
}

function getMargedGeometries(features, tuples, max) {
  const newGeometries = [];

  _.forEach(features, (geoJson) => {
    const newGeoJson = JSON.parse(JSON.stringify(geoJson));
    const tuple = _.find(tuples, newGeoJson.properties);

    if (tuple) {
      newGeoJson.properties.count = tuple.count;
      newGeoJson.properties.propotion = tuple.count / max;
      newGeometries.push(newGeoJson);
    } else {
      newGeoJson.properties.count = 0;
      newGeoJson.properties.propotion = 0;
      newGeometries.push(newGeoJson);
    }
  });

  return newGeometries;
}
