import cookie from 'js-cookie';
import { MutationTree } from 'vuex';
import {
  SET_USER,
  CLEAR_USER,
  SET_TOKEN,
  SET_MOBILE,
  SET_BACKGROUND,
  SET_ACTIVETABNAME,
  SET_LOADING,
  SET_UID,
  SET_UID_AND_TOKEN,
} from './mutation-types';
import { AppState } from './state';

const mutations: MutationTree<AppState> = {
  [SET_USER](state, payload: User) {
    state.user = Object.assign(state.user, payload);
  },

  [CLEAR_USER](state) {
    state.user = {
      uid: '',
      userId: '',
      username: '',
      password: '',
      account: '',
      avatar: '',
      createTime: 0,
    };
    cookie.remove(state.uidKey);
    cookie.remove(state.tokenKey);
    window.sessionStorage.removeItem(state.uidKey);
    window.sessionStorage.removeItem(state.uidKey);
  },

  [SET_TOKEN](state, payload: string) {
    state.token = payload;
    // cookie.set('token', payload, { expires: 3 });
    cookie.set(state.tokenKey, payload, { expires: 3 });
  },

  [SET_MOBILE](state, payload: boolean) {
    state.mobile = payload;
  },

  [SET_BACKGROUND](state, payload: string) {
    state.background = payload;
    localStorage.setItem('background', payload);
  },

  [SET_ACTIVETABNAME](state, payload: 'message' | 'contacts') {
    state.activeTabName = payload;
  },

  [SET_LOADING](state, payload: boolean) {
    state.loading = payload;
  },

  [SET_UID](state, payload: string) {
    state.uid = payload;
    cookie.set(state.uidKey, payload, { expires: 3 });
  },
  [SET_UID_AND_TOKEN](state) {
    const uid = cookie.get(state.uidKey);
    const token = cookie.get(state.tokenKey);
    state.uid = uid!;
    state.token = token!;
  },
};

export default mutations;
