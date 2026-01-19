<template>
  <div class="weapon-list-item" @click="copyCode">
    <!-- Left: Build Name + Tier -->
    <div class="item-left">
      <span class="build-name">{{ code.build }}</span>
      <span v-if="isValidTier(code.tier)" class="tier-tag-mini">{{ code.tier }}</span>
    </div>

    <!-- Middle: Code -->
    <div class="item-middle">
      <code class="code-text">{{ shortCode }}</code>
    </div>

    <!-- Right: Stats -->
    <div class="item-right">
      <span v-if="code.price !== null" class="stat-item stat-price">
        {{ code.price }}万
      </span>
      <span v-if="code.range !== null" class="stat-item stat-range">
        射程{{ code.range }}m
      </span>
      <span v-if="code.update_time" class="stat-item stat-time">
        更新{{ formatDate(code.update_time) }}
      </span>
    </div>

    <!-- Copy Button (hidden by default, shown on hover) -->
    <button class="copy-btn" @click.stop="copyCode">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useMessage } from 'naive-ui'
import type { WeaponCode } from '../stores/weaponStore'

const props = defineProps<{
  code: WeaponCode
}>()

const message = useMessage()

const shortCode = computed(() => {
  const parts = props.code.code.split('-')
  return parts[parts.length - 1] || props.code.code
})

const isValidTier = (tier: string) => {
  return tier && tier !== '-' && tier !== 'null' && tier !== ''
}

const formatDate = (dateStr: string) => {
  // Format like "2024.01" or "2024-01" to "24年1月"
  const match = dateStr.match(/(\d{4})[-.](\d{1,2})/)
  if (match) {
    const year = match[1].slice(-2) // Get last 2 digits
    const month = match[2].replace(/^0/, '') // Remove leading zero
    return `${year}年${month}月`
  }
  return dateStr
}

const copyCode = async () => {
  try {
    await navigator.clipboard.writeText(shortCode.value)
    message.success('已复制到剪贴板')
  } catch {
    message.error('复制失败，请手动复制')
  }
}
</script>

<style scoped>
.weapon-list-item {
  display: grid;
  grid-template-columns: 2fr 2fr 1.5fr auto;
  gap: 1rem;
  align-items: center;
  padding: 0.625rem 0.875rem;
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.weapon-list-item:hover {
  background: #F1F5F9;
  border-color: #2563EB;
}

.item-left {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  min-width: 0;
}

.build-name {
  font-weight: 500;
  font-size: 0.8125rem;
  color: #1E293B;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tier-tag-mini {
  padding: 0.125rem 0.375rem;
  font-size: 0.65rem;
  font-weight: 600;
  background: #FEF3C7;
  color: #D97706;
  border: 1px solid #FDE68A;
  border-radius: 0.25rem;
  white-space: nowrap;
  flex-shrink: 0;
}

.item-middle {
  flex-shrink: 0;
}

.code-text {
  font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
  font-size: 0.75rem;
  color: #475569;
  background: #E2E8F0;
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  white-space: nowrap;
}

.item-right {
  display: flex;
  align-items: center;
  gap: 0.625rem;
  flex-shrink: 0;
}

.stat-item {
  font-size: 0.7rem;
  font-weight: 500;
  white-space: nowrap;
}

.stat-price {
  color: #DC2626;
}

.stat-range {
  color: #2563EB;
}

.stat-time {
  color: #94A3B8;
  font-size: 0.65rem;
}

.copy-btn {
  padding: 0.375rem;
  background: transparent;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  color: #94A3B8;
  opacity: 0;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.weapon-list-item:hover .copy-btn {
  opacity: 1;
}

.copy-btn:hover {
  background: #EFF6FF;
  color: #2563EB;
}
</style>
