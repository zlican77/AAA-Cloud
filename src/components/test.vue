<template>
  
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue';
export default {
    setup() {
    const userNum = ref(0);
    const userList = ref('');
    const messageList = ref('');
    const message = ref('');
    const uname = ref(prompt('请输入用户名', 'user' + uuid(8, 16)));
    let ws = null;

    onMounted(() => {
      ws = new WebSocket('ws://127.0.0.1:8080/ws');

      ws.onopen = () => {
        listMsg('系统消息：建立连接成功');
      };

      ws.onmessage = (e) => {
        const msg = JSON.parse(e.data);
        let sender, user_name, name_list, change_type;

        switch (msg.type) {
          case 'system':
            sender = '系统消息: ';
            break;
          case 'user':
            sender = msg.from + ': ';
            break;
          case 'handshake':
            sendMsg({ type: 'login', content: uname.value });
            return;
          case 'login':
          case 'logout':
            user_name = msg.content;
            name_list = msg.user_list;
            change_type = msg.type;
            dealUser(user_name, change_type, name_list);
            return;
        }

        const data = sender + msg.content;
        listMsg(data);
      };

      ws.onerror = () => {
        listMsg('系统消息 : 出错了,请退出重试.');
      };
    });

    onUnmounted(() => {
      if (ws) {
        ws.close();
      }
    });

    const sendMessage = () => {
      const content = message.value.trim();
      if (content) {
        const msg = { content, type: 'user' };
        sendMsg(msg);
        message.value = '';
      }
    };

    const listMsg = (data) => {
      messageList.value += `<p>${data}</p>`;
    };

    const dealUser = (user_name, type, name_list) => {
      userList.value = '';
      name_list.forEach((name) => {
        userList.value += `<p>${name}</p>`;
      });
      userNum.value = name_list.length;
      const change = type === 'login' ? '上线' : '下线';
      listMsg(`系统消息: ${user_name} 已${change}`);
    };

    const sendMsg = (msg) => {
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify(msg));
      } else {
        listMsg('系统消息: 连接未建立，无法发送消息');
      }
    };

    const uuid = (len, radix) => {
      const chars = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'.split('');
      const uuid = [];
      radix = radix || chars.length;
      if (len) {
        for (let i = 0, r; i < len; i++) {
          r = 0 | Math.random() * radix;
          uuid[i] = chars[r];
        }
      } else {
        let r;
        uuid[8] = uuid[13] = uuid[18] = uuid[23] = '-';
        uuid[14] = '4';
        for (let i = 0; i < 36; i++) {
          if (!uuid[i]) {
            r = 0 | Math.random() * 16;
            uuid[i] = chars[(i === 19) ? (r & 0x3) | 0x8 : r];
          }
        }
      }
      return uuid.join('');
    };

    return {
      userNum,
      userList,
      messageList,
      message,
      sendMessage,
    };
  },
};
</script>

<style>

</style>