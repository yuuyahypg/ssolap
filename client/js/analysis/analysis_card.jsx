import * as React from 'react';
import { connect } from 'react-redux';
import { Tabs, Tab } from 'material-ui/Tabs';
import injectTapEventPlugin from 'react-tap-event-plugin';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';

import DataTable from './data_table.jsx';
import AnalysisMapCard from './analysis_map_card.jsx';
import * as actions from '../actions/data.jsx';

class AnalysisCard extends React.Component {
  componentWillMount() {
    this.props.getData();
  }

  handleChange(value) {
    this.props.switchTab(value);
  }

  render() {
    return (
      <MuiThemeProvider muiTheme={ getMuiTheme() }>
        <Tabs
          value={ this.props.tabState }
          onChange={ this.handleChange.bind(this) }>
          <Tab label="Table" value="table">
            <DataTable/>
          </Tab>
          <Tab label="Map" value="map">
            <AnalysisMapCard/>
          </Tab>
        </Tabs>
      </MuiThemeProvider>
    );
  }
}

AnalysisCard.propTypes = {};

function mapStateToProps(state) {
  const { data } = state;
  return {
    tuples: data.tuples,
    tabState: data.tabState,
  };
}

function mapDispatchToProps(dispatch) {
  return {
    getData: () => dispatch(actions.fetchRequest()),
    switchTab: state => dispatch(actions.fetchSwitchTab(state)),
  };
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(AnalysisCard);
