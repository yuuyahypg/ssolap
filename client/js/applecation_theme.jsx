import React from 'react';
import { tealA400 } from 'material-ui/styles/colors';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import getMuiTheme from 'material-ui/styles/getMuiTheme';

const ApplecationTheme = getMuiTheme({
  userAgent: (typeof navigator !== 'undefined' && navigator.userAgent) || 'all',
  appBar: {
    "color": tealA400,
  },
});

export default ApplecationTheme;
