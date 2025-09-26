<template>
  <div class="header" @dblclick="toggleFullscreen">
    <div
      class="logo-container"
      :class="{
        'padding-left-small': !isDarwin,
        'padding-left-large': isDarwin,
        'center-position': !hasConnections && !isFullscreen,
        'left-position': hasConnections || isFullscreen,
        fullscreen: isFullscreen,
      }"
    >
      <div class="logo-wrapper" @click="toConnection">
        <img src="@/assets/images/icon.png" alt="Logo" class="logo-image" />
        <span class="app-title" :class="{ hidden: hasConnections }">GTerm</span>
      </div>
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
const router = useRouter();

const isDarwin = ref(false);
const isFullscreen = ref(false);
const windowIsMaximised = ref(false);

const toConnection = () => {
  router.push({ name: 'Connection' });
};

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
}

.logo-container {
  position: absolute;
  display: flex;
  align-items: center;
  transition: all 0.3s ease-in-out;

  &.padding-left-small {
    padding-left: 8px;
  }

  &.padding-left-large {
    padding-left: 80px;

    &.fullscreen {
      padding-left: 8px;
    }
  }

  &.center-position {
    left: 50%;
    transform: translateX(-50%);
  }

  &.left-position {
    left: 0;
    transform: translateX(0);
  }
}

.logo-wrapper {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.logo-image {
  width: 24px;
  height: 24px;
}

.app-title {
  padding-left: 8px;
  font-size: 14px;
  font-weight: 600;

  &.hidden {
    display: none;
  }
}

.connection-tabs-container {
  margin-left: 16px;
}

.window-controls {
  display: flex;
  align-items: center;
  margin-left: auto;
}

.window-control-btn {
  width: 38px;
  height: 38px;
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
