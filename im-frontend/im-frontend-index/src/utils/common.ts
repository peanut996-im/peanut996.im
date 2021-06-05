import Vue from 'vue';
import sha1 from 'js-sha1';
import UUID from 'uuid-int';

export function handleResponse(res: Response) {
  const { code, message, data } = res;
  if (code) {
    Vue.prototype.$message.error(message);
    return;
  }
  if (message) {
    Vue.prototype.$message.success(message);
  }
  return data;
}

// 判断一个字符串是否包含另外一个字符串
export function isContainStr(str1: string, str2: string) {
  return str2.indexOf(str1) >= 0;
}

/**
 * 屏蔽词
 * @param text 文本
 */
export function parseText(text: string) {
  return text;
}

/**
 * 判断是否URL
 * @param text 文本
 */
export function isUrl(text: string) {
  // 解析网址
  // eslint-disable-next-line no-useless-escape
  const UrlReg = new RegExp(/http(s)?:\/\/([\w-]+\.)+[\w-]+(\/[\w- .\/?%&=]*)?/);
  return UrlReg.test(text);
}

/**
 * 消息时间格式化
 * @param time
 */
export function formatTime(time: number) {
  const moment = Vue.prototype.$moment;
  // 大于昨天
  if (moment().add(-1, 'days').startOf('day') > time) {
    return moment(time).format('M/D HH:mm');
  }
  // 昨天
  if (moment().startOf('day') > time) {
    return `昨天 ${moment(time).format('HH:mm')}`;
  }
  // 大于五分钟不显示秒
  // if (new Date().valueOf() > time + 300000) {
  //   return moment(time).format('HH:mm');
  // }
  return moment(time).format('HH:mm');
}

/**
 * 转换时间字符串为数字 unix
 * @param time
 */
export function convertTime(time: string): number {
  return parseInt(time.substring(0, 13), 10);
}

/**
 * 群名/用户名校验
 * @param name
 */
export function nameVerify(name: string): boolean {
  const nameReg = /^(?!_)(?!.*?_$)[a-zA-Z0-9_\u4e00-\u9fa5]+$/;
  if (name.length === 0) {
    Vue.prototype.$message.error('请输入名字');
    return false;
  }
  if (!nameReg.test(name)) {
    Vue.prototype.$message.error('名字只含有汉字、字母、数字和下划线 不能以下划线开头和结尾');
    return false;
  }
  if (name.length > 16) {
    Vue.prototype.$message.error('名字太长');
    return false;
  }
  return true;
}

/**
 * 密码校验
 * @param password
 */
export function passwordVerify(password: string): boolean {
  const passwordReg = /^\w+$/gis;
  if (password.length === 0) {
    Vue.prototype.$message.error('请输入密码');
    return false;
  }
  if (!passwordReg.test(password)) {
    Vue.prototype.$message.error('密码只含有字母、数字和下划线');
    return false;
  }
  if (password.length > 16) {
    Vue.prototype.$message.error('密码最多16位,请重新输入');
    return false;
  }
  return true;
}

/**
 * 生成签名
 * @param obj
 */
export function makeSign(obj: object): string {
  let plain = '';
  Object.keys(obj)
    .filter((key: string) => key !== 'sign' && key !== 'EIO' && key !== 'transport' && key !== 't')
    .sort()
    .forEach((key: string) => {
      // @ts-ignore
      plain += `${key}${obj[key]}`;
    });

  plain += process.env.VUE_APP_KEY;
  console.log('before sha1: ', plain);
  console.log(plain);
  return sha1(plain).toUpperCase();
}

/**
 * 对象数组按照time属性排序
 * @param objArray 对象数组
 * @param reverse 是否降序
 */
export function sortByTime(objArray: Object[], reverse: boolean): Object[] {
  const newObjArray = objArray.sort((a: any, b: any) => {
    // time 越大 越后面
    if (a.time > b.time) {
      return reverse ? -1 : 1;
    }
    if (a.time < b.time) {
      return reverse ? 1 : -1;
    }
    return 0;
  });
  return newObjArray;
}

// 设定一个int类型的id, 范围在 [0, 511]之间
const id = 0;

// 使用id初始化
const generator = UUID(id);

/**
 * return a fake snowFlake ID
 */
export function newSnowFake(): string {
  return String(generator.uuid());
}

export function getFileExtension(fileName: string): string {
  return fileName.substring(fileName.lastIndexOf('.') + 1);
}

/**
 * 格式化输出文件大小
 * @param value
 */
export function renderSize(fileSize: number): string {
  const value = String(fileSize);
  if (value == null || value === '') {
    return '0 Bytes';
  }
  const unitArr = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
  let index = 0;
  const origin = parseFloat(value);
  index = Math.floor(Math.log(origin) / Math.log(1024));
  let size = origin / 1024 ** index;
  size = Number(size.toFixed(2)); // 保留的小数位数
  return size + unitArr[index];
}
