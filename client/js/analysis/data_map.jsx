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
let Circle;

class DataMap extends React.Component {
  componentDidMount() {
    Map = require('react-leaflet').Map;
    TileLayer = require('react-leaflet').TileLayer;
    GeoJSON = require('react-leaflet').GeoJSON;
    Circle = require('react-leaflet').Circle;
    this.forceUpdate(() => {
      const bounds = this.refs.map.leafletElement.getBounds();
      this.props.getGeometry({
        southWestLon: bounds._southWest.lng,
        southWestLat: bounds._southWest.lat,
        northEastLon: bounds._northEast.lng,
        northEastLat: bounds._northEast.lat,
        region: this.props.region,
      });
      this.props.getRoad();
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
        onEachFeature={ this.onEachFeature }
        filter={ this.filter }/>
    );
  }

  renderRoad() {
    return (
      <GeoJSON
        key={ `${_.now()}${this.props.road.features[0].geometry.coordinates[0][0][0]}` }
        data={ this.props.road.features }
        style={ this.roadStyle } />
    );
  }

  renderPoints() {
    const features = {
      "type": "Feature",
      "geometry": {
        "type": "MultiPoint",
        "coordinates": this.props.coordinates,
      },
    };

    return (
      <GeoJSON
        data={ features }/>
    );
  }

  style(feature) {
    const separater = [0, 0.05, 0.1, 0.15, 0.2, 0.25, 0.3, 0.35, 0.4, 0.45, 0.5, 0.55, 0.6, 0.65, 0.7, 0.75, 0.8, 0.85, 0.9, 0.95, 1];
    const getColor = chroma.scale(["blue", "aqua", "lime", "green", "yellow", "orange", "orangered", "red"]).classes(separater);
    const color = getColor(feature.properties.propotion).toString();

    return {
      color,
      opacity: 0.5,
      fillOpacity: 0.5,
    };
  }

  roadStyle(feature) {
    return {
      color: "black",
      opacity: 1,
      fillOpacity: 1,
    };
  }

  pointStyle(feature, latlng) {
    return  <Circle center={ latlng } radius={ 2 }/>;
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
              center={ [36.084219, 140.114217] }
              zoom={ 13 }
              onMoveend={ this.updateGeometry.bind(this) }
              ref="map"
              className="col-xs-10">
              <TileLayer
                url='http://cyberjapandata.gsi.go.jp/xyz/blank/{z}/{x}/{y}.png'
                attribution='&copy; <a href="http://maps.gsi.go.jp/development/ichiran.html">地理院タイル</a> contributors'/>
              { _.size(this.props.geometries.features) > 0 ? this.renderGeometry() : null }
              { _.size(this.props.road.features) > 0 ? this.renderRoad() : null }
              { _.size(this.props.coordinates) > 0 ? this.renderPoints() : null }
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
    road: map.road,
    bounds: map.bounds,
    region: dimensions.stateDimensions.region,
    tuples: map.tuples,
    coordinates: map.coordinates,
  };
}

function mapDispatchToProps(dispatch) {
  return {
    getGeometry: param => dispatch(actions.fetchSetGeo(param)),
    getRoad: param => dispatch(actions.fetchSetRoad()),
    updateGeo: param => dispatch(actions.fetchUpdateGeo(param)),
  };
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(DataMap);
