/* eslint-disable no-unused-vars */
declare module 'socket.io-client';

// 图片尺寸
interface ImageSize {
  width: number;
  height: number;
}

// 未读消息对象
interface UnReadGather {
  [key: string]: number;
}

// 右键菜单操作烈性
declare enum ContextMenuType {
  COPY = 'COPY', // 复制
  TOP_REVERT = 'TOP_REVERT', // 取消置顶
  TOP = 'TOP', // 置顶
  READ = 'READ', // 一键已读
  DELETE = 'DELETE', // 删除
}

interface ChatMessage {
  id?: string;
  from: string;
  to: string;
  time?: string;
  type: string;
  content: string;
  file?: File;
  fileName?: string;
  width?: number;
  height?: number;
  size?: number;
}
// 好友
interface Friend {
  uid: string;
  account: string;
  roomID: string;
  avatar?: string;
  messages?: ChatMessage[];
  createTime?: number;
  isTop?: boolean; // 是否置顶聊天
  online?: 1 | 0; // 是否在线
  isManager?: 1 | 0; // 是否为群主
}
// 群组
interface Group {
  groupID: string;
  groupAdmin: string; // 群主id
  groupName: string;
  notice?: string;
  messages?: ChatMessage[];
  createTime?: number;
  isTop?: boolean; // 是否置顶聊天
  members?: Friend[]; // 群成员列表
}

// 服务端返回值格式
interface Response {
  code: number;
  message: string;
  data: any;
}

// 所有好友的好友信息
interface FriendGather {
  [uid: string]: Friend;
}
// 所有群的群信息
interface GroupGather {
  [groupID: string]: Group;
}

interface FriendMap {
  friendA: string;
  friendB: string;
}

interface GroupMap {
  uid: string;
  groupID: string;
}

interface RoomGather {
  [to: string]: Friend & Group;
}

interface UserGather {
  [uid: string]: User;
}
