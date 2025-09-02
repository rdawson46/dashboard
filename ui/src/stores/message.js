import { defineStore } from "pinia";

export const useMessageStore = defineStore('mess', {
    state: () => ({
        messages: null,
        messageId: null,
    }),

    getters: {
        messages: (state) => state.messages ? state.messages : [],
        messagesId: (state) => state.messageId ? state.messageId : null,
    },

    actions: {
        addMessage: () => {},
        newMessage: () => {
            // set message to []
            // set message id to ID
        },
    }
})
