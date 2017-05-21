import * as React from 'react';
import { connect } from 'react-redux';
import { Route, Switch } from 'react-router';
import Helmet from 'react-helmet';
import { AppBar } from 'material-ui';
import injectTapEventPlugin from 'react-tap-event-plugin';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import routes from './routes.jsx';
import ApplecationTheme from './applecation_theme.jsx';

injectTapEventPlugin();

class App extends React.Component {
  render() {
    return (
      <div>
        <MuiThemeProvider muiTheme={ ApplecationTheme }>
          <AppBar
            title="Title"
            iconClassNameRight="muidocs-icon-navigation-expand-more"/>
        </MuiThemeProvider>
        { this.props.children }
      </div>
    );
  }
}

App.propTypes = {};

function mapStateToProps(state) {
  return {};
}

function mapDispatchToProps(dispatch) {
  return {};
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(App);
