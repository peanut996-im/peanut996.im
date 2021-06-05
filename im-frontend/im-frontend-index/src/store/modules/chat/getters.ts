import { GetterTree } from 'vuex';
import { RootState } from '@/store';
import { ChatState } from './state';

const getters: GetterTree<ChatState, RootState> = {
  socket(state) {
    return state.socket;
  },
  dropped(state) {
    return state.dropped;
  },
  activeRoom(state) {
    return state.activeRoom;
  },
  roomGather(state) {
    return state.roomGather;
  },
  friendGather(state) {
    return state.friendGather;
  },
  groupGather(state) {
    return state.groupGather;
  },
  userGather(state) {
    return state.userGather;
  },
  unReadGather(state) {
    return state.unReadGather;
  },
  userSearchResult(state) {
    return state.userSearchResult;
  },
  groupSearchResult(state) {
    return state.groupSearchResult;
  },
  pullMessageResult(state) {
    return state.pullMessageResult;
  },
};

export default getters;
