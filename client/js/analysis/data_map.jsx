import * as React from 'react';
import { connect } from 'react-redux';
import injectTapEventPlugin from 'react-tap-event-plugin';

import * as actions from '../actions/map.jsx';

let Map;
let TileLayer;
let GeoJSON;

class DataMap extends React.Component {
  componentDidMount() {
    Map = require('react-leaflet').Map;
    TileLayer = require('react-leaflet').TileLayer;
    GeoJSON = require('react-leaflet').GeoJSON;
    this.forceUpdate(() => {
      const bounds = this.refs.map.leafletElement.getBounds();
      this.props.getGeometry({
        southWestLon: bounds._southWest.lng,
        southWestLat: bounds._southWest.lat,
        northEastLon: bounds._northEast.lng,
        northEastLat: bounds._northEast.lat,
      });
    });
  }

  test() {
    console.log(this.props);
  }

  renderGeometry() {
    console.log(this.props.geometries);
    return (
      <GeoJSON data={ this.props.geometries.features } />
    );
  }

  render() {
    return (
      <div>
        {
          (Map) && (TileLayer)
          ? (
            <Map
              center={ [35.68, 139.7] }
              zoom={ 13 }
              onMoveend={ this.test.bind(this) }
              ref="map">
              <TileLayer
                url='https://cyberjapandata.gsi.go.jp/xyz/std/{z}/{x}/{y}.png'
                attribution='&copy; <a href="http://maps.gsi.go.jp/development/ichiran.html">地理院タイル</a> contributors'/>
              { this.props.geometries.features ? this.renderGeometry() : null }
            </Map>
          )
          : (null)
        }
      </div>
    );
  }
}

function mapStateToProps(state) {
  const { map } = state;
  return {
    geometries: map.geometries,
    bounds: map.bounds,
  };
}

function mapDispatchToProps(dispatch) {
  return {
    getGeometry: param => dispatch(actions.fetchSetGeo(param)),
  };
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(DataMap);
