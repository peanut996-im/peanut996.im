<template>
  <div class="room" ref="room">
    <div v-for="chat in myChatArr" :key="chat.groupID || chat.uid">
      <div
        v-if="chat.groupID"
        class="room-card"
        :class="{ active: activeRoom && activeRoom.groupID === chat.groupID }"
        @click="changeActiveRoom(chat)"
        v-contextmenu="'groupmenu' + chat.groupID"
      >
        <!-- 自定义右键菜单 -->
        <v-contextmenu :ref="'groupmenu' + chat.groupID">
          <v-contextmenu-item v-if="chat.isTop === true" @click="handleMyCommand('TOP_REVERT', chat)">取消置顶</v-contextmenu-item>
          <v-contextmenu-item v-else @click="handleMyCommand('TOP', chat)">置顶</v-contextmenu-item>
          <v-contextmenu-item @click="handleMyCommand('READ', chat)">标记已读</v-contextmenu-item>
          <v-contextmenu-item divider></v-contextmenu-item>
          <v-contextmenu-item @click="handleMyCommand('DELETE', chat)">删除</v-contextmenu-item>
        </v-contextmenu>
        <a-badge class="room-card-badge" v-if="unReadGather[chat.groupID]" :count="unReadGather[chat.groupID]" />
        <img class="room-card-type" src="~@/assets/group.png" alt="" />
        <div class="room-card-message">
          <div class="room-card-info">
            <div class="room-card-name">{{ chat.groupName }}</div>
            <!-- 显示最后一次聊天时间 -->
            <div
              class="room-card-time"
              v-if="chat.messages && chat.messages[chat.messages.length - 1]"
              v-text="_formatTime(chat.messages[chat.messages.length - 1])"
            ></div>
          </div>
          <div class="room-card-new" v-if="chat.messages.length > 0">
            <!-- 消息列表未读信息简述考虑撤回情况 -->
            <template v-if="chat.messages[chat.messages.length - 1].isRevoke">
              <div>{{ chat.messages[chat.messages.length - 1].revokeUserName }}撤回了一条消息</div>
            </template>
            <template v-else>
              <div
                v-text="_parseText(chat.messages[chat.messages.length - 1])"
                v-if="chat.messages[chat.messages.length - 1].type === 'text'"
              ></div>
              <div class="image" v-if="chat.messages[chat.messages.length - 1].type === 'image'">[图片]</div>
              <div class="image" v-if="chat.messages[chat.messages.length - 1].type === 'file'">[文件]</div>
            </template>
          </div>
        </div>
      </div>
      <div
        v-else
        class="room-card"
        :class="{ active: activeRoom && !activeRoom.groupID && activeRoom.uid === chat.uid }"
        @click="changeActiveRoom(chat)"
        v-contextmenu="'contextmenu' + chat.uid"
      >
        <!-- 自定义右键菜单 -->
        <v-contextmenu :ref="'contextmenu' + chat.uid">
          <v-contextmenu-item v-if="chat.isTop === true" @click="handleMyCommand('TOP_REVERT', chat)">取消置顶</v-contextmenu-item>
          <v-contextmenu-item v-else @click="handleMyCommand('TOP', chat)">置顶</v-contextmenu-item>
          <v-contextmenu-item @click="handleMyCommand('READ', chat)">标记已读</v-contextmenu-item>
          <v-contextmenu-item divider></v-contextmenu-item>
          <v-contextmenu-item @click="handleMyCommand('DELETE', chat)">删除</v-contextmenu-item>
        </v-contextmenu>
        <a-badge class="room-card-badge" :count="unReadGather[chat.roomID]" />
        <!--三元表达式: 删除后 先验证friendGather 再验证 avatar是否存在-->
        <img
          class="room-card-type"
          :src="!friendGather[chat.uid] ? randomAvatar : friendGather[chat.uid].avatar ? friendGather[chat.uid].avatar : randomAvatar"
          alt=""
        />
        <div class="room-card-message">
          <div class="room-card-info">
            <div class="room-card-name">{{ chat.account }}</div>
            <!-- 显示最后一次聊天时间 -->
            <div
              class="room-card-time"
              v-if="chat.messages.length && chat.messages[chat.messages.length - 1]"
              v-text="_formatTime(chat.messages[chat.messages.length - 1])"
            ></div>
          </div>
          <div class="room-card-new" v-if="chat.messages.length > 0">
            <template v-if="chat.messages">
              <div
                v-text="_parseText(chat.messages[chat.messages.length - 1])"
                v-if="chat.messages[chat.messages.length - 1].type === 'text'"
              ></div>
              <div class="image" v-if="chat.messages[chat.messages.length - 1].type === 'image'">[图片]</div>
              <div class="image" v-if="chat.messages[chat.messages.length - 1].type === 'file'">[文件]</div>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator';
import { namespace } from 'vuex-class';
import { parseText, formatTime } from '@/utils/common';
import { getRandomAvatar } from '@/common';

const chatModule = namespace('chat');
const appModule = namespace('app');
@Component
export default class Room extends Vue {
  @chatModule.State('activeRoom') activeRoom: Group & Friend;

  @chatModule.Getter('groupGather') groupGather: GroupGather;

  @chatModule.Getter('friendGather') friendGather: FriendGather;

  @chatModule.Getter('roomGather') roomGather: RoomGather;

  @chatModule.Getter('userGather') userGather: UserGather;

  @chatModule.Getter('unReadGather') unReadGather: UnReadGather;

  @chatModule.Mutation('lose_unread_gather') lose_unread_gather: Function;

  @appModule.Getter('user') user: User;

  myChatArr: Array<Group | Friend> = [];

  created() {
    this.renderChat();
  }

  mounted() {
    // hack方法 页面初始化时定位到当前room
    setTimeout(() => {
      this.setRoomScrollTop();
    }, 100);
  }

  // 重置滚动条至当前activeRoom位置
  setRoomScrollTop() {
    const { offsetHeight: roomHeight, scrollTop: roomTop } = (document.querySelector('.room') as HTMLElement)!;
    const activeRommDom = document.querySelector('.room-card.active') as HTMLElement;
    if (activeRommDom) {
      const { offsetTop: domTop } = activeRommDom!;
      if (domTop - roomHeight >= roomTop) {
        document.querySelector('.room')!.scrollTop = domTop - roomHeight;
      }
    }
  }

  @Watch('groupGather', { deep: true })
  changeMyGroupGather() {
    // this.sortChat();
    this.renderChat();
  }

  @Watch('friendGather', { deep: true })
  changeMyFriendGather() {
    // this.sortChat();
    this.renderChat();
  }

  @Watch('activeRoom')
  changeMyActiveRoomEvent() {
    this.$nextTick(() => {
      this.setRoomScrollTop();
    });
  }

  get currentUID() {
    return this.user.uid;
  }

  get randomAvatar() {
    return getRandomAvatar();
  }

  async handleMyCommand(type: ContextMenuType, chat: Group | Friend) {
    //   // 消息ID
    const chatID = (chat as Group).groupID || (chat as Friend).uid;
    if (type === 'TOP') {
      // 修复重复置顶bug,在已置顶某个窗口的情况下 直接置顶另外一个,需要先取消第一个置顶的窗口
      const topID = await this.$localforage.getItem(`${this.currentUID}-topChatID`);
      // continue here.
      if (topID) {
        const topRoom = this.myChatArr.find((room) => ((room as Group).groupID || (room as Friend).roomID) === topID);
        if (topRoom) {
          delete topRoom.isTop;
        }
      }
      await this.$localforage.setItem(`${this.currentUID}-topChatID`, chatID);
      this.renderChat().then(() => {
        this.$message.success('置顶成功');
      });
    } else if (type === 'TOP_REVERT') {
      await this.$localforage.removeItem(`${this.currentUID}-topChatID`);
      // 删除isTop属性,取消置顶
      // eslint-disable-next-line no-param-reassign
      delete chat.isTop;
      this.renderChat().then(() => {
        this.$message.info('取消置顶');
      });
    } else if (type === 'READ') {
      this.lose_unread_gather(chatID);
    } else if (type === 'DELETE') {
      // 如果聊天列表仅有一个消息不允许删除
      console.debug(this.myChatArr);
      // 先查询本地时候有删除记录
      const deletedChat = (await this.$localforage.getItem(`${this.currentUID}-deletedChatID`)) as string[];
      if (Array.isArray(deletedChat)) {
        if (!deletedChat.includes(chatID)) {
          deletedChat.push(chatID);
        }
        await this.$localforage.setItem(`${this.currentUID}-deletedChatID`, deletedChat);
      } else {
        // 本地删除聊天(非删除好友,本地记录)
        await this.$localforage.setItem(`${this.currentUID}-deletedChatID`, [chatID]);
      }
      // 删除聊天窗口后默认激活第一个聊天窗口
      this.renderChat().then(() => {
        this.$message.success(`已删除${(chat as Group).groupName || (chat as Friend).account}聊天窗口`);
      });

      this.changeActiveRoom(this.myChatArr[0] as Friend | Group);
    }
  }

  // 用户是否离线状态
  avatarOffLine() {
    // 机器人默认在线
    // return chat.uid === DEFAULT_ROBOT ? false : !chat.online;
    // 默认保持在线
    return false;
  }

  /**
   * 渲染聊天列表
   *
   */
  async renderChat() {
    console.debug('render chats now!');
    const groups = Object.values(this.groupGather);
    const friends = Object.values(this.friendGather);

    let roomArr = [...groups, ...friends];
    // 过滤已删除
    const deletedChat = (await this.$localforage.getItem(`${this.currentUID}-deletedChatID`)) as string[];
    if (Array.isArray(deletedChat)) {
      roomArr = roomArr.filter((chat) => !deletedChat.includes((chat as Group).groupID || (chat as Friend).uid));
    }

    // sort by messages
    roomArr = roomArr.sort((a: Group | Friend, b: Group | Friend) => {
      if (a.messages?.length && b.messages?.length) {
        // @ts-ignore
        return parseInt(b.messages[b.messages.length - 1].time, 10) - parseInt(a.messages[a.messages.length - 1].time, 10);
      }
      if (a.messages?.length) {
        return -1;
      }
      return 1;
    });
    // pin to top
    const topChatID = (await this.$localforage.getItem(`${this.currentUID}-topChatID`)) as string;
    if (topChatID) {
      // 找到需要置顶的窗口
      const chat = roomArr.find((c) => ((c as Group).groupID || (c as Friend).uid) === topChatID);
      if (chat) {
        // 移动至第一位
        roomArr = roomArr.filter((k) => ((k as Group).groupID || (k as Friend).uid) !== topChatID);
        chat.isTop = true;
        roomArr.unshift(chat);
      }
    }
    // 此处避免Await造成v-for页面闪烁问题,所以在最后才赋值this.chatArr = roomArr;
    this.myChatArr = roomArr;
  }

  changeActiveRoom(activeRoom: Friend | Group) {
    console.debug('set my active room in Room');
    this.$emit('setActiveRoom', activeRoom);
    this.lose_unread_gather((activeRoom as Group).groupID || (activeRoom as Friend).roomID);
  }

  _parseText(chat: ChatMessage) {
    if (chat.to) {
      const unReadCount = this.unReadGather[chat.to];
      if (unReadCount && unReadCount > 1) {
        return `[${this.unReadGather[chat.to]}条] ${this.userGather[chat.from].account}:${parseText(chat.content)}`;
      }
      return `${this.userGather[chat.from].account}:${parseText(chat.content)}`;
    }
    return parseText(chat.content);
  }

  _formatTime(chat: ChatMessage) {
    return formatTime(parseInt(chat?.time!.substring(0, 13), 10));
  }
}
</script>
<style lang="scss" scoped>
@import '@/styles/theme';

@mixin button($bcolor, $url, $x1, $y1, $bor, $col) {
  background: $bcolor;
  -webkit-mask: url($url);
  mask: url($url);
  -webkit-mask-size: $x1 $y1;
  mask-size: $x1 $y1;
  border: $bor;
  color: $col;
}

.room {
  height: calc(100% - 60px);
  overflow: auto;
  background: $room-bg-color;
  .room-card {
    position: relative;
    min-height: 65px;
    display: flex;
    align-items: center;
    border-bottom: 1px solid #e8e8e8 !important;
    // background-color: rgba(0, 0, 0, 0.2);
    padding: 5px 10px;
    text-align: left;
    transition: all 0.2s linear;
    cursor: pointer;
    &:hover,
    &.active {
      background-color: #d6d6d6;
    }
    .room-card-badge {
      position: absolute;
      left: 40px;
      top: 10px;
      ::v-deep.ant-badge-count {
        box-shadow: none;
        width: 10px;
      }
    }
    .room-card-type {
      width: 40px;
      height: 40px;
      margin-right: 10px;
      // border-radius: 50%;
      object-fit: cover;
      &.offLine {
        filter: grayscale(100%);
      }
    }
    .room-card-message {
      flex: 1;
      display: flex;
      width: 75%;
      flex-direction: column;
      .room-card-info {
        .room-card-name {
          overflow: hidden; //超出的文本隐藏
          text-overflow: ellipsis; //溢出用省略号显示
          white-space: nowrap; //溢出不换行
          color: #474747;
          font-weight: bold;
          font-size: 16px;
          display: inline-block;
          max-width: 110px;
        }
        .room-card-time {
          overflow: hidden; //超出的文本隐藏
          text-overflow: ellipsis; //溢出用省略号显示
          white-space: nowrap; //溢出不换行
          color: #a9a9a9;
          font-size: 14px;
          float: right;
        }
      }

      .room-card-new {
        > * {
          color: #a9a9a9;
          display: block;
          overflow: hidden; //超出的文本隐藏
          text-overflow: ellipsis; //溢出用省略号显示
          white-space: nowrap; //溢出不换行
        }
        color: rgb(255, 255, 255, 0.6);
        font-size: 14px;
      }
    }
  }
}

@keyframes ani {
  from {
    -webkit-mask-position: 100% 0;
    mask-position: 100% 0;
  }

  to {
    -webkit-mask-position: 0 0;
    mask-position: 0 0;
  }
}
</style>
