import * as ActionTypes from '../actions/map.jsx';

const initialState = {
  geometries: {},
  margedGeometries: {},
  road: {},
  bounds: {
    northEast: {
      lon: 0,
      lat: 0,
    },
    southWest: {
      lon: 0,
      lat: 0,
    },
  },
  tuples: [],
  coordinates: [],
  max: 0,
};

const map = (state = initialState, action) => {
  switch (action.type) {
    case ActionTypes.SET_GEO:
      return {
        ...state,
        geometries: action.geometries,
        bounds: {
          northEast: {
            lon: action.bounds.northEastLon,
            lat: action.bounds.northEastLat,
          },
          southWest: {
            lon: action.bounds.southWestLon,
            lat: action.bounds.southWestLat,
          },
        },
      };
    case ActionTypes.SET_ROAD:
      return {
        ...state,
        road: action.road,
      };
    case ActionTypes.UPDATE_GEOMETRIES:
      return {
        ...state,
        geometries: action.geometries,
        margedGeometries: {
          type: "FeatureCollection",
          features:action.margedGeometries,
        },
      };
    case ActionTypes.GET_TUPLES:
      return {
        ...state,
        tuples: action.tuples,
        max: action.max,
        margedGeometries: {
          type: "FeatureCollection",
          features:action.geometries,
        },
        coordinates: action.coordinates,
      };
    default:
      return state;
  }
};

export default map;
