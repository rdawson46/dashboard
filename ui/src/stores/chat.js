import { defineStore } from 'pinia';

export const useChatStore = defineStore('chat', {
  state: () => ({
    history: [],
  }),
  actions: {
    addChatToHistory(chat) {
      this.history.unshift(chat);
    },
    setHistory(history) {
      this.history = history;
    }
  },
});
