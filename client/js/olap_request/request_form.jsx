import * as React from 'react';
import { connect } from 'react-redux';
import _ from 'lodash';
import SelectField from 'material-ui/SelectField';
import MenuItem from 'material-ui/MenuItem';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import injectTapEventPlugin from 'react-tap-event-plugin';

class RequestForm extends React.Component {
  getReferences(dimension, array) {
    array.push(dimension.name);
    _.forEach(dimension.references, (value) => {
      this.getReferences(value, array);
    });
  }

  handleOnChange(value, e) {
    this.props.onSelectMenu(value, e.target.textContent);
  }

  renderField(value, key) {
    const array = [];
    this.getReferences(this.props.rootDimensions[key], array);
    return (
      <div key={ `sf-${value}-${key}` }>
        <SelectField
          floatingLabelText={ value }
          value={ this.props.stateDimensions[value] }
          onChange={ this.handleOnChange.bind(this, value) }>
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
            this.props.fact.dimensions ? this.props.fact.dimensions.map((value, key) => this.renderField(value, key)) : null
          }
        </div>
      </MuiThemeProvider>
    );
  }
}

export default RequestForm;
