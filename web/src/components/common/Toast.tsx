import { useStore } from '../../store/uploadStore'
import { X, CheckCircle, AlertCircle, Info } from 'lucide-react'

export function Toast() {
  const { toasts, removeToast } = useStore()

  if (toasts.length === 0) return null

  return (
    <div className="fixed top-4 right-4 z-[100] flex flex-col gap-2 max-w-sm">
      {toasts.map((toast) => {
        const Icon =
          toast.type === 'success'
            ? CheckCircle
            : toast.type === 'error'
              ? AlertCircle
              : Info
        const colors =
          toast.type === 'success'
            ? 'bg-emerald-50 dark:bg-emerald-950/50 border-emerald-200 dark:border-emerald-800 text-emerald-800 dark:text-emerald-200'
            : toast.type === 'error'
              ? 'bg-red-50 dark:bg-red-950/50 border-red-200 dark:border-red-800 text-red-800 dark:text-red-200'
              : 'bg-blue-50 dark:bg-blue-950/50 border-blue-200 dark:border-blue-800 text-blue-800 dark:text-blue-200'

        return (
          <div
            key={toast.id}
            className={`flex items-start gap-2 px-4 py-3 rounded-lg border shadow-lg animate-in slide-in-from-right ${colors}`}
          >
            <Icon size={18} className="mt-0.5 shrink-0" />
            <p className="text-sm flex-1">{toast.message}</p>
            <button
              onClick={() => removeToast(toast.id)}
              className="shrink-0 opacity-60 hover:opacity-100"
            >
              <X size={14} />
            </button>
          </div>
        )
      })}
    </div>
  )
}
