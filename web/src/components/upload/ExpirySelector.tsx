import { useStore } from '../../store/uploadStore'

const EXPIRY_OPTIONS = [
  { value: 0, labelKey: 'never' as const },
  { value: 1, labelKey: '1h' as const },
  { value: 6, labelKey: '6h' as const },
  { value: 24, labelKey: '24h' as const },
  { value: 168, labelKey: '7d' as const },
  { value: 720, labelKey: '30d' as const },
]

export function ExpirySelector() {
  const { expireHours, setExpireHours, t } = useStore()

  const getLabel = (opt: (typeof EXPIRY_OPTIONS)[number]) => {
    if (opt.value === 0) return t.expiryNever
    if (opt.value < 24) return t.expiryHour(opt.value)
    return t.expiryDay(opt.value / 24)
  }

  return (
    <div className="flex items-center gap-2 mt-3">
      <span className="text-xs text-zinc-500 dark:text-zinc-500 shrink-0">
        {t.expiryLabel}
      </span>
      <div className="flex items-center gap-1 flex-wrap">
        {EXPIRY_OPTIONS.map((opt) => (
          <button
            key={opt.value}
            onClick={() => setExpireHours(opt.value)}
            className={`px-2.5 py-1 text-xs rounded-lg transition-colors ${
              expireHours === opt.value
                ? 'bg-violet-100 dark:bg-violet-900/40 text-violet-700 dark:text-violet-300 font-medium'
                : 'text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-800'
            }`}
          >
            {getLabel(opt)}
          </button>
        ))}
      </div>
    </div>
  )
}
