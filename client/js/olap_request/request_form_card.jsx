import * as React from 'react';
import { connect } from 'react-redux';
import RaisedButton from 'material-ui/RaisedButton';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import { Link } from 'react-router';
import injectTapEventPlugin from 'react-tap-event-plugin';

import RequestForm from './request_form.jsx';
import * as actions from '../actions/dimensions.jsx';

class RequestFormCard extends React.Component {
  componentDidMount() {
    this.props.getDimensionsInfo();
    console.log(this.props);
  }

  render() {
    return (
      <div style={ {
        paddingLeft: '5%',
      } }>
        <RequestForm
          rootDimensions={ this.props.rootDimensions }
          fact={ this.props.fact }
          stateDimensions={ this.props.stateDimensions }
          onSelectMenu={ this.props.onSelectMenu }/>
        <MuiThemeProvider muiTheme={ getMuiTheme() }>
          <RaisedButton label="Request" primary={ true } containerElement={ <Link to="/analysis" /> }/>
        </MuiThemeProvider>
      </div>
    );
  }
}

function mapStateToProps(state) {
  const { dimensions } = state;
  return {
    rootDimensions: dimensions.rootDimensions,
    fact: dimensions.fact,
    stateDimensions: dimensions.stateDimensions,
  };
}

function mapDispatchToProps(dispatch) {
  return {
    getDimensionsInfo: () => dispatch(actions.fetchDimensions()),
    onSelectMenu: (dimension, level) => dispatch(actions.selectLevel(dimension, level)),
  };
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(RequestFormCard);
