import OSS from 'ali-oss';

const OSSClient = new OSS({
  region: process.env.VUE_APP_OSS_REGION,
  accessKeyId: process.env.VUE_APP_OSS_ACCESS_KEY_ID!,
  accessKeySecret: process.env.VUE_APP_OSS_ACCESS_KEY_SECRET!,
  bucket: process.env.VUE_APP_OSS_BUCKET,
});

export default OSSClient;
