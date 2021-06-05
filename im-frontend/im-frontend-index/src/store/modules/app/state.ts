import cookie from 'js-cookie';

export interface AppState {
  user: User;
  uid: string;
  token: string;
  mobile: boolean;
  background: string;
  activeTabName: 'message' | 'contacts';
  loading: boolean;
  loginUrl: string;
  tokenKey: string;
  uidKey: string;
  env: string;
  appKey: string;
}

const appState: AppState = {
  user: {
    uid: '',
    userId: '',
    username: '',
    account: '',
    password: '',
    avatar: '',
    createTime: 0,
  },
  uid: '',
  token: cookie.get('token') as string,
  mobile: false,
  background: '',
  activeTabName: 'message',
  loginUrl: process.env.VUE_APP_LOGIN_URL,
  loading: false, // 全局Loading状态
  tokenKey: process.env.VUE_APP_TOKEN_KEY,
  uidKey: process.env.VUE_APP_UID_KEY,
  env: process.env.NODE_ENV,
  appKey: process.env.VUE_APP_KEY,
};

export default appState;
