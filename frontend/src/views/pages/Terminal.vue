<template>
  <div class="xterm-container">
    <template v-for="conn in connectionStore.connections" :key="conn.id">
      <NResult
        v-if="conn.errorCausedClosed"
        :status="conn.errorCausedClosed ? 'error' : 'warning'"
        :title="conn.message"
        :class="{ 'terminal-hidden': isTerminalHidden(conn.id) }"
      >
        <template #footer>
          <NButton type="primary" @click="() => reconnect(conn.id)">{{ $t('frontend.terminal.reconnect') }}</NButton>
        </template>
        <template #default>
          <div v-if="conn.details">
            <NCollapse>
              <NCollapseItem :title="$t('frontend.terminal.error.details')" name="details">
                <NCode :code="conn.details" language="bash" :word-wrap="true" />
              </NCollapseItem>
            </NCollapse>
          </div>
        </template>
      </NResult>
      <NResult
        v-else-if="conn.isFingerprintConfirm"
        status="warning"
        :title="$t('frontend.terminal.fingerprint.confirm')"
        :description="$t('frontend.terminal.fingerprint.description')"
        :class="{ 'terminal-hidden': isTerminalHidden(conn.id) }"
      >
        <template #icon>
          <NIcon size="48">
            <Icon icon="ph:fingerprint" />
          </NIcon>
        </template>
        <template #footer>
          <NSpace>
            <NButton type="error" @click="() => rejectFingerprint(conn.id)">
              {{ $t('frontend.terminal.fingerprint.reject') }}
            </NButton>
            <NButton type="primary" @click="() => acceptFingerprint(conn.id)">
              {{ $t('frontend.terminal.fingerprint.accept') }}
            </NButton>
          </NSpace>
        </template>
        <template #default>
          <div class="fingerprint-details">
            <p>
              <strong>{{ $t('frontend.terminal.fingerprint.host') }}:</strong> {{ conn.hostAddress }}
            </p>
            <p>
              <strong>{{ $t('frontend.terminal.fingerprint.fingerprint') }}:</strong>
            </p>
            <NCode :code="conn.hostFingerprint || ''" language="bash" :word-wrap="true" />
          </div>
        </template>
      </NResult>
      <NResult
        v-else-if="conn.isConnecting || !connectedTerminals[conn.id]"
        status="info"
        :title="$t('frontend.terminal.connecting')"
        :description="$t('frontend.terminal.connecting_desc')"
        :class="{ 'terminal-hidden': isTerminalHidden(conn.id) }"
      >
        <template #icon>
          <NSpin size="large" />
        </template>
      </NResult>
      <div
        v-show="isTerminalVisible(conn)"
        :ref="el => el && (terminalRefs[conn.id] = el as HTMLElement)"
        :class="{ 'terminal-hidden': isTerminalHidden(conn.id) }"
        class="xterm-wrapper"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import '@xterm/xterm/css/xterm.css';
import { Icon } from '@iconify/vue';
import { TerminalType } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/enums';
import { WebsocketPort } from '@wailsApp/github.com/MisakaTAT/GTerm/backend/services/terminalsrv';
import { FitAddon } from '@xterm/addon-fit';
import { WebLinksAddon } from '@xterm/addon-web-links';
import { WebglAddon } from '@xterm/addon-webgl';
import { CanvasAddon } from '@xterm/addon-canvas';
import { Terminal } from '@xterm/xterm';
import { debounce } from 'lodash';
import { NButton, NCode, NCollapse, NCollapseItem, NIcon, NResult, NSpace, NSpin } from 'naive-ui';
import { onActivated, onMounted, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useConnectionStore } from '@/stores/connection';
import { loadTheme } from '@/themes/xtermjs';
import { getTranslated } from '@/utils/call';

defineOptions({ name: 'Terminal' });

const { t } = useI18n();

const connectionStore = useConnectionStore();
const connectionTabs = inject<any>('connectionTabs');
const activeConn = computed(() => connectionStore.activeConnection);
const terminalRefs = ref<Record<number, HTMLElement | null>>({});
const terminals = ref<Record<number, Terminal | undefined>>({});
const sockets = ref<Record<number, WebSocket | undefined>>({});
const fitAddons = ref<Record<number, FitAddon | undefined>>({});
const webLinksAddons = ref<Record<number, WebLinksAddon | undefined>>({});
const webglAddons = ref<Record<number, WebglAddon | undefined>>({});
const canvasAddons = ref<Record<number, CanvasAddon | undefined>>({});
const connectedTerminals = ref<Record<number, boolean>>({});

// 流控制
class FlowControl {
  private blocked = false;
  private pendingCallbacks = 0;
  private lowWatermark = 5;
  private highWatermark = 10;
  private bytesWritten = 0;
  private bytesThreshold = 1024 * 128; // 128KB

  constructor(
    private xterm: Terminal,
    private onPause: (paused: boolean) => void,
  ) {}

  async write(data: string): Promise<void> {
    if (this.blocked) {
      await new Promise<void>(resolve => {
        const checkUnblocked = () => {
          if (!this.blocked) {
            resolve();
          } else {
            setTimeout(checkUnblocked, 10);
          }
        };
        checkUnblocked();
      });
    }

    this.bytesWritten += data.length;

    if (this.bytesWritten > this.bytesThreshold) {
      this.pendingCallbacks++;
      this.bytesWritten = 0;

      if (!this.blocked && this.pendingCallbacks > this.highWatermark) {
        this.blocked = true;
        this.onPause(true);
      }

      this.xterm.write(data, () => {
        this.pendingCallbacks--;
        if (this.blocked && this.pendingCallbacks < this.lowWatermark) {
          this.blocked = false;
          this.onPause(false);
        }
      });
    } else {
      this.xterm.write(data);
    }
  }
}

// 流控制实例
const flowControls = ref<Record<number, FlowControl | undefined>>({});

const isTerminalHidden = (connId: number) => {
  return connId !== activeConn.value?.id;
};

const isTerminalVisible = (conn: any) => {
  return !conn.connectionError && !conn.isConnecting && !conn.isFingerprintConfirm && connectedTerminals.value[conn.id];
};

const updateStatus = (
  id: number,
  status: Partial<{
    isConnecting: boolean;
    errorCausedClosed: boolean;
    isFingerprintConfirm: boolean;
    message: string;
    details: string;
    hostAddress: string;
    hostFingerprint: string;
  }>,
) => {
  connectionStore.updateConnectionStatus(id, status);
  if (connectionTabs?.value) {
    let tabStatus = 'connecting';
    if (status.errorCausedClosed) {
      tabStatus = 'error';
    } else if (status.isFingerprintConfirm) {
      tabStatus = 'warning';
    } else if (!status.isConnecting && connectedTerminals.value[id]) {
      tabStatus = 'connected';
    }
    connectionTabs.value.updateTabStatus(id, tabStatus);
  }
};

const acceptFingerprint = (id: number) => {
  const conn = connectionStore.connections.find(c => c.id === id);
  if (!conn || !sockets.value[id]) return;

  sockets.value[id]?.send(
    JSON.stringify({
      type: TerminalType.TerminalTypeFingerprintConfirm,
      accept: true,
    }),
  );

  updateStatus(id, {
    isFingerprintConfirm: false,
    isConnecting: true,
  });
};

const rejectFingerprint = (id: number) => {
  const conn = connectionStore.connections.find(c => c.id === id);
  if (!conn || !sockets.value[id]) return;

  sockets.value[id]?.send(
    JSON.stringify({
      type: TerminalType.TerminalTypeFingerprintConfirm,
      accept: false,
    }),
  );

  updateStatus(id, {
    isFingerprintConfirm: false,
    isConnecting: true,
  });
};

const setXtermDomSize = (id: number) => {
  const terminalEl = terminalRefs.value[id];
  const fitAddon = fitAddons.value[id];
  if (fitAddon && connectedTerminals.value[id] && terminalEl) {
    const xtermElement = terminalEl.querySelector('.xterm') as HTMLElement;
    if (xtermElement) {
      const hight = terminalEl.clientHeight - 16;
      const weight = terminalEl.clientWidth - 16;
      xtermElement.style.height = `${hight}px`;
      xtermElement.style.width = `${weight}px`;
    }
    fitAddon.fit();
  }
};

const initializeTerminal = async (id: number) => {
  if (terminals.value[id]) return;

  const conn = connectionStore.connections.find(c => c.id === id);
  if (!conn) return;

  const theme = await loadTheme(conn.theme || 'Default');
  terminals.value[id] = new Terminal({
    convertEol: true,
    disableStdin: false,
    fontSize: 16,
    cursorBlink: true,
    cursorStyle: 'bar',
    theme,
    scrollback: 1000,
  });

  fitAddons.value[id] = new FitAddon();
  webLinksAddons.value[id] = new WebLinksAddon();
};

const initializeXterm = async (id: number) => {
  const terminalEl = terminalRefs.value[id] as HTMLElement;
  const terminal = terminals.value[id];
  const fitAddon = fitAddons.value[id];
  if (!terminalEl || !terminal || !fitAddon) return;

  fitAddon.activate(terminal);
  webLinksAddons.value[id]?.activate(terminal);

  // WebGL 优先，失败回退 Canvas
  try {
    const webglAddon = new WebglAddon();
    terminal.loadAddon(webglAddon);
    webglAddons.value[id] = webglAddon;
  } catch (error) {
    console.warn('WebGL addon failed, falling back to Canvas:', error);
    const canvasAddon = new CanvasAddon();
    terminal.loadAddon(canvasAddon);
    canvasAddons.value[id] = canvasAddon;
  }
  terminal.attachCustomKeyEventHandler(arg => {
    if (arg.code === 'PageUp' && arg.type === 'keydown') {
      terminal?.scrollPages(-1);
      return false;
    } else if (arg.code === 'PageDown' && arg.type === 'keydown') {
      terminal?.scrollPages(1);
      return false;
    }
    return true;
  });

  terminal.open(terminalEl);

  // 初始化流控制
  flowControls.value[id] = new FlowControl(terminal, (paused: boolean) => {
    const socket = sockets.value[id];
    if (socket?.readyState === WebSocket.OPEN) {
      socket.send(
        JSON.stringify({
          type: TerminalType.TerminalTypeFlowControl,
          pause: paused,
        }),
      );
    }
  });

  terminal.onData(data => sockets.value[id]?.send(JSON.stringify({ type: TerminalType.TerminalTypeCMD, cmd: data })));
  terminal.onResize(({ cols, rows }) => {
    if (sockets.value[id]?.readyState === WebSocket.OPEN) {
      sockets.value[id]?.send(JSON.stringify({ type: TerminalType.TerminalTypeResize, cols, rows }));
    }
  });
};

const initializeWebsocket = async (id: number, hostId: number) => {
  try {
    if (sockets.value[id]?.readyState === WebSocket.OPEN) {
      updateStatus(id, { isConnecting: false });
      return;
    }

    if (sockets.value[id]) {
      sockets.value[id].close();
      sockets.value[id] = undefined;
    }

    updateStatus(id, {
      isConnecting: true,
      errorCausedClosed: false,
      isFingerprintConfirm: false,
      message: '',
      details: '',
      hostAddress: '',
      hostFingerprint: '',
    });

    const port = await WebsocketPort();
    sockets.value[id] = new WebSocket(`ws://localhost:${port}/ws/terminal?hostId=${hostId}`);
    const socket = sockets.value[id];
    if (!socket) return;

    socket.onopen = () => {
      updateStatus(id, { isConnecting: false });
    };

    socket.onmessage = async (event: MessageEvent) => {
      const data = JSON.parse(event.data);
      switch (data.type) {
        case TerminalType.TerminalTypeError:
          updateStatus(id, {
            isConnecting: false,
            errorCausedClosed: true,
            message: getTranslated(data.code, data.message),
            details: data.details,
          });
          connectedTerminals.value[id] = false;
          socket?.close();
          break;
        case TerminalType.TerminalTypeFingerprintConfirm:
          updateStatus(id, {
            isConnecting: false,
            isFingerprintConfirm: true,
            hostAddress: data.host,
            hostFingerprint: data.fingerprint,
          });
          break;
        case TerminalType.TerminalTypeConnected:
          updateStatus(id, { isConnecting: false });
          connectedTerminals.value[id] = true;
          if (connectionTabs?.value) {
            connectionTabs.value.updateTabStatus(id, 'connected');
          }
          nextTick(() => {
            setXtermDomSize(id);
            terminals.value[id]?.focus();
            fitAddons.value[id]?.fit();
          });
          break;
        case TerminalType.TerminalTypeData:
          await handleTerminalData(id, data.content);
          break;
      }
    };

    socket.onerror = () => {
      updateStatus(id, {
        isConnecting: false,
        errorCausedClosed: true,
        message: t('frontend.terminal.error.connection'),
      });
      connectedTerminals.value[id] = false;
      socket?.close();
    };

    socket.onclose = event => {
      if (!connectionStore.connections.find(c => c.id === id)?.errorCausedClosed) {
        updateStatus(id, {
          isConnecting: false,
          errorCausedClosed: true,
          // TODO: reason 需要返回 message 原文。用于在 code 不存在的时候 fallback
          message: getTranslated(event.reason, t('frontend.terminal.error.disconnected')),
        });
      }
      connectedTerminals.value[id] = false;
      sockets.value[id] = undefined;
    };
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : t('frontend.terminal.error.failed');
    updateStatus(id, {
      isConnecting: false,
      errorCausedClosed: true,
      message: errorMessage,
    });
    sockets.value[id]?.close();
    sockets.value[id] = undefined;
  }
};

const reconnect = async (id: number) => {
  const conn = connectionStore.connections.find(c => c.id === id);
  if (!conn) return;

  updateStatus(id, {
    errorCausedClosed: false,
    isConnecting: true,
    message: '',
  });
  await initializeWebsocket(id, conn.connId);
};

const handleTerminalData = async (id: number, content: string) => {
  const flowControl = flowControls.value[id];
  if (!flowControl) return;

  await flowControl.write(content);
};

const closeTerminal = (id: number) => {
  sockets.value[id]?.close();
  sockets.value[id] = undefined;
  connectedTerminals.value[id] = false;

  if (terminals.value[id]) {
    terminals.value[id].onData(() => {});
    terminals.value[id].onResize(() => {});
    terminals.value[id].dispose();
    terminals.value[id] = undefined;
  }

  fitAddons.value[id]?.dispose();
  fitAddons.value[id] = undefined;

  // WebLinksAddon 无 dispose
  delete webLinksAddons.value[id];

  webglAddons.value[id]?.dispose();
  webglAddons.value[id] = undefined;

  canvasAddons.value[id]?.dispose();
  canvasAddons.value[id] = undefined;

  // 清理流控
  delete flowControls.value[id];

  updateStatus(id, {
    isConnecting: false,
    errorCausedClosed: false,
    message: '',
    details: '',
  });
};

const handleResize = debounce(() => {
  Object.keys(terminals.value).forEach(id => {
    setXtermDomSize(Number(id));
  });
}, 100);

onMounted(() => {
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  Object.keys(terminals.value).forEach(id => closeTerminal(Number(id)));
  window.removeEventListener('resize', handleResize);
});

watchEffect(async () => {
  const connections = connectionStore.connections;
  const activeConnId = activeConn.value?.id;

  if (!activeConnId) return;

  const conn = connections.find(c => c.id === activeConnId);
  if (conn && !terminals.value[activeConnId]) {
    await initializeTerminal(activeConnId);
    await initializeXterm(activeConnId);
    if (!sockets.value[activeConnId]) {
      await initializeWebsocket(activeConnId, conn.connId);
    }
  }
});

onActivated(() => {
  const activeId = activeConn.value?.id;
  if (activeId && terminals.value[activeId]) {
    nextTick(() => {
      terminals.value[activeId]?.focus();
    });
  }
});

defineExpose({ closeTerminal });
</script>

<style lang="less" scoped>
.xterm-container {
  height: 100%;
  position: relative;

  .terminal-hidden {
    position: absolute;
    opacity: 0;
    pointer-events: none;
    z-index: -1;
  }

  .xterm-wrapper {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
  }

  :deep(.xterm) {
    padding: 8px 16px 8px 8px;

    .xterm-screen {
      width: 100% !important;
      height: 100% !important;
    }
    canvas {
      width: 100% !important;
      height: 100% !important;
    }

    .xterm-viewport {
      &::-webkit-scrollbar {
        width: 8px;
      }
      &::-webkit-scrollbar-thumb {
        background: rgba(255, 255, 255, 0.2);
        border-radius: 4px;

        &:hover {
          background: rgba(255, 255, 255, 0.3);
        }
      }
    }
  }

  .n-result {
    margin: auto;
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }

  .n-button {
    min-width: 120px;
    margin-top: 16px;
  }

  .fingerprint-details {
    text-align: left;
  }
}
</style>
