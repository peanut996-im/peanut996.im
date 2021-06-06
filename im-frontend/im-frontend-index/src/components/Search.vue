<template>
  <div class="search">
    <div class="search-select">
      <a-select
        show-search
        size="small"
        placeholder="搜索聊天组"
        :default-active-first-option="false"
        :show-arrow="false"
        :filter-option="false"
        :not-found-content="null"
        @search="handleMySearch"
      >
        <a-icon slot="suffixIcon" type="search" />
        <a-select-option v-for="(chat, index) in mySearchData" :key="index" @click="selectMyChat(chat)">
          <div v-if="chat.account">{{ chat.account }}</div>
          <div v-if="chat.groupName">{{ chat.groupName }}</div>
        </a-select-option>
      </a-select>

      <a-dropdown class="search-dropdown">
        <a-icon type="plus-square" class="search-dropdown-button" />
        <a-menu slot="overlay">
          <a-menu-item>
            <div @click="() => (visibleAddGroup = !visibleAddGroup)">创建群</div>
          </a-menu-item>
          <a-menu-item>
            <div @click="() => (visibleJoinGroup = !visibleJoinGroup)">搜索群</div>
          </a-menu-item>
          <a-menu-item>
            <div @click="() => (visibleAddFriend = !visibleAddFriend)">搜索用户</div>
          </a-menu-item>
        </a-menu>
      </a-dropdown>
    </div>

    <a-modal v-model="visibleAddGroup" footer="" title="创建群聊">
      <div style="display: flex">
        <a-input v-model="groupName" placeholder="请输入群名字"></a-input>
        <a-button @click="createGroup" :loadig="loading" type="primary">确定</a-button>
      </div>
    </a-modal>
    <a-modal v-model="visibleJoinGroup" footer="" title="搜索群组">
      <div style="display: flex" v-if="visibleJoinGroup">
        <a-select
          show-search
          placeholder="请输入群名字"
          style="width: 90%"
          :default-active-first-option="false"
          :show-arrow="false"
          :filter-option="false"
          :not-found-content="null"
          @search="handleMyGroupSearch"
          @change="handleMyGroupChange"
        >
          <a-select-option v-for="(group, index) in groupSearchResult" :key="index" @click="handleMyGroupSelect(group)">
            <div>{{ group.groupName }}</div>
          </a-select-option>
        </a-select>
        <a-button @click="joinGroup" type="primary" :loading="loading">加入群</a-button>
      </div>
    </a-modal>
    <a-modal v-model="visibleAddFriend" footer="" title="创建聊天/搜索用户">
      <div style="display: flex" v-if="visibleAddFriend">
        <a-select
          show-search
          placeholder="请输入用户名"
          style="width: 90%"
          :default-active-first-option="false"
          :show-arrow="false"
          :filter-option="false"
          :not-found-content="null"
          @search="handleMyUserSearch"
          @change="handleMyUserChange"
        >
          <a-select-option v-for="(user, index) in userSearchResult" :key="index" @click="handleMyUserSelect(user)">
            <div>{{ user.account }}</div>
          </a-select-option>
        </a-select>
        <a-button @click="addFriend" type="primary" :loading="loading">添加好友</a-button>
      </div>
    </a-modal>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator';
import { namespace } from 'vuex-class';
import { isContainStr, nameVerify } from '@/utils/common';

const chatModule = namespace('chat');
const appModule = namespace('app');

@Component
export default class Search extends Vue {
  @appModule.Getter('loading') loading: boolean;

  @chatModule.Getter('groupGather') groupGather: GroupGather;

  @chatModule.Getter('friendGather') friendGather: FriendGather;

  @chatModule.Getter('groupSearchResult') groupSearchResult: Array<Group>;

  @chatModule.Getter('userSearchResult') userSearchResult: Array<Friend>;

  @appModule.Mutation('set_loading') setLoading: Function;

  visibleAddGroup: boolean = false;

  visibleJoinGroup: boolean = false;

  visibleAddFriend: boolean = false;

  groupName: string = '';

  mySearchData: Array<Group | Friend> = [];

  mySearchResult: Array<Group | User> = [];

  groupID: string = '';

  myGroupArr: Array<Group> = [];

  selected: User | Group | null;

  userArr: Array<User> = [];

  myUserArr: Array<User> = [];

  created() {
    this.getMySearchData();
  }

  // 监控群组列表
  @Watch('groupGather')
  changeMyGroupGather() {
    this.getMySearchData();
  }

  // 监控群组列表
  @Watch('friendGather')
  changeMyFriendGather() {
    this.getMySearchData();
  }

  getMySearchData() {
    this.mySearchData = [...Object.values(this.friendGather), ...Object.values(this.groupGather)];
  }

  // 消息列表查询
  handleMySearch(value: string) {
    const newSearchData: Array<Friend | Group> = [];
    this.mySearchData = [...Object.values(this.groupGather), ...Object.values(this.friendGather)];
    // eslint-disable-next-line no-restricted-syntax
    for (const chat of this.mySearchData) {
      if ((chat as Friend).account) {
        if (isContainStr(value, (chat as Friend).account)) {
          newSearchData.push(chat);
        }
      } else if (isContainStr(value, (chat as Group).groupName)) {
        newSearchData.push(chat);
      }
    }
    this.mySearchData = newSearchData;
  }

  handleMyGroupSelect(group: Group) {
    this.selected = group;
  }

  handleMyGroupChange() {
    this.mySearchResult = [];
  }

  async handleMyUserSearch(value: string) {
    if (!value) {
      return;
    }
    const data = {
      account: value,
    };
    this.$emit('findUser', data);
  }

  // 群组名查询
  async handleMyGroupSearch(value: string) {
    if (!value) {
      return;
    }
    const data = {
      groupName: value,
    };
    this.$emit('findGroup', data);
  }

  handleMyUserSelect(u: User) {
    // this.friend.uid = friend.uid;
    this.selected = u;
  }

  handleMyUserChange() {
    this.mySearchResult = [];
  }

  selectMyChat(activeRoom: Friend & Group) {
    this.$emit('setActiveRoom', activeRoom);
  }

  @Watch('groupSearchResult')
  changeGroupSearchResult() {
    this.mySearchData = this.groupSearchResult;
  }

  @Watch('userSearchResult')
  changeUserSearchResult() {
    this.mySearchData = this.userSearchResult;
  }

  // 创建群组
  createGroup() {
    this.setLoading(true);
    this.visibleAddGroup = false;
    if (!nameVerify(this.groupName)) {
      this.visibleAddGroup = true;
      return;
    }
    console.debug('start emit create group to chat');
    this.$emit('createGroup', this.groupName);
    this.groupName = '';
  }

  joinGroup() {
    this.setLoading(true);
    this.visibleJoinGroup = false;
    this.$emit('joinGroup', this.selected);
    this.selected = null;
  }

  addFriend() {
    this.setLoading(true);
    this.visibleAddFriend = false;
    this.$emit('addFriend', this.selected);
    this.selected = null;
  }
}
</script>
<style lang="scss" scoped>
@import '@/styles/theme';

.search {
  background: $room-bg-color;
  position: relative;
  height: 60px;
  padding: 10px;
  display: flex;
  align-items: center;
  .search-select {
    width: 80%;
    .ant-select {
      width: 100%;
    }
  }
  .search-dropdown {
    position: absolute;
    right: 15px;
    top: 13px;
    font-size: 20px;
    padding: 0;
    cursor: pointer;
    line-height: 40px;
    color: gray;
    transition: 0.2s all linear;
    border-radius: 4px;
    &:hover {
      color: $primary-color;
    }
  }
}
</style>
