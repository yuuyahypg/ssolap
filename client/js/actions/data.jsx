import _ from 'lodash';

const prefix = 'data';
const request = require('superagent');

export const SET_TUPLES = `${prefix}/SET_TUPLES`;
export function setTuples(json) {
  console.log(json);
  return {
    type: SET_TUPLES,
    // tuples: json,
  };
}

export function fetchRequest() {
  return (dispatch, getState) => {
    return request
      .get('./api/request')
      .query({
        gender: "gender",
        age: "age",
        work: "none",
        transport: "none",
        purpose: "purpose",
        region: "region2",
        time: "minute",
      })
      .end((err, res) => {
        dispatch(setTuples(res.body));
      });
  };
}
