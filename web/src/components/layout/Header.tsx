import { Sun, Moon, Monitor, Languages } from 'lucide-react'
import { useTheme } from '../../hooks/useTheme'
import { useStore } from '../../store/uploadStore'

export function Header() {
  const { theme, setTheme } = useTheme()
  const { locale, setLocale } = useStore()

  /** 循环切换主题：light -> dark -> system */
  const nextTheme = () => {
    const order: Array<'light' | 'dark' | 'system'> = ['light', 'dark', 'system']
    const idx = order.indexOf(theme)
    setTheme(order[(idx + 1) % order.length])
  }

  /** 切换语言：zh <-> en */
  const toggleLocale = () => {
    setLocale(locale === 'zh' ? 'en' : 'zh')
  }

  const ThemeIcon = theme === 'dark' ? Moon : theme === 'light' ? Sun : Monitor

  return (
    <header className="w-full border-b border-zinc-200 dark:border-zinc-800 bg-white/80 dark:bg-zinc-950/80 backdrop-blur-sm sticky top-0 z-50">
      <div className="max-w-4xl mx-auto px-4 h-14 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-violet-500 to-indigo-600 flex items-center justify-center">
            <span className="text-white font-bold text-sm">D</span>
          </div>
          <h1 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            doyo-img
          </h1>
        </div>
        <div className="flex items-center gap-1">
          {/* 语言切换按钮 */}
          <button
            onClick={toggleLocale}
            className="flex items-center gap-1 px-2 py-1.5 rounded-lg hover:bg-zinc-100 dark:hover:bg-zinc-800 text-zinc-600 dark:text-zinc-400 transition-colors text-sm"
            title={locale === 'zh' ? 'Switch to English' : '切换到中文'}
          >
            <Languages size={16} />
            <span className="text-xs font-medium">{locale === 'zh' ? '中' : 'EN'}</span>
          </button>
          {/* 主题切换按钮 */}
          <button
            onClick={nextTheme}
            className="p-2 rounded-lg hover:bg-zinc-100 dark:hover:bg-zinc-800 text-zinc-600 dark:text-zinc-400 transition-colors"
            title={`${theme}`}
          >
            <ThemeIcon size={18} />
          </button>
        </div>
      </div>
    </header>
  )
}
