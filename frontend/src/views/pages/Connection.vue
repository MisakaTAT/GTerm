<template>
  <div class="page-container">
    <div ref="sidebarRef" class="sidebar" :style="{ width: `${sidebarWidth}px` }">
      <div class="groups-list">
        <div class="list-header">
          <NInput
            v-model:value="searchText"
            size="small"
            clearable
            :placeholder="$t('frontend.connection.search')"
            :allow-input="value => !/\s/.test(value)"
          >
            <template #prefix>
              <Icon icon="ph:magnifying-glass" />
            </template>
          </NInput>
          <n-divider vertical />
          <div class="header-right">
            <NTooltip trigger="hover">
              <template #trigger>
                <NButton text size="large" @click="handleAddGroup">
                  <template #icon>
                    <Icon icon="ph:folder-plus" />
                  </template>
                </NButton>
              </template>
              {{ $t('frontend.connection.add.group') }}
            </NTooltip>
            <NTooltip trigger="hover">
              <template #trigger>
                <NButton text size="large" @click="handleAddConn">
                  <template #icon>
                    <Icon icon="ph:plus" />
                  </template>
                </NButton>
              </template>
              {{ $t('frontend.connection.add.conn') }}
            </NTooltip>
          </div>
        </div>
        <div class="list-content">
          <!-- 资产列表部分 -->
          <div class="section-header">
            <div class="section-title">{{ $t('frontend.connection.sections.assetList') }}</div>
            <NButton text size="tiny" @click="assetListCollapsed = !assetListCollapsed">
              <template #icon>
                <Icon :icon="assetListCollapsed ? 'ph:caret-right' : 'ph:caret-down'" />
              </template>
            </NButton>
          </div>
          <div v-show="!assetListCollapsed" class="asset-list">
            <div
              v-for="conn in filteredAssets"
              :key="conn.id"
              class="asset-item"
              @click="handleSelectConn(conn)"
              @contextmenu="handleConnContextMenu($event, conn)"
            >
              <div class="asset-icon">
                <Icon :icon="getSessionIcon(conn).icon" :class="{ 'text-logo': getSessionIcon(conn).isText }" />
              </div>
              <div class="asset-info">
                <div class="asset-name">{{ conn.label }}</div>
              </div>
            </div>
          </div>

          <!-- 分组列表部分 -->
          <div class="section-header">
            <div class="section-title">{{ $t('frontend.connection.sections.groupList') }}</div>
            <NButton text size="tiny" @click="groupListCollapsed = !groupListCollapsed">
              <template #icon>
                <Icon :icon="groupListCollapsed ? 'ph:caret-right' : 'ph:caret-down'" />
              </template>
            </NButton>
          </div>
          <div v-show="!groupListCollapsed" class="group-list">
            <div
              v-for="group in filteredGroups"
              :key="group.id"
              class="group-item"
              :class="{ active: selectedGroup?.id === group.id }"
              @click="handleSelectGroup(group)"
              @contextmenu="handleGroupContextMenu($event, group)"
            >
              <div class="group-icon">
                <Icon :icon="getGroupConnCount(group) > 0 ? 'ph:folders-duotone' : 'ph:folder-dashed-duotone'" />
              </div>
              <div class="group-info">
                <div class="group-name">{{ group.name }}</div>
              </div>
              <div class="group-count">{{ getGroupConnCount(group) }}</div>
            </div>
          </div>

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
        </div>
      </div>
    </div>

    <div class="resize-handle" :style="{ left: `${sidebarWidth}px` }" @mousedown="startResize"></div>

    <div class="main-content">
      <div class="content-header">
        <div class="header-left">
          <h2>{{ selectedGroup ? selectedGroup.name : $t('frontend.sider.assets') }}</h2>
          <NBadge :value="filteredConns.length" show-zero type="success" />
        </div>
      </div>

      <div v-if="filteredConns.length > 0" class="conns-grid">
        <div v-for="conn in filteredConns" :key="conn.id" class="conn-card" @click="toTerminal(conn)">
          <div class="card-header">
            <div class="card-left">
              <div class="os-icon" :class="{ 'text-logo': getSessionIcon(conn).isText }">
                <Icon :icon="getSessionIcon(conn).icon" />
              </div>
              <div class="card-info">
                <div class="conn-name">{{ conn.label }}</div>
                <div v-if="conn.connProtocol === ConnProtocol.SSH" class="conn-info">
                  {{ conn.credential?.username }}@{{ conn.host }}
                </div>
                <div v-if="conn.connProtocol === ConnProtocol.Serial" class="conn-info">
                  {{ conn.serialPort }}
                </div>
              </div>
            </div>
            <NButton circle text size="small" class="edit-btn" @click.stop="handleEditConn(conn)">
              <template #icon>
                <Icon icon="ph:pencil-simple" class="edit-icon" />
              </template>
            </NButton>
          </div>

          <div class="card-footer">
            <div class="protocol-info">
              <Icon :icon="getProtocolIcon(conn)" />
              <span>{{ conn.connProtocol }}</span>
            </div>
            <div class="connection-tags" style="margin-left: auto">
              <NTooltip v-if="getConnCount(conn) > 0" trigger="hover">
                <template #trigger>
                  <div>
                    <NTag size="tiny" type="success">
                      {{ getConnCount(conn) }}
                    </NTag>
                  </div>
                </template>
                {{ $t('frontend.connection.connection.active') }}
              </NTooltip>
              <NTooltip v-if="getErrorConnCount(conn) > 0" trigger="hover">
                <template #trigger>
                  <div>
                    <NTag size="tiny" type="error">
                      {{ getErrorConnCount(conn) }}
                    </NTag>
                  </div>
                </template>
                {{ $t('frontend.connection.connection.disconnected') }}
              </NTooltip>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="empty-state">
        <NResult
          status="404"
          :title="$t('frontend.connection.empty.title')"
          :description="
            selectedGroup
              ? $t('frontend.connection.empty.group_desc', { name: selectedGroup.name })
              : $t('frontend.connection.empty.all_desc')
          "
        />
      </div>
    </div>

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
import { NBadge, NButton, NDropdown, NInput, NResult, NTag, NTooltip, useDialog, useThemeVars } from 'naive-ui';
import type { DropdownOption } from 'naive-ui';
import { computed, h, onMounted, onUnmounted, ref } from 'vue';
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

const groups = ref<Group[]>();
const conns = ref<Connection[]>();
const selectedGroup = ref<Group | null>(null);

const showDropdown = ref(false);
const dropdownX = ref(0);
const dropdownY = ref(0);
const dropdownOptions = ref<DropdownOption[]>([]);

const searchText = ref('');

const assetListCollapsed = ref(false);
const groupListCollapsed = ref(false);

const showConnModal = ref(false);
const showGroupModal = ref(false);
const isEditConn = ref(false);
const isEditGroup = ref(false);
const connectionId = ref<number>(0);
const editGroup = ref<Group | undefined>(undefined);

const currentContextNode = ref<any>(null);

const dialog = useDialog();

const updateDropdownOptions = (type: 'group' | 'conn') => {
  if (type === 'group') {
    dropdownOptions.value = [
      {
        label: t('frontend.connection.menu.edit_group'),
        key: 'edit-group',
        icon: () => h(Icon, { icon: 'ph:pencil-simple' }),
      },
      {
        label: t('frontend.connection.menu.delete_group'),
        key: 'delete-group',
        icon: () => h(Icon, { icon: 'ph:trash' }),
      },
    ];
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

const filteredConns = computed(() => {
  if (!selectedGroup.value) return conns.value || [];
  return conns.value?.filter(conn => conn.groupID === selectedGroup.value?.id) || [];
});

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
  const groupId = Number.parseInt(currentContextNode.value.key.replace('group-', ''));
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

const getConnCount = (conn: Connection) => {
  return connStore.connections.filter(c => c.host === conn.host && !c.errorCausedClosed).length;
};

const getErrorConnCount = (conn: Connection) => {
  return connStore.connections.filter(c => c.host === conn.host && c.errorCausedClosed).length;
};

const getProtocolIcon = (conn: Connection) => {
  const protocol = conn.connProtocol;
  switch (protocol) {
    case ConnProtocol.SSH:
      return 'ph:terminal-duotone';
    case ConnProtocol.Telnet:
      return 'ph:broadcast-duotone';
    case ConnProtocol.RDP:
      return 'ph:desktop-duotone';
    case ConnProtocol.VNC:
      return 'ph:monitor-duotone';
    case ConnProtocol.Serial:
      return 'ph:plug-duotone';
    default:
      return 'ph:gconn-duotone';
  }
};

const themeVars = useThemeVars();

const sidebarRef = ref<HTMLElement | null>(null);
const sidebarWidth = ref(Number(localStorage.getItem('sidebarWidth')) || 260);
const minWidth = 260;
const maxWidth = 380;
const isResizing = ref(false);

const startResize = (e: MouseEvent) => {
  e.preventDefault();
  isResizing.value = true;
  document.body.style.cursor = 'col-resize';
  document.body.style.userSelect = 'none';

  const startX = e.clientX;
  const startWidth = sidebarWidth.value;

  const handleMouseMove = (e: MouseEvent) => {
    if (!isResizing.value) return;

    const delta = e.clientX - startX;
    const newWidth = Math.min(Math.max(startWidth + delta, minWidth), maxWidth);
    sidebarWidth.value = newWidth;
    localStorage.setItem('sidebarWidth', newWidth.toString());
  };

  const handleMouseUp = () => {
    isResizing.value = false;
    document.body.style.cursor = '';
    document.body.style.userSelect = '';
    document.removeEventListener('mousemove', handleMouseMove);
    document.removeEventListener('mouseup', handleMouseUp);
  };

  document.addEventListener('mousemove', handleMouseMove);
  document.addEventListener('mouseup', handleMouseUp);
};

onMounted(async () => {
  await fetchData();

  window.addEventListener('sidebar-width-change', ((e: CustomEvent) => {
    sidebarWidth.value = e.detail;
  }) as EventListener);
});

onUnmounted(() => {
  window.removeEventListener('sidebar-width-change', ((e: CustomEvent) => {
    sidebarWidth.value = e.detail;
  }) as EventListener);
});

const filteredAssets = computed(() => {
  if (!searchText.value) return conns.value || [];
  return (
    conns.value?.filter(
      conn =>
        conn.label.toLowerCase().includes(searchText.value.toLowerCase()) ||
        conn.host?.toLowerCase().includes(searchText.value.toLowerCase()),
    ) || []
  );
});

const filteredGroups = computed(() => {
  if (!searchText.value) return groups.value || [];
  return groups.value?.filter(group => group.name.toLowerCase().includes(searchText.value.toLowerCase())) || [];
});

const getGroupConnCount = (group: Group) => {
  return conns.value?.filter(conn => conn.groupID === group.id).length || 0;
};

const handleSelectConn = (conn: Connection) => {
  toTerminal(conn);
};

const handleSelectGroup = (group: Group | null) => {
  if (selectedGroup.value?.id === group?.id) {
    selectedGroup.value = null;
  } else {
    selectedGroup.value = group;
  }
};

const handleConnContextMenu = (event: MouseEvent, conn: Connection) => {
  event.preventDefault();
  dropdownX.value = event.clientX;
  dropdownY.value = event.clientY;
  showDropdown.value = true;
  currentContextNode.value = { key: `conn-${conn.id}` };
  updateDropdownOptions('conn');
};

const handleGroupContextMenu = (event: MouseEvent, group: Group) => {
  event.preventDefault();
  dropdownX.value = event.clientX;
  dropdownY.value = event.clientY;
  showDropdown.value = true;
  currentContextNode.value = { key: `group-${group.id}` };
  updateDropdownOptions('group');
};

const handleDeleteGroup = async () => {
  if (!currentContextNode.value) return;
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
</script>

<style lang="less" scoped>
.page-container {
  height: 100%;
  display: flex;
  position: relative;
}

.sidebar {
  border-right: none;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.groups-list {
  flex: 1;
  display: flex;
  flex-direction: column;

  .list-header {
    padding: 6px;
    display: flex;
    align-items: center;
    border-bottom: 1px solid v-bind('themeVars.borderColor');

    .title {
      font-size: 13px;
      font-weight: 600;
      color: v-bind('themeVars.textColorBase');
    }

    .header-right {
      display: flex;
      align-items: center;
      gap: 4px;
      margin-left: auto;
    }
  }

  .list-content {
    padding: 2px 4px;
    height: calc(100vh - 80px);
    overflow-y: auto;
    scroll-behavior: smooth;
    position: relative;

    &::-webkit-scrollbar {
      display: none;
    }
    -ms-overflow-style: none;
    scrollbar-width: none;

    :deep(.n-tree) {
      .n-tree-node {
        &:hover {
          background: v-bind('`${themeVars.primaryColor}10`');
        }

        &.n-tree-node--selected {
          background: v-bind('`${themeVars.primaryColor}20`');

          .n-tree-node-content {
            color: v-bind('themeVars.primaryColor');
          }
        }

        .n-tree-node-content {
          .n-tree-node-content__text {
            border-bottom: none !important;
            text-decoration: none !important;
          }
        }
      }

      .n-tree__empty {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
      }
    }
  }
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  padding: 24px;
  position: relative;

  .content-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 24px;

    .header-left {
      display: flex;
      align-items: center;
      gap: 12px;

      h2 {
        font-size: 20px;
        font-weight: 600;
        color: v-bind('themeVars.textColorBase');
        margin: 0;
      }
    }
  }
}

.conns-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
  overflow-y: auto;
  padding-right: 4px;
}

.conn-card {
  background: v-bind('themeVars.cardColor');
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
  border: 1px solid v-bind('themeVars.borderColor');

  &:hover {
    border-color: v-bind('themeVars.primaryColor');

    .card-header {
      .edit-btn {
        opacity: 1;
      }
    }
  }

  .card-header {
    padding: 12px;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .card-left {
      display: flex;
      align-items: center;
      gap: 10px;
      flex: 1;
      min-width: 0;

      .os-icon {
        width: 38px;
        height: 38px;
        border-radius: 6px;
        background: v-bind('`${themeVars.primaryColor}20`');
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
        min-width: 0;

        .conn-name {
          font-size: 14px;
          font-weight: 600;
          color: v-bind('themeVars.textColorBase');
          margin-bottom: 2px;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }

        .conn-info {
          font-size: 12px;
          color: v-bind('themeVars.textColor3');
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }
      }
    }

    .edit-btn {
      opacity: 0;
      transition: all 0.2s ease;
      margin-left: 8px;
      background: v-bind('`${themeVars.primaryColor}20`');
      width: 32px;
      height: 32px;
      border-radius: 6px !important;

      :deep(.edit-icon) {
        font-size: 16px;
      }

      &:hover {
        color: v-bind('themeVars.primaryColor');
        background: v-bind('`${themeVars.primaryColor}30`');
      }
    }
  }

  .card-footer {
    padding: 8px 12px;
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: auto;
    border-top: 1px dashed v-bind('themeVars.borderColor');

    .protocol-info {
      display: flex;
      align-items: center;
      gap: 4px;
      font-size: 12px;
      color: v-bind('themeVars.textColor3');
      flex-shrink: 0;

      :deep(svg) {
        font-size: 16px;
        color: v-bind('themeVars.primaryColor');
      }
    }

    .connection-tags {
      display: flex;
      gap: 4px;
      align-items: center;
      margin-left: auto;
    }

    :deep(.n-tag) {
      padding: 0;
      width: 18px;
      height: 18px;
      line-height: 18px;
      text-align: center;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: calc(100% - 72px);
  width: 100%;
}

.resize-handle {
  width: 4px;
  cursor: col-resize;
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  z-index: 10;
  margin-left: -2px;
  border-right: 1px solid v-bind('themeVars.borderColor');
  background: transparent;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 8px 4px;

  .section-title {
    padding: 0;
    font-size: 12px;
    font-weight: 600;
    color: v-bind('themeVars.textColor3');
  }

  :deep(.n-button) {
    width: 16px;
    height: 16px;
    color: v-bind('themeVars.textColor3');

    .n-button__icon {
      font-size: 14px;
    }
  }
}

.asset-list,
.group-list {
  .asset-item,
  .group-item {
    display: flex;
    align-items: center;
    padding: 4px 8px;
    margin: 0 4px;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s ease;

    &:hover {
      background: v-bind('`${themeVars.primaryColor}10`');
    }

    .asset-icon,
    .group-icon {
      width: 20px;
      height: 20px;
      border-radius: 50%;
      background: v-bind('`${themeVars.primaryColor}20`');
      color: v-bind('themeVars.primaryColor');
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 16px;
      flex-shrink: 0;
    }

    .asset-info,
    .group-info {
      flex: 1;
      min-width: 0;
      margin-left: 6px;

      .asset-name,
      .group-name {
        font-size: 14px;
        color: v-bind('themeVars.textColorBase');
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
    }

    .group-count {
      font-size: 12px;
      color: v-bind('themeVars.textColor3');
      margin-left: 6px;
      min-width: 16px;
      text-align: right;
    }
  }

  .group-item {
    &.active {
      background: v-bind('`${themeVars.primaryColor}20`');

      .group-name {
        color: v-bind('themeVars.primaryColor');
      }
    }
  }
}
</style>
