import { useState } from 'react'
import { Copy, Check, Trash2, ChevronDown, Clock } from 'lucide-react'
import type { UploadResult, LinkFormat } from '../../types'
import { generateLink, formatFileSize } from '../../utils/format'
import { useCopy } from '../../hooks/useCopy'
import { useStore } from '../../store/uploadStore'
import { deleteImage } from '../../api/client'

/** 链接格式选项，label 不需要翻译（URL/Markdown/HTML/BBCode 是通用名称） */
const FORMATS: { key: LinkFormat; label: string }[] = [
  { key: 'url', label: 'URL' },
  { key: 'markdown', label: 'Markdown' },
  { key: 'html', label: 'HTML' },
  { key: 'bbcode', label: 'BBCode' },
]

export function UploadResultCard({ result }: { result: UploadResult }) {
  const [format, setFormat] = useState<LinkFormat>('url')
  const { copy, copied } = useCopy()
  const { removeResult, addToast, t } = useStore()

  const linkText = generateLink(result.url, format, result.name)

  /** 复制链接到剪贴板 */
  const handleCopy = (e: React.MouseEvent) => {
    e.stopPropagation()
    copy(linkText)
    addToast(t.copiedToClipboard, 'success')
  }

  /** 使用 delete token 删除图片 */
  const handleDelete = async (e: React.MouseEvent) => {
    e.stopPropagation()
    try {
      await deleteImage(result.id, result.delete_token)
      removeResult(result.id)
      addToast(t.imageDeleted, 'info')
    } catch {
      addToast(t.deleteImageFailed, 'error')
    }
  }

  return (
    <div className="group rounded-xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 overflow-hidden shadow-sm hover:shadow-md transition-shadow">
      {/* 图片预览 */}
      <div className="relative aspect-video bg-zinc-100 dark:bg-zinc-800 overflow-hidden">
        <img
          src={result.thumbnail_url || result.url}
          alt={result.name}
          className="w-full h-full object-contain"
          loading="lazy"
        />
        <button
          onClick={handleDelete}
          className="absolute top-2 right-2 p-1.5 rounded-lg bg-black/50 text-white opacity-0 group-hover:opacity-100 hover:bg-red-500 transition-all"
          title={t.deleteImage}
        >
          <Trash2 size={14} />
        </button>
      </div>

      <div className="p-3">
        {/* 图片信息 */}
        <div className="flex items-center justify-between mb-2">
          <p className="text-sm font-medium text-zinc-700 dark:text-zinc-300 truncate max-w-[60%]" title={result.name}>
            {result.name}
          </p>
          <div className="flex items-center gap-2 text-xs text-zinc-500">
            {result.width > 0 && (
              <span>{result.width}x{result.height}</span>
            )}
            <span>{formatFileSize(result.size)}</span>
          </div>
        </div>

        {/* 过期时间信息 */}
        {result.expires_at && (
          <ExpiryBadge expiresAt={result.expires_at} />
        )}

        {/* 链接格式选择器 */}
        <div className="flex items-center gap-1 mb-2">
          {FORMATS.map((f) => (
            <button
              key={f.key}
              onClick={() => setFormat(f.key)}
              className={`px-2 py-0.5 text-xs rounded-md transition-colors ${
                format === f.key
                  ? 'bg-violet-100 dark:bg-violet-900/40 text-violet-700 dark:text-violet-300 font-medium'
                  : 'text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-800'
              }`}
            >
              {f.label}
            </button>
          ))}
        </div>

        {/* 链接展示 + 复制按钮 */}
        <div className="flex items-center gap-2">
          <input
            type="text"
            readOnly
            value={linkText}
            className="flex-1 text-xs bg-zinc-50 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-lg px-3 py-2 text-zinc-600 dark:text-zinc-400 font-mono select-all outline-none focus:border-violet-400"
            onClick={(e) => (e.target as HTMLInputElement).select()}
          />
          <button
            onClick={handleCopy}
            className={`shrink-0 px-3 py-2 rounded-lg text-xs font-medium transition-all ${
              copied
                ? 'bg-emerald-100 dark:bg-emerald-900/40 text-emerald-700 dark:text-emerald-300'
                : 'bg-violet-600 hover:bg-violet-700 text-white'
            }`}
          >
            {copied ? <Check size={14} /> : <Copy size={14} />}
          </button>
        </div>
      </div>
    </div>
  )
}

export function UploadResults() {
  const { results, clearResults, t } = useStore()
  const [expanded, setExpanded] = useState(true)

  if (results.length === 0) return null

  return (
    <div className="mt-6">
      <div className="flex items-center justify-between mb-3">
        <button
          onClick={() => setExpanded(!expanded)}
          className="flex items-center gap-1.5 text-sm font-medium text-zinc-700 dark:text-zinc-300"
        >
          <ChevronDown
            size={16}
            className={`transition-transform ${expanded ? '' : '-rotate-90'}`}
          />
          {t.uploadResults} ({results.length})
        </button>
        <button
          onClick={clearResults}
          className="text-xs text-zinc-500 hover:text-red-500 transition-colors"
        >
          {t.clearAll}
        </button>
      </div>

      {expanded && (
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          {results.map((result) => (
            <UploadResultCard key={result.id} result={result} />
          ))}
        </div>
      )}
    </div>
  )
}

/** 过期时间徽章组件 */
function ExpiryBadge({ expiresAt }: { expiresAt: string }) {
  const { t } = useStore()
  const expiry = new Date(expiresAt)
  const now = new Date()
  const isExpired = now > expiry

  if (isExpired) {
    return (
      <div className="flex items-center gap-1 mb-2 text-xs text-red-500">
        <Clock size={12} />
        <span>{t.expired}</span>
      </div>
    )
  }

  const diffMs = expiry.getTime() - now.getTime()
  const diffH = Math.floor(diffMs / (1000 * 60 * 60))
  const diffD = Math.floor(diffH / 24)
  const remaining = diffD > 0 ? t.expiryDay(diffD) : t.expiryHour(Math.max(1, diffH))

  return (
    <div className="flex items-center gap-1 mb-2 text-xs text-amber-500 dark:text-amber-400">
      <Clock size={12} />
      <span>{t.expiresIn} {remaining}</span>
    </div>
  )
}
