/* eslint-disable no-unused-vars */
import { ActionTree } from 'vuex';
import io from 'socket.io-client';
import Vue from 'vue';
import { makeSign, sortByTime } from '@/utils/common';
import { ErrorHttpInnerError, ErrorTokenInvalid, ErrorSignInvalid } from '@/api/constants';
import { RootState } from '@/store';
import { SET_LOADING, CLEAR_USER, SET_USER } from '../app/mutation-types';
import { ChatState } from './state';
import {
  SET_SOCKET,
  SET_DROPPED,
  DEL_GROUP,
  DEL_FRIEND,
  ADD_UNREAD_GATHER,
  SET_MY_FRIEND_GATHER,
  SET_MY_GROUP_GATHER,
  SET_ROOM_GATHER,
  SAVE_MESSAGE,
  SET_ACTIVE_ROOM,
  SET_USER_GATHER,
  SET_USER_SEARCH_RESULT,
  SET_GROUP_SEARCH_RESULT,
  SET_PULL_MESSAGE_RESULT,
  CLEAR_FRIEND_GATHER,
  CLEAR_GROUP_GATHER,
} from './mutation-types';

const actions: ActionTree<ChatState, RootState> = {
  // 初始化socket连接和监听socket事件
  async connectSocket({ commit, state, dispatch, rootState }) {
    const { uid, token } = rootState.app;
    const query = {
      uid,
      token,
      sign: '',
    };
    query.sign = makeSign(query);
    const socket: SocketIOClient.Socket = io.connect(`${process.env.VUE_APP_GATE_URL}`, {
      reconnection: true,
      reconnectionDelay: process.env.NODE_ENV === 'production' ? 2000 : 10000,
      query,
    });

    socket.on('connect', async () => {
      console.log(`socket connect success: socket.id: ${socket.id}`);
      commit(SET_SOCKET, socket);
    });

    socket.on('disconnect', async () => {
      console.log(`socket disconnect: socket.id: ${socket.id}`);
    });

    socket.on('auth', (res: any) => {
      console.debug(`auth response: `);
      console.debug(res);
      switch (res.code) {
        case ErrorHttpInnerError:
          Vue.prototype.$message.error('鉴权服务无响应!', 1.5);
          break;
        case ErrorSignInvalid:
          Vue.prototype.$message.error('签名错误!', 1.5).then(() => {
            window.location.reload();
          });
          commit(`app/${CLEAR_USER}`, false, { root: true });
          break;
        case ErrorTokenInvalid:
          Vue.prototype.$message.error('鉴权失败,请重新登录!', 1.5).then(() => {
            console.debug(`go tot loginUrl: ${rootState.app.loginUrl}`);
            window.location.href = rootState.app.loginUrl;
          });
          commit(`app/${CLEAR_USER}`, false, { root: true });
          break;
        default:
          console.debug('auth done!');
      }
    });

    socket.on('load', (res: Response | any) => {
      console.debug('get load data:');
      console.debug(res);
      dispatch('handleLoadData', res).then(() => {
        console.debug('load data complete.');
      });
      commit(SET_DROPPED, false);
    });

    socket.on('chat', (res: Response | any) => {
      if (res.code) {
        return Vue.prototype.$message.error('消息发送失败', 1.5);
      }
      if (res.content) {
        commit(SAVE_MESSAGE, res);
        const { activeRoom } = state;
        if (activeRoom) {
          // 未读消息计数
          if (activeRoom.groupID && activeRoom.groupID !== res.to) {
            commit(ADD_UNREAD_GATHER, res.to);
          } else if (activeRoom.roomID && activeRoom.roomID !== res.to) {
            commit(ADD_UNREAD_GATHER, res.to);
          }
        }
      }
    });

    socket.on('findUser', (res: Response) => {
      if (res.code) {
        console.error(res.message);
        return Vue.prototype.$message.error('模糊查找失败', 0.5);
      }
      const { data } = res;
      commit(SET_USER_SEARCH_RESULT, data);
      commit(`app/${SET_LOADING}`, false, { root: true });
    });
    socket.on('findGroup', (res: Response) => {
      if (res.code) {
        console.error(res.message);
        return Vue.prototype.$message.error('模糊查找失败', 0.5);
      }
      const { data } = res;
      commit(SET_GROUP_SEARCH_RESULT, data);
      commit(`app/${SET_LOADING}`, false, { root: true });
    });

    socket.on('createGroup', (res: Response) => {
      if (res.code) {
        console.error(res.message);
        return Vue.prototype.$message.error('创建群组失败', 0.5);
      }
      const { data } = res;
      const sortMessageGroup = data;
      sortMessageGroup.message = sortByTime(sortMessageGroup.messages, false);
      // 更新本地群组信息
      commit(SET_MY_GROUP_GATHER, sortMessageGroup);
      // 更新本地房间信息
      commit(SET_ROOM_GATHER, sortMessageGroup);
      commit(`app/${SET_LOADING}`, false, { root: true });
    });

    socket.on('addFriend', (res: Response) => {
      if (res.code) {
        console.error(res.message);
        return Vue.prototype.$message.error('添加好友失败', 0.5);
      }
      const { data, message } = res;
      data.message = sortByTime(data.messages, false);
      commit(SET_MY_FRIEND_GATHER, data);
      commit(SET_ROOM_GATHER, data);
      commit(`app/${SET_LOADING}`, false, { root: true });
      Vue.prototype.$message.success(message);
    });

    socket.on('deleteFriend', (res: Response) => {
      if (res.code) {
        console.error(res.message);
        return Vue.prototype.$message.error('删除好友失败', 0.5);
      }
      const { data, message } = res;
      commit(DEL_FRIEND, data);
      // 返回数组首项
      commit(SET_ACTIVE_ROOM, state.friendGather[Object.keys(state.friendGather)[0]].uid);
      Vue.prototype.$message.success(message);
    });

    socket.on('joinGroup', (res: Response) => {
      if (res.code) {
        console.error(res.message);
        return Vue.prototype.$message.error('添加群组失败', 0.5);
      }
      const { data, message } = res;
      data.message = sortByTime(data.messages, false);
      commit(SET_MY_GROUP_GATHER, data);
      commit(SET_ROOM_GATHER, data);
      commit(`app/${SET_LOADING}`, false, { root: true });
      Vue.prototype.$message.success(message);
    });

    socket.on('leaveGroup', (res: Response) => {
      if (res.code) {
        console.error(res.message);
        return Vue.prototype.$message.error('删除群组失败', 0.5);
      }
      const { data, message } = res;
      commit(DEL_GROUP, data);
      // 返回数组首项
      commit(SET_ACTIVE_ROOM, state.groupGather[Object.keys(state.groupGather)[0]].groupID);
      Vue.prototype.$message.success(message);
    });

    socket.on('pullMessage', (res: Response) => {
      if (res.code) {
        console.error(res.message);
        return Vue.prototype.$message.error('拉取信息失败', 0.5);
      }
      const { data } = res;
      commit(SET_PULL_MESSAGE_RESULT, data);
    });

    socket.on('updateUser', (res: Response) => {
      if (res.code) {
        console.error(res.message);
        return Vue.prototype.$message.error('更新信息失败', 0.5);
      }
      const { data, message } = res;
      commit(`app/${SET_USER}`, data, { root: true });
      Vue.prototype.$message.success(message);
    });
  },

  async handleLoadData({ commit, dispatch, state, rootState }, payload) {
    // clear
    commit(CLEAR_FRIEND_GATHER);
    commit(CLEAR_GROUP_GATHER);
    const { user, friends, groups } = payload;
    console.debug('process on handleLoadData');
    // TODO Save to State
    console.debug('start set user');
    commit(`app/${SET_USER}`, user, { root: true });
    commit(SET_USER_GATHER, user);
    console.debug('start set my_friend_gather');
    friends.forEach((friend) => {
      const sortMessageFriend = friend;
      // 消息排序
      sortMessageFriend.messages = sortByTime(friend.messages, false);
      commit(SET_MY_FRIEND_GATHER, sortMessageFriend);
      commit(SET_ROOM_GATHER, sortMessageFriend);
      commit(SET_USER_GATHER, sortMessageFriend);
    });
    console.debug('start set my_group_gather');
    // 触发watch操作
    groups.forEach((group) => {
      const sortMessageGroup = group;
      // 消息排序
      sortMessageGroup.messages = sortByTime(group.messages, false);
      commit(SET_MY_GROUP_GATHER, sortMessageGroup);
      commit(SET_ROOM_GATHER, sortMessageGroup);
      sortMessageGroup.members.forEach((member) => {
        commit(SET_USER_GATHER, member);
      });
    });
    const { friendGather, groupGather } = state;
    if (!state.activeRoom) {
      // 随机
      console.debug('set active room');
      commit(SET_ACTIVE_ROOM, groupGather[groups[0].groupID] || friendGather[friends[0].uid]);
    } else {
      console.debug('flush active room');
      commit(SET_ACTIVE_ROOM, state.activeRoom.uid ? friendGather[state.activeRoom.uid] : groupGather[state.activeRoom.groupID]);
    }
  },
};

export default actions;
