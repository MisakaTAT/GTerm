import type { CancellablePromise } from '@wailsio/runtime/types/cancellable';
import { useMessage } from 'naive-ui';
import { i18n } from './i18n';

interface BackendResp {
  ok: boolean;
  msg?: string;
  data?: any;
  code?: string;
}

interface CallResult<T = any> {
  ok: boolean;
  msg?: string;
  data?: T;
}

export function getTranslated(code: string | undefined, fallback: string | undefined): string {
  code = `backend.${code}`;
  const { t } = i18n.global;
  if (!code) return fallback || '';
  const translated = t(code);
  return translated === code ? fallback || '' : translated;
}

export async function callBackendFunction<T = any>(
  backendFunction: (...args: any[]) => CancellablePromise<BackendResp | null>,
  ...args: any[]
): Promise<CallResult<T>> {
  try {
    const resp = await backendFunction(...args);
    return handleCallResp<T>(resp);
  } catch (error: any) {
    return {
      ok: false,
      msg: error.message,
    };
  }
}

export function handleCallResp<T = any>(resp: BackendResp | null): CallResult<T> {
  if (!resp) {
    return {
      ok: false,
      msg: 'No data returned',
    };
  }
  if (resp.ok) {
    return {
      ok: true,
      data: resp.data as T,
      msg: getTranslated(resp.code || '', resp?.msg),
    };
  }
  return {
    ok: false,
    msg: getTranslated(resp.code || '', resp?.msg),
  };
}

export function useCall() {
  const message = useMessage();

  async function call<T = any>(
    backendFunction: (...args: any[]) => CancellablePromise<BackendResp | null>,
    options: {
      args?: any[];
    } = {},
  ): Promise<CallResult<T>> {
    const { args = [] } = options;
    const result = await callBackendFunction<T>(backendFunction, ...args);

    if (result.msg) {
      result.ok ? message.success(result.msg) : message.error(result.msg);
    }

    return result;
  }

  return {
    call,
  };
}
