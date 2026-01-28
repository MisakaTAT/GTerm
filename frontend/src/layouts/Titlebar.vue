<template>
  <div class="titlebar-container" @dblclick="toggleFullscreen">
    <div class="titlebar-left" :class="{ 'darwin-left': isDarwin }">
      <NTooltip placement="bottom" trigger="hover">
        <template #trigger>
          <div
            class="menu-item top-menu-item"
            :class="{ active: selectedKey === 'Connection' }"
            @click="handleSelect('Connection')"
          >
            <NIcon size="16">
              <Icon icon="ph:hard-drives" />
            </NIcon>
          </div>
        </template>
        {{ $t('frontend.sider.assets') }}
      </NTooltip>

      <NTooltip placement="bottom" trigger="hover">
        <template #trigger>
          <div
            class="menu-item top-menu-item"
            :class="{ active: selectedKey === 'FileTransfer' }"
            @click="handleSelect('FileTransfer')"
          >
            <NIcon size="16">
              <Icon icon="ph:folders" />
            </NIcon>
          </div>
        </template>
        {{ $t('frontend.sider.file_transfer') }}
      </NTooltip>

      <NTooltip placement="bottom" trigger="hover">
        <template #trigger>
          <div
            class="menu-item top-menu-item"
            :class="{ active: selectedKey === 'Credential' }"
            @click="handleSelect('Credential')"
          >
            <NIcon size="16">
              <Icon icon="ph:vault" />
            </NIcon>
          </div>
        </template>
        {{ $t('frontend.sider.credentials') }}
      </NTooltip>
    </div>

    <div class="titlebar-center">
      <ConnectionTabs v-if="hasConnections" ref="connectionTabsRef" class="connection-tabs-container" />
    </div>

    <div class="titlebar-right">
      <NTooltip placement="bottom" trigger="hover">
        <template #trigger>
          <div class="menu-item bottom-menu-item" @click="toggleTheme">
            <NIcon size="16">
              <Icon :icon="prefStore.isDark ? 'ph:sun' : 'ph:moon'" />
            </NIcon>
          </div>
        </template>
        {{ prefStore.isDark ? $t('frontend.sider.theme.toggle_light') : $t('frontend.sider.theme.toggle_dark') }}
      </NTooltip>

      <NDropdown
        :options="[
          {
            label: $t('frontend.sider.menu.preferences'),
            key: 'preferences',
            icon: () => h(NIcon, { size: 'large' }, { default: () => h(Icon, { icon: 'ph:sliders-horizontal' }) }),
          },
          {
            label: $t('frontend.sider.menu.check_update'),
            key: 'check-update',
            icon: () => h(NIcon, { size: 'large' }, { default: () => h(Icon, { icon: 'ph:arrow-circle-up' }) }),
          },
          {
            label: $t('frontend.sider.menu.about'),
            key: 'about',
            icon: () => h(NIcon, { size: 'large' }, { default: () => h(Icon, { icon: 'ph:info' }) }),
          },
        ]"
        trigger="click"
        :width="180"
        placement="bottom"
        @select="handleSettingsSelect"
      >
        <NTooltip placement="bottom" trigger="hover">
          <template #trigger>
            <div class="menu-item bottom-menu-item">
              <NIcon size="16">
                <Icon icon="ph:gear-six" />
              </NIcon>
            </div>
          </template>
          {{ $t('frontend.sider.settings') }}
        </NTooltip>
      </NDropdown>

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
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { Events, Window } from '@wailsio/runtime';
import { NDropdown, NIcon, NTooltip, useThemeVars } from 'naive-ui';
import { h, onMounted, onUnmounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import ConnectionTabs from '@/layouts/ConnectionTabs.vue';
import { useConnectionStore } from '@/stores/connection';
import { useDialogStore } from '@/stores/dialog';
import { usePreferencesStore } from '@/stores/preferences';

const props = defineProps<{
  isDarwin?: boolean;
}>();

const router = useRouter();
const route = useRoute();
const prefStore = usePreferencesStore();
const dialogStore = useDialogStore();
const connectionStore = useConnectionStore();
const selectedKey = ref('Connection');
const themeVars = useThemeVars();
const connectionTabsRef = ref();
const isDarwin = computed(() => props.isDarwin || false);
const isFullscreen = ref(false);
const windowIsMaximised = ref(false);

const hasConnections = computed(() => connectionStore.hasConnections);

defineExpose({
  connectionTabsRef,
});

const updateSelectedKey = () => {
  const routeName = route.name as string;
  selectedKey.value = routeName;
};

watch(() => route.name, updateSelectedKey);

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
  updateSelectedKey();
  await updateWindowState();
  Events.On('window:state-changed', updateWindowState);
  window.addEventListener('resize', updateWindowState);
});

onUnmounted(() => {
  Events.Off('window:state-changed');
  window.removeEventListener('resize', updateWindowState);
});

const handleSelect = (key: string) => {
  selectedKey.value = key;
  router.push({ name: key });
};

const handleSettingsSelect = (key: string) => {
  switch (key) {
    case 'preferences':
      router.push({ name: 'Preferences' });
      break;
    case 'about':
      dialogStore.openAboutDialog();
      break;
    case 'check-update':
      // TODO: 实现检查更新功能
      break;
  }
};

const toggleTheme = () => {
  prefStore.isDark ? prefStore.toLight() : prefStore.toDark();
};
</script>

<style lang="less" scoped>
.titlebar-container {
  --wails-draggable: drag;
  display: flex;
  flex-direction: row;
  height: 100%;
  align-items: center;
  width: 100%;
  flex: 1;
}

.titlebar-left {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 0;
  flex-shrink: 0;
  height: 100%;

  &.darwin-left {
    margin-left: 82px;
  }
}

.titlebar-center {
  display: flex;
  flex-direction: row;
  align-items: center;
  flex: 1;
  overflow: hidden;
  min-width: 0;
  height: 100%;
}

.titlebar-right {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 0;
  flex-shrink: 0;
  margin-left: auto;
  padding-right: 8px;
  height: 100%;
}

.connection-tabs-container {
  margin-left: 8px;
}

.window-controls {
  display: flex;
  align-items: center;
  gap: 0;
  margin-left: 4px;
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

.menu-item {
  width: 28px;
  height: 28px;
  padding: 3px;
  box-sizing: border-box;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 4px;
  cursor: pointer;
  position: relative;

  &::before {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 22px;
    height: 22px;
    border-radius: 4px;
    opacity: 0;
    transition:
      opacity 0.2s,
      background-color 0.2s;
  }

  &:hover::before {
    opacity: 1;
    background-color: v-bind('themeVars.hoverColor');
  }
}

.top-menu-item {
  &:hover {
    color: v-bind('themeVars.primaryColor');
  }

  &.active {
    color: v-bind('themeVars.primaryColor');

    &::before {
      opacity: 1;
      background-color: v-bind('themeVars.hoverColor');
    }
  }
}

.bottom-menu-item {
  &:hover {
    color: v-bind('themeVars.primaryColor');
  }
}
</style>
