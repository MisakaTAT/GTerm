<template>
  <div ref="widthRef" class="header">
    <div
      class="logo-container"
      :class="{
        'padding-left-small': !isDarwin,
        'padding-left-large': isDarwin,
        'center-position': !hasConnections,
        'left-position': hasConnections,
      }"
    >
      <div class="logo-wrapper" @click="toConnection">
        <img src="@/assets/images/icon.png" alt="Logo" class="logo-image" />
        <span class="app-title" :class="{ hidden: hasConnections }">GTerm</span>
      </div>
      <ConnectionTabs v-if="hasConnections" ref="connectionTabsRef" class="connection-tabs-container" />
    </div>

    <div v-if="!isDarwin" class="window-controls">
      <div class="window-control-btn" @click="WindowMinimise">
        <n-icon size="16"><icon icon="ph:minus-bold" /></n-icon>
      </div>
      <div class="window-control-btn" @click="reduction">
        <n-icon size="16">
          <span v-if="windowIsMaximised"><icon icon="ph:corners-in-bold" /></span>
          <span v-else><icon icon="ph:corners-out-bold" /></span>
        </n-icon>
      </div>
      <div class="window-control-btn close-btn" @click="Quit">
        <n-icon size="16"><icon icon="ph:x-bold" /></n-icon>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { WindowMinimise, Quit, WindowIsMaximised, WindowMaximise, WindowUnmaximise } from '@wailsApp/runtime';
import { NIcon } from 'naive-ui';
import { Icon } from '@iconify/vue';
import { IsDarwin } from '@wailsApp/go/services/PreferencesSrv';
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

const toConnection = () => {
  router.push({ name: 'Connection' });
};

const windowIsMaximised = ref(false);

const reduction = () => {
  WindowIsMaximised().then(isMaximised => {
    if (isMaximised) {
      windowIsMaximised.value = false;
      WindowUnmaximise();
    } else {
      windowIsMaximised.value = true;
      WindowMaximise();
    }
  });
};

onMounted(async () => {
  isDarwin.value = await IsDarwin();
});
</script>

<style lang="less" scoped>
.header {
  --wails-draggable: drag;
  width: 100%;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-top-left-radius: 8px;
  border-top-right-radius: 8px;
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
