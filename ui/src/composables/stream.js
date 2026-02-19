import { Marked } from 'marked';
import { markedHighlight } from "marked-highlight";
import hljs from 'highlight.js';
import 'highlight.js/styles/atom-one-dark.css';
import 'vue3-toastify/dist/index.css';
import { useChatStore } from '@/stores/chat.js';
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
    const chatStore = useChatStore();
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
                                    messageId.value = m.messageId;
                                    const description = m.desc
                                    chatStore.addChatToHistory({ id: messageId.value, description: description })
                                    router.push({ name: "Existing Chat", params: { id: messageId.value } })
                                    break

                                case "response":
                                    data = data.data

                                    if (data.message.tool_calls) {
                                        assistantMessage.tool_calls = data.message.tool_calls;
                                        assistantMessage.loading = false; // Stop loading animation on message with tool calls
                                    }

                                    if (data.message.content) {
                                        let token = data.message.content;
                                        fullResponse += token;
                                        assistantMessage.content = marked.parse(fullResponse);
                                    }

                                    if (data.done) {
                                        if (assistantMessage.content) {
                                            apiMessages.push({ 'role': 'assistant', 'content': assistantMessage.content });
                                        }
                                        assistantMessage.loading = false;

                                        chatMessages.value.push({ role: 'info', content: assistantMessage.content, details: data });
                                        return;
                                    }
                                    break;

                                case "message":
                                    const toolResult = data.data;

                                    // Find the assistant message that contains the tool calls
                                    const assistantMessageWithTools = chatMessages.value.find(
                                        msg => msg.role === 'assistant' && msg.tool_calls && !msg.content
                                    );

                                    if (assistantMessageWithTools && assistantMessageWithTools.tool_calls) {
                                        const toolCall = assistantMessageWithTools.tool_calls.find(
                                            tc => tc.id === toolResult.tool_call_id
                                        );

                                        if (toolCall) {
                                            toolCall.result = toolResult.content;
                                        } else {
                                            // Fallback to prevent data loss if no match is found
                                            chatMessages.value.push(toolResult);
                                        }
                                    } else {
                                        chatMessages.value.push(toolResult);
                                    }
                                    
                                    // Ensure there's a placeholder for the next assistant response, but don't create duplicates.
                                    const lastMessage = chatMessages.value[chatMessages.value.length - 1];
                                    if (!lastMessage || lastMessage.role !== 'assistant' || lastMessage.content || lastMessage.tool_calls) {
                                        chatMessages.value.push({ role: 'assistant', content: '', loading: true });
                                    }
                                    break;
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
