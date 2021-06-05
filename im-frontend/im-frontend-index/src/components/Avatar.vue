<template>
  <div class="avatar">
    <a-popover v-if="data.uid !== user.uid || data.from !== user.uid" trigger="click">
      <div slot="content" class="avatar-card">
        <a-card :bordered="false" style="width: 300px">
          <template slot="title">
            <h2>{{ data.account || (data.from && userGather[data.from].account) }}</h2>
            <a-avatar :size="60" style="float: right" :src="data.avatar || (data.from && userGather[data.from].avatar) || randomAvatar" />
          </template>
          <!--          <a-button-->
          <!--            v-if="user.uid === 'admin'"-->
          <!--            style="margin-bottom: 5px"-->
          <!--            @click="deleteUser(data.from)"-->
          <!--            :loading="loading"-->
          <!--            type="primary"-->
          <!--          >-->
          <!--            删除用户-->
          <!--          </a-button>-->
          <a-button @click="_setMyActiveRoom(data.from || data.uid)" type="primary" v-if="friendGather[data.from || data.uid]"
            >发消息</a-button
          >
          <a-button
            @click="addFriend(data.from || data.uid)"
            :loading="loading"
            type="primary"
            v-else-if="user.uid !== (data.from || data.uid)"
            >添加好友</a-button
          >
        </a-card>
      </div>
      <!--      左边-->
      <a-avatar
        :style="{ order: data.from === user.uid && highLight ? '3' : '1' }"
        class="avatar-img"
        :class="{ offLine: !data.online && highLight === false }"
        :src="data.avatar || (data.from && userGather[data.from].avatar) || randomAvatar"
      />
    </a-popover>
    <a-avatar
      v-else
      class="avatar-img"
      :style="{ order: data.from === user.uid && highLight ? '3' : '1' }"
      :class="{ offLine: !data.online && highLight === false }"
      :src="(userGather[data.from] && userGather[data.from].avatar) || randomAvatar"
    />
    <div class="avatar-name" style="order: 2">{{ (userGather[data.from] && userGather[data.from].account) || data.account }}</div>
    <div class="avatar-time" :style="{ order: data.from === user.uid && highLight ? '1' : '3' }" v-if="showTime">
      {{ _formatTime(data.time) }}
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';
import { namespace } from 'vuex-class';
import { formatTime } from '@/utils/common';
import { getRandomAvatar } from '@/common';

const chatModule = namespace('chat');
const appModule = namespace('app');

@Component
export default class Avatar extends Vue {
  @Prop() data: User & ChatMessage; // 列表或者用户

  @Prop({ default: true }) showTime: boolean; // 是否显示时间

  @Prop({ type: Boolean, default: false }) highLight: boolean; // 头像是否常亮

  @appModule.Getter('user') user: User;

  @appModule.Getter('loading') loading: boolean;

  @chatModule.Getter('userGather') userGather: UserGather;

  @chatModule.Getter('friendGather') friendGather: FriendGather;

  @chatModule.Getter('socket') socket: SocketIOClient.Socket;

  @chatModule.Mutation('set_active_room') setActiveRoom: Function;

  @appModule.Mutation('set_loading') setLoading: Function;

  addFriend(friendID: string) {
    // 设置按钮loading,避免网络延迟重复点击造成多次执行
    this.setLoading(true);
    this.socket.emit('addFriend', {
      friendA: this.user.uid,
      friendB: friendID,
    });
  }

  _formatTime(time: string) {
    return formatTime(parseInt(time.substring(0, 13), 10));
  }

  _setMyActiveRoom(uid: string) {
    this.setActiveRoom(this.friendGather[uid]);
  }

  get randomAvatar() {
    return getRandomAvatar();
  }
}
</script>
<style lang="scss" scoped>
.avatar {
  display: flex;
  align-items: center;
  height: 37px;
  margin-bottom: 6px;
  .avatar-img {
    cursor: pointer;
    width: 40px;
    height: 40px;
    border-radius: 0;
  }
  .offLine {
    filter: grayscale(100%);
  }
  .avatar-name {
    margin: 0 12px;
    color: #080808;
    margin: 0 12px;
    max-width: 160px;
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
    color: #080808;
  }
  .avatar-time {
    font-size: 12px;
    color: #080808;
  }
}
.avatar-card {
  display: flex;
  font-size: 18px;
  flex-direction: column;
  align-items: center;
  .ant-card-body {
    text-align: right;
  }
  h2 {
    display: inline-block;
    line-height: 60px;
    max-width: 190px;
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
  }
  > div {
    margin: 4px;
  }
}
</style>
