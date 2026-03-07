export interface ImageMeta {
  id: string
  name: string
  format: string
  mime_type: string
  size: number
  width: number
  height: number
  storage_path: string
  delete_token: string
  created_at: string
}

export interface UploadResult {
  id: string
  name: string
  url: string
  thumbnail_url?: string
  size: number
  width: number
  height: number
  format: string
  delete_token: string
  created_at: string
  expires_at?: string
}

export interface ApiResponse<T> {
  code: number
  data: T
  message: string
}

export interface UploadResponse {
  images: UploadResult[]
  errors?: string[]
}

export interface RecentResponse {
  images: ImageMeta[]
}

export interface PublicConfig {
  max_file_size: number
  max_batch_size: number
  allowed_formats: string[]
  compress_enabled: boolean
  base_url: string
  watermark_enabled: boolean
  default_expire_hours: number
  max_expire_days: number
}

export type LinkFormat = 'url' | 'markdown' | 'html' | 'bbcode'

export interface UploadItem {
  id: string
  file: File
  progress: number
  status: 'pending' | 'uploading' | 'success' | 'error'
  result?: UploadResult
  error?: string
}
