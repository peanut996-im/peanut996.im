import Vue from 'vue';
import { MutationTree } from 'vuex';
import {
  SET_SOCKET,
  SET_DROPPED,
  DEL_GROUP,
  DEL_FRIEND,
  ADD_UNREAD_GATHER,
  LOSE_UNREAD_GATHER,
  SET_MY_FRIEND_GATHER,
  SET_MY_GROUP_GATHER,
  SET_ROOM_GATHER,
  SET_ACTIVE_ROOM,
  SAVE_MESSAGE,
  SET_USER_GATHER,
  SET_USER_SEARCH_RESULT,
  SET_GROUP_SEARCH_RESULT,
  SET_PULL_MESSAGE_RESULT,
  SET_MESSAGES,
  CLEAR_FRIEND_GATHER,
  CLEAR_GROUP_GATHER,
} from './mutation-types';
import { ChatState } from './state';

const mutations: MutationTree<ChatState> = {
  // 保存socket
  [SET_SOCKET](state, payload: SocketIOClient.Socket) {
    state.socket = payload;
  },

  // 设置用户是否处于掉线重连状态
  [SET_DROPPED](state, payload: boolean) {
    state.dropped = payload;
  },

  [SAVE_MESSAGE](state, payload: ChatMessage) {
    // get friend/group from room map
    if (state.roomGather[payload.to]) {
      if (state.roomGather[payload.to].messages) {
        state.roomGather[payload.to].messages!.push(payload);
      } else {
        Vue.set(state.roomGather[payload.to], 'messages', [payload]);
      }
    }
  },

  [SET_PULL_MESSAGE_RESULT](state, payload: Array<ChatMessage>) {
    state.pullMessageResult = payload;
  },

  [SET_MESSAGES](state, payload: ChatMessage[]) {
    if (payload && payload.length) {
      if (state.activeRoom?.uid) {
        Vue.set(state.friendGather[state.roomGather[payload[0].to].uid], 'messages', payload);
      }
      if (state.activeRoom?.groupID) {
        Vue.set(state.groupGather[payload[0].to], 'messages', payload);
      }
    }
  },

  [SET_ACTIVE_ROOM](state, payload: Friend & Group) {
    state.activeRoom = payload;
  },

  [SET_MY_GROUP_GATHER](state, payload: Group) {
    Vue.set(state.groupGather, payload.groupID, payload);
  },

  [SET_USER_GATHER](state, payload: User) {
    Vue.set(state.userGather, payload.uid, payload);
  },

  [CLEAR_FRIEND_GATHER](state) {
    state.friendGather = {};
  },
  [CLEAR_GROUP_GATHER](state) {
    state.groupGather = {};
  },
  [SET_MY_FRIEND_GATHER](state, payload: Friend) {
    Vue.set(state.friendGather, payload.uid, payload);
  },

  [SET_ROOM_GATHER](state, payload: Friend & Group) {
    if (payload.roomID) {
      Vue.set(state.roomGather, payload.roomID, payload);
    }
    if (payload.groupID) {
      Vue.set(state.roomGather, payload.groupID, payload);
    }
  },

  // 退群
  [DEL_GROUP](state, payload: GroupMap) {
    Vue.delete(state.groupGather, payload.groupID);
  },

  // 删好友
  [DEL_FRIEND](state, payload: FriendMap) {
    Vue.delete(state.friendGather, payload.friendB);
  },

  // 给某个聊天组添加未读消息
  [ADD_UNREAD_GATHER](state, payload: string) {
    document.title = '🔴Peanut996.IM';
    if (!state.unReadGather[payload]) {
      Vue.set(state.unReadGather, payload, 1);
    } else {
      ++state.unReadGather[payload];
    }
  },

  // 给某个聊天组清空未读消息
  [LOSE_UNREAD_GATHER](state, payload: string) {
    document.title = 'Peanut996.IM';
    Vue.set(state.unReadGather, payload, 0);
  },

  [SET_USER_SEARCH_RESULT](state, payload: Array<User>) {
    state.userSearchResult = payload;
  },

  [SET_GROUP_SEARCH_RESULT](state, payload: Array<Group>) {
    state.groupSearchResult = payload;
  },
};

export default mutations;
