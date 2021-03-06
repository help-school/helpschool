import './mock-db';
import React from 'react';
import ReactDOM from 'react-dom';
import {Router } from 'react-router-dom';
import history from './history';
import Auth from "./auth/auth";

// Importing the Bootstrap CSS
import "bootstrap/dist/css/bootstrap.min.css";

import "slick-carousel/slick/slick.css";
import "slick-carousel/slick/slick-theme.css";

import './index.css';
import App from './App';
import * as serviceWorker from './serviceWorker';

import { createStore, applyMiddleware } from 'redux';
import { Provider } from 'react-redux';
import logger from 'redux-logger';
import reducers from './store/reducers';
import ReduxThunk from 'redux-thunk';

const store = createStore(reducers, {}, applyMiddleware(logger, ReduxThunk));

ReactDOM.render(
    <Auth>
        <Provider store={store}>
            <Router history={history}>
                <App />
            </Router>
        </Provider>
    </Auth>,
    document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
// serviceWorker.unregister();
