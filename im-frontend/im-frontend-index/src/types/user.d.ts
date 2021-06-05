interface User {
  uid: string;
  userId: string;
  username: string;
  account: string;
  password: string;
  avatar: string;
  role?: string;
  tag?: string;
  createTime: number;
  online?: 1 | 0; // 是否在线
}
