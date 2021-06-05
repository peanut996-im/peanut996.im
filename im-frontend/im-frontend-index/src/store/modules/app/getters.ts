import { GetterTree } from 'vuex';
import cookie from 'js-cookie';
import { RootState } from '@/store';
import { AppState } from './state';

const getters: GetterTree<AppState, RootState> = {
  user(state) {
    return state.user;
  },
  uid(state) {
    const uid = cookie.get(state.uidKey);
    if (!uid) {
      return '';
    }
    state.uid = uid;
    return state.uid;
  },
  mobile(state) {
    return state.mobile;
  },
  background(state) {
    // eslint-disable-next-line no-unused-expressions
    state.background;
    return localStorage.getItem('background');
  },
  activeTabName(state) {
    return state.activeTabName;
  },
  token(state) {
    const token = cookie.get(state.tokenKey);
    if (!token) {
      return '';
    }
    state.token = token;
    return state.token;
  },
  loading(state) {
    return state.loading;
  },
  isLogined(state) {
    return cookie.get(state.tokenKey) !== '' || window.sessionStorage.getItem(state.tokenKey) !== '';
  },
};

export default getters;
