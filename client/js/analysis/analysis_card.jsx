import * as React from 'react';
import { connect } from 'react-redux';

import * as actions from '../actions/data.jsx';

class AnalysisCard extends React.Component {
  componentWillMount() {
    this.props.getData();
  }

  render() {
    console.log(this.props);
    return (
      <div/>
    );
  }
}

AnalysisCard.propTypes = {};

function mapStateToProps(state) {
  return {};
}

function mapDispatchToProps(dispatch) {
  return {
    getData: () => dispatch(actions.fetchRequest()),
  };
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(AnalysisCard);
