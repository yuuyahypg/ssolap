import * as React from 'react';
import { connect } from 'react-redux';
import _ from 'lodash';
import SelectField from 'material-ui/SelectField';
import MenuItem from 'material-ui/MenuItem';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import injectTapEventPlugin from 'react-tap-event-plugin';

class RequestForm extends React.Component {
  getLevels(dimension, array) {
    _.forEach(dimension.rollUp, (branches) => {
      _.forEach(branches, (level) => {
        array.push(level);
      });
    });
  }

  handleOnChange(value, e) {
    this.props.onSelectMenu(value, e.target.textContent);
  }

  renderField(dimension, key) {
    const array = [];
    this.getLevels(dimension, array);
    return (
      <div key={ `sf-${dimension.name}-${key}` }>
        <SelectField
          floatingLabelText={ dimension.name }
          value={ this.props.stateDimensions[dimension.name] }
          onChange={ this.handleOnChange.bind(this, dimension.name) }>
          {
            array.map((level, levelKey) => {
              return (
                <MenuItem
                  key={ `menu-${level}-${levelKey}` }
                  value={ level }
                  primaryText={ level }/>
              );
            })
          }
          <MenuItem value={ "none" } primaryText="none" />
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
            this.props.rootDimensions ? this.props.rootDimensions.map((dimension, key) => this.renderField(dimension, key)) : null
          }
        </div>
      </MuiThemeProvider>
    );
  }
}

export default RequestForm;
