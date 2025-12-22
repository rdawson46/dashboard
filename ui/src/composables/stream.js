import { Marked } from 'marked';
import { markedHighlight } from "marked-highlight";
import hljs from 'highlight.js';
import 'highlight.js/styles/atom-one-dark.css';
import 'vue3-toastify/dist/index.css';
import { useNotify } from '@/composables/notify.js'


const marked = new Marked(
    markedHighlight({
        langPrefix: 'hljs language-',
        highlight(code, lang) {
            const language = hljs.getLanguage(lang) ? lang : 'plaintext';
            return hljs.highlight(code, { language }).value;
        }
    })
);

export async function useStream(
    body,
    messageId,
    apiMessages,
    chatMessages,
    router
) {
    const url = '/api/stream'

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(body),
            credentials: 'include'
        });

        if (!response.ok) {
            if (response.status === 401) {
                await router.replace('/login')
                return
            }
            useNotify(`Error: ${response.status} ${response.statusText}`);
            chatMessages.value.pop();
            return;
        }

        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        let buffer = '';
        let fullResponse = '';


        while (true) {
            const { done, value } = await reader.read();
            if (done) break;

            buffer += decoder.decode(value, { stream: true });
            const lines = buffer.split('\n');
            buffer = lines.pop() || '';

            for (const line of lines) {
                let assistantMessage = chatMessages.value[chatMessages.value.length - 1];

                if (line.startsWith("data: ")) {
                    const jsonStr = line.substring(6);
                    if (jsonStr.trim()) {
                        try {
                            let data = JSON.parse(jsonStr);

                            switch (data.type) {
                                case "Message ID":
                                    if (messageId.value) break;
                                    const m = data.data;
                                    messageId.value = m;
                                    router.push({ name: "Existing Chat", params: { id: messageId.value } })
                                    break

                                case "response":
                                    data = data.data

                                    if (data.done) {
                                        if (data.message.content.length) {
                                            fullResponse += data.message.content
                                            assistantMessage.content = marked.parse(fullResponse);
                                        }
                                        assistantMessage.loading = false;
                                        apiMessages.push({ 'role': 'assistant', 'content': fullResponse });

                                        chatMessages.value.push({ role: 'info', content: fullResponse, details: data });
                                        return;
                                    }

                                    if (data.message.tool_calls) {
                                        console.log(data)
                                        for (const toolCall of data.message.tool_calls) {
                                            let { name } = toolCall;
                                            console.log(toolCall)
                                        }
                                        assistantMessage.tool_calls = data.message.tool_calls
                                        continue
                                    }

                                    let token = data.message.content;
                                    fullResponse += token;
                                    assistantMessage.content = marked.parse(fullResponse);
                                    break;

                                case "message":
                                    data = data.data

                                    // will need to append this message to the chat
                                    chatMessages.value.push(data);
                                    chatMessages.value.push({ role: 'assistant', content: '' });
                                    break
                                default:
                                    console.log(data)
                                    break
                            }

                        } catch (e) {
                            console.error('Error parsing JSON:', e);
                            useNotify('Error processing server response.');
                        }
                    }
                }
            }
        }
    } catch (error) {
        console.error('Error during fetch:', error);
        useNotify('Error when querying agent.');
        chatMessages.value.pop();
    } 
}
