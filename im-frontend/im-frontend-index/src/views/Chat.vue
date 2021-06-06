<template>
  <div
    class="chat"
    :style="{
      '--bg-image': `url('${background}')`,
    }"
  >
    <!-- 左侧导航栏 -->
    <div class="chat-part1" v-if="visibleNav">
      <Nav @logout="goToSSO" @updateUser="updateUser"></Nav>
    </div>
    <!-- 消息列表/通讯人列表 -->
    <div class="chat-part2">
      <template v-if="activeTabName === 'message'">
        <Search
          @findUser="findUser"
          @findGroup="findGroup"
          @createGroup="createGroup"
          @joinGroup="joinGroup"
          @addFriend="addFriend"
          @setActiveRoom="setActiveRoom"
        >
        </Search>
        <Room @setActiveRoom="setActiveRoom"></Room>
      </template>
      <template v-else>
        <Contact @addFriend="addFriend" @setActiveRoom="setActiveRoom"></Contact>
      </template>
    </div>
    <!-- 右侧聊天窗口 -->
    <div class="chat-part3">
      <a-icon class="chat-team" type="message" @click="toggleDrawer" />
      <div class="chat-nav">
        <a-icon type="menu-fold" @click="toggleNav" v-if="visibleNav" />
        <a-icon type="menu-unfold" @click="toggleNav" v-else />
      </div>
      <Message
        v-if="activeRoom"
        @pullMessage="pullMessage"
        @deleteFriend="deleteFriend"
        @leaveGroup="leaveGroup"
        @inviteFriend="inviteFriend"
      ></Message>
    </div>

    <!-- 移动端兼容 -->
    <a-drawer placement="left" :closable="false" :visible="visibleDrawer" @close="toggleDrawer" style="height: 100%">
      <div class="chat-drawer">
        <template v-if="activeTabName === 'message'">
          <Search
            @findUser="findUser"
            @findGroup="findGroup"
            @createGroup="createGroup"
            @joinGroup="joinGroup"
            @addFriend="addFriend"
            @setActiveRoom="setActiveRoom"
            @addGroup="createGroup"
          >
          </Search>
          <Room @setActiveRoom="setActiveRoom"></Room>
        </template>
        <template v-else>
          <Contact @addFriend="addFriend"></Contact>
        </template>
      </div>
    </a-drawer>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Nav from '@/components/Nav.vue';
import Room from '@/components/Room.vue';
import Message from '@/components/Message.vue';
import Search from '@/components/Search.vue';
import Contact from '@/components/Contact.vue';
import { namespace } from 'vuex-class';

const appModule = namespace('app');
const chatModule = namespace('chat');

@Component({
  components: {
    Nav,
    Room,
    Message,
    Search,
    Contact,
  },
})
export default class Chat extends Vue {
  @appModule.Getter('user') user: User;

  @appModule.Getter('uid') uid: string;

  @appModule.Getter('token') token: string;

  @appModule.Getter('isLogined') isLogined: boolean;

  @appModule.Mutation('clear_user') clearUser: Function;

  @appModule.Action('goToSSO') goToSSO: Function;

  @appModule.Mutation('set_uid') setUID: Function;

  @appModule.Mutation('set_token') setToken: Function;

  @appModule.Mutation('set_uid_and_token') setUIDAndToken: Function;

  @appModule.Getter('background') background: string;

  @appModule.Getter('activeTabName') activeTabName: string;

  @appModule.Mutation('set_activeTabName') _setActiveTabName: Function;

  @chatModule.Getter('socket') socket: SocketIOClient.Socket;

  @chatModule.Getter('activeRoom') activeRoom: Friend | Group;

  @chatModule.Mutation('set_active_room') _setActiveRoom: Function;

  @chatModule.Action('connectSocket') connectSocket: Function;

  visibleDrawer: boolean = false;

  visibleNav: boolean = true;

  created() {
    // 单点登陆/获取url链接中传递的userName,直接后台默认注册登陆
    const { uid, token } = this.$route.query;
    if (uid && token) {
      this.setUID(uid);
      this.setToken(token);
      this.handleConnect();
    } else if (!this.isLogined) {
      this.goToSSO();
    } else {
      this.setUIDAndToken();
      this.handleConnect();
    }
  }

  // 登录
  async handleLogin(user: User) {
    const res = await this.goToSSO(user);
    if (res) {
      // 进入系统事件
      await this.handleConnect();
    }
  }

  // 进入系统初始化事件
  async handleConnect() {
    this.connectSocket();
  }

  // 创建群组
  createGroup(groupName: string) {
    console.debug(`create group ${this.user.uid} ${groupName}`);
    this.socket.emit('createGroup', {
      uid: this.user.uid,
      groupName,
    });
  }

  // 加入群组
  joinGroup(group: Group) {
    this.socket.emit('joinGroup', {
      uid: this.user.uid,
      groupID: group.groupID,
    });
  }

  leaveGroup(group: Group) {
    this.socket.emit('leaveGroup', {
      uid: this.user.uid,
      groupID: group.groupID,
    });
  }

  findUser(user: User) {
    this.socket.emit('findUser', {
      account: user.account,
    });
  }

  findGroup(group: Group) {
    this.socket.emit('findGroup', {
      groupName: group.groupName,
    });
  }

  // 添加好友/发起私聊窗口
  addFriend(friend: User) {
    this.socket.emit('addFriend', {
      friendA: this.user.uid,
      friendB: friend.uid,
    });
  }

  inviteFriend(data: any) {
    this.socket.emit('inviteFriend', {
      friends: data.friends,
      groupID: data.groupID,
    });
  }

  pullMessage(data: any) {
    this.socket.emit('pullMessage', data);
  }

  deleteFriend(friend: User) {
    this.socket.emit('deleteFriend', {
      friendA: this.user.uid,
      friendB: friend.uid,
    });
  }

  updateUser(data: any) {
    this.socket.emit('updateUser', data);
  }

  setActiveRoom(room: Friend | Group) {
    console.debug('change setMyActiveRoom On Chat');
    this._setActiveRoom(room);
  }

  toggleDrawer() {
    this.visibleDrawer = !this.visibleDrawer;
  }

  toggleNav() {
    this.visibleNav = !this.visibleNav;
  }
}
</script>
<style lang="scss" scoped>
@import '@/styles/theme';

.chat {
  font-size: 16px;
  z-index: 999;
  max-width: 1000px;
  min-width: 300px;
  width: 100%;
  height: 80%;
  max-height: 800px;
  min-height: 470px;
  position: relative;
  margin: auto 20px;
  box-shadow: 10px 20px 80px rgba(0, 0, 0, 0.8);
  display: flex;
  border-radius: 8px;
  overflow: hidden;
  .chat-part1 {
    width: 74px;
    height: 100%;
    background-color: rgb(0, 0, 0, 0.7);
  }
  .chat-part2 {
    width: 260px;
    height: 100%;
    background-color: rgb(0, 0, 0, 0.3);
    border-right: 1px solid #d6d6d6;
    background: $room-bg-color;
    overflow: auto;
  }
  .chat-part3 {
    flex: 1;
    height: 100%;
    background: $room-bg-color;
    // background-color: rgb(0, 0, 0, 0.2);
    overflow-y: hidden;
    position: relative;
    .chat-group {
      height: 53px;
      border-bottom: 1px solid #ccc;
      line-height: 50px;
      font-weight: bold;
    }
  }
  .chat-team {
    display: none;
  }
  .chat-nav {
    display: none;
  }
}
.chat::after {
  content: '';
  background: var(--bg-image) 0 / cover fixed;
  position: absolute;
  object-fit: cover;
  width: 100%;
  height: 100%;
  filter: blur(10px);
  transform: scale(1.08);
  z-index: -1;
}

@media screen and (max-width: 768px) {
  .chat {
    margin: 0;
    height: 100%;
    .chat-part2 {
      display: none;
    }
    .chat-team {
      display: block !important;
      position: absolute;
      font-size: 20px;
      top: 20px;
      color: #080808;
      right: 100px;
      z-index: 999;
      &:active {
        color: $primary-color;
      }
    }
    .chat-nav {
      display: block !important;
      position: absolute;
      font-size: 20px;
      top: 16px;
      color: #080808;
      right: 62px;
      z-index: 999;
      &:active {
        color: $primary-color;
      }
    }
  }
}
</style>
