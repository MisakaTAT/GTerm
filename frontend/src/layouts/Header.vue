<template>
  <div class="header" @dblclick="toggleFullscreen">
    <div class="header-content">
      <ConnectionTabs v-if="hasConnections" ref="connectionTabsRef" class="connection-tabs-container" />
    </div>

    <div v-if="!isDarwin" class="window-controls">
      <div class="window-control-btn" @click="Window.Minimise">
        <NIcon size="16"><Icon icon="ph:minus-bold" /></NIcon>
      </div>
      <div class="window-control-btn" @click="toggleMaximize">
        <NIcon size="16">
          <Icon :icon="windowIsMaximised ? 'ph:corners-in-bold' : 'ph:corners-out-bold'" />
        </NIcon>
      </div>
      <div class="window-control-btn close-btn" @click="Window.Close">
        <NIcon size="16"><Icon icon="ph:x-bold" /></NIcon>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Icon } from '@iconify/vue';
import { IsDarwin } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/services/preferencessrv';
import { Events, Window } from '@wailsio/runtime';
import { NIcon } from 'naive-ui';
import ConnectionTabs from '@/layouts/ConnectionTabs.vue';
import { useConnectionStore } from '@/stores/connection';

const connectionTabsRef = ref();

defineExpose({
  connectionTabsRef,
});

const connectionStore = useConnectionStore();
const hasConnections = computed(() => connectionStore.hasConnections);

const isDarwin = ref(false);
const isFullscreen = ref(false);
const windowIsMaximised = ref(false);

const toggleMaximize = async () => {
  const isMaximised = await Window.IsMaximised();
  windowIsMaximised.value = !isMaximised;
  isMaximised ? Window.UnMaximise() : Window.Maximise();
};

const checkFullscreenStatus = async () => {
  isFullscreen.value = await Window.IsFullscreen();
};

const toggleFullscreen = async () => {
  const fullscreen = await Window.IsFullscreen();
  if (fullscreen) {
    Window.UnFullscreen();
  } else {
    Window.Fullscreen();
  }
  await checkFullscreenStatus();
};

const updateWindowState = async () => {
  await checkFullscreenStatus();
  windowIsMaximised.value = await Window.IsMaximised();
};

onMounted(async () => {
  isDarwin.value = await IsDarwin();
  await updateWindowState();
  Events.On('window:state-changed', updateWindowState);
  window.addEventListener('resize', updateWindowState);
});

onUnmounted(() => {
  Events.Off('window:state-changed');
  window.removeEventListener('resize', updateWindowState);
});
</script>

<style lang="less" scoped>
.header {
  --wails-draggable: drag;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  border-top-left-radius: 8px;
  border-top-right-radius: 8px;
  justify-content: space-between;
  flex: 1;
}

.header-content {
  display: flex;
  align-items: center;
  flex: 1;
}

.connection-tabs-container {
  margin-left: 8px;
}

.window-controls {
  display: flex;
  align-items: center;
  margin-left: auto;
}

.window-control-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: rgba(0, 0, 0, 0.06);
  }

  &.close-btn {
    border-top-right-radius: 8px;

    &:hover {
      background-color: #e54d42;
      color: #fff;
    }
  }
}
</style>
