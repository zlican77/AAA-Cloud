<template>
  <div class="container">
    <h1>云上私密聊天室</h1>
    <div class="chat-area">
      <div class="user-list">
        <p>当前在线: <span>{{ userNum }}</span></p>
        <div id="user_list" ref="userList"></div>
      </div>
      <div id="msg_list" ref="msgList" class="message-list">
        <p v-for="(message, index) in messages" :key="index" ref="message"></p>
      </div>
    </div>
    <textarea id="msg_box" v-model="message" @keydown.enter.prevent="sendMessage"></textarea>
    <button @click="sendMessage">发送</button>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, onUpdated } from 'vue';

export default {
  setup() {
    const messages = ref([]);
    const msgList = ref(null);
    const message = ref('');

    const userNum = ref(0);
    const userList = ref('');
    const messageList = ref('');
    const uname = ref(prompt('请输入用户名', 'user' + uuid(8, 16)));
    let ws = null;

    onMounted(() => {
      ws = new WebSocket('ws://127.0.0.1:8080/ws');

      ws.onopen = () => {
        listMsg('系统消息：建立连接成功');
      };

      ws.onmessage = (e) => {
        const msg = JSON.parse(e.data);
        let sender;

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
        }

        const data = sender + msg.content;
        listMsg(data);
      };

      ws.onerror = () => {
        listMsg('系统消息 : 出错了,请退出重试.');
      };
    });

    onUpdated(() => {
      if (msgList.value) {
        msgList.value.scrollTop = msgList.value.scrollHeight;
      }
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
      messages.value.push(data);
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

<style scoped>
.container {
  width: 800px;
  height: 600px;
  margin: 30px auto;
  text-align: center;
}

.chat-area {
  width: 800px;
  border: 1px solid gray;
  height: 300px;
  display: flex;
}

.user-list {
  width: 200px;
  height: 300px;
  text-align: left;
  border-right: 1px solid gray;
  overflow: auto;
}

.message-list {
  width: 598px;
  border-left: 1px solid gray;
  height: 300px;
  overflow-y: auto;
}

textarea {
  width: 500px;
  height: 100px;
  margin-top: 10px;
}

button {
  margin-left: 10px;
}
</style>