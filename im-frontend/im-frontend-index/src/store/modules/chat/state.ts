export interface ChatState {
  socket: SocketIOClient.Socket;
  dropped: boolean;
  activeRoom: (Group & Friend) | null;
  userGather: UserGather;
  friendGather: FriendGather;
  groupGather: GroupGather;
  roomGather: RoomGather;
  unReadGather: UnReadGather;
  userSearchResult: Array<User>;
  groupSearchResult: Array<Group>;
  pullMessageResult: Array<ChatMessage>;
}

const chatState: ChatState = {
  // @ts-ignore
  socket: null, // ws实例
  dropped: false, // 是否断开连接
  activeRoom: null,
  userGather: {},
  unReadGather: {}, // 所有会话未读消息集合
  friendGather: {},
  groupGather: {},
  roomGather: {},
  userSearchResult: [],
  groupSearchResult: [],
  pullMessageResult: [],
};

export default chatState;
