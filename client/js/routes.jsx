import React from 'react';
import { Route, IndexRoute } from 'react-router';
import App from './application.jsx';
import RequestFormCard from './olap_request/request_form_card.jsx';
import AnalysisCard from './analysis/analysis_card.jsx';
import NotFound from './not_found.jsx';

const routes = (
  <Route path="/" component={ App }>
    <IndexRoute component={ RequestFormCard } />
    <Route path="analysis" component={ AnalysisCard } />
    <Route path="*" component={ NotFound }/>
  </Route>
);

export default routes;
