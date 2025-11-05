<script setup>
import { ref, watch, nextTick } from 'vue';

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
  if (message === '') return;

  emit('send-message', {
    content: message,
    webSearch: searchActive.value,
    code: codeActive.value,
  });

  inputValue.value = '';
  historyIndex.value = -1;
  pristineInput = '';

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
    <div class="input-area-main glass-card" :class="{ 'middle': isCentered }">
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
        <button id="send-btn" @click="sendMessage"><i class="fa-solid fa-paper-plane"></i></button>
      </div>
    </div>
</template>

<style scoped>
.input-area-main {
    position: absolute;
    bottom: 1rem;
    left: 20%;
    width: 60%;
    border-radius: 25px;
    padding: 5px 15px;
    transition: all 500ms ease-out;
    display: flex;
    flex-direction: column;
}

.input-area-main.middle {
  top: 50%;
  transform: translateY(-50%);
  bottom: auto;
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
