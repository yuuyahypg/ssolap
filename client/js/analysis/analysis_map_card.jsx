import * as React from 'react';
import { connect } from 'react-redux';

import DataMapValueSelector from './data_map_value_selector.jsx';
import * as actions from '../actions/data.jsx';

class AnalysisMapCard extends React.Component {
  render() {
    console.log(this.props);
    return (
      <div>
        <DataMapValueSelector/>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
  };
}

function mapDispatchToProps(dispatch) {
  return {
  };
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(AnalysisMapCard);
