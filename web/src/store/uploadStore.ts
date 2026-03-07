import { create } from 'zustand'
import type { PublicConfig, UploadResult } from '../types'
import type { Locale } from '../i18n/messages'
import { getMessages } from '../i18n/messages'
import { fetchConfig } from '../api/client'

type Theme = 'light' | 'dark' | 'system'

const HISTORY_KEY = 'doyo-history'
const MAX_HISTORY = 100

/** 从 localStorage 加载上传历史记录 */
function loadHistory(): UploadResult[] {
  try {
    const raw = localStorage.getItem(HISTORY_KEY)
    if (!raw) return []
    const parsed = JSON.parse(raw)
    if (Array.isArray(parsed)) return parsed
    return []
  } catch {
    return []
  }
}

/** 将上传历史保存到 localStorage，最多保留 MAX_HISTORY 条 */
function saveHistory(results: UploadResult[]) {
  try {
    const trimmed = results.slice(0, MAX_HISTORY)
    localStorage.setItem(HISTORY_KEY, JSON.stringify(trimmed))
  } catch {
    // localStorage 满或不可用，静默忽略
  }
}

/** 检测浏览器默认语言，返回 zh 或 en */
function detectLocale(): Locale {
  const saved = localStorage.getItem('doyo-locale') as Locale | null
  if (saved === 'zh' || saved === 'en') return saved
  const lang = navigator.language || ''
  return lang.startsWith('zh') ? 'zh' : 'en'
}

interface UploadStore {
  // 服务端公开配置
  config: PublicConfig | null
  loadConfig: () => Promise<void>

  // 上传结果（持久化到 localStorage）
  results: UploadResult[]
  addResults: (results: UploadResult[]) => void
  removeResult: (id: string) => void
  clearResults: () => void

  // 上传状态
  uploading: boolean
  progress: number
  setUploading: (uploading: boolean) => void
  setProgress: (progress: number) => void

  // 过期时间（小时）
  expireHours: number
  setExpireHours: (hours: number) => void

  // 主题
  theme: Theme
  setTheme: (theme: Theme) => void

  // 多语言
  locale: Locale
  setLocale: (locale: Locale) => void
  t: ReturnType<typeof getMessages>

  // 消息提示
  toasts: { id: string; message: string; type: 'success' | 'error' | 'info' }[]
  addToast: (message: string, type: 'success' | 'error' | 'info') => void
  removeToast: (id: string) => void
}

const initialLocale = detectLocale()

export const useStore = create<UploadStore>((set, get) => ({
  config: null,
  loadConfig: async () => {
    try {
      const config = await fetchConfig()
      set({ config })
    } catch (err) {
      console.error('加载配置失败:', err)
    }
  },

  results: loadHistory(),
  addResults: (results) =>
    set((state) => {
      // 按 id 去重，避免重复记录
      const existingIds = new Set(state.results.map((r) => r.id))
      const newResults = results.filter((r) => !existingIds.has(r.id))
      const updated = [...newResults, ...state.results]
      saveHistory(updated)
      return { results: updated }
    }),
  removeResult: (id) =>
    set((state) => {
      const updated = state.results.filter((r) => r.id !== id)
      saveHistory(updated)
      return { results: updated }
    }),
  clearResults: () => {
    saveHistory([])
    set({ results: [] })
  },

  uploading: false,
  progress: 0,
  setUploading: (uploading) => set({ uploading }),
  setProgress: (progress) => set({ progress }),

  expireHours: 0,
  setExpireHours: (hours) => set({ expireHours: hours }),

  theme: (localStorage.getItem('doyo-theme') as Theme) || 'system',
  setTheme: (theme) => {
    localStorage.setItem('doyo-theme', theme)
    set({ theme })
  },

  locale: initialLocale,
  t: getMessages(initialLocale),
  setLocale: (locale: Locale) => {
    localStorage.setItem('doyo-locale', locale)
    set({ locale, t: getMessages(locale) })
  },

  toasts: [],
  addToast: (message, type) => {
    const id = Date.now().toString(36) + Math.random().toString(36).slice(2)
    set((state) => ({ toasts: [...state.toasts, { id, message, type }] }))
    setTimeout(() => get().removeToast(id), 3000)
  },
  removeToast: (id) =>
    set((state) => ({ toasts: state.toasts.filter((t) => t.id !== id) })),
}))
