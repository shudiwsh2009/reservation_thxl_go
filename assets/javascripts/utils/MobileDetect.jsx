import MobileDetect from 'mobile-detect';

const md = new MobileDetect(window.navigator.userAgent);

export default {
  isWechat() {
    return !!md.match('MicroMessenger');
  },
};
