<template>
  <div class="file-transfer-container">
    <div class="integrated-header">
      <div class="left-section">
        <div class="connection-dropdown">
          <NSelect
            v-model:value="selectedConnectionId"
            :options="connectionOptions"
            :placeholder="t('frontend.file_transfer.select_connection')"
            :disabled="connecting"
            class="connection-select"
            @update:value="handleConnectionChange"
          />

          <NTooltip placement="bottom">
            <template #trigger>
              <NButton
                quaternary
                circle
                :disabled="!selectedConnectionId"
                :loading="connecting"
                class="connection-action-btn"
                @click="isConnected ? disconnectSftp() : connectSftp()"
              >
                <template #icon>
                  <NIcon>
                    <Icon :icon="isConnected ? 'ph:link-break-bold' : 'ph:link-bold'" />
                  </NIcon>
                </template>
              </NButton>
            </template>
            {{ isConnected ? t('frontend.file_transfer.disconnect') : t('frontend.file_transfer.connect') }}
          </NTooltip>
        </div>

        <NDivider v-if="isConnected" vertical style="height: 20px; margin: 0 8px" />

        <div v-if="isConnected" class="path-navigator">
          <div class="path-control">
            <NTooltip placement="bottom">
              <template #trigger>
                <NButton quaternary circle class="path-nav-btn" @click="navigateHome">
                  <template #icon>
                    <NIcon><Icon icon="ph:house-bold" /></NIcon>
                  </template>
                </NButton>
              </template>
              {{ t('frontend.file_transfer.navigate_home') }}
            </NTooltip>

            <NTooltip placement="bottom" :disabled="!canNavigateUp">
              <template #trigger>
                <NButton quaternary circle class="path-nav-btn" :disabled="!canNavigateUp" @click="navigateUp">
                  <template #icon>
                    <NIcon><Icon icon="ph:arrow-up-bold" /></NIcon>
                  </template>
                </NButton>
              </template>
              {{ t('frontend.file_transfer.navigate_up') }}
            </NTooltip>

            <NInput v-model:value="remotePathInput" class="path-input" @keydown.enter="navigateToPath">
              <template #prefix>
                <NIcon><Icon icon="ph:folder-open-bold" /></NIcon>
              </template>
            </NInput>

            <NTooltip placement="bottom">
              <template #trigger>
                <NButton quaternary class="path-nav-go-btn" @click="navigateToPath">
                  <template #icon>
                    <NIcon><Icon icon="ph:arrow-circle-right-bold" /></NIcon>
                  </template>
                  {{ t('frontend.file_transfer.goto') }}
                </NButton>
              </template>
              {{ t('frontend.file_transfer.navigate_to_path') }}
            </NTooltip>
          </div>
        </div>
      </div>

      <div v-if="isConnected" class="right-section">
        <NButtonGroup size="small">
          <NTooltip placement="bottom">
            <template #trigger>
              <NButton quaternary @click="refreshRemoteFiles">
                <template #icon>
                  <NIcon><Icon icon="ph:arrows-clockwise-bold" /></NIcon>
                </template>
              </NButton>
            </template>
            {{ t('frontend.file_transfer.refresh') }}
          </NTooltip>

          <NTooltip placement="bottom">
            <template #trigger>
              <NButton quaternary @click="uploadToRemote">
                <template #icon>
                  <NIcon><Icon icon="ph:arrow-up-bold" /></NIcon>
                </template>
              </NButton>
            </template>
            {{ t('frontend.file_transfer.upload') }}
          </NTooltip>

          <NTooltip placement="bottom" :disabled="!remoteSelectedFiles.length">
            <template #trigger>
              <NButton quaternary :disabled="!remoteSelectedFiles.length" @click="downloadFromRemote">
                <template #icon>
                  <NIcon><Icon icon="ph:arrow-down-bold" /></NIcon>
                </template>
              </NButton>
            </template>
            {{ t('frontend.file_transfer.download') }}
          </NTooltip>
        </NButtonGroup>

        <NTooltip placement="bottom">
          <template #trigger>
            <NPopover
              trigger="click"
              placement="bottom"
              :width="380"
              :show="showTasksPopover"
              @update:show="showTasksPopover = $event"
            >
              <template #trigger>
                <NButton quaternary circle @click="toggleTasksPopover">
                  <template #icon>
                    <NIcon>
                      <Icon icon="ph:arrows-down-up-bold" />
                    </NIcon>
                    <NBadge type="info" :value="activeTasks.length" :show="activeTasks.length > 0" :max="99" />
                  </template>
                </NButton>
              </template>

              <div class="tasks-popover-content">
                <div v-if="transferTasks.length === 0" class="no-tasks">
                  <span>{{ $t('frontend.file_transfer.no_tasks') }}</span>
                </div>
                <div v-else class="tasks-list">
                  <div
                    v-for="task in transferTasks"
                    :key="task.id"
                    class="task-item"
                    :class="{ 'is-completed': task.status === FileTransferTaskState.Completed }"
                  >
                    <div class="task-info">
                      <div class="task-title">
                        <NIcon class="task-type-icon">
                          <Icon :icon="task.isUpload ? 'ph:arrow-up-bold' : 'ph:arrow-down-bold'" />
                        </NIcon>
                        <span class="task-name">{{ getTaskFilename(task) }}</span>
                      </div>
                      <div class="task-details">
                        <span class="task-progress-text">
                          {{ formatFileSize(task.transferred) }} / {{ formatFileSize(task.size) }}
                        </span>
                        <span class="task-status-text">{{ getTaskStatusText(task) }}</span>
                      </div>
                    </div>
                    <NProgress
                      type="line"
                      :percentage="getTaskPercentage(task)"
                      :processing="task.status === FileTransferTaskState.Progress"
                      :status="getTaskProgressStatus(task)"
                      :show-indicator="false"
                      :height="4"
                      class="task-progress-bar"
                    />
                  </div>
                </div>
                <div class="tasks-actions">
                  <NButton v-if="hasCompletedTasks" size="small" @click="clearCompletedTasks">
                    {{ $t('frontend.file_transfer.clear_completed') }}
                  </NButton>
                  <NButton v-if="transferTasks.length > 0" size="small" @click="clearAllTasks">
                    {{ $t('frontend.file_transfer.clear_all') }}
                  </NButton>
                </div>
              </div>
            </NPopover>
          </template>
          {{ t('frontend.file_transfer.transfer_tasks') }}
        </NTooltip>
      </div>
    </div>

    <div v-if="!isConnected" class="empty-state">
      <div class="empty-content">
        <NIcon size="64" class="empty-icon">
          <Icon icon="ph:cloud-slash-bold" />
        </NIcon>
        <div class="empty-text">{{ $t('frontend.file_transfer.no_connection') }}</div>
      </div>
    </div>

    <div v-else class="file-browser">
      <div class="file-list-container">
        <div class="file-list-header">
          <div class="file-grid-cell checkbox-cell">
            <NCheckbox
              :checked="allRemoteSelected"
              :indeterminate="someRemoteSelected && !allRemoteSelected"
              @update:checked="handleCheckAll"
            />
          </div>
          <div class="file-grid-cell name-cell">{{ t('frontend.file_transfer.name') }}</div>
          <div class="file-grid-cell size-cell">{{ t('frontend.file_transfer.size') }}</div>
          <div class="file-grid-cell time-cell">{{ t('frontend.file_transfer.modified') }}</div>
          <div class="file-grid-cell owner-cell">{{ t('frontend.file_transfer.owner') }}</div>
          <div class="file-grid-cell group-cell">{{ t('frontend.file_transfer.group') }}</div>
          <div class="file-grid-cell perm-cell">{{ t('frontend.file_transfer.permissions') }}</div>
        </div>

        <div class="file-list-body">
          <NScrollbar>
            <div class="file-grid">
              <div
                v-for="file in remoteFiles"
                :key="file.name"
                class="file-grid-row"
                :class="{ 'is-selected': remoteSelectedFiles.includes(file.name) }"
                @click="handleRowClick(file)"
              >
                <div class="file-grid-cell checkbox-cell" @click.stop>
                  <NCheckbox
                    :checked="remoteSelectedFiles.includes(file.name)"
                    @update:checked="checked => handleFileSelect(file, checked)"
                  />
                </div>
                <div class="file-grid-cell name-cell">
                  <div class="file-name">
                    <NIcon class="file-icon" :class="{ 'is-folder': file.isDir }">
                      <Icon :icon="file.isDir ? 'ph:folder-simple-bold' : 'ph:file-bold'" />
                    </NIcon>
                    <span class="file-label">{{ file.name }}</span>
                  </div>
                </div>
                <div class="file-grid-cell size-cell">
                  {{ file.isDir ? '-' : formatFileSize(file.size) }}
                </div>
                <div class="file-grid-cell time-cell">
                  {{ formatDate(file.modTime) }}
                </div>
                <div class="file-grid-cell owner-cell">
                  {{ file.owner }}
                </div>
                <div class="file-grid-cell group-cell">
                  {{ file.group }}
                </div>
                <div class="file-grid-cell perm-cell">
                  {{ file.permissions }}
                </div>
              </div>
            </div>
          </NScrollbar>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { FileTransferTaskState } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/enums';
import { ListConnection } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/services/connectionsrv';
import {
  ConnectSFTP,
  DisconnectSFTP,
  DownloadFiles,
  ListRemoteFiles,
  SelectDownloadDirectory,
  SelectUploadFiles,
  UploadFiles,
} from '@wailsApp/github.com/MisakaTAT/GTerm/backend/services/filetransfersrv';
import type {
  FileList,
  FileTransferItemInfo,
  FileTransferTask,
} from '@wailsApp/github.com/MisakaTAT/GTerm/backend/types';
import { Events } from '@wailsio/runtime';
import {
  NBadge,
  NButton,
  NButtonGroup,
  NCheckbox,
  NDivider,
  NIcon,
  NInput,
  NPopover,
  NProgress,
  NScrollbar,
  NSelect,
  NTooltip,
  useThemeVars,
} from 'naive-ui';
import { v4 as uuidv4 } from 'uuid';
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useCall } from '@/utils/call';

const { t } = useI18n();
const { call } = useCall();
const themeVars = useThemeVars();

const selectedConnectionId = ref<number | null>(null);
const connecting = ref(false);
const isConnected = ref(false);

const remoteAbsolutePath = ref('');
const remotePathInput = ref('');
const remoteFiles = ref<FileTransferItemInfo[]>([]);
const remoteSelectedFiles = ref<string[]>([]);
const showTasksPopover = ref(false);

const transferTasks = ref<FileTransferTask[]>([]);
const connections = ref<any[]>([]);

const pathHistory = ref<string[]>([]);
const currentHistoryIndex = ref(-1);

const connectionOptions = computed(() => {
  if (!connections.value || !Array.isArray(connections.value)) {
    return [];
  }
  return connections.value.map(conn => ({
    label: conn.label,
    value: conn.id,
  }));
});

const fetchConns = async () => {
  const result = await call(ListConnection);
  connections.value = result.data || [];
};

const allRemoteSelected = computed(() => {
  return remoteFiles.value.length > 0 && remoteSelectedFiles.value.length === remoteFiles.value.length;
});

const someRemoteSelected = computed(() => {
  return remoteSelectedFiles.value.length > 0;
});

const activeTasks = computed(() => {
  return transferTasks.value.filter(
    task => task.status === FileTransferTaskState.Pending || task.status === FileTransferTaskState.Progress,
  );
});

const hasCompletedTasks = computed(() => {
  return transferTasks.value.some(
    task => task.status === FileTransferTaskState.Completed || FileTransferTaskState.Error,
  );
});

const toggleTasksPopover = () => {
  showTasksPopover.value = !showTasksPopover.value;
};

const getTaskFilename = (task: FileTransferTask): string => {
  const pathParts = task.source.split('/');
  return pathParts[pathParts.length - 1];
};

const getTaskPercentage = (task: FileTransferTask): number => {
  if (task.size === 0) return 0;
  return Math.min(Math.floor((task.transferred / task.size) * 100), 100);
};

const getTaskProgressStatus = (task: FileTransferTask): 'success' | 'error' | 'warning' | undefined => {
  switch (task.status) {
    case FileTransferTaskState.Completed:
      return 'success';
    case FileTransferTaskState.Error:
      return 'error';
    case FileTransferTaskState.Pending:
      return 'warning';
    default:
      return undefined;
  }
};

const getTaskStatusText = (task: FileTransferTask): string => {
  switch (task.status) {
    case FileTransferTaskState.Completed:
      return t('frontend.file_transfer.completed');
    case FileTransferTaskState.Error:
      return task.error || t('frontend.file_transfer.failed');
    case FileTransferTaskState.Pending:
      return t('frontend.file_transfer.pending');
    case FileTransferTaskState.Progress:
      return `${getTaskPercentage(task)}%`;
    default:
      return '';
  }
};

const addToHistory = (path: string) => {
  if (pathHistory.value[pathHistory.value.length - 1] !== path) {
    pathHistory.value.push(path);
    currentHistoryIndex.value = pathHistory.value.length - 1;
  }
};

const refreshRemoteFiles = async () => {
  const resp = await call(ListRemoteFiles, {
    args: [remoteAbsolutePath.value],
  });

  if (resp.ok) {
    const response = resp.data as FileList;
    remoteFiles.value = (response.files || []).filter(Boolean) as FileTransferItemInfo[];
    remoteAbsolutePath.value = response.absolutePath || '';
    remotePathInput.value = remoteAbsolutePath.value;
    remoteSelectedFiles.value = [];
  }
};

const navigateTo = async (path: string, addHistory = true) => {
  remoteAbsolutePath.value = path;

  if (addHistory) {
    addToHistory(path);
  }

  await refreshRemoteFiles();
};

const clearCompletedTasks = () => {
  transferTasks.value = transferTasks.value.filter(
    task => task.status === FileTransferTaskState.Pending || task.status === FileTransferTaskState.Progress,
  );
};

const clearAllTasks = () => {
  transferTasks.value = [];
};

const connectSftp = async () => {
  if (!selectedConnectionId.value) return;

  connecting.value = true;
  const resp = await call(ConnectSFTP, {
    args: [selectedConnectionId.value],
  });

  if (resp.ok) {
    isConnected.value = true;
    pathHistory.value = [];
    currentHistoryIndex.value = -1;
    await navigateTo('', false);
  }
  connecting.value = false;
};

const canNavigateUp = computed(() => {
  return remoteAbsolutePath.value !== '' && remoteAbsolutePath.value !== '/';
});

const navigateHome = () => {
  navigateTo('');
};

const navigateUp = () => {
  if (!canNavigateUp.value) return;

  const path = remoteAbsolutePath.value;
  const lastSlashIndex = path.lastIndexOf('/');

  if (lastSlashIndex <= 0) {
    navigateTo('/');
  } else {
    let parentPath = path.substring(0, lastSlashIndex);
    if (parentPath === '') {
      parentPath = '/';
    }
    navigateTo(parentPath);
  }
};

const navigateToPath = () => {
  if (!remotePathInput.value) return;
  navigateTo(remotePathInput.value);
};

const handleFileClick = (file: FileTransferItemInfo) => {
  if (file.isDir) {
    const path = remoteAbsolutePath.value;
    let newPath;

    if (path === '' || path === '/') {
      newPath = `/${file.name}`;
    } else {
      newPath = `${path}/${file.name}`;
    }
    navigateTo(newPath);
  }
};

const handleFileSelect = (file: FileTransferItemInfo, checked: boolean) => {
  if (checked) {
    if (!remoteSelectedFiles.value.includes(file.name)) {
      remoteSelectedFiles.value.push(file.name);
    }
  } else {
    const index = remoteSelectedFiles.value.indexOf(file.name);
    if (index !== -1) {
      remoteSelectedFiles.value.splice(index, 1);
    }
  }
};

const handleRowClick = (file: FileTransferItemInfo) => {
  if (file.isDir) {
    handleFileClick(file);
  } else {
    handleFileSelect(file, !remoteSelectedFiles.value.includes(file.name));
  }
};

const handleCheckAll = (checked: boolean) => {
  if (checked) {
    remoteSelectedFiles.value = remoteFiles.value.map(file => file.name);
  } else {
    remoteSelectedFiles.value = [];
  }
};

const disconnectSftp = async () => {
  const resp = await call(DisconnectSFTP);
  if (resp.ok) {
    isConnected.value = false;
    remoteFiles.value = [];
  }
};

const uploadToRemote = async () => {
  const filesResponse = await call(SelectUploadFiles, {
    args: [t('frontend.file_transfer.select_files_title')],
  });
  if (!filesResponse.ok) {
    transferTasks.value = transferTasks.value.map(task => {
      if (task.isUpload && task.status === FileTransferTaskState.Pending) {
        return { ...task, status: FileTransferTaskState.Error, error: filesResponse.msg || null };
      }
      return task;
    });

    return;
  }

  const sources = filesResponse.data || [];
  if (sources.length === 0) {
    return;
  }

  for (const source of sources) {
    const filename = source.split('/').pop() || t('frontend.file_transfer.unknown_file');
    const fileSize = 0;

    const task: FileTransferTask = {
      id: uuidv4(),
      source,
      destination: `${remoteAbsolutePath.value === '' ? '/' : remoteAbsolutePath.value}/${filename}`,
      size: fileSize,
      transferred: 0,
      isUpload: true,
      status: FileTransferTaskState.Pending,
      error: null,
    };

    transferTasks.value.push(task);
  }

  if (!showTasksPopover.value && sources.length > 0) {
    showTasksPopover.value = true;
  }

  const response = await call(UploadFiles, {
    args: [sources, remoteAbsolutePath.value],
  });

  if (response.ok) {
    await refreshRemoteFiles();
  } else {
    transferTasks.value = transferTasks.value.map(task => {
      if (task.isUpload && task.status === FileTransferTaskState.Pending) {
        return { ...task, status: FileTransferTaskState.Error, error: response.msg || null };
      }
      return task;
    });
  }
};

const downloadFromRemote = async () => {
  if (!remoteSelectedFiles.value.length) return;

  const files = remoteFiles.value.filter(f => remoteSelectedFiles.value.includes(f.name) && !f.isDir);
  const sources = files.map(f => `${remoteAbsolutePath.value === '' ? '/' : remoteAbsolutePath.value}/${f.name}`);

  const dirResponse = await call(SelectDownloadDirectory, {
    args: [t('frontend.file_transfer.select_directory_title')],
  });
  if (!dirResponse.ok) {
    transferTasks.value = transferTasks.value.map(task => {
      if (!task.isUpload && task.status === FileTransferTaskState.Pending) {
        return { ...task, status: FileTransferTaskState.Error, error: dirResponse.msg || null };
      }
      return task;
    });

    return;
  }

  const downloadPath = dirResponse.data;

  for (let i = 0; i < files.length; i++) {
    const file = files[i];
    const source = sources[i];
    const destination = `${downloadPath}/${file.name}`;

    const task: FileTransferTask = {
      id: uuidv4(),
      source,
      destination,
      size: file.size,
      transferred: 0,
      isUpload: false,
      status: FileTransferTaskState.Pending,
      error: null,
    };

    transferTasks.value.push(task);
  }

  if (!showTasksPopover.value && files.length > 0) {
    showTasksPopover.value = true;
  }

  const response = await call(DownloadFiles, {
    args: [sources, downloadPath],
  });

  if (!response.ok) {
    transferTasks.value = transferTasks.value.map(task => {
      if (!task.isUpload && task.status === FileTransferTaskState.Pending) {
        return { ...task, status: FileTransferTaskState.Error, error: response.msg || null };
      }
      return task;
    });
  }
};

const formatFileSize = (size: number): string => {
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`;
  if (size < 1024 * 1024 * 1024) return `${(size / 1024 / 1024).toFixed(2)} MB`;
  return `${(size / 1024 / 1024 / 1024).toFixed(2)} GB`;
};

const formatDate = (dateStr: string): string => {
  return new Date(dateStr).toLocaleString();
};

const handleTransferProgress = (data: any) => {
  const taskType = data.type === 'upload';
  const fileName = data.fileName || '';
  let task = transferTasks.value.find(t => t.isUpload === taskType && t.source.endsWith(fileName));
  if (!task && fileName) {
    const newTask: FileTransferTask = {
      id: uuidv4(),
      source: taskType ? fileName : `${remoteAbsolutePath.value}/${fileName}`,
      destination: taskType ? `${remoteAbsolutePath.value}/${fileName}` : fileName,
      size: data.total || 0,
      transferred: data.transferred || 0,
      isUpload: taskType,
      status: FileTransferTaskState.Progress,
      error: null,
    };

    transferTasks.value.push(newTask);
    task = newTask;
  }

  if (task) {
    task.status = FileTransferTaskState.Progress;
    task.transferred = data.transferred || 0;
    task.size = data.total || task.size;

    if (data.progress >= 100) {
      setTimeout(() => {
        task!.status = FileTransferTaskState.Completed;
      }, 500);
    }
  }
};

const handleConnectionChange = () => {
  if (isConnected.value) {
    disconnectSftp();
  }
};

onMounted(() => {
  fetchConns().then(() => {
    if (connections.value.length > 0) {
      selectedConnectionId.value = connections.value[0].id;
    }
  });

  Events.On('transfer:progress', handleTransferProgress);
});

onBeforeUnmount(() => {
  if (isConnected.value) {
    disconnectSftp();
  }

  Events.Off('transfer:progress');
});
</script>

<style lang="less" scoped>
.file-transfer-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 0;
  background-color: v-bind('themeVars.bodyColor');
}

.integrated-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px;
  background-color: v-bind('themeVars.cardColor');
  border-bottom: 1px solid v-bind('themeVars.borderColor');
  height: 48px;
  flex-shrink: 0;
}

.left-section {
  display: flex;
  align-items: center;
  flex: 1;
  overflow: hidden;
}

.right-section {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.connection-dropdown {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.connection-select {
  width: 180px;
}

.connection-action-btn {
  margin-left: 4px;
}

.path-navigator {
  flex: 1;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.path-control {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.path-nav-btn {
  margin: 0;
  flex-shrink: 0;
}

.path-input {
  flex: 1;
  min-width: 100px;
}

.path-nav-go-btn {
  margin: 0;
  flex-shrink: 0;
}

.tasks-popover-content {
  max-height: 300px;
  overflow-y: auto;
}

.no-tasks {
  padding: 16px;
  text-align: center;
  color: v-bind('themeVars.textColor3');
}

.tasks-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.task-item {
  background-color: v-bind('themeVars.tableColorStriped');
  border-radius: 4px;
  padding: 8px;

  &.is-completed {
    opacity: 0.8;
  }
}

.task-info {
  margin-bottom: 6px;
}

.task-title {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 4px;
}

.task-type-icon {
  font-size: 14px;
  color: v-bind('themeVars.primaryColor');
}

.task-name {
  font-size: 13px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.task-details {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: v-bind('themeVars.textColor3');
}

.task-progress-bar {
  margin-top: 2px;
}

.tasks-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 8px;
  border-top: 1px solid v-bind('themeVars.borderColor');
  margin-top: 8px;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: v-bind('themeVars.cardColor');
}

.empty-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.empty-icon {
  color: v-bind('themeVars.textColor3');
}

.empty-text {
  font-size: 16px;
  color: v-bind('themeVars.textColor3');
}

.file-browser {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: v-bind('themeVars.cardColor');
  height: 0;
  min-height: 0;
}

.file-list-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0;
}

.file-list-header {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  background-color: v-bind('themeVars.tableColorStriped');
  border-bottom: 1px solid v-bind('themeVars.borderColor');
  font-weight: 500;
  flex-shrink: 0;
}

.file-list-body {
  flex: 1;
  overflow: hidden;
  min-height: 0;
  position: relative;
}

.file-list-body .n-scrollbar {
  height: 100%;
}

.file-grid-row {
  display: flex;
  align-items: center;
  padding: 6px 16px;
  border-bottom: 1px solid v-bind('themeVars.borderColor');
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: v-bind('themeVars.tableColorHover');
  }

  &.is-selected {
    background-color: v-bind('themeVars.primaryColorSuppl');
  }
}

.file-grid-cell {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding: 0 4px;
}

.checkbox-cell {
  width: 40px;
  flex-shrink: 0;
}

.name-cell {
  flex: 3;
  min-width: 200px;
}

.size-cell {
  flex: 1;
  min-width: 100px;
}

.time-cell {
  flex: 2;
  min-width: 180px;
}

.owner-cell {
  flex: 1;
  min-width: 100px;
}

.group-cell {
  flex: 1;
  min-width: 100px;
}

.perm-cell {
  flex: 1;
  min-width: 100px;
}

.file-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-icon {
  font-size: 18px;

  &.is-folder {
    color: v-bind('themeVars.warningColor');
  }
}

.file-label {
  font-weight: 500;
}
</style>
