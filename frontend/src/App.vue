<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useWeaponStore } from './stores/weaponStore'
import SearchBar from './components/SearchBar.vue'
import WeaponGroup from './components/WeaponGroup.vue'
import { NMessageProvider } from 'naive-ui'

const weaponStore = useWeaponStore()
const {
  loading,
  searchQuery,
  selectedMode,
  selectedWeaponType,
  selectedDataSource,
  filteredCodes,
  groupedCodes,
} = storeToRefs(weaponStore)

// Watch for data source changes and reload codes
watch(selectedDataSource, () => {
  weaponStore.loadCodes()
})

onMounted(() => {
  const checkWailsReady = () => {
    if ((window as any).go && (window as any).go.app && (window as any).go.app.App) {
      weaponStore.loadCodes()
    } else {
      setTimeout(checkWailsReady, 100)
    }
  }
  checkWailsReady()
})

const weaponTypes = [
  { label: '全部枪械', value: 'all' },
  { label: '突击步枪', value: '突击步枪' },
  { label: '冲锋枪', value: '冲锋枪' },
  { label: '射手步枪', value: '射手步枪' },
  { label: '手枪', value: '手枪' },
  { label: '栓狙', value: '狙击' },
  { label: '机枪', value: '机枪' },
  { label: '霰弹枪', value: '霰弹枪' },
  { label: '弓弩', value: '弓弩' },
]

const dataSources = [
  { label: '刀仔', value: '刀仔' },
  { label: '武器大师', value: '武器大师' },
]
</script>

<template>
  <n-message-provider>
    <div class="app-container">
    <!-- Layout Container -->
    <div class="layout">
      <!-- Sidebar -->
      <aside class="sidebar">
        <!-- Logo Section -->
        <div class="logo-section">
          <h1 class="logo-title">三角洲改枪码工具</h1>
          <p class="logo-subtitle">Weapon Codes Database</p>
        </div>

        <!-- Scrollable Content -->
        <div class="sidebar-content">
          <!-- Search Section -->
          <section class="sidebar-section">
            <SearchBar v-model="searchQuery" />
          </section>

          <!-- Game Mode Filter -->
          <section class="sidebar-section">
            <h2 class="section-title">游戏模式</h2>
            <div class="mode-list">
              <button
                v-for="mode in ['烽火地带', '全面战场']"
                :key="mode"
                @click="selectedMode = mode"
                :class="['mode-button', { active: selectedMode === mode }]"
              >
                <span class="mode-label">{{ mode }}</span>
              </button>
            </div>
          </section>

          <!-- Weapon Type Filter -->
          <section class="sidebar-section">
            <h2 class="section-title">枪械类型</h2>
            <div class="type-grid">
              <button
                v-for="type in weaponTypes"
                :key="type.value"
                @click="selectedWeaponType = type.value"
                :class="['type-button', { active: selectedWeaponType === type.value }]"
              >
                {{ type.label }}
              </button>
            </div>
          </section>
        </div>

        <!-- Footer Stats -->
        <div class="sidebar-footer">
          <div class="stats-text">
            <span class="stats-label">共</span>
            <span class="stats-number">{{ filteredCodes.length }}</span>
            <span class="stats-label">条改枪码</span>
          </div>
        </div>
      </aside>

      <!-- Main Content -->
      <main class="main-content">
        <!-- Header with Data Source Tabs -->
        <header class="content-header">
          <div class="header-tabs">
            <button
              v-for="source in dataSources"
              :key="source.value"
              @click="selectedDataSource = source.value"
              :class="['tab-button', { active: selectedDataSource === source.value }]"
            >
              {{ source.label }}
            </button>
          </div>
          <div class="header-info">
            <h2 class="breadcrumb-item">
              {{ selectedMode }}
            </h2>
            <span class="breadcrumb-separator">·</span>
            <h2 class="breadcrumb-item">
              {{ selectedWeaponType === 'all' ? '全部枪械' : weaponTypes.find(t => t.value === selectedWeaponType)?.label }}
            </h2>
          </div>
        </header>

        <!-- Content Area -->
        <div class="content-area">
          <!-- Loading State -->
          <div
            v-if="loading"
            class="state-container"
          >
            <div class="loading-spinner">
              <div class="spinner-ring"></div>
              <div class="spinner-dot"></div>
            </div>
            <p class="state-text">加载中...</p>
          </div>

          <!-- Empty State -->
          <div
            v-else-if="filteredCodes.length === 0"
            class="state-container"
          >
            <div class="empty-icon">
              <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <h3 class="state-title">暂无数据</h3>
            <p class="state-desc">未找到匹配的改枪码</p>
            <button
              @click="weaponStore.loadCodes()"
              class="reload-button"
            >
              重新加载
            </button>
          </div>

          <!-- Weapon Groups -->
          <div
            v-else
            class="weapons-container"
          >
            <WeaponGroup
              v-for="[weaponName, codes] in groupedCodes"
              :key="weaponName"
              :weapon-name="weaponName"
              :codes="codes"
            />
          </div>
        </div>
      </main>
    </div>
    </div>
  </n-message-provider>
</template>

<style scoped>
.app-container {
  min-height: 100vh;
  background: #F8FAFC;
}

.layout {
  display: flex;
  height: 100vh;
}

/* Sidebar Styles */
.sidebar {
  width: 18rem;
  flex-shrink: 0;
  background: #FFFFFF;
  border-right: 1px solid #E2E8F0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.logo-section {
  padding: 1.5rem;
  border-bottom: 1px solid #E2E8F0;
}

.logo-title {
  font-family: 'Space Grotesk', sans-serif;
  font-size: 1.5rem;
  font-weight: 700;
  background: linear-gradient(to right, #2563EB, #7C3AED);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
}

.logo-subtitle {
  font-size: 0.875rem;
  color: #64748B;
  margin: 0.25rem 0 0 0;
}

.sidebar-content {
  flex: 1;
  overflow-y: auto;
  padding: 1rem 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.sidebar-section {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.section-title {
  font-size: 0.75rem;
  font-weight: 600;
  color: #64748B;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin: 0;
}

/* Mode Buttons */
.mode-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.mode-button {
  width: 100%;
  padding: 0.625rem 1rem;
  background: #F1F5F9;
  color: #475569;
  border: 1px solid transparent;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 500;
  font-size: 0.875rem;
}

.mode-button:hover {
  background: #E2E8F0;
  color: #1E293B;
}

.mode-button.active {
  background: #2563EB;
  color: white;
  border-color: #2563EB;
}

.mode-label {
  flex: 1;
  text-align: left;
}

/* Type Grid */
.type-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.5rem;
}

.type-button {
  padding: 0.5rem 0.75rem;
  background: #F1F5F9;
  color: #475569;
  border: 1px solid transparent;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
  font-weight: 500;
  font-size: 0.875rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.type-button:hover {
  background: #E2E8F0;
  color: #1E293B;
}

.type-button.active {
  background: #EFF6FF;
  color: #2563EB;
  border-color: #2563EB;
}

/* Sidebar Footer */
.sidebar-footer {
  padding: 1rem;
  border-top: 1px solid #E2E8F0;
  background: #F8FAFC;
}

.stats-text {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.25rem;
  font-size: 0.875rem;
}

.stats-label {
  color: #64748B;
}

.stats-number {
  font-family: 'Space Grotesk', sans-serif;
  font-weight: 700;
  color: #2563EB;
}

/* Main Content */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.content-header {
  height: 5rem;
  border-bottom: 1px solid #E2E8F0;
  background: #FFFFFF;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 2rem;
  gap: 2rem;
}

.header-tabs {
  display: flex;
  gap: 0.5rem;
  background: #F1F5F9;
  padding: 0.25rem;
  border-radius: 0.5rem;
}

.tab-button {
  padding: 0.5rem 1.25rem;
  background: transparent;
  color: #64748B;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: all 0.2s;
  font-weight: 500;
  font-size: 0.875rem;
  white-space: nowrap;
}

.tab-button:hover {
  color: #475569;
  background: rgba(255, 255, 255, 0.5);
}

.tab-button.active {
  background: #FFFFFF;
  color: #2563EB;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.header-info {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex: 1;
  justify-content: flex-end;
}

.header-breadcrumb {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.breadcrumb-item {
  font-family: 'Space Grotesk', sans-serif;
  font-size: 1.125rem;
  font-weight: 600;
  color: #1E293B;
  margin: 0;
}

.breadcrumb-separator {
  color: #CBD5E1;
}

/* Content Area */
.content-area {
  flex: 1;
  overflow-y: auto;
  padding: 2rem;
}

/* States */
.state-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 1rem;
}

.loading-spinner {
  position: relative;
  width: 3rem;
  height: 3rem;
}

.spinner-ring {
  position: absolute;
  width: 100%;
  height: 100%;
  border: 4px solid #E2E8F0;
  border-radius: 50%;
}

.spinner-dot {
  position: absolute;
  width: 100%;
  height: 100%;
  border: 4px solid #2563EB;
  border-radius: 50%;
  border-top-color: transparent;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.state-text {
  color: #64748B;
  margin: 0;
}

.empty-icon {
  width: 5rem;
  height: 5rem;
  border-radius: 50%;
  background: #F1F5F9;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94A3B8;
  margin-bottom: 0.5rem;
}

.state-title {
  font-family: 'Space Grotesk', sans-serif;
  font-size: 1.25rem;
  font-weight: 600;
  color: #1E293B;
  margin: 0;
}

.state-desc {
  color: #64748B;
  margin: 0;
}

.reload-button {
  padding: 0.625rem 1.5rem;
  background: #2563EB;
  color: white;
  font-weight: 500;
  border-radius: 0.5rem;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
}

.reload-button:hover {
  background: #1D4ED8;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
}

/* Weapon Groups Container */
.weapons-container {
  display: flex;
  flex-direction: column;
}

/* Custom scrollbar for sidebar */
.sidebar-content::-webkit-scrollbar {
  width: 0.375rem;
}

.sidebar-content::-webkit-scrollbar-track {
  background: transparent;
}

.sidebar-content::-webkit-scrollbar-thumb {
  background: #CBD5E1;
  border-radius: 9999px;
}

.sidebar-content::-webkit-scrollbar-thumb:hover {
  background: #94A3B8;
}
</style>
