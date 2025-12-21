<script setup>
import { ref, watch, nextTick } from 'vue';
import { useNotify } from '@/composables/notify.js';

const props = defineProps({
  userMessageHistory: {
    type: Array,
    default: () => []
  },
  isCentered: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['send-message']);

const searchActive = ref(false);
const codeActive = ref(false);
const isUploading = ref(false);
const fileInput = ref(null);
const uploadedFiles = ref([]);

const inputContainer = ref(null);
const inputValue = ref('');
const historyIndex = ref(-1);
let pristineInput = '';

const handleInput = (e) => {
    inputValue.value = e.target.innerText;
    if (historyIndex.value !== -1) {
        historyIndex.value = -1;
        pristineInput = '';
    }
};

const triggerFileUpload = () => {
  fileInput.value.click();
};

const uploadFile = async (event) => {
    const files = event.target.files;
    if (!files.length) return;

    isUploading.value = true;

    const uploadPromises = Array.from(files).map(file => {
        const formData = new FormData();
        formData.append('file', file);
        return fetch('/api/uploadFile', {
            method: 'POST',
            body: formData,
            credentials: 'include'
        }).then(async response => {
            if (response.ok) {
                const result = await response.json();
                return { id: result.fileID, name: file.name };
            } else {
                const error = await response.text();
                throw new Error(`Failed to upload ${file.name}: ${error}`);
            }
        });
    });

    try {
        const results = await Promise.all(uploadPromises);
        uploadedFiles.value.push(...results);
        useNotify(`${results.length} file(s) uploaded successfully!`);
    } catch (error) {
        useNotify(`Error uploading file(s): ${error.message}`);
    } finally {
        isUploading.value = false;
        if (fileInput.value) {
            fileInput.value.value = '';
        }
    }
};

const removeUploadedFile = (fileIdToRemove) => {
    uploadedFiles.value = uploadedFiles.value.filter(file => file.id !== fileIdToRemove);
};

function navigateHistoryUp() {
    if (!props.userMessageHistory.length) return;

    if (historyIndex.value === -1) {
        pristineInput = inputValue.value;
        historyIndex.value = props.userMessageHistory.length - 1;
    } else if (historyIndex.value > 0) {
        historyIndex.value--;
    } else {
        return;
    }
    inputValue.value = props.userMessageHistory[historyIndex.value].content;
}

function navigateHistoryDown() {
    if (historyIndex.value === -1) return;

    if (historyIndex.value < props.userMessageHistory.length - 1) {
        historyIndex.value++;
        inputValue.value = props.userMessageHistory[historyIndex.value].content;
    } else {
        historyIndex.value = -1;
        inputValue.value = pristineInput;
    }
}

watch(inputValue, (newValue) => {
    if (inputContainer.value && newValue !== inputContainer.value.innerText) {
        inputContainer.value.innerText = newValue;
        
        nextTick(() => {
            if (document.activeElement === inputContainer.value) {
                const selection = window.getSelection();

                if (selection) {
                    const range = document.createRange();
                    range.selectNodeContents(inputContainer.value);
                    range.collapse(false)
                    selection.removeAllRanges()
                    selection.addRange(range)
                }
            }
        });
    }
});

function sendMessage() {
  const message = inputValue.value.trim();
  if (message === '' && uploadedFiles.value.length === 0) return;

  const messagePayload = {
    content: message,
    webSearch: searchActive.value,
    code: codeActive.value,
  };

  if (uploadedFiles.value.length > 0) {
    messagePayload.fileIds = uploadedFiles.value.map(file => file.id);
  }

  emit('send-message', messagePayload);

  inputValue.value = '';
  historyIndex.value = -1;
  pristineInput = '';
  uploadedFiles.value = [];

  nextTick(() => {
    if(inputContainer.value) {
      inputContainer.value.focus();
    }
  });
}

defineExpose({
  focus: () => {
    inputContainer.value?.focus();
  }
});

</script>

<template>
    <div class="chat-input-container" :class="{ 'middle': isCentered }">
        <div class="file-pills-container">
            <div v-for="file in uploadedFiles" :key="file.id" class="file-pill">
                <span><i class="fa-solid fa-paperclip"></i> {{ file.name }}</span>
                <button @click="removeUploadedFile(file.id)" class="close-btn">&times;</button>
            </div>
        </div>
        <div class="input-area-main glass-card">
            <div class="input-area-sub">
                <div
                ref="inputContainer"
                id="message-input"
                placeholder="Type your message..."
                contenteditable="true"
                @input="handleInput"
                @keydown.down.exact.prevent="navigateHistoryDown()"
                @keydown.up.exact.prevent="navigateHistoryUp()"
                @keydown.enter.exact.prevent="sendMessage"
                ></div>
            </div>

            <div class="input-area-buttons">
                <button @click='searchActive = !searchActive' :class="{ active: searchActive }">
                <i class="fa-solid fa-globe"></i>
                <span>Web Search</span>
                </button>
                <button @click='codeActive = !codeActive' :class="{ active: codeActive }">
                <i class="fa-solid fa-code"></i>
                <span>Code</span>
                </button>
                <button @click="triggerFileUpload" :disabled="isUploading">
                <template v-if="isUploading">
                    <i class="fa-solid fa-spinner fa-spin"></i>
                    <span>Uploading...</span>
                </template>
                <template v-else>
                    <i class="fa-solid fa-file-arrow-up"></i>
                    <span>File Upload</span>
                </template>
                </button>
                <input type="file" ref="fileInput" @change="uploadFile" style="display: none" multiple />
                <button id="send-btn" @click="sendMessage"><i class="fa-solid fa-paper-plane"></i></button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.chat-input-container {
    position: absolute;
    bottom: 1rem;
    left: 20%;
    width: 60%;
    display: flex;
    flex-direction: column;
    transition: all 500ms ease-out;
}

.chat-input-container.middle {
  top: 50%;
  transform: translateY(-50%);
  bottom: auto;
}

.input-area-main {
    border-radius: 25px;
    padding: 5px 15px;
    display: flex;
    flex-direction: column;
    width: 100%;
}

.file-pills-container {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
    max-width: 100%;
    justify-content: flex-start;
}

.file-pill {
    display: flex;
    align-items: center;
    background: var(--bg-color-light);
    border-radius: 20px;
    padding: 0.2rem 0.6rem;
    font-size: 0.8rem;
    color: var(--text-color);
    border: 1px solid var(--border-color);
    max-width: 200px;
}

.file-pill span {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    margin-right: 0.5rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

.file-pill .close-btn {
    background: none;
    border: none;
    color: var(--text-color);
    cursor: pointer;
    font-size: 1rem;
    padding: 0;
    line-height: 1;
}

.file-pill .close-btn:hover {
    color: var(--primary-color);
}

.input-area-sub {
    display: flex;
    align-items: center;
}

.active {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.input-area-buttons {
  display: flex;
  gap: 1rem;
  margin: 0.25rem 0;
}

.input-area-buttons button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: var(--bg-color-light);
  color: var(--text-color);
  border: 1px solid var(--border-color);
  padding: 0.5rem 1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.input-area-buttons button.active {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.input-area-buttons button:last-child {
    margin-left: auto;
}

#message-input {
    flex-grow: 1;
    padding: 5px 5px 0 5px;
    border: none;
    resize: none;
    font-family: 'Sqgoe UI', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    font-size: 1rem;
    line-height: 1.5;
    max-height: 200px;
    min-height: 40px;
    color: #f0f0f0;
    background-color: transparent;
    text-align: left;
}

#message-input:focus {
    outline: none;
}

#send-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    padding: 10px;
    transition: color 0.2s ease;
}

#send-btn:hover {
    color: #81b2f3;
}
</style>
