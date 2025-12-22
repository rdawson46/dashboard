<script setup>
defineProps({
  open: Boolean,
  rag: Boolean
});

defineEmits(['update:rag', 'trigger-upload', 'trigger-select-files']);
</script>

<template>
  <div v-if="open" class="options-menu glass-card">
    <button class="menu-button" @click="$emit('trigger-upload')">
      <i class="fa-solid fa-upload"></i>
      <span>Upload File</span>
    </button>

    <div class="rag-toggle">
      <span class="rag-label">
          <i class="fa-solid fa-circle-nodes"></i>
          RAG
      </span>
      <label class="switch">
        <input type="checkbox" :checked="rag" @change="$emit('update:rag', $event.target.checked)">
        <span class="slider round"></span>
      </label>
    </div>

    <button class="menu-button" @click="$emit('trigger-select-files')">
      <i class="fa-solid fa-file-circle-check"></i>
      <span>Select Files</span>
    </button>
  </div>
</template>

<style scoped>
.options-menu {
  position: absolute;
  bottom: calc(100% + 10px); /* Position above the button with a 10px gap */
  left: 0;
  width: 250px;
  padding: 1rem;
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  z-index: 10;
  border: 1px solid var(--border-color);
}

.menu-button {
  display: flex;
  align-items: center;
  gap: 1rem;
  background: transparent;
  color: var(--text-color);
  border: none;
  padding: 0.5rem;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.9rem;
  text-align: left;
}

.menu-button:hover:not(:disabled) {
  background: var(--bg-color-light);
}

.menu-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.rag-toggle {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: var(--text-color);
  font-size: 0.9rem;
  padding: 0.5rem;
  border-radius: 6px;
  background: transparent;
  transition: all 0.2s ease;
}

.rag-toggle:hover {
    background: var(--bg-color-light);
}

.rag-label {
  display: flex;
  align-items: center;
  gap: 1rem;
}

/* The switch - the box around the slider */
.switch {
  position: relative;
  display: inline-block;
  width: 50px;
  height: 28px;
}

/* Hide default HTML checkbox */
.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

/* The slider */
.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  -webkit-transition: .4s;
  transition: .4s;
}

.slider:before {
  position: absolute;
  content: "";
  height: 20px;
  width: 20px;
  left: 4px;
  bottom: 4px;
  background-color: white;
  -webkit-transition: .4s;
  transition: .4s;
}

input:checked + .slider {
  background-color: var(--primary-color);
}

input:checked + .slider:before {
  -webkit-transform: translateX(22px);
  -ms-transform: translateX(22px);
  transform: translateX(22px);
}

/* Rounded sliders */
.slider.round {
  border-radius: 28px;
}

.slider.round:before {
  border-radius: 50%;
}
</style>
