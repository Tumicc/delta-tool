<template>
  <div class="weapon-group">
    <!-- Weapon Header -->
    <div class="group-header" @click="isExpanded = !isExpanded">
      <div class="weapon-info">
        <div class="weapon-title-row">
          <h3 class="weapon-name">{{ weaponName }}</h3>
          <!-- Tier Tag -->
          <span
            v-if="tier && tier !== '-'"
            class="tier-tag"
          >
            {{ tier }}
          </span>
        </div>
        <span class="code-count">{{ codes.length }} 个改装方案</span>
      </div>
      <svg
        :class="['chevron', { expanded: isExpanded }]"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
      </svg>
    </div>

    <!-- Codes List -->
    <transition
      name="expand"
      @enter="onEnter"
      @after-enter="onAfterEnter"
      @leave="onLeave"
    >
      <div v-show="isExpanded" class="codes-list">
        <WeaponCard
          v-for="code in codes"
          :key="code.id"
          :code="code"
        />
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import WeaponCard from './WeaponCard.vue'
import type { WeaponCode } from '../stores/weaponStore'

const props = defineProps<{
  weaponName: string
  codes: WeaponCode[]
}>()

const isExpanded = ref(true)

// Get tier from the first code (all codes for the same weapon should have the same tier)
const tier = computed(() => {
  return props.codes[0]?.tier || null
})

// Animation helpers
const onEnter = (el: Element) => {
  const element = el as HTMLElement
  element.style.height = '0'
  element.style.overflow = 'hidden'
}

const onAfterEnter = (el: Element) => {
  const element = el as HTMLElement
  element.style.height = 'auto'
  element.style.overflow = 'visible'
}

const onLeave = (el: Element) => {
  const element = el as HTMLElement
  const height = element.scrollHeight
  element.style.height = height + 'px'
  element.style.overflow = 'hidden'

  // Force reflow
  element.offsetHeight

  element.style.height = '0'
}
</script>

<style scoped>
.weapon-group {
  margin-bottom: 1rem;
}

.group-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem;
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
  user-select: none;
}

.group-header:hover {
  border-color: #2563EB;
}

.weapon-info {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.weapon-title-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.weapon-name {
  font-family: 'Space Grotesk', sans-serif;
  font-size: 1rem;
  font-weight: 600;
  color: #1E293B;
  margin: 0;
}

.tier-tag {
  padding: 0.125rem 0.375rem;
  font-size: 0.65rem;
  font-weight: 600;
  background: #FEF3C7;
  color: #D97706;
  border: 1px solid #FDE68A;
  border-radius: 0.25rem;
}

.code-count {
  font-size: 0.7rem;
  color: #94A3B8;
}

.chevron {
  width: 1rem;
  height: 1rem;
  color: #CBD5E1;
  transition: transform 0.3s ease;
  flex-shrink: 0;
}

.chevron.expanded {
  transform: rotate(180deg);
}

.codes-list {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
  margin-top: 0.75rem;
  padding: 0 0.5rem;
}

/* Expand transition */
.expand-enter-active,
.expand-leave-active {
  transition: all 0.3s ease;
  overflow: hidden;
}

.expand-enter-from,
.expand-leave-to {
  height: 0 !important;
  opacity: 0;
}
</style>
