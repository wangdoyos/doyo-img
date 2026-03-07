import { useEffect } from 'react'
import { useStore } from '../store/uploadStore'

export function useTheme() {
  const { theme, setTheme } = useStore()

  useEffect(() => {
    const root = document.documentElement

    const applyTheme = (t: 'light' | 'dark') => {
      if (t === 'dark') {
        root.classList.add('dark')
      } else {
        root.classList.remove('dark')
      }
    }

    if (theme === 'system') {
      const mq = window.matchMedia('(prefers-color-scheme: dark)')
      applyTheme(mq.matches ? 'dark' : 'light')
      const handler = (e: MediaQueryListEvent) => applyTheme(e.matches ? 'dark' : 'light')
      mq.addEventListener('change', handler)
      return () => mq.removeEventListener('change', handler)
    } else {
      applyTheme(theme)
    }
  }, [theme])

  return { theme, setTheme }
}
