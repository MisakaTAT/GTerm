<template>
  <div class="tabs">
    <div
      v-for="tab in tabs"
      :key="tab.id"
      class="tab-item"
      :class="{ active: tab.id === activeTab }"
      @click="switchTab(tab.id)"
    >
      <div class="status-dot" :class="tabStatus[tab.id]" />
      <span class="tab-label">{{ tab.label }}</span>
      <NButton circle text size="tiny" class="close-btn" @click.stop="closeTab(tab.id)">
        <template #icon>
          <Icon icon="ph:x-bold" />
        </template>
      </NButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { NButton, useThemeVars } from 'naive-ui';
import { useRouter } from 'vue-router';
import { useConnectionStore } from '@/stores/connection';

const router = useRouter();
const connectionStore = useConnectionStore();
const activeTab = computed(() => connectionStore.activeConnectionId);
const themeVars = useThemeVars();
const tabs = computed(() => connectionStore.connections);
const terminalRefs = ref<Map<number, any>>(new Map());
const tabStatus = ref<Record<number, string>>({});

const updateTabStatus = (id: number, status: 'connected' | 'error' | 'connecting' | 'warning') => {
  tabStatus.value[id] = status;
};

const registerTerminal = (id: number, terminal: any) => {
  terminalRefs.value.set(id, terminal);
  if (terminal.status) {
    tabStatus.value[id] = terminal.status;
  }
};

const switchTab = (id: number) => {
  connectionStore.setActiveConnection(id);
  if (router.currentRoute.value.name !== 'Terminal') {
    router.push({ name: 'Terminal' });
  }
};

const closeTab = (id: number) => {
  const terminal = terminalRefs.value.get(id);
  if (terminal?.closeTerminal) {
    terminal.closeTerminal();
  }
  terminalRefs.value.delete(id);
  connectionStore.removeConnection(id);

  if (connectionStore.connections.length === 0) {
    router.push({ name: 'Connection' });
  }
};

defineExpose({
  registerTerminal,
  updateTabStatus,
});
</script>

<style lang="less" scoped>
.tabs {
  display: flex;
  align-items: center;
  gap: 4px;
  height: 100%;
  overflow-x: auto;
  overflow-y: hidden;
  scrollbar-width: none;
  -ms-overflow-style: none;

  &::-webkit-scrollbar {
    display: none;
  }
}

.tab-item {
  display: flex;
  align-items: center;
  gap: 4px;
  height: 20px;
  padding: 0 8px;
  border-radius: 10px;
  border: 1px solid v-bind('themeVars.borderColor');
  background-color: v-bind('themeVars.cardColor');
  cursor: pointer;
  white-space: nowrap;
  flex-shrink: 0;
  transition: all 0.2s;

  .tab-label {
    font-size: 11px;
    color: v-bind('themeVars.textColorBase');
    max-width: 100px;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .status-dot {
    width: 5px;
    height: 5px;
    border-radius: 50%;
    flex-shrink: 0;

    &.connected {
      background-color: v-bind('themeVars.successColor');
    }

    &.error {
      background-color: v-bind('themeVars.errorColor');
    }

    &.connecting {
      background-color: v-bind('themeVars.infoColor');
      animation: pulse 1.5s infinite;
    }
  }

  .close-btn {
    opacity: 0;
    width: 14px;
    height: 14px;
    min-width: 14px;
    transition: opacity 0.2s;
  }

  &:hover .close-btn {
    opacity: 1;
  }

  &.active {
    background-color: v-bind('themeVars.primaryColor');
    border-color: v-bind('themeVars.primaryColor');

    .tab-label {
      color: #fff;
    }

    .status-dot {
      background-color: #fff;
    }

    .close-btn {
      color: #fff;
    }
  }
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
</style>
