import { Github } from 'lucide-react'
import { useStore } from '../../store/uploadStore'

export function Footer() {
  const { t } = useStore()

  return (
    <footer className="w-full border-t border-zinc-200 dark:border-zinc-800 py-4 mt-auto">
      <div className="max-w-4xl mx-auto px-4 flex items-center justify-between text-sm text-zinc-500 dark:text-zinc-500">
        <span>{t.footerText}</span>
        <a
          href="https://github.com/doyo-img/doyo-img"
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center gap-1 hover:text-zinc-700 dark:hover:text-zinc-300 transition-colors"
        >
          <Github size={16} />
          <span>GitHub</span>
        </a>
      </div>
    </footer>
  )
}
