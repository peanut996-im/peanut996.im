import { ActionTree } from 'vuex';
import Vue from 'vue';
import { RootState } from '@/store';
import { CLEAR_USER } from './mutation-types';
import { AppState } from './state';

const actions: ActionTree<AppState, RootState> = {
  goToSSO({ state, commit }) {
    commit(CLEAR_USER);
    Vue.prototype.$message.info('前往登录页面', 0.3).then(() => {
      window.location.href = state.loginUrl;
    });
  },
};

export default actions;
