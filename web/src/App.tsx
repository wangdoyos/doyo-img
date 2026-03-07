import { useEffect, useCallback } from 'react'
import { Header } from './components/layout/Header'
import { Footer } from './components/layout/Footer'
import { UploadZone } from './components/upload/UploadZone'
import { UploadResults } from './components/result/UploadResult'
import { Toast } from './components/common/Toast'
import { useStore } from './store/uploadStore'
import { useUpload } from './hooks/useUpload'

function App() {
  const { loadConfig, t } = useStore()
  const { upload } = useUpload()

  useEffect(() => {
    loadConfig()
  }, [loadConfig])

  /** 全局粘贴事件处理：监听 Ctrl+V 粘贴图片 */
  const handlePaste = useCallback(
    (e: ClipboardEvent) => {
      const items = e.clipboardData?.items
      if (!items) return

      const files: File[] = []
      for (let i = 0; i < items.length; i++) {
        if (items[i].type.startsWith('image/')) {
          const file = items[i].getAsFile()
          if (file) files.push(file)
        }
      }
      if (files.length > 0) {
        e.preventDefault()
        upload(files)
      }
    },
    [upload],
  )

  useEffect(() => {
    document.addEventListener('paste', handlePaste)
    return () => document.removeEventListener('paste', handlePaste)
  }, [handlePaste])

  return (
    <>
      <Header />
      <main className="flex-1 max-w-4xl w-full mx-auto px-4 py-8">
        <div className="text-center mb-8">
          <h2 className="text-2xl font-bold text-zinc-900 dark:text-zinc-100 mb-2">
            {t.pageTitle}
          </h2>
          <p className="text-zinc-500 dark:text-zinc-500">
            {t.pageSubtitle}
          </p>
        </div>
        <UploadZone />
        <UploadResults />
      </main>
      <Footer />
      <Toast />
    </>
  )
}

export default App
