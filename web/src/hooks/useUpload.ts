import { useCallback } from 'react'
import { uploadFiles } from '../api/client'
import { useStore } from '../store/uploadStore'

export function useUpload() {
  const { setUploading, setProgress, addResults, addToast, config, t, expireHours } = useStore()

  const upload = useCallback(
    async (files: File[]) => {
      if (files.length === 0) return

      // 客户端预校验
      if (config) {
        const maxSize = config.max_file_size
        const maxBatch = config.max_batch_size
        const allowed = config.allowed_formats

        if (files.length > maxBatch) {
          addToast(t.tooManyFiles(maxBatch), 'error')
          return
        }

        for (const file of files) {
          if (file.size > maxSize) {
            addToast(t.fileTooLarge(file.name, Math.round(maxSize / 1024 / 1024)), 'error')
            return
          }
          const ext = file.name.split('.').pop()?.toLowerCase() || ''
          const isImage = file.type.startsWith('image/') || allowed.includes(ext)
          if (!isImage) {
            addToast(t.formatNotAllowed(file.name), 'error')
            return
          }
        }
      }

      setUploading(true)
      setProgress(0)

      try {
        const response = await uploadFiles(files, (loaded, total) => {
          setProgress(Math.round((loaded / total) * 100))
        }, expireHours)

        if (response.images && response.images.length > 0) {
          addResults(response.images)
          addToast(t.uploadSuccess(response.images.length), 'success')
        }

        if (response.errors && response.errors.length > 0) {
          response.errors.forEach((err) => addToast(err, 'error'))
        }
      } catch (err) {
        addToast(err instanceof Error ? err.message : t.uploadFailed, 'error')
      } finally {
        setUploading(false)
        setProgress(0)
      }
    },
    [config, t, expireHours, setUploading, setProgress, addResults, addToast],
  )

  return { upload }
}
