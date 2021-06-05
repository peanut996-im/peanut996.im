/* eslint-disable no-unused-vars */

// EventName
const EventAuth = 'auth';
const EventLoad = 'load';
const EventAddFriend = 'addFriend';
const EventDeleteFriend = 'deleteFriend';
const EventCreateGroup = 'createGroup';
const EventJoinGroup = 'joinGroup';
const EventLeaveGroup = 'leaveGroup';
const EventChat = 'chat';
const EventGetUserInfo = 'getUserInfo';

// status code
const ErrorCodeOK = 0;
const ErrorSignInvalid = 1001;
const ErrorTokenInvalid = 1002;
const ErrorAuthFailed = 1003;
const ErrorHttpInnerError = 1004;
const ErrorHttpParamInvalid = 1005;
const ErrorHttpResourceExists = 1006;
const ErrorHttpResourceNotFound = 1007;

export {
  EventAuth,
  EventChat,
  EventLoad,
  EventAddFriend,
  EventDeleteFriend,
  EventCreateGroup,
  EventJoinGroup,
  EventLeaveGroup,
  EventGetUserInfo,
};

export {
  ErrorHttpInnerError,
  ErrorAuthFailed,
  ErrorHttpParamInvalid,
  ErrorHttpResourceExists,
  ErrorHttpResourceNotFound,
  ErrorSignInvalid,
  ErrorTokenInvalid,
  ErrorCodeOK,
};
