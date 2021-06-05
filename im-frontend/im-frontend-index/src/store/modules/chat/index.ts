import { Module } from 'vuex';
import { RootState } from '@/store';
import actions from './actions';
import mutations from './mutations';
import getters from './getters';
import state, { ChatState } from './state';

const chat: Module<ChatState, RootState> = {
  namespaced: true,
  state,
  mutations,
  actions,
  getters,
};

export default chat;
