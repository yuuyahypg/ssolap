import * as React from 'react';
import { connect } from 'react-redux';
import SelectField from 'material-ui/SelectField';
import MenuItem from 'material-ui/MenuItem';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import injectTapEventPlugin from 'react-tap-event-plugin';
import _ from 'lodash';

import DataTable from './data_table.jsx';
import DataMap from './data_map.jsx';
import * as actions from '../actions/data.jsx';

class DataMapValueSelector extends React.Component {
  handleOnChange(level, e) {
    this.props.onSelectMenu(level, e.target.textContent);
  }

  renderSelectField(level, dimension) {
    if (dimension === "region" || level === "none") {
      return null;
    }

    return (
      <div key={ `sf-${level}-${dimension}` }>
        <SelectField
          floatingLabelText={ level }
          value={ this.props.selectedValue[level] }
          onChange={ this.handleOnChange.bind(this, level) }>
          {
            this.props.valueList[level].map((value, valueKey) => {
              return (
                <MenuItem
                  key={ `menu-${value}-${valueKey}` }
                  value={ value }
                  primaryText={ value }/>
              );
            })
          }
        </SelectField>
        <br/>
      </div>
    );
  }

  render() {
    return (
      <MuiThemeProvider muiTheme={ getMuiTheme() }>
        <div>
          {
            _.size(this.props.stateDimensions) > 0 && _.size(this.props.valueList) > 0 ? _.map(this.props.stateDimensions, (value, key) => { return this.renderSelectField(value, key); }) : null
          }
        </div>
      </MuiThemeProvider>
    );
  }
}

function mapStateToProps(state) {
  const { data, dimensions } = state;
  return {
    stateDimensions: dimensions.stateDimensions,
    valueList: data.valueList,
    selectedValue: data.selectedValue,
  };
}

function mapDispatchToProps(dispatch) {
  return {
    onSelectMenu: (level, value) => dispatch(actions.fetchSelectValue(level, value)),
  };
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(DataMapValueSelector);
