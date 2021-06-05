// default background
const bg = `${process.env.VUE_APP_OSS_BUCKET_URL}/catalina-background.jpg`;

export const DEFAULT_BACKGROUND = bg;

// MIME类型
export const MIME_TYPE = ['xls', 'xlsx', 'doc', 'docx', 'exe', 'pdf', 'ppt', 'txt', 'zip', 'img', 'rar'];
// 图片类型
export const IMAGE_TYPE = ['png', 'jpg', 'jpeg', 'gif'];

/**
 * 返回随机头像
 */
export function getRandomAvatar() {
  const baseOSSUrl = process.env.VUE_APP_OSS_BUCKET_URL;
  const random = Math.floor(Math.random() * 20) + 1;
  return `${baseOSSUrl}/avatar/avatar${random}.png`;
}
