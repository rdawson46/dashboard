import { Marked } from 'marked';
import { markedHighlight } from "marked-highlight";
import hljs from 'highlight.js';
import 'highlight.js/styles/atom-one-dark.css';
import { toast } from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';

const marked = new Marked(
  markedHighlight({
    langPrefix: 'hljs language-',
    highlight(code, lang) {
      const language = hljs.getLanguage(lang) ? lang : 'plaintext';
      return hljs.highlight(code, { language }).value;
    }
  })
);

function notify(message) {
  toast(message, {
    autoClose: 1000,
    theme: 'dark',
  });
}

export async function useStream(
  body,
  messageId,
  apiMessages,
  chatMessages
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
      notify(`Error: ${response.status} ${response.statusText}`);
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
                  // TODO: swap to path param
                  const m = data.messageId;
                  messageId.value = m;
                  break
                
                case "response":
                  data = data.data

                  if (data.done) {
                    if (data.message.content.length) {
                      fullResponse += data.message.content
                      assistantMessage.content = marked.parse(fullResponse);
                    }
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
              notify('Error processing server response.');
            }
          }
        }
      }
    }
  } catch (error) {
    console.error('Error during fetch:', error);
    notify('Error when querying agent.');
    chatMessages.value.pop();
  } 
}
