import * as React from 'react';
import { connect } from 'react-redux';
import { Table, TableBody, TableHeader, TableHeaderColumn, TableRow, TableRowColumn } from 'material-ui/Table';
import injectTapEventPlugin from 'react-tap-event-plugin';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import _ from 'lodash';

import * as actions from '../actions/data.jsx';

class DataTable extends React.Component {
  renderHeader() {
    return (
      <TableRow>
        {
          this.props.columns.map((value) => {
            return (
              <TableHeaderColumn key={ `header-${value}` }>{ value }</TableHeaderColumn>
            );
          })
        }
        {
          this.props.measures.map((measure) => {
            return (
              <TableHeaderColumn key={ `header-${measure.name}` }>{ measure.name }</TableHeaderColumn>
            );
          })
        }
        <TableHeaderColumn>count</TableHeaderColumn>
      </TableRow>
    );
  }

  renderRow(value, key) {
    return (
      <TableRow key={ `row-${key}` }>
        {
          this.props.columns.map((column) => {
            return value[column] ? (
              <TableRowColumn key={ `row-${key}-${column}` }>{ value[column] }</TableRowColumn>
            ) : (
              <TableRowColumn key={ `row-${key}-${column}` }></TableRowColumn>
            );
          })
        }
        {
          this.props.measures.map((measure) => {
            return value[measure.name] ? (
              <TableRowColumn key={ `row-${key}-${measure.name}` }>{ value[measure.name] }</TableRowColumn>
            ) : (
              <TableRowColumn key={ `row-${key}-${measure.name}` }></TableRowColumn>
            );
          })
        }
        <TableRowColumn>{ value.count }</TableRowColumn>
      </TableRow>
    );
  }

  render() {
    return (
      <MuiThemeProvider muiTheme={ getMuiTheme() }>
        <Table>
          <TableHeader>
            { this.renderHeader() }
          </TableHeader>
          <TableBody>
            {
              this.props.tuples.map((value, key) => {
                if (key < 20) {
                  return this.renderRow(value, key);
                }
                return null;
              })
            }
          </TableBody>
        </Table>
      </MuiThemeProvider>
    );
  }
}

DataTable.propTypes = {};

function mapStateToProps(state) {
  const { dimensions, data } = state;
  return {
    columns: dimensions.columns,
    measures: dimensions.fact.measures,
    tuples: data.tuples,
  };
}

function mapDispatchToProps(dispatch) {
  return {
  };
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(DataTable);
