<template>
  <NModal
    v-model:show="dialogStore.aboutDialogVisible"
    preset="dialog"
    :close-on-esc="true"
    :title="$t('frontend.about.name')"
    transform-origin="center"
    :show-icon="false"
    style="width: 460px"
    :on-close="dialogStore.closeAboutDialog"
    :auto-focus="false"
  >
    <NSpace :size="10" :wrap="false" :wrap-item="false" align="center" vertical>
      <NAvatar :size="80" :src="iconUrl" color="#0000" />
      <div class="about-app-title">GTerm</div>
      <n-button text @click="onOpenVersionURL">{{ version }}</n-button>
      <NSpace :size="1" :wrap="false" :wrap-item="false" align="center">
        <n-button text @click="onOpenSource">{{ $t('frontend.about.github') }}</n-button>
        <NDivider vertical />
        <n-button text @click="onOpenWebsite">{{ $t('frontend.about.website') }}</n-button>
      </NSpace>
      <div :style="{ color: themeVars.textColor3 }" class="about-copyright">
        <span>{{ copyright }}</span>
      </div>
    </NSpace>
  </NModal>
</template>

<script setup lang="ts">
import { Copyright, Version, VersionURL } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/services/preferencessrv';
import { Browser } from '@wailsio/runtime';
import { NAvatar, NDivider, NModal, NSpace, useThemeVars } from 'naive-ui';
import iconUrl from '@/assets/images/icon.png';
import { useDialogStore } from '@/stores/dialog';

const version = ref('');
const copyright = ref('');
const versionUrl = ref('');

onMounted(async () => {
  version.value = await Version();
  copyright.value = await Copyright();
  versionUrl.value = await VersionURL();
});

const dialogStore = useDialogStore();
const themeVars = useThemeVars();

const onOpenVersionURL = () => {
  Browser.OpenURL(versionUrl.value);
};

const onOpenSource = () => {
  Browser.OpenURL('https://github.com/MisakaTAT/GTerm');
};

const onOpenWebsite = () => {
  Browser.OpenURL('https://github.com/MisakaTAT');
};
</script>

<style lang="less" scoped>
.about-app-title {
  font-weight: bold;
  font-size: 18px;
}

.about-copyright {
  font-size: 12px;
  text-align: center;
  white-space: pre-wrap;
}
</style>
