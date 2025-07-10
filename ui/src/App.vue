<script module>
</script>
<script setup>
    import { ref, onMounted } from 'vue';
    import { Marked } from 'marked';
    import { markedHighlight } from "marked-highlight";
    import hljs from 'highlight.js';
    import 'highlight.js/styles/rose-pine-moon.css';

    /*
     TODO:

     * add a UI blocker to prevent overload
     * add animations
     * add a loader/spinner somewhere

    */

    const chatContainer = ref(null)
    const inputContainer = ref(null)

    const marked = new Marked(
        markedHighlight({
            langPrefix: 'hljs language-',
            highlight(code, lang) {
                const language = hljs.getLanguage(lang) ? lang : 'plaintext';
                return hljs.highlight(code, { language }).value;
            }
        })
    );

    async function stream(url, body, message) {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(body)
        });

        const reader = response.body.getReader();
        const decoder = new TextDecoder();

        let buffer = '';
        let fullResponse = '';

        try {
            while (true) {
                const { done, value } = await reader.read();

                if (done) break;

                buffer += decoder.decode(value, { stream: true })

                const lines = buffer.split('\n')
                buffer = lines.pop() || '';

                for (const line of lines) {
                    if (line.startsWith("data: ")) {
                        const jsonStr = line.substring(6)
                        if (jsonStr.trim()) {
                            try {
                                const data = JSON.parse(jsonStr)

                                if (data.done) return

                                let token = data.message.content
                                fullResponse += token;
                                message.innerHTML = marked.parse(fullResponse);
                                chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
                            } catch (e) {
                                return
                            }
                        }
                    }
                }
            }
        } catch (error) {
            console.log(error)
        } finally {
            reader.releaseLock();
        }
    }

    function addMessage(role, content) {
        const messageElement = document.createElement('div');
        messageElement.classList.add(role, 'message');

        messageElement.innerHTML = content;
        chatContainer.value.appendChild(messageElement);

        chatContainer.value.scrollTop = chatContainer.value.scrollHeight;

        return messageElement
    }

    async function query() {
        let chat = inputContainer.value.value.trim()
        if (chat === '') return;

        addMessage('user', chat);
        inputContainer.value.value = '';

        let messElem = addMessage('assistant', '')
        await stream('/api/stream', {"query": chat}, messElem);
    }

    onMounted(() => {
        // console.log(chatContainer.value)
    })
</script>

<template>
    <div id="wrapper">
        <h3>Query Machine</h3>
        <div id="chat-container" ref="chatContainer"></div>
        <div class="input-area">
            <input ref="inputContainer" type="text" id="message-input" placeholder="Type your message..." @keyup.enter="query">

            <!--<button id="send-btn" @click='stream("/api/stream", {"query": "why is the sky blue"})'>Send</button>-->
            <button id="send-btn" @click='query'>Send</button>
            <button id="clear-btn">Clear Chat</button>
        </div>
    </div>
</template>

<style scoped>
/* TODO: add transitions */
#wrapper {
    height: 100vh;
    max-width: 100%;
    margin: 0 auto;
    padding: 20px;
    overflow-y: hidden;
}

#chat-container {
    height: 75%;
    width: 90%;
    border-radius: 5px;
    padding: 10px;
    overflow-y: auto;
    margin-bottom: 10px;
    scroll-behavior: smooth;
}

#message-input {
    flex-grow: 1;
    padding: 10px;
    border: 1px solid #000000;
    border-radius: 5px;
    margin-right: 10px;
    transition: border-color 0.2s ease;
}

button {
    padding: 10px 15px;
    margin: 0 10px;
    border: none;
    border-radius: 5px;
    cursor: pointer;
}

</style>

<style>

.message {
    padding: 8px 25px;
    margin-bottom: 10px;
    border-radius: 20px;
    max-width: 85%;
    width: fit-content;
    word-wrap: break-word;
    text-align: left;
    animation: fadeIn 0.3s ease-out;
}

.user {
    color: #303030;
    background-color: #cf6f6a;
    margin-left: auto;
    margin-right: 0;
}

.assistant {
    /*background-color: #2b5b6f;*/
}

@keyframes fadeIn {
    from {
        opacity: 0;
    }
    to {
        opacity: 1;
    }
}

code {
    border-radius: 15px;
}

</style>
