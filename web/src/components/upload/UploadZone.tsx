import { useCallback, useRef, useState } from 'react'
import { Upload, ImagePlus, Clipboard } from 'lucide-react'
import { useUpload } from '../../hooks/useUpload'
import { useStore } from '../../store/uploadStore'
import { ExpirySelector } from './ExpirySelector'

export function UploadZone() {
  const [isDragging, setIsDragging] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const { upload } = useUpload()
  const { uploading, progress, t, config } = useStore()

  /** 过滤出图片文件并触发上传 */
  const handleFiles = useCallback(
    (files: FileList | File[]) => {
      const fileArray = Array.from(files).filter((f) => f.type.startsWith('image/'))
      if (fileArray.length > 0) {
        upload(fileArray)
      }
    },
    [upload],
  )

  /** 处理拖拽放置 */
  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault()
      setIsDragging(false)
      if (e.dataTransfer.files.length > 0) {
        handleFiles(e.dataTransfer.files)
      }
    },
    [handleFiles],
  )

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(true)
  }, [])

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)
  }, [])

  /** 处理粘贴事件，提取图片文件 */
  const handlePaste = useCallback(
    (e: React.ClipboardEvent) => {
      const items = e.clipboardData.items
      const files: File[] = []
      for (let i = 0; i < items.length; i++) {
        if (items[i].type.startsWith('image/')) {
          const file = items[i].getAsFile()
          if (file) files.push(file)
        }
      }
      if (files.length > 0) {
        e.preventDefault()
        handleFiles(files)
      }
    },
    [handleFiles],
  )

  const handleClick = () => {
    fileInputRef.current?.click()
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      handleFiles(e.target.files)
      e.target.value = ''
    }
  }

  return (
    <div
      className={`
        relative rounded-2xl border-2 border-dashed transition-all duration-200 cursor-pointer
        ${
          isDragging
            ? 'border-violet-500 bg-violet-50/50 dark:bg-violet-950/20 scale-[1.02]'
            : 'border-zinc-300 dark:border-zinc-700 hover:border-violet-400 dark:hover:border-violet-600 bg-white dark:bg-zinc-900'
        }
        ${uploading ? 'pointer-events-none opacity-70' : ''}
      `}
      onDrop={handleDrop}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onPaste={handlePaste}
      onClick={handleClick}
      tabIndex={0}
    >
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        multiple
        className="hidden"
        onChange={handleInputChange}
      />

      <div className="flex flex-col items-center justify-center py-16 px-4">
        {uploading ? (
          <>
            <div className="w-16 h-16 rounded-full bg-violet-100 dark:bg-violet-900/30 flex items-center justify-center mb-4">
              <Upload size={28} className="text-violet-600 dark:text-violet-400 animate-bounce" />
            </div>
            <p className="text-lg font-medium text-zinc-700 dark:text-zinc-300 mb-2">
              {t.uploading}
            </p>
            <div className="w-64 h-2 bg-zinc-200 dark:bg-zinc-700 rounded-full overflow-hidden">
              <div
                className="h-full bg-gradient-to-r from-violet-500 to-indigo-500 rounded-full transition-all duration-300"
                style={{ width: `${progress}%` }}
              />
            </div>
            <p className="text-sm text-zinc-500 mt-2">{progress}%</p>
          </>
        ) : (
          <>
            <div className="w-16 h-16 rounded-full bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center mb-4">
              <ImagePlus size={28} className="text-zinc-400 dark:text-zinc-500" />
            </div>
            <p className="text-lg font-medium text-zinc-700 dark:text-zinc-300 mb-1">
              {t.uploadDropText}
            </p>
            <div className="flex items-center gap-1.5 text-sm text-zinc-500 dark:text-zinc-500">
              <Clipboard size={14} />
              <span>{t.uploadPasteHint}</span>
            </div>
            <p className="text-xs text-zinc-400 dark:text-zinc-600 mt-3">
              {t.uploadFormatHint}
            </p>
            {config?.watermark_enabled && (
              <p className="text-xs text-amber-500 dark:text-amber-400 mt-1">
                {t.watermarkHint}
              </p>
            )}
          </>
        )}
      </div>

      {isDragging && (
        <div className="absolute inset-0 rounded-2xl bg-violet-500/10 flex items-center justify-center">
          <p className="text-violet-600 dark:text-violet-400 font-medium text-lg">
            {t.uploadReleaseHint}
          </p>
        </div>
      )}

      {/* 过期时间选择器 — 阻止点击冒泡，避免触发文件选择 */}
      {!uploading && (
        <div onClick={(e) => e.stopPropagation()}>
          <ExpirySelector />
        </div>
      )}
    </div>
  )
}
