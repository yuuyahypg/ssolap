import * as React from 'react';
import { AppBar } from 'material-ui';
import injectTapEventPlugin from 'react-tap-event-plugin';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import ApplecationTheme from './applecation_theme.jsx';

export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      message: '',
    };
  }

  componentDidMount() {
    fetch('./api/home')
    .then(x => x.json())
    .then((json) => {
      this.setState({
        message: json.message,
      });
    });
  }

  render() {
    const { message } = this.state;
    return (
      <MuiThemeProvider muiTheme={ ApplecationTheme }>
        <AppBar
          title="Title"
          iconClassNameRight="muidocs-icon-navigation-expand-more"
        />
      </MuiThemeProvider>
    );
  }
}
