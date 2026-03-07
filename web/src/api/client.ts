import type { ApiResponse, PublicConfig, UploadResponse } from '../types'

const API_BASE = '/api'

export async function fetchConfig(): Promise<PublicConfig> {
  const res = await fetch(`${API_BASE}/config`)
  const json: ApiResponse<PublicConfig> = await res.json()
  if (json.code !== 0) throw new Error(json.message)
  return json.data
}

export function uploadFiles(
  files: File[],
  onProgress?: (loaded: number, total: number) => void,
  expireHours?: number,
): Promise<UploadResponse> {
  return new Promise((resolve, reject) => {
    const formData = new FormData()
    files.forEach((file) => formData.append('file', file))
    if (expireHours && expireHours > 0) {
      formData.append('expire_hours', String(expireHours))
    }

    const xhr = new XMLHttpRequest()
    xhr.open('POST', `${API_BASE}/upload`)

    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable && onProgress) {
        onProgress(e.loaded, e.total)
      }
    }

    xhr.onload = () => {
      try {
        const json: ApiResponse<UploadResponse> = JSON.parse(xhr.responseText)
        if (json.code !== 0) {
          reject(new Error(json.message))
        } else {
          resolve(json.data)
        }
      } catch {
        reject(new Error('Failed to parse response'))
      }
    }

    xhr.onerror = () => reject(new Error('Upload failed'))
    xhr.send(formData)
  })
}

export async function deleteImage(id: string, token: string): Promise<void> {
  const res = await fetch(`${API_BASE}/image/${id}`, {
    method: 'DELETE',
    headers: { 'X-Delete-Token': token },
  })
  const json: ApiResponse<unknown> = await res.json()
  if (json.code !== 0) throw new Error(json.message)
}

export async function fetchRecent(limit = 20): Promise<ApiResponse<{ images: unknown[] }>> {
  const res = await fetch(`${API_BASE}/recent?limit=${limit}`)
  return res.json()
}
