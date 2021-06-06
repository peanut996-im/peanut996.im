import {makeSign, sha1} from "./tool.js";
import * as Cookies from './js-cookie.js';

let tokenKey = '';
let uidKey = '';
let baseUrl = '';
const inSixHours = 0.25;
const isLogined = (tokenKey, uidKey) => {
    return window.Cookies.get(tokenKey) != undefined ||
        window.sessionStorage.getItem(uidKey) != undefined;
}

const saveToken = (token, uid) => {
    window.sessionStorage.setItem(tokenKey, token);
    window.sessionStorage.setItem(uidKey, uid);
    window.Cookies.set(tokenKey, token, {expires: inSixHours});
    window.Cookies.set(uidKey, uid, {expires: inSixHours});
}

const getToken = () => {
    if (sessionStorage.getItem(tokenKey) != undefined) {
        return {token: sessionStorage.getItem(tokenKey), uid: sessionStorage.getItem(uidKey)}
    }
    return {token: window.Cookies.get(tokenKey), uid: window.Cookies.get(uidKey)}
}

const app = new Vue({
    el: "#app",
    data: {
        account: '',
        password: '',
        indexUrl: '/index',
        env: '',
        ssoLogin: '',
    },
    methods: {
        nameVerify: function (name) {
            const nameReg = /^(?!_)(?!.*?_$)[a-zA-Z0-9_\u4e00-\u9fa5]+$/;
            if (name.length === 0) {
                this.$message.error('请输入名字');
                return false;
            }
            if (!nameReg.test(name)) {
                this.$message.error('名字只含有汉字、字母、数字和下划线 不能以下划线开头和结尾');
                return false;
            }
            if (name.length > 16) {
                this.$message.error('名字太长');
                return false;
            }
            return true;
        },

        validateData: function (account, password) {
            if (!this.nameVerify(account) ) {
                return false;
            }
            if (password == undefined || password === '') {
                this.$message.error('请输入密码');
                return false;
            }
            return true;
        },

        login: function () {
            console.log(`account: ${this.account}\npassword: ${this.password}`);
            if(!this.validateData(this.account,this.password)){
                console.debug("data validate failed");
                return;
            }
            let loginOption = {
                account: this.account,
                password: sha1(this.password).toUpperCase(),
            }
            loginOption.sign = makeSign(loginOption);
            console.log(`loginOption:`, loginOption);
            this.$message.loading('login in progress...wait', 0.5);
            axios.post(baseUrl + '/login', loginOption)
                .then(response => {
                    console.log(response.data.data);
                    let {code, message, data} = response.data;
                    if (code === 0) {
                        saveToken(data.token, data.uid);
                        this.$message.success('login succeed! Redirecting', 1.5)
                            .then(() => {
                                if (this.ssoLogin === 'on') {
                                    window.location.href = `${this.indexUrl}?token=${data.token}&uid=${data.uid}`
                                } else {
                                    if (this.env === 'production') {
                                        window.location.href = `${this.indexUrl}`;
                                    }
                                    if (this.env === 'development') {
                                        this.$message.info('开发环境强制使用SSO登录', 1).then(() => {
                                            window.location.href = `${this.indexUrl}?token=${data.token}&uid=${data.uid}`;
                                        });
                                    }
                                }
                            });
                    } else {
                        this.$message.error('login failed!', 2.5);
                    }

                })
                .catch(error => {
                    console.log(error);
                    this.$message.error('login failed!', 2.5);
                });
        }
        ,
        register: function () {
            if(!this.validateData(this.account,this.password)){
                console.debug("data validate failed");
                return;
            }
            let registerOption = {
                account: this.account,
                password: sha1(this.password).toUpperCase(),
            }
            registerOption.sign = makeSign(registerOption);
            console.log('registerOption: ', registerOption);
            this.$message.loading('register in progress...wait', 0.5);
            axios.post(baseUrl + '/register', registerOption)
                .then(response => {
                    let {code, message, data} = response.data;
                    if (code === 0) {
                        saveToken(data.token, data.uid);
                        this.$message.success('register succeed! Redirecting', 1.5)
                            .then(() => {
                                if (this.ssoLogin === 'on') {
                                    window.location.href = `${this.indexUrl}?token=${data.token}&uid=${data.uid}`
                                } else {
                                    if (this.env === 'production') {
                                        window.location.href = `${this.indexUrl}`;
                                    }
                                    if (this.env === 'development') {
                                        this.$message.info('开发环境强制使用SSO登录', 1).then(() => {
                                            window.location.href = `${this.indexUrl}?token=${data.token}&uid=${data.uid}`;
                                        });
                                    }
                                }
                            });
                    } else {
                        this.$message.error('register failed!', 2.5);
                    }

                })
                .catch(error => {
                    console.log(error);
                    this.$message.error('register failed!', 2.5);
                });
        }
    },
    mounted() {
        tokenKey = this.$refs.tokenKey.value;
        uidKey = this.$refs.uidKey.value;
        baseUrl = this.$refs.ssoUrl.value;
        this.env = this.$refs.nodeEnv.value;
        this.indexUrl = this.$refs.indexUrl.value;
        this.ssoLogin = this.$refs.ssoLogin.value;

        console.log(`ssoLogin: ${this.$refs.ssoLogin.value}`);
        if (tokenKey === "" || uidKey === "") {
            this.$message.error('Cannot get dotenv config');
            return
        }
        if (isLogined(tokenKey, uidKey)) {
            console.debug(`after: ${this.ssoLogin}`)
            console.debug(getToken());
            const {token, uid} = getToken();
            if (this.env === 'production') {
                this.$message.success('Have been login, Redirecting', 1.5)
                    .then(() => {
                        if (this.ssoLogin === 'on') {
                            window.location.href = `${this.indexUrl}?token=${token}&uid=${uid}`;
                        } else {
                            console.debug('go to index')
                            window.location.href = `${this.indexUrl}`;
                        }
                    });
            }

        }
    }

});
