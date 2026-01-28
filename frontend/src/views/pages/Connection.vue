<template>
  <div class="page-container">
    <!-- 顶部工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <NInput
          v-model:value="searchText"
          size="medium"
          clearable
          :placeholder="$t('frontend.connection.search')"
          :allow-input="value => !/\s/.test(value)"
          style="max-width: 300px"
        >
          <template #prefix>
            <Icon icon="ph:magnifying-glass" />
          </template>
        </NInput>
      </div>
      <div class="toolbar-right">
        <NTooltip trigger="hover">
          <template #trigger>
            <NButton type="primary" @click="handleAddGroup">
              <template #icon>
                <Icon icon="ph:folder-plus" />
              </template>
              {{ $t('frontend.connection.add.group') }}
            </NButton>
          </template>
        </NTooltip>
        <NTooltip trigger="hover">
          <template #trigger>
            <NButton type="primary" @click="handleAddConn">
              <template #icon>
                <Icon icon="ph:plus" />
              </template>
              {{ $t('frontend.connection.add.conn') }}
            </NButton>
          </template>
        </NTooltip>
      </div>
    </div>

    <!-- 分组列表 -->
    <div class="groups-container">
      <!-- 默认分组（未分组的连接） -->
      <div v-if="defaultGroupConns.length > 0 || !searchText" class="group-section">
        <div
          class="group-header"
          @click="toggleGroupCollapse('default')"
          @contextmenu="handleGroupContextMenu($event, null)"
        >
          <Icon
            :icon="defaultGroupConns.length > 0 ? 'ph:folders-duotone' : 'ph:folder-dashed-duotone'"
            class="group-icon"
          />
          <span class="group-title">{{ $t('frontend.connection.defaultGroup') }}</span>
          <span class="group-count">{{ defaultGroupConns.length }}</span>
        </div>
        <div v-show="!collapsedGroups.has('default')" class="group-content">
          <div v-if="defaultGroupConns.length > 0" class="conns-grid">
            <div
              v-for="conn in defaultGroupConns"
              :key="conn.id"
              class="conn-card"
              @click="toTerminal(conn)"
              @contextmenu="handleConnContextMenu($event, conn)"
            >
              <div class="card-content">
                <div class="os-icon" :class="{ 'text-logo': getSessionIcon(conn).isText }">
                  <Icon :icon="getSessionIcon(conn).icon" />
                </div>
                <div class="card-info">
                  <div class="conn-name">{{ conn.label }}</div>
                  <div class="conn-meta">
                    <span v-if="conn.connProtocol === ConnProtocol.SSH" class="conn-info">
                      {{ conn.credential?.username }}@{{ conn.host }}
                    </span>
                    <span v-if="conn.connProtocol === ConnProtocol.Serial" class="conn-info">
                      {{ conn.serialPort }}
                    </span>
                    <span class="protocol-badge">{{ conn.connProtocol }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="empty-group">
            <NEmpty
              :description="
                $t('frontend.connection.empty.group_desc', { name: $t('frontend.connection.defaultGroup') })
              "
            />
          </div>
        </div>
      </div>

      <!-- 其他分组 -->
      <div v-for="group in filteredGroups" :key="group.id" class="group-section">
        <div
          class="group-header"
          @click="toggleGroupCollapse(group.id)"
          @contextmenu="handleGroupContextMenu($event, group)"
        >
          <Icon
            :icon="getGroupConnCount(group) > 0 ? 'ph:folders-duotone' : 'ph:folder-dashed-duotone'"
            class="group-icon"
          />
          <span class="group-title">{{ group.name }}</span>
          <span class="group-count">{{ getGroupConnCount(group) }}</span>
        </div>
        <div v-show="!collapsedGroups.has(group.id)" class="group-content">
          <div v-if="getGroupConnCount(group) > 0" class="conns-grid">
            <div
              v-for="conn in getGroupConns(group)"
              :key="conn.id"
              class="conn-card"
              @click="toTerminal(conn)"
              @contextmenu="handleConnContextMenu($event, conn)"
            >
              <div class="card-content">
                <div class="os-icon" :class="{ 'text-logo': getSessionIcon(conn).isText }">
                  <Icon :icon="getSessionIcon(conn).icon" />
                </div>
                <div class="card-info">
                  <div class="conn-name">{{ conn.label }}</div>
                  <div class="conn-meta">
                    <span v-if="conn.connProtocol === ConnProtocol.SSH" class="conn-info">
                      {{ conn.credential?.username }}@{{ conn.host }}
                    </span>
                    <span v-if="conn.connProtocol === ConnProtocol.Serial" class="conn-info">
                      {{ conn.serialPort }}
                    </span>
                    <span class="protocol-badge">{{ conn.connProtocol }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="empty-group">
            <NEmpty :description="$t('frontend.connection.empty.group_desc', { name: group.name })" />
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="groupsWithConns.length === 0 && defaultGroupConns.length === 0" class="empty-state">
        <NResult
          status="404"
          :title="$t('frontend.connection.empty.title')"
          :description="$t('frontend.connection.empty.all_desc')"
        />
      </div>
    </div>

    <!-- 右键菜单 -->
    <NDropdown
      trigger="manual"
      placement="bottom-start"
      :show="showDropdown"
      :options="dropdownOptions"
      :x="dropdownX"
      :y="dropdownY"
      @select="handleDropdownSelect"
      @clickoutside="handleClickoutside"
    />

    <!-- 模态框 -->
    <ConnectionModal
      v-model:show="showConnModal"
      :is-edit="isEditConn"
      :connection-id="connectionId"
      @success="handleConnSuccess"
    />
    <GroupModal v-model:show="showGroupModal" :is-edit="isEditGroup" :group="editGroup" @success="handleGroupSuccess" />
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import type { Connection, Group } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/dal/model';
import { ConnProtocol } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/enums';
import { DeleteConnection, ListConnection } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/services/connectionsrv';
import { DeleteGroup, ListGroup } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/services/groupsrv';
import { NBadge, NButton, NDropdown, NEmpty, NInput, NResult, NTooltip, useDialog, useThemeVars } from 'naive-ui';
import type { DropdownOption } from 'naive-ui';
import { computed, h, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useConnectionStore } from '@/stores/connection';
import { useCall } from '@/utils/call';
import ConnectionModal from '@/views/modals/ConnectionModal.vue';
import GroupModal from '@/views/modals/GroupModal.vue';

const router = useRouter();
const connStore = useConnectionStore();
const { t } = useI18n();
const { call } = useCall();

const groups = ref<Group[]>([]);
const conns = ref<Connection[]>([]);
const collapsedGroups = ref<Set<number | 'default'>>(new Set());

const showDropdown = ref(false);
const dropdownX = ref(0);
const dropdownY = ref(0);
const dropdownOptions = ref<DropdownOption[]>([]);

const searchText = ref('');

const showConnModal = ref(false);
const showGroupModal = ref(false);
const isEditConn = ref(false);
const isEditGroup = ref(false);
const connectionId = ref<number>(0);
const editGroup = ref<Group | undefined>(undefined);

const currentContextNode = ref<any>(null);

const dialog = useDialog();
const themeVars = useThemeVars();

// 默认分组（未分组的连接）
const defaultGroupConns = computed(() => {
  let filtered = conns.value?.filter(conn => !conn.groupID || conn.groupID === 0) || [];
  if (searchText.value) {
    filtered = filtered.filter(
      conn =>
        conn.label.toLowerCase().includes(searchText.value.toLowerCase()) ||
        conn.host?.toLowerCase().includes(searchText.value.toLowerCase()),
    );
  }
  return filtered;
});

// 过滤后的分组列表
const filteredGroups = computed(() => {
  if (!searchText.value) return groups.value || [];
  return groups.value?.filter(group => group.name.toLowerCase().includes(searchText.value.toLowerCase())) || [];
});

// 有连接的分组列表（用于判断是否显示空状态）
const groupsWithConns = computed(() => {
  return filteredGroups.value.filter(group => getGroupConnCount(group) > 0);
});

// 获取分组下的连接
const getGroupConns = (group: Group) => {
  let filtered = conns.value?.filter(conn => conn.groupID === group.id) || [];
  if (searchText.value) {
    filtered = filtered.filter(
      conn =>
        conn.label.toLowerCase().includes(searchText.value.toLowerCase()) ||
        conn.host?.toLowerCase().includes(searchText.value.toLowerCase()),
    );
  }
  return filtered;
};

// 获取分组连接数量
const getGroupConnCount = (group: Group) => {
  return conns.value?.filter(conn => conn.groupID === group.id).length || 0;
};

// 切换分组折叠状态
const toggleGroupCollapse = (groupId: number | 'default') => {
  if (collapsedGroups.value.has(groupId)) {
    collapsedGroups.value.delete(groupId);
  } else {
    collapsedGroups.value.add(groupId);
  }
};

const updateDropdownOptions = (type: 'group' | 'conn', isDefaultGroup = false) => {
  if (type === 'group') {
    const options: DropdownOption[] = [];
    if (!isDefaultGroup) {
      options.push({
        label: t('frontend.connection.menu.edit_group'),
        key: 'edit-group',
        icon: () => h(Icon, { icon: 'ph:pencil-simple' }),
      });
      options.push({
        label: t('frontend.connection.menu.delete_group'),
        key: 'delete-group',
        icon: () => h(Icon, { icon: 'ph:trash' }),
      });
    }
    dropdownOptions.value = options;
  } else {
    dropdownOptions.value = [
      {
        label: t('frontend.connection.menu.edit_conn'),
        key: 'edit-conn',
        icon: () => h(Icon, { icon: 'ph:pencil-simple' }),
      },
      {
        label: t('frontend.connection.menu.delete_conn'),
        key: 'delete-conn',
        icon: () => h(Icon, { icon: 'ph:trash' }),
      },
    ];
  }
};

const handleClickoutside = () => {
  showDropdown.value = false;
};

const toTerminal = (conn: Connection) => {
  const connection = {
    id: Date.now(),
    connId: conn.id,
    label: `${conn.label} (${connStore.connections.filter(c => c.host === conn.host).length + 1})`,
    host: conn.host,
    username: conn.credential?.username || '',
    theme: conn.theme || 'Default',
  };
  connStore.addConnection(connection);
  router.push({ name: 'Terminal' });
};

const fetchGroups = async () => {
  const result = await call(ListGroup);
  return result.data || [];
};

const fetchConns = async () => {
  const result = await call(ListConnection);
  return result.data || [];
};

const fetchData = async () => {
  const [groupsData, connsData] = await Promise.all([fetchGroups(), fetchConns()]);
  groups.value = groupsData;
  conns.value = connsData;

  // 默认折叠没有连接的分组
  if (defaultGroupConns.value.length === 0) {
    collapsedGroups.value.add('default');
  }

  groups.value.forEach(group => {
    if (getGroupConnCount(group) === 0) {
      collapsedGroups.value.add(group.id);
    }
  });
};

const handleEditConn = (conn: Connection) => {
  isEditConn.value = true;
  connectionId.value = conn.id;
  showConnModal.value = true;
};

const handleConnSuccess = () => {
  fetchData();
};

const handleAddConn = () => {
  isEditConn.value = false;
  connectionId.value = 0;
  showConnModal.value = true;
};

const handleAddGroup = () => {
  isEditGroup.value = false;
  editGroup.value = undefined;
  showGroupModal.value = true;
};

const handleEditGroup = () => {
  if (!currentContextNode.value) return;
  const groupId =
    currentContextNode.value.key === 'default-group'
      ? null
      : Number.parseInt(currentContextNode.value.key.replace('group-', ''));
  if (groupId === null) return; // 默认分组不能编辑
  const group = groups.value?.find(g => g.id === groupId);
  if (group) {
    isEditGroup.value = true;
    editGroup.value = group;
    showGroupModal.value = true;
  }
};

const handleGroupSuccess = () => {
  fetchData();
};

const vendorIconMap: Record<string, string> = {
  cisco: 'simple-icons:cisco',
  huawei: 'simple-icons:huawei',
  fortinet: 'simple-icons:fortinet',
  mikrotik: 'simple-icons:mikrotik',
  pfsense: 'simple-icons:pfsense',
  juniper: 'simple-icons:junipernetworks',
  hp: 'simple-icons:hp',
  dell: 'simple-icons:dell',
  redhat: 'simple-icons:redhat',
  ubuntu: 'simple-icons:ubuntu',
  centos: 'simple-icons:centos',
  debian: 'simple-icons:debian',
  opensuse: 'simple-icons:opensuse',
  fedora: 'simple-icons:fedora',
  almalinux: 'simple-icons:almalinux',
  kalilinux: 'simple-icons:kalilinux',
  archlinux: 'simple-icons:archlinux',
  rockylinux: 'simple-icons:rockylinux',
  alpinelinux: 'simple-icons:alpinelinux',
  gentoo: 'simple-icons:gentoo',
  raspberrypi: 'simple-icons:raspberrypi',
  linuxmint: 'simple-icons:linuxmint',
  elementary: 'simple-icons:elementary',
  zorinos: 'simple-icons:zorinos',
  popos: 'simple-icons:popos',
  linux: 'simple-icons:linux',
  vmware: 'simple-icons:vmware',
};

const isTextLogo = (vendor: string) => {
  const textLogoVendors = ['cisco', 'juniper', 'vmware'];
  return textLogoVendors.includes(vendor.toLowerCase());
};

const getSessionIcon = (conn: Connection) => {
  const vendor = conn.metadata?.vendor || '';
  if (vendor && vendor in vendorIconMap) {
    return {
      icon: vendorIconMap[vendor],
      isText: isTextLogo(vendor),
    };
  }

  const type = conn.metadata?.type || '';
  if (type === 'server') {
    return {
      icon: 'ph:hard-drives',
      isText: false,
    };
  } else if (type === 'network') {
    return {
      icon: 'ph:network',
      isText: false,
    };
  }
  return {
    icon: 'ph:ghost',
    isText: false,
  };
};

const handleConnContextMenu = (event: MouseEvent, conn: Connection) => {
  event.preventDefault();
  dropdownX.value = event.clientX;
  dropdownY.value = event.clientY;
  showDropdown.value = true;
  currentContextNode.value = { key: `conn-${conn.id}` };
  updateDropdownOptions('conn');
};

const handleGroupContextMenu = (event: MouseEvent, group: Group | null) => {
  event.preventDefault();
  dropdownX.value = event.clientX;
  dropdownY.value = event.clientY;
  showDropdown.value = true;
  const isDefaultGroup = group === null;
  currentContextNode.value = { key: group ? `group-${group.id}` : 'default-group' };
  updateDropdownOptions('group', isDefaultGroup);
};

const handleDeleteGroup = async () => {
  if (!currentContextNode.value) return;
  // 默认分组不能删除
  if (currentContextNode.value.key === 'default-group') return;
  const groupId = Number.parseInt(currentContextNode.value.key.replace('group-', ''));
  const group = groups.value?.find(g => g.id === groupId);
  if (!group) return;

  dialog.warning({
    title: t('frontend.connection.delete.group.title'),
    content: t('frontend.connection.delete.group.content', { name: group.name }),
    positiveText: t('frontend.connection.delete.group.confirm'),
    negativeText: t('frontend.connection.delete.group.cancel'),
    onPositiveClick: async () => {
      const result = await call(DeleteGroup, {
        args: [groupId],
      });

      if (result.ok) {
        await fetchData();
      }
    },
  });
};

const handleDropdownSelect = async (key: string) => {
  showDropdown.value = false;
  if (!currentContextNode.value) return;

  const connId = Number.parseInt(currentContextNode.value.key.replace('conn-', ''));

  switch (key) {
    case 'edit-group': {
      handleEditGroup();
      break;
    }
    case 'delete-group': {
      await handleDeleteGroup();
      break;
    }
    case 'edit-conn': {
      const conn = conns.value?.find(h => h.id === connId);
      if (conn) handleEditConn(conn);
      break;
    }
    case 'delete-conn': {
      const result = await call(DeleteConnection, {
        args: [connId],
      });

      if (result.ok) {
        await fetchData();
      }
      break;
    }
  }
  currentContextNode.value = null;
};

onMounted(async () => {
  await fetchData();
});
</script>

<style lang="less" scoped>
.page-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 12px;
  overflow: hidden;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  flex-shrink: 0;

  .toolbar-left {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .toolbar-right {
    display: flex;
    align-items: center;
    gap: 12px;
  }
}

.groups-container {
  flex: 1;
  overflow-y: auto;
  padding-right: 4px;

  &::-webkit-scrollbar {
    width: 8px;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
  }

  &::-webkit-scrollbar-thumb {
    background: v-bind('themeVars.borderColor');
    border-radius: 4px;

    &:hover {
      background: v-bind('themeVars.textColor3');
    }
  }
}

.group-section {
  margin-bottom: 16px;

  &:last-child {
    margin-bottom: 0;
  }
}

.group-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 8px;
  border-radius: 6px;

  &:hover {
    background: v-bind('`${themeVars.primaryColor}08`');
  }

  .group-icon {
    font-size: 18px;
    color: v-bind('themeVars.primaryColor');
    flex-shrink: 0;
  }

  .group-title {
    font-size: 14px;
    font-weight: 500;
    color: v-bind('themeVars.textColorBase');
    flex: 1;
  }

  .group-count {
    font-size: 12px;
    color: v-bind('themeVars.textColor3');
    background: v-bind('themeVars.borderColor');
    padding: 2px 8px;
    border-radius: 10px;
    flex-shrink: 0;
  }
}

.group-content {
  padding-left: 8px;
}

.empty-group {
  padding: 24px;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 120px;
}

.conns-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
}

.conn-card {
  background: v-bind('themeVars.cardColor');
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid v-bind('themeVars.borderColor');
  padding: 8px;

  &:hover {
    border-color: v-bind('themeVars.primaryColor');
    background: v-bind('`${themeVars.primaryColor}05`');
  }

  .card-content {
    display: flex;
    align-items: center;
    gap: 10px;

    .os-icon {
      width: 40px;
      height: 40px;
      border-radius: 4px;
      background: v-bind('`${themeVars.primaryColor}15`');
      color: v-bind('themeVars.primaryColor');
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 22px;
      flex-shrink: 0;

      &.text-logo {
        :deep(svg) {
          width: 80%;
          height: 100%;
        }
      }

      :deep(svg) {
        width: 70%;
        height: 100%;
        object-fit: contain;
      }
    }

    .card-info {
      flex: 1;
      min-width: 0;

      .conn-name {
        font-size: 14px;
        font-weight: 500;
        color: v-bind('themeVars.textColorBase');
        margin-bottom: 2px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .conn-meta {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 12px;

        .conn-info {
          color: v-bind('themeVars.textColor3');
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
          flex: 1;
          min-width: 0;
        }

        .protocol-badge {
          color: v-bind('themeVars.textColor3');
          background: v-bind('themeVars.borderColor');
          padding: 2px 6px;
          border-radius: 4px;
          font-size: 11px;
          flex-shrink: 0;
        }
      }
    }
  }
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  width: 100%;
}
</style>
