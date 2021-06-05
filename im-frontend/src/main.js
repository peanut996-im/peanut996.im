const express = require('express');
const http = require('http');
const path = require('path');
const morgan = require('morgan');
const app = express();
const server = http.createServer(app);

// template view engine settings
app.set('views', path.join(__dirname, 'views'));
app.engine('.html', require('ejs').__express);
app.set('view engine', 'html');

app.use(morgan(`short`));

let codeEnv = '';
switch (process.env.NODE_ENV) {
    case 'production':
        codeEnv = 'prod';
        break;
    case 'development':
        codeEnv = 'dev';
        break;
    default:
        codeEnv = '';
}

const dotenv = require('dotenv');
dotenv.config({
    path: `${__dirname}/../.env.${codeEnv}`
});

app.get('/', (req, res) => {
    return res.redirect("/login");
});

// set index.html for vue dist
app.use(express.static(path.join(__dirname, '/../im-frontend-index/dist')));
app.use('/static', express.static(path.join(__dirname, 'assets')));

//set for vue index.html
app.get('/index', (req, res) => {
    res.sendFile('index.html', {'root': __dirname + '/../im-frontend-index/dist'});
});

// set for origin folder
app.get('/index-dev', (req, res) => {
    res.sendFile('/assets/index-native.html', {'root': __dirname});
});


app.get('/login', (req, res) => {
    res.render('login', {
        tokenKey: process.env.TOKEN_KEY,
        uidKey: process.env.UID_KEY,
        ssoUrl: process.env.SSO_URL,
        indexUrl: process.env.INDEX_URL,
        nodeEnv: process.env.NODE_ENV,
        ssoLogin: process.env.SSO_LOGIN,
    });
});

server.listen(process.env.APP_PORT, () => {
    console.log(`express listening on :${process.env.APP_PORT}`);
});

