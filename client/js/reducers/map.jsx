import * as ActionTypes from '../actions/map.jsx';

const initialState = {
  geometries: {},
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
    default:
      return state;
  }
};

export default map;
