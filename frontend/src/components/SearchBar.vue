<template>
  <div class="relative">
    <div class="relative flex items-center">
      <!-- Search Icon -->
      <svg
        class="absolute left-4 w-5 h-5 text-slate-400 pointer-events-none"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
        />
      </svg>

      <!-- Input -->
      <input
        :value="modelValue"
        @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
        type="text"
        placeholder="搜索武器名称、改装类型、改枪码..."
        class="search-input"
      />

      <!-- Clear Button -->
      <button
        v-if="modelValue"
        @click="$emit('update:modelValue', '')"
        class="absolute right-4 p-1 text-slate-400 hover:text-slate-200 transition-colors duration-200"
        aria-label="Clear search"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: string
}>()

defineEmits<{
  'update:modelValue': [value: string]
}>()
</script>

<style scoped>
.search-input {
  width: 100%;
  padding-left: 3rem;
  padding-right: 3rem;
  padding-top: 0.75rem;
  padding-bottom: 0.75rem;
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  border-radius: 0.75rem;
  color: #1E293B;
  font-size: 0.875rem;
  transition: all 0.2s;
}

.search-input::placeholder {
  color: #94A3B8;
}

.search-input:focus {
  outline: none;
  border-color: #2563EB;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.search-input:hover {
  border-color: #CBD5E1;
}
</style>
