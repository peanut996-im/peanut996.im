import { makeSign,sha1} from './tool.js'

var log = {};
log.Info = (msg)=>{
    console.log(new Date().toLocaleTimeString()+ ' '+ msg);
}

var token = window.localStorage.getItem('peanut996.im.token');
var uid = window.localStorage.getItem('peanut996.im.uid');

console.log(`current uid: ${uid}, token: ${token}`);
if( token == undefined || uid == undefined ){
    console.log('user not login');
    window.location.href = '/login';
}
let query = {
    token: token,
    uid: uid
}

query.sign = makeSign(query);

const socket = io('ws://localhost:9000/', {
    transport: ['websocket'],
    query: query
});

// socket.emit("getUserInfo",{
//     uid
// });

const form = document.getElementById('form');
const input = document.getElementById('input');
const addFriendBtn = document.getElementById('add-friend-button');
const addFriendInput = document.getElementById('add-friend-input');
const deleteFriendBtn = document.getElementById('delete-friend-button');
const deleteFriendInput = document.getElementById('delete-friend-input');
const createGroupBtn = document.getElementById('create-group-button');
const createGroupInput = document.getElementById('create-group-input');
const joinGroupBtn = document.getElementById('join-group-button');
const joinGroupInput = document.getElementById('join-group-input');
const leaveGroupBtn = document.getElementById('leave-group-button');
const leaveGroupInput = document.getElementById('leave-group-input');

form.addEventListener('submit', function (e) {
    e.preventDefault();
    if (input.value) {
        let from = uid;
        let to = "1396470649915969536";
        let type = "text"
        let content = input.value;
        let data = {
            from,to,type,content
        }
        socket.emit('chat', data);
        log.Info("/chat emit data:"+input.value);
        input.value = '';
    }
});



addFriendBtn.addEventListener('click', function (e) {
    e.preventDefault();
    let friendA = uid;
    if (addFriendInput.value) {
        let friendB = addFriendInput.value;
        let data = {
            friendA,friendB
        };
        socket.emit('addFriend',data);
        log.Info("/addFriend emit data:");
        console.log(data);

    }
});

deleteFriendBtn.addEventListener('click', function (e) {
    e.preventDefault();
    let friendA = uid;
    if (deleteFriendInput.value) {
        let friendB = deleteFriendInput.value;
        let data = {
            friendA,friendB
        };
        socket.emit('deleteFriend',data);
        log.Info("/deleteFriend emit data:");
        console.log(data);

    }
});

createGroupBtn.addEventListener('click', function (e) {
    e.preventDefault();
    let groupAdmin = uid;

    if (createGroupInput.value) {
        let groupName = createGroupInput.value;
        let data = {
            groupAdmin,groupName
        };
        socket.emit('createGroup',data);
        log.Info("/createGroup emit data:");
        console.log(data);

    }
});

joinGroupBtn.addEventListener('click', function (e) {
    e.preventDefault();

    if (joinGroupInput.value) {
        let groupID = joinGroupInput.value;
        let data = {
            groupID,uid
        };
        socket.emit('joinGroup',data);
        log.Info("/joinGroup emit data:");
        console.log(data);

    }
});

leaveGroupBtn.addEventListener('click', function (e) {
    e.preventDefault();

    if (leaveGroupInput.value) {
        let groupID = leaveGroupInput.value;
        let data = {
            groupID,uid
        };
        socket.emit('leaveGroup',data);
        log.Info("/leaveGroup emit data:");
        console.log(data);
    }
})


/**
 * Socket 事件监听
 */
socket.on('chat', function(data) {
    const msg = data.content;
    if (data.content!=undefined){
        log.Info("new chat message here")
        var item = document.createElement('li');
        console.log(new Date().toLocaleTimeString()+ ' /chat:'+ msg);
        item.textContent = msg;
        messages.appendChild(item);
        window.scrollTo(0, document.body.scrollHeight);
    }

});

socket.on('connect', () => {
    console.log(`socket.id: [${socket.id}], socket.remoteAddr: [${socket.io.engine.hostname}]`);
});

socket.on('console', msg => {
    console.log(new Date().toLocaleTimeString()+ ' /console:'+ msg);
})
socket.on('pingpong', msg => {
    console.log(new Date().toLocaleTimeString()+ ' /pingpong:'+ msg);
});

socket.on('userinfo', msg => {

});

socket.on('auth', msg => {
    log.Info('/auth: auth status '+ msg.message);
    console.log(msg);
});

socket.on('load', res => {
    log.Info("/load :get data: ");
    console.log(res);
});

export {
    socket,
    uid,
    token
}
