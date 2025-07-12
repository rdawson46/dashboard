<script module>
</script>
<script setup>
    import { ref, onMounted } from 'vue';
    import { Marked } from 'marked';
    import { markedHighlight } from "marked-highlight";
    import hljs from 'highlight.js';
    import 'highlight.js/styles/rose-pine-moon.css';
    import Sidebar from './components/sidebar.vue';

    // HACK: storing this in the front end 
    let messages = [];

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
                    console.log(line)
                    if (line.startsWith("data: ")) {
                        const jsonStr = line.substring(6)
                        if (jsonStr.trim()) {
                            try {
                                const data = JSON.parse(jsonStr)

                                if (data.done) {
                                    // TODO: add to UI

                                    console.log(data.eval_duration)
                                    messages.push({
                                        'role': 'assistant',
                                        'content': fullResponse,
                                    })

                                    return
                                }

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
        let chat = inputContainer.value.innerText.trim()
        if (chat === '') return;

        addMessage('user', chat);
        messages.push({
            'role': 'user',
            'content': chat,
        })
        
        inputContainer.value.innerText = '';

        let messElem = addMessage('assistant', '')
        await stream('/api/stream', {"messages": messages}, messElem);
    }

    onMounted(() => {
        // console.log(chatContainer.value)
        inputContainer.value.focus()
    })
</script>

<template>
    <Sidebar />
    <div id="wrapper">
        <h1>Query Machine</h1>
        <div id="chat-container" ref="chatContainer"></div>
        <div class="input-area">
            <div ref="inputContainer" id="message-input" contenteditable="true" @keydown.enter.exact.prevent="query"></div>
            <button id="send-btn" @click='query'>Send</button>
        </div>
    </div>
</template>

<style scoped>
/* TODO: add transitions */
#wrapper {
    height: 100vh;
    width: 100%;
    margin: 0 auto;
    padding: 20px;
    overflow-y: hidden;
}

#chat-container {
    height: 75%;
    width: 80%;
    border-radius: 5px;
    padding: 10px;
    overflow-y: auto;
    margin-bottom: 10px;
    scroll-behavior: smooth;
    margin-left: auto;
    margin-right: auto;
}

.input-area {
    display: flex;
    align-items: center;
    border: 1px solid #ccc;
    border-radius: 20px;
    padding: 5px;
    max-width: 60%;
    place-items: left;
    margin-left: auto;
    margin-right: auto;
}

#message-input {
    flex-grow: 1;
    padding: 10px;
    border: none;
    resize: none;
    font-family: inherit;
    font-size: 1rem;
    line-height: 1.5;
    min-height: 50px;
    max-height: 200px;
    overflow-y: auto;
    text-align: left;
}

#message-input:focus {
    outline: none;
}

button {
    padding: 10px 15px;
    margin: 0 10px;
    border: none;
    border-radius: 20px;
    cursor: pointer;
    background-color: #2b5b6f;
    color: white;
    transition: background-color 0.2s ease;
}

button:hover {
    background-color: #1e4258;
}

</style>

<style>

.message {
    padding: 8px 25px;
    margin-bottom: 10px;
    border-radius: 20px;
    width: fit-content;
    word-wrap: break-word;
    text-align: left;
    animation: fadeIn 0.3s ease-out;
}

.user {
    color: #303030;
    max-width: 85%;
    background-color: #cf6f6a;
    margin-left: auto;
    margin-right: 0;
}

.assistant {
    /*background-color: #2b5b6f;*/
    margin-left: auto;
    margin-right: auto;
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
