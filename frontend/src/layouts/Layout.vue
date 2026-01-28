<template>
  <NLayout class="layout-container">
    <NLayoutHeader bordered class="header">
      <Header ref="titlebarRef" :is-darwin="isDarwin" />
    </NLayoutHeader>

    <NLayout class="content">
      <NLayoutContent>
        <router-view v-slot="{ Component }">
          <keep-alive>
            <component :is="Component" />
          </keep-alive>
        </router-view>
      </NLayoutContent>
    </NLayout>
  </NLayout>
</template>

<script lang="ts" setup>
import { IsDarwin } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/services/preferencessrv';
import { NLayout, NLayoutContent, NLayoutHeader } from 'naive-ui';
import { computed, onMounted, provide, ref } from 'vue';
import Header from '@/layouts/Header.vue';

const headerRef = ref();
const titlebarRef = ref();
const isDarwin = ref(false);

onMounted(async () => {
  isDarwin.value = await IsDarwin();
});

provide(
  'connectionTabs',
  computed(() => titlebarRef.value?.connectionTabsRef || headerRef.value?.connectionTabsRef),
);
</script>

<style lang="less" scoped>
.layout-container {
  height: 100vh;
  .header {
    height: 35px;
    display: flex;
    align-items: center;
  }
  .content {
    height: calc(100vh - 35px);
  }
}
</style>
