import React from 'react';
import { tealA400 } from 'material-ui/styles/colors';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import getMuiTheme from 'material-ui/styles/getMuiTheme';

const ApplecationTheme = getMuiTheme({
  appBar: {
    "background-color": tealA400,
  },
});

export default ApplecationTheme;
