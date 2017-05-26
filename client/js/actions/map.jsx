import _ from 'lodash';

const prefix = 'map';
const request = require('superagent');

export const SET_GEO = `${prefix}/SET_GEO`;
export function setGeo(param, geojson) {
  return {
    type: SET_GEO,
    geometries: geojson.geojson,
    bounds: param,
  };
}

export function fetchSetGeo(param) {
  return (dispatch) => {
    return request
      .get('./api/geometry')
      .query(param)
      .end((err, res) => {
        dispatch(setGeo(param, res.body));
      });
  };
}
