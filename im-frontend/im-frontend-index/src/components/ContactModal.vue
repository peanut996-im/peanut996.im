<template>
  <a-modal title="邀请好友加入本群" app :visible="showContactDialog" footer="" @cancel="cancelEvent">
    <a-transfer
      class="tree-transfer"
      :data-source="dataSource"
      :target-keys="targetKeys"
      show-search
      :filter-option="filterOption"
      :render="(item) => item.title"
      :show-select-all="false"
      @change="onChange"
    >
      <template slot="children" slot-scope="{ props: { direction, selectedKeys }, on: { itemSelect } }">
        <a-tree
          v-if="direction === 'left'"
          blockNode
          checkable
          checkStrictly
          defaultExpandAll
          :checkedKeys="[...selectedKeys, ...targetKeys]"
          :treeData="treeData"
          show-icon
          @check="
            (_, props) => {
              onChecked(_, props, [...selectedKeys, ...targetKeys], itemSelect);
            }
          "
          @select="
            (_, props) => {
              onChecked(_, props, [...selectedKeys, ...targetKeys], itemSelect);
            }
          "
        >
        </a-tree>
      </template>
      <template slot="footer" slot-scope="props">
        <template v-if="props.direction === 'right'">
          <a-button type="primary" slot="footer" size="small" style="float: right; margin: 5px" @click="onSubmit"> 添加 </a-button>
          <a-button slot="footer" size="small" style="float: right; margin: 5px" @click="showContactDialog = false"> 取消 </a-button>
        </template>
      </template>
    </a-transfer>
  </a-modal>
</template>
<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { namespace } from 'vuex-class';
import cnchar from 'cnchar';

const chatModule = namespace('chat');
const appModule = namespace('app');

function isChecked(selectedKeys: any, eventKey: any) {
  return selectedKeys.indexOf(eventKey) !== -1;
}

function handleTreeData(data: any, targetKeys: any = []) {
  data.forEach((item: any) => {
    // eslint-disable-next-line no-param-reassign
    item.disabled = targetKeys.includes(item.key);
    if (item.children) {
      handleTreeData(item.children, targetKeys);
    }
  });
  return data;
}

@Component
export default class ContactModal extends Vue {
  @chatModule.Getter('friendGather') friendGather: FriendGather;

  @chatModule.Getter('socket') socket: SocketIOClient.Socket;

  @chatModule.Action('inviteFriendsIntoGroup') inviteFriendsIntoGroup: Function;

  @chatModule.State('activeRoom') activeRoom: Group & Friend;

  @appModule.Getter('user') user: User;

  targetKeys: string[] = [];

  // 添加成员dialog
  showContactDialog: boolean = false;

  filterOption(inputValue: any, option: any) {
    return option.title.indexOf(inputValue) > -1;
  }

  showDialog() {
    this.showContactDialog = true;
  }

  // 获取联系人列表,按A-Z字母排序
  get myContactList() {
    const list = Object.values(this.friendGather).filter((friend) => !this.activeRoom.members!.some((member) => member.uid === friend.uid));
    const charList = list.map((k) => cnchar.spell(k.account).toString().charAt(0).toUpperCase()).sort();
    const myContactList = [] as any;
    // eslint-disable-next-line no-restricted-syntax
    for (const char of Array.from(new Set(charList))) {
      // eslint-disable-next-line no-restricted-syntax
      myContactList.push({
        key: char,
        title: char,
        disabled: true,
        children: list
          .filter((k) => cnchar.spell(k.account).toString().charAt(0).toUpperCase() === char)
          .map((t) => ({
            key: t.uid,
            title: t.account,
            avatar: t.avatar,
          })),
      });
    }
    return myContactList;
  }

  get dataSource() {
    const transferDataSource: any = [];
    // 数组扁平化
    function flatten(list: any = []) {
      list.forEach((item: any) => {
        // eslint-disable-next-line no-param-reassign
        delete item.disabled;
        if (!item.children) {
          transferDataSource.push(item);
        }
        if (item.children) {
          flatten(item.children);
        }
      });
    }
    flatten(JSON.parse(JSON.stringify(this.treeData)));
    return transferDataSource;
  }

  get treeData() {
    return handleTreeData(this.myContactList, this.targetKeys);
  }

  onSubmit() {
    if (this.targetKeys.length === 0) {
      this.$message.warning('请选择联系人');
      return;
    }
    this.$emit('inviteFriend', {
      groupID: this.activeRoom.groupID,
      friends: this.targetKeys,
    });
    this.cancelEvent();
  }

  onChange(targetKeys: string[]) {
    console.log('Target Keys:', targetKeys);
    // eslint-disable-next-line no-restricted-syntax
    for (const targetKey of targetKeys) {
      if (targetKey.length === 1 && targetKey.charCodeAt(0) > 64 && targetKey.charCodeAt(0) < 91) {
        return this.$message.error('请勿选择标题', 0.5);
      }
    }
    this.targetKeys = targetKeys;
  }

  onChecked(_: any, e: any, checkedKeys: any, itemSelect: any) {
    const { eventKey } = e.node;
    itemSelect(eventKey, !isChecked(checkedKeys, eventKey));
  }

  cancelEvent() {
    this.showContactDialog = false;
    this.targetKeys = [];
  }
}
</script>
<style scoped>
.tree-transfer .ant-transfer-list:first-child {
  width: 50%;
  flex: none;
}
</style>
