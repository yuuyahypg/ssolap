import * as React from 'react';
import { connect } from 'react-redux';
import RaisedButton from 'material-ui/RaisedButton';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import injectTapEventPlugin from 'react-tap-event-plugin';

import _ from 'lodash';

import * as actions from '../actions/map.jsx';

const chroma = require('chroma-js');

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
        region: this.props.region,
      });
    });
  }

  updateGeometry() {
    const bounds = this.refs.map.leafletElement.getBounds();
    this.props.updateGeo({
      southWestLon: bounds._southWest.lng,
      southWestLat: bounds._southWest.lat,
      northEastLon: bounds._northEast.lng,
      northEastLat: bounds._northEast.lat,
      region: this.props.region,
    });
  }

  renderGeometry() {
    return (
      <GeoJSON
        key={ `${_.now()}${this.props.geometries.features[0].geometry.coordinates[0][0][0]}` }
        data={ this.props.geometries.features }
        style={ this.style }
        onEachFeature={ this.onEachFeature }/>
    );
  }

  style(feature) {
    const separater = [0, 0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 0.07, 0.08, 0.09, 0.1, 0.15, 0.2, 0.25, 0.3, 0.35, 0.4, 0.45, 0.5, 0.6, 0.7, 0.8, 1];
    const getColor = chroma.scale(["blue", "aqua", "lime", "green", "yellow", "orange", "orangered", "red"]).classes([0, 0.01, 0.05, 0.1, 0.15, 0.2, 0.3, 0.4, 0.5, 0.7, 1]);
    const color = getColor(feature.properties.propotion).toString();

    return {
      color,
      opacity: 0.5,
      fillOpacity: 0.5,
    };
  }

  onEachFeature(feature, layer) {
    const popupContent = `<p>${feature.properties.prefecture ? feature.properties.prefecture : ""}
                          ${feature.properties.city ? feature.properties.city : ""}
                          ${feature.properties.region1 ? feature.properties.region1 : ""}
                          ${feature.properties.region2 ? feature.properties.region2 : ""}</p>
                          <p>count:${feature.properties.count}</p>`;
    layer.bindPopup(popupContent);
  }

  filter(feature) {
    if (feature.properties.count > 0) {
      return true;
    }

    return false;
  }

  render() {
    return (
      <div className={ this.props.className }>
        {
          (Map) && (TileLayer)
          ? (
            <Map
              center={ [35.68, 139.7] }
              zoom={ 13 }
              onMoveend={ this.updateGeometry.bind(this) }
              ref="map"
              className="col-xs-10">
              <TileLayer
                url='http://cyberjapandata.gsi.go.jp/xyz/blank/{z}/{x}/{y}.png'
                attribution='&copy; <a href="http://maps.gsi.go.jp/development/ichiran.html">地理院タイル</a> contributors'/>
              { _.size(this.props.geometries.features) > 0 ? this.renderGeometry() : null }
            </Map>
          )
          : (null)
        }
      </div>
    );
  }
}

function mapStateToProps(state) {
  const { map, dimensions } = state;
  return {
    geometries: map.margedGeometries,
    bounds: map.bounds,
    region: dimensions.stateDimensions.region,
  };
}

function mapDispatchToProps(dispatch) {
  return {
    getGeometry: param => dispatch(actions.fetchSetGeo(param)),
    updateGeo: param => dispatch(actions.fetchUpdateGeo(param)),
  };
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(DataMap);
